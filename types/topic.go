package types

type SubscribeTopicReq struct {
	Name string `json:"name"`
}

type GetMessageList struct {
	BasePage
	Name string `form:"name" json:"name"`
}

type DeleteMessage struct {
	Id uint `json:"id"`
}

type MessageResp struct {
	Info            string `json:"info" gorm:"size:1000"`
	Title           string `json:"title" gorm:"comment:商品标题"`
	ImgPath         string `json:"img_path" gorm:"comment:商品图片"`
	ProductName     string `json:"product_name" gorm:"comment:商品名字"`
	ProductCategory string `json:"product_category" gorm:"comment:商品分类"`
}
