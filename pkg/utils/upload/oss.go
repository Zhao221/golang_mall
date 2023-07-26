package upload

import (
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"golang_mall/conf"
	"mime/multipart"
	"path/filepath"
	"time"
)

// UploadToQiNiu 封装上传图片到七牛云然后返回状态和图片的url，单张
func UploadToQiNiu(file multipart.File, fileSize int64, name string) (path string, err error) {
	qConfig := conf.Config.Oss
	var AccessKey = qConfig.AccessKeyId
	var SerectKey = qConfig.AccessKeySecret
	var Bucket = qConfig.BucketName
	var ImgUrl = qConfig.QiNiuServer
	putPlicy := storage.PutPolicy{
		Scope: Bucket,
	}
	mac := qbox.NewMac(AccessKey, SerectKey)
	upToken := putPlicy.UploadToken(mac)
	cfg := storage.Config{
		Zone:          &storage.ZoneHuanan,
		UseCdnDomains: false,
		UseHTTPS:      false,
	}
	putExtra := storage.PutExtra{}
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	uniqueFileName := generateUniqueFileName( name)
	err = formUploader.Put(context.Background(), &ret, upToken, uniqueFileName, file, fileSize, &putExtra)
	if err != nil {
		return "", err
	}
	url := ImgUrl + "/" + ret.Key
	return url, err
}

func UpToQiNiu(file multipart.File, fileSize int64) (path string, err error) {
	qConfig := conf.Config.Oss
	var AccessKey = qConfig.AccessKeyId
	var SerectKey = qConfig.AccessKeySecret
	var Bucket = qConfig.BucketName
	var ImgUrl = qConfig.QiNiuServer
	putPlicy := storage.PutPolicy{
		Scope: Bucket,
	}
	mac := qbox.NewMac(AccessKey, SerectKey)
	upToken := putPlicy.UploadToken(mac)
	cfg := storage.Config{
		Zone:          &storage.ZoneHuanan,
		UseCdnDomains: false,
		UseHTTPS:      false,
	}
	putExtra := storage.PutExtra{}
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	err = formUploader.PutWithoutKey(context.Background(), &ret, upToken, file, fileSize, &putExtra)
	if err != nil {
		return "", err
	}
	url := ImgUrl + "/" + ret.Key
	return url, nil
}

func generateUniqueFileName(originalName string) string {
	// 使用当前时间戳和原始文件名生成唯一文件名
	timestamp := time.Now().Unix()
	ext := filepath.Ext(originalName)
	nameWithoutExt := originalName[:len(originalName)-len(ext)]
	uniqueName := fmt.Sprintf("%s_%d%s", nameWithoutExt, timestamp, ext)
	return uniqueName
}