package service

import (
	"context"
	"fmt"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/vishnusunil243/OrderService/adapters"
	helperstruct "github.com/vishnusunil243/OrderService/helperStruct"
	"github.com/vishnusunil243/proto-files/pb"
)

var (
	Tracer     opentracing.Tracer
	CartClient pb.CartServiceClient
)

func RetrieveTracer(tr opentracing.Tracer) {
	Tracer = tr
}

type OrderService struct {
	Adapter adapters.AdapterInterface
	pb.UnimplementedOrderServiceServer
}

func NewOrderService(adapter adapters.AdapterInterface) *OrderService {
	return &OrderService{
		Adapter: adapter,
	}
}
func (order *OrderService) OrderAll(req *pb.OrderRequest) (*pb.OrderResponse, error) {
	span := Tracer.StartSpan("orderAll grpc")
	defer span.Finish()
	cartItems, err := CartClient.GetAllCartItems(context.TODO(), &pb.UserCartCreate{
		UserId: req.UserId,
	})
	if err != nil {
		return &pb.OrderResponse{}, fmt.Errorf("unable to get cart items")
	}
	var cart []helperstruct.OrderAll
	for {
		items, err := cartItems.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return &pb.OrderResponse{}, err
		}
		item := helperstruct.OrderAll{
			ProductId: uint(items.ProductId),
			Quantity:  float64(items.Quantity),
			Total:     uint(items.Total),
		}
		cart = append(cart, item)
	}
	orderId, err := order.Adapter.OrderAll(cart, uint(req.UserId))
	if err != nil {
		return &pb.OrderResponse{}, err
	}
	return &pb.OrderResponse{OrderId: uint32(orderId)}, nil

}
