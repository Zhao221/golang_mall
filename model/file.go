package model

import "github.com/jinzhu/gorm"

type FileUpload struct {
	gorm.Model
	Size int64 `json:"size"`
	UserId    uint   `json:"user_id" gorm:"comment:用户id"` 	// 上传到阿里云后阿里云的视频id ，我们有这个id，可以后续做删除功能
	Name      string `json:"name" gorm:"comment:文件名"`           // 文件名
	Url       string `json:"url" gorm:"comment:文件地址;size:256" ` // 文件地址
	Tag       string `json:"tag" gorm:"comment:文件标签"`           // 文件标签
	Key       string `json:"key" gorm:"comment:编号"`             // 编号
	LocalPath string `json:"local_path" form:"local_path" gorm:"comment:用户上传文件的本地路径"`
}
