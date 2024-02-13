package adapters

import (
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
	for _, item := range items {
		insertProductItems := `INSERT INTO order_items (product_id,quantity,total,order_id) VALUES ($1,$2,$3,$4)`
		if err := tx.Exec(insertProductItems, item.ProductId, item.Quantity, item.Total, orderId).Error; err != nil {
			tx.Rollback()
			return 1, err
		}
		updateTotal := `UPDATE orders SET total=total+$1 WHERE orderId=$2`
		if err := tx.Exec(updateTotal, item.Total, orderId).Error; err != nil {
			tx.Rollback()
			return 1, err
		}
	}
	return orderId, nil
}
