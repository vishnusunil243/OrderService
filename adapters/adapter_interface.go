package adapters

import (
	helperstruct "github.com/vishnusunil243/OrderService/helperStruct"
)

type AdapterInterface interface {
	OrderAll(items []helperstruct.OrderAll, userId uint) (int, error)
	UserCancelOrder(orderId int) error
	ChangeOrderStatus(orderId, orderStatusId int) error
	GetAllOrdersUser(userId int) ([]helperstruct.GetAllOrder, error)
	GetAllOrders() ([]helperstruct.GetAllOrder, error)
	GetOrder(orderId int) (helperstruct.GetAllOrder, error)
}
