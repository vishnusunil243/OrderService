package initializer

import (
	"github.com/vishnusunil243/OrderService/adapters"
	"github.com/vishnusunil243/OrderService/service"
	"gorm.io/gorm"
)

func Initializer(db *gorm.DB) *service.OrderService {
	adapter := adapters.NewOrderAdapter(db)
	service := service.NewOrderService(adapter)
	return service
}
