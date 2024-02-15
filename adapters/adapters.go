package adapters

import (
	"fmt"

	"github.com/vishnusunil243/OrderService/entities"
	helperstruct "github.com/vishnusunil243/OrderService/helperStruct"
	"gorm.io/gorm"
)

type OrderAdapter struct {
	DB *gorm.DB
}

func NewOrderAdapter(db *gorm.DB) *OrderAdapter {
	return &OrderAdapter{
		DB: db,
	}
}
func (order *OrderAdapter) OrderAll(items []helperstruct.OrderAll, userId uint) (int, error) {
	var orderId int
	tx := order.DB.Begin()
	query := `INSERT INTO orders (user_id,payment_type_id,total) VALUES ($1,$2,0) RETURNING id`
	if err := tx.Raw(query, userId, 1).Scan(&orderId).Error; err != nil {
		tx.Rollback()
		return 1, err
	}
	if orderId == 0 {
		return -1, fmt.Errorf("order not found")
	}
	for _, item := range items {
		insertProductItems := `INSERT INTO order_items (product_id,quantity,total,order_id) VALUES ($1,$2,$3,$4)`
		if err := tx.Exec(insertProductItems, item.ProductId, item.Quantity, item.Total, orderId).Error; err != nil {
			tx.Rollback()
			return 1, err
		}
		updateTotal := `UPDATE orders SET total=total+$1 WHERE id=$2`
		if err := tx.Exec(updateTotal, item.Total, orderId).Error; err != nil {
			tx.Rollback()
			return 1, err
		}
	}
	if err := tx.Commit().Error; err != nil {
		return -1, fmt.Errorf("error running transactions")
	}
	return orderId, nil
}
func (order *OrderAdapter) UserCancelOrder(orderId int) error {
	tx := order.DB.Begin()
	deleteOrderItems := `DELETE FROM order_items WHERE order_id=?`
	if err := tx.Exec(deleteOrderItems, orderId).Error; err != nil {
		return err
	}
	deleteOrder := `UPDATE orders SET order_status_id=$1 WHERE id=$2`
	if err := tx.Exec(deleteOrder, 5, orderId).Error; err != nil {
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}
func (order *OrderAdapter) ChangeOrderStatus(orderId, orderStatusId int) error {
	updateOrderStatus := `UPDATE orders SET order_status_id=$1 WHERE id=$2`
	if err := order.DB.Exec(updateOrderStatus, orderStatusId, orderId).Error; err != nil {
		return err
	}
	return nil
}
func (order *OrderAdapter) GetAllOrdersUser(userId int) ([]helperstruct.GetAllOrder, error) {
	tx := order.DB.Begin()
	var orders []entities.Order
	var res []helperstruct.GetAllOrder
	getAllOrder := `SELECT * FROM orders WHERE user_id=?`
	if err := tx.Raw(getAllOrder, userId).Scan(&orders).Error; err != nil {
		return []helperstruct.GetAllOrder{}, err
	}
	for _, order := range orders {
		var orderItems []entities.OrderItems
		getAllOrderItems := `SELECT * from order_items WHERE order_id=?`
		if err := tx.Raw(getAllOrderItems, order.Id).Scan(&orderItems).Error; err != nil {
			return []helperstruct.GetAllOrder{}, err
		}
		response := helperstruct.GetAllOrder{
			OrderId:       int(order.Id),
			AddressId:     int(order.AddressId),
			PaymentTypeId: int(order.PaymentTypeId),
			OrderStatusId: int(order.OrderStatusId),
			OrderItems:    orderItems,
		}
		res = append(res, response)

	}
	if err := tx.Commit().Error; err != nil {
		return []helperstruct.GetAllOrder{}, err
	}
	return res, nil
}
func (order *OrderAdapter) GetAllOrders() ([]helperstruct.GetAllOrder, error) {
	var res []helperstruct.GetAllOrder
	tx := order.DB.Begin()
	var orders []entities.Order
	orderQuery := `SELECT * FROM orders`
	if err := tx.Raw(orderQuery).Scan(&orders).Error; err != nil {
		tx.Rollback()
		return []helperstruct.GetAllOrder{}, err
	}
	for _, order := range orders {
		var orderItems []entities.OrderItems
		orderItemsQuery := `SELECT * FROM order_items WHERE order_id=?`
		if err := tx.Raw(orderItemsQuery, order.Id).Scan(&orderItems).Error; err != nil {
			tx.Rollback()
			return []helperstruct.GetAllOrder{}, err
		}
		response := helperstruct.GetAllOrder{
			OrderId:       int(order.Id),
			AddressId:     int(order.AddressId),
			PaymentTypeId: int(order.PaymentTypeId),
			OrderStatusId: int(order.OrderStatusId),
			OrderItems:    orderItems,
		}
		res = append(res, response)
	}
	if err := tx.Commit().Error; err != nil {
		return []helperstruct.GetAllOrder{}, err
	}
	return res, nil
}
func (order *OrderAdapter) GetOrder(orderId int) (helperstruct.GetAllOrder, error) {
	tx := order.DB.Begin()
	var orderData entities.Order
	orderQuery := `SELECT * FROM orders WHERE id=?`
	if err := tx.Raw(orderQuery, orderId).Scan(&orderData).Error; err != nil {
		tx.Rollback()
		return helperstruct.GetAllOrder{}, err
	}
	var orderItems []entities.OrderItems
	orderItemQuery := `SELECT * FROM order_items WHERE order_id=?`
	if err := tx.Raw(orderItemQuery, orderId).Scan(&orderItems).Error; err != nil {
		return helperstruct.GetAllOrder{}, err
	}
	res := helperstruct.GetAllOrder{
		OrderId:       int(orderData.Id),
		AddressId:     int(orderData.AddressId),
		PaymentTypeId: int(orderData.PaymentTypeId),
		OrderStatusId: int(orderData.OrderStatusId),
		OrderItems:    orderItems,
	}
	return res, nil
}
