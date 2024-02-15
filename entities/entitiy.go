package entities

type Order struct {
	Id            uint `gorm:"primaryKey"`
	UserId        uint
	PaymentTypeId uint
	AddressId     uint
	OrderStatusId uint
	OrderStatus   `gorm:"ForeignKey:OrderStatusId"`
	Total         float64
}
type OrderItems struct {
	Id        uint
	OrderId   uint
	Order     `gorm:"ForeignKey:OrderId"`
	ProductId uint
	Quantity  int
	Total     float64
}
type OrderStatus struct {
	Id     uint
	Status string
}
