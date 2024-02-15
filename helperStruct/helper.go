package helperstruct

import "github.com/vishnusunil243/OrderService/entities"

type OrderAll struct {
	ProductId uint
	Quantity  float64
	Total     uint
}
type GetAllOrder struct {
	OrderId       int
	AddressId     int
	PaymentTypeId int
	OrderStatusId int
	OrderItems    []entities.OrderItems
}
