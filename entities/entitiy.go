package entities

type Order struct {
	Id            uint `gorm:"primaryKey"`
	UserId        uint
	PaymentTypeId uint
	AddressId     uint
	Total         float64
}
type OrderItems struct {
	Id        uint
	OrderId   uint
	ProductId uint
	quantity  int
	total     float64
}
