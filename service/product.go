package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	errors2 "github.com/pkg/errors"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"golang_mall/conf"
	"golang_mall/consts"
	"golang_mall/dao/mysql"
	"golang_mall/global"
	"golang_mall/model"
	"golang_mall/pkg/utils/ctl"
	"golang_mall/pkg/utils/upload"
	"golang_mall/repository/kafka"
	"golang_mall/types"
	"log"
	"mime/multipart"
	"strconv"
	"sync"
	"time"
)

var ProductSrvIns *ProductSrv
var ProductSrvOnce sync.Once
var ctx context.Context

type ProductSrv struct{}

func GetProductSrv() *ProductSrv {
	ProductSrvOnce.Do(func() {
		ProductSrvIns = &ProductSrv{}
	})
	return ProductSrvIns
}

func (s *ProductSrv) ProductCreate(c context.Context, files []*multipart.FileHeader, req types.ProductCreateReq) (err error) {
	u, err := ctl.GetUserInfo(c)
	if err != nil {
		return err
	}
	uId := u.Id
	boss, _ := mysql.NewUserDao(c).GetUserById(uId)
	// 以第一张图作为封面图
	tmp, _ := files[0].Open()
	var path string
	path, err = upload.UploadToQiNiu(tmp, files[0].Size, files[0].Filename)
	onSale, _ := strconv.ParseBool(req.OnSale)
	if err != nil {
		return err
	}
	product := &model.Product{
		Name:          req.Name,
		CategoryID:    req.CategoryID,
		Title:         req.Title,
		Info:          req.Info,
		ImgPath:       path,
		Price:         req.Price,
		DiscountPrice: req.DiscountPrice,
		Num:           req.Num,
		OnSale:        onSale,
		BossID:        uId,
		BossName:      boss.UserName,
		BossAvatar:    boss.Avatar,
	}
	productDao := mysql.NewProductDao(c)
	err = productDao.CreateProduct(product)
	if err != nil {
		return err
	}
	// 需要上传大量图片且希望加快上传速度时
	// 可以考虑使用 WaitGroup 和 goroutine 控制并发上传。
	// 设置最大并发任务数量
	maxConcurrentTasks := 5
	semaphore := make(chan struct{}, maxConcurrentTasks)
	wg := new(sync.WaitGroup)
	wg.Add(len(files))
	for index, file := range files {
		key := index
		value := file
		// 从信号量中获取一个空位，如果信号量已满，将会阻塞等待
		semaphore <- struct{}{}
		go func() {
			defer wg.Done()                // 当上传完成时，通知 WaitGroup 已完成一个任务
			defer func() { <-semaphore }() // 释放信号量空位，使其他任务可以执行
			num := strconv.Itoa(key)
			tmp, _ = value.Open()
			path, err = upload.UploadToQiNiu(tmp, value.Size, value.Filename+num)
			if err != nil {
				return
			}
			productImg := &model.ProductImg{
				ProductID: product.ID,
				ImgPath:   path,
				Name:      value.Filename,
			}
			err = mysql.NewProductImgDao(c).CreateProductImg(productImg)
			if err != nil {
				return
			}
		}()
	}
	wg.Wait()
	// 将商品信息存入kafka
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = mysql.NewProductDao(c).NewProduct(product)
		if err != nil {
			return
		}
	}()
	wg.Wait()
	return err
}

func (s *ProductSrv) ProductUpdate(c context.Context, req types.ProductUpdateReq) (err error) {
	product := &model.Product{
		Name:          req.Name,
		CategoryID:    req.CategoryID,
		Title:         req.Title,
		Info:          req.Info,
		Price:         req.Price,
		DiscountPrice: req.DiscountPrice,
		OnSale:        req.OnSale,
	}
	err = mysql.NewProductDao(c).UpdateProduct(req.ID, product)
	if err != nil {
		return err
	}
	return err
}

func (s *ProductSrv) ProductDelete(ctx context.Context, req *types.ProductDeleteReq) (err error) {
	u, _ := ctl.GetUserInfo(ctx)
	err = mysql.NewProductDao(ctx).DeleteProduct(req.ID, u.Id)
	if err != nil {
		return err
	}
	return err
}

func (s *ProductSrv) ProductList(c context.Context, req types.ProductListReq) (resp interface{}, err error) {
	var total int64
	condition := make(map[string]interface{})
	if req.CategoryID != 0 {
		condition["category_id"] = req.CategoryID
	}
	productDao := mysql.NewProductDao(c)
	products, _ := productDao.ListProductByCondition(condition, req.BasePage)
	total, err = productDao.CountProductByCondition(condition)
	if err != nil {
		return nil, err
	}
	pRespList := make([]*types.ProductResp, 0)
	for _, p := range products {
		pResp := &types.ProductResp{
			ID:            p.ID,
			Name:          p.Name,
			CategoryID:    p.CategoryID,
			Title:         p.Title,
			Info:          p.Info,
			ImgPath:       p.ImgPath,
			Price:         p.Price,
			DiscountPrice: p.DiscountPrice,
			View:          p.View(),
			CreatedAt:     p.CreatedAt.Unix(),
			Num:           p.Num,
			OnSale:        p.OnSale,
			BossID:        p.BossID,
			BossName:      p.BossName,
			BossAvatar:    p.BossAvatar,
		}
		var filename string
		mysql.NewUserDao(c).Model(&model.ProductImg{}).Select("name").Where("id = ?", p.ID).First(&filename)
		pResp.BossAvatar = GetDownloadURL(conf.Config.Oss.AccessKeyId, conf.Config.Oss.AccessKeySecret, conf.Config.Oss.BucketName, filename)
		pRespList = append(pRespList, pResp)
	}

	resp = &types.DataListResp{
		Item:  pRespList,
		Total: total,
	}

	return resp, err
}

func GetDownloadURL(accessKey, secretKey, bucket, fileName string) string {
	// 初始化 Mac 对象
	mac := qbox.NewMac(accessKey, secretKey)
	// 初始化七牛云存储配置
	cfg := storage.Config{}
	cfg.UseHTTPS = true

	// 创建 BucketManager 对象
	bucketManager := storage.NewBucketManager(mac, &cfg)

	// 获取文件的下载地址
	domain, err := bucketManager.ListBucketDomains(bucket)
	if err != nil {
		return fmt.Sprintf("%s", errors2.New("获取空间域名失败"))
	}
	deadline := time.Now().Add(time.Second * 3600).Unix() // 设置 URL 有效期为 1 小时
	privateURL := storage.MakePrivateURL(mac, domain[0].Domain, fileName, deadline)
	return privateURL
}

func (s *ProductSrv) ProductShow(c context.Context, req types.ProductShowReq) (resp interface{}, err error) {
	p, err := mysql.NewProductDao(c).ShowProductById(req.ID)
	if err != nil {
		return nil, err
	}
	pResp := &types.ProductResp{
		ID:            p.ID,
		Name:          p.Name,
		CategoryID:    p.CategoryID,
		Title:         p.Title,
		Info:          p.Info,
		ImgPath:       p.ImgPath,
		Price:         p.Price,
		DiscountPrice: p.DiscountPrice,
		View:          p.View(),
		CreatedAt:     p.CreatedAt.Unix(),
		Num:           p.Num,
		OnSale:        p.OnSale,
		BossID:        p.BossID,
		BossName:      p.BossName,
		BossAvatar:    p.BossAvatar,
	}
	var filename string
	global.GVA_DB.Table("file_upload").Select("name").Where("user_id", p.BossID).First(&filename)
	pResp.BossAvatar = GetDownloadURL(conf.Config.Oss.AccessKeyId, conf.Config.Oss.AccessKeySecret, conf.Config.Oss.BucketName, filename)
	pResp.ImgPath = GetDownloadURL(conf.Config.Oss.AccessKeyId, conf.Config.Oss.AccessKeySecret, conf.Config.Oss.BucketName, p.Name)
	resp = pResp
	return resp, err
}

func (s *ProductSrv) ProductSearch(c context.Context, req types.ProductSearchReq) (resp interface{}, err error) {
	products, count, err := mysql.NewProductDao(c).SearchProduct(req.Info, req.BasePage)
	if err != nil {
		return nil, err
	}
	pRespList := make([]*types.ProductResp, 0)
	for _, p := range products {
		pResp := &types.ProductResp{
			ID:            p.ID,
			Name:          p.Name,
			CategoryID:    p.CategoryID,
			Title:         p.Title,
			Info:          p.Info,
			ImgPath:       p.ImgPath,
			Price:         p.Price,
			DiscountPrice: p.DiscountPrice,
			View:          p.View(),
			CreatedAt:     p.CreatedAt.Unix(),
			Num:           p.Num,
			OnSale:        p.OnSale,
			BossID:        p.BossID,
			BossName:      p.BossName,
			BossAvatar:    p.BossAvatar,
		}
		var filename string
		global.GVA_DB.Table("file_upload").Select("name").Where("user_id", p.BossID).First(&filename)
		pResp.BossAvatar = GetDownloadURL(conf.Config.Oss.AccessKeyId, conf.Config.Oss.AccessKeySecret, conf.Config.Oss.BucketName, filename)
		pResp.ImgPath = GetDownloadURL(conf.Config.Oss.AccessKeyId, conf.Config.Oss.AccessKeySecret, conf.Config.Oss.BucketName, p.Name)
		pRespList = append(pRespList, pResp)
	}
	resp = &types.DataListResp{
		Item:  pRespList,
		Total: count,
	}
	return resp, err
}

func (s *ProductSrv) ProductImgList(c context.Context, req types.ListProductImgReq) (resp interface{}, err error) {
	productImgs, _ := mysql.NewProductImgDao(c).ListProductImgByProductId(req.ID)
	var filename string
	global.GVA_DB.Table("product_img").Select("name").Where("id", req.ID).First(&filename)
	for i := range productImgs {
		if conf.Config.System.UploadModel == consts.UploadModelLocal {
			productImgs[i].ImgPath = GetDownloadURL(conf.Config.Oss.AccessKeyId, conf.Config.Oss.AccessKeySecret, conf.Config.Oss.BucketName, filename)
		}
	}

	resp = &types.DataListResp{
		Item:  productImgs,
		Total: int64(len(productImgs)),
	}

	return resp, err
}

func Consume() error {
	// 初始化 Kafka 消费者
	brokers := []string{"localhost:9092"}
	consumer, err := sarama.NewConsumer(brokers, kafka.Kcfg)
	partitionConsumer, err := consumer.ConsumePartition("newProduct", 0, sarama.OffsetNewest)
	if err != nil {
		log.Printf("Error consuming partition: %v", err)
		return err
	}
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			var product model.Product
			err = json.Unmarshal(msg.Value, &product)
			if err != nil {
				log.Printf("Error unmarshaling product: %v", err)
				return err
			} else {
				fmt.Printf("New product: %+v\n", product)
				var Ids []uint
				var category string
				mysql.NewCategoryDao(ctx).Model(&model.Category{}).Select("category_name").Where("id=?", product.CategoryID).First(&category)
				mysql.NewTopicDao(ctx).Table("topic").Select("user_id").Where("name=?", "newProduct").Find(&Ids)
				var messages []model.Message
				var wg sync.WaitGroup
				wg.Add(len(Ids))
				for _, value := range Ids {
					go func() {
						defer wg.Done()
						message := model.Message{
							UserId:          value,
							ProductId:       product.ID,
							Info:            product.Info,
							Title:           product.Title,
							ImgPath:         product.ImgPath,
							TopicName:       "newProduct",
							ProductName:     product.Name,
							ProductCategory: category,
						}
						messages = append(messages, message)
					}()
				}
				wg.Wait()
				mysql.NewMessageDao(ctx).Table("message").Create(&messages)
			}
		case err = <-partitionConsumer.Errors():
			log.Printf("Error consuming message: %v", err)
			return err
		}
	}
}
