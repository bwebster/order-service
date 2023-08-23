package main

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	ordersv1 "github.com/bwebster/order-service/gen/orders/v1"
	"github.com/bwebster/order-service/gen/orders/v1/ordersv1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const address = "localhost:8080"

func main() {
	mux := http.NewServeMux()
	path, handler := ordersv1connect.NewOrderServiceHandler(&orderStoreServiceServer{})
	mux.Handle(path, handler)
	fmt.Println("... Listening on", address)
	http.ListenAndServe(
		address,
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}

// orderStoreServiceServer implements the OrderService API.
type orderStoreServiceServer struct {
	ordersv1connect.UnimplementedOrderServiceHandler
}

func (s *orderStoreServiceServer) GetOrder(
	ctx context.Context,
	req *connect.Request[ordersv1.GetOrderRequest],
) (*connect.Response[ordersv1.GetOrderResponse], error) {
	orderID := req.Msg.GetOrderId()

	return connect.NewResponse(&ordersv1.GetOrderResponse{
		Order: &ordersv1.Order{
			OrderId:    orderID,
			CustomerId: "abc123",
		},
	}), nil
}
