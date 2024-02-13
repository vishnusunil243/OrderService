package adapters

import (
	helperstruct "github.com/vishnusunil243/OrderService/helperStruct"
)

type AdapterInterface interface {
	OrderAll(items []helperstruct.OrderAll, userId uint) (int, error)
}
