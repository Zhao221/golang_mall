package types

type PaymentDownReq struct {
	OrderId   uint    `json:"order_id"`
	Num       uint    `json:"num"`
	ProductID uint    `json:"product_id"`
	BossID    uint    `json:"boss_id"`
	Money     float64 `json:"money"`
	Key       string  `json:"key"`
}
