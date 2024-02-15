package service

import (
	"context"
	"fmt"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/vishnusunil243/OrderService/adapters"
	helperstruct "github.com/vishnusunil243/OrderService/helperStruct"
	"github.com/vishnusunil243/proto-files/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
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
func (order *OrderService) OrderAll(ctx context.Context, req *pb.OrderRequest) (*pb.OrderResponse, error) {
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
	if len(cart) == 0 {
		return &pb.OrderResponse{}, fmt.Errorf("cart is empty please add products to the cart to complete order")
	}
	if _, err := CartClient.TruncateCart(context.TODO(), &pb.UserCartCreate{UserId: req.UserId}); err != nil {
		return &pb.OrderResponse{}, err
	}
	orderId, err := order.Adapter.OrderAll(cart, uint(req.UserId))
	if err != nil {
		return &pb.OrderResponse{}, err
	}
	return &pb.OrderResponse{OrderId: uint32(orderId)}, nil

}
func (order *OrderService) UserCancelOrder(ctx context.Context, req *pb.OrderResponse) (*pb.OrderResponse, error) {
	err := order.Adapter.UserCancelOrder(int(req.OrderId))
	if err != nil {
		return &pb.OrderResponse{}, err
	}
	return &pb.OrderResponse{OrderId: req.OrderId}, nil
}
func (order *OrderService) ChangeOrderStatus(ctx context.Context, req *pb.ChangeOrderStatusRequest) (*pb.OrderResponse, error) {
	err := order.Adapter.ChangeOrderStatus(int(req.OrderId), int(req.StatusId))
	if err != nil {
		return &pb.OrderResponse{}, err
	}
	return &pb.OrderResponse{OrderId: req.OrderId}, nil
}
func (order *OrderService) GetAllOrdersUser(req *pb.OrderRequest, srv pb.OrderService_GetAllOrdersUserServer) error {
	orders, err := order.Adapter.GetAllOrdersUser(int(req.UserId))
	if err != nil {
		return err
	}
	for _, ordr := range orders {
		var orderItems []*pb.OrderItems
		for _, ordrItem := range ordr.OrderItems {
			itm := &pb.OrderItems{
				OrderId:  uint32(ordrItem.OrderId),
				Id:       uint32(ordrItem.ProductId),
				Quantity: int32(ordrItem.Quantity),
				Total:    ordrItem.Total,
			}
			orderItems = append(orderItems, itm)
		}
		res := &pb.GetAllOrderResponse{
			OrderId:       uint32(ordr.OrderId),
			AddressId:     uint32(ordr.AddressId),
			PaymentTypeId: uint32(ordr.PaymentTypeId),
			OrderStatusId: uint32(ordr.OrderStatusId),
			OrderItems:    orderItems,
		}
		fmt.Println(ordr.OrderStatusId)
		if err := srv.Send(res); err != nil {
			return err
		}
	}
	return nil
}
func (order *OrderService) GetAllOrders(empty *pb.NoParam, srv pb.OrderService_GetAllOrdersServer) error {
	orders, err := order.Adapter.GetAllOrders()
	if err != nil {
		return err
	}
	for _, ordr := range orders {
		var orderItems []*pb.OrderItems
		for _, ordrItem := range ordr.OrderItems {
			itm := &pb.OrderItems{
				OrderId:  uint32(ordrItem.OrderId),
				Id:       uint32(ordrItem.ProductId),
				Quantity: int32(ordrItem.Quantity),
				Total:    ordrItem.Total,
			}
			orderItems = append(orderItems, itm)
		}
		res := &pb.GetAllOrderResponse{
			OrderId:       uint32(ordr.OrderId),
			AddressId:     uint32(ordr.AddressId),
			PaymentTypeId: uint32(ordr.PaymentTypeId),
			OrderStatusId: uint32(ordr.OrderStatusId),
			OrderItems:    orderItems,
		}
		if err := srv.Send(res); err != nil {
			return err
		}
	}
	return nil
}
func (order *OrderService) GetOrder(ctx context.Context, req *pb.OrderResponse) (*pb.GetAllOrderResponse, error) {
	orderData, err := order.Adapter.GetOrder(int(req.OrderId))
	if err != nil {
		return &pb.GetAllOrderResponse{}, err
	}
	var orderItems []*pb.OrderItems
	for _, item := range orderData.OrderItems {
		ordrItem := &pb.OrderItems{
			Id:       uint32(item.ProductId),
			OrderId:  uint32(item.OrderId),
			Quantity: int32(item.Quantity),
			Total:    item.Total,
		}
		orderItems = append(orderItems, ordrItem)
	}
	res := &pb.GetAllOrderResponse{
		OrderId:       req.OrderId,
		AddressId:     uint32(orderData.AddressId),
		PaymentTypeId: uint32(orderData.PaymentTypeId),
		OrderStatusId: uint32(orderData.OrderStatusId),
		OrderItems:    orderItems,
	}
	return res, nil
}

type HealthChecker struct {
	grpc_health_v1.UnimplementedHealthServer
}

func (s *HealthChecker) Check(ctx context.Context, in *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	fmt.Println("check called")
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (s *HealthChecker) Watch(in *grpc_health_v1.HealthCheckRequest, srv grpc_health_v1.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "Watching is not supported")
}
