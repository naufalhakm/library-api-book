package handlers

import (
	"context"
	"library-api-book/internal/services"
	pb "library-api-book/proto/book"
)

type BookHandler struct {
	service services.BookService
	pb.UnimplementedBookServiceServer
}

func NewBookHandler(service services.BookService) *BookHandler {
	return &BookHandler{service: service}
}

func (handler *BookHandler) DecreaseStock(ctx context.Context, req *pb.DecreaseStockRequest) (*pb.DecreaseStockResponse, error) {
	err := handler.service.DecreaseStock(ctx, req.BookId)
	if err != nil {
		return &pb.DecreaseStockResponse{Success: false, Message: err.Message}, nil
	}
	return &pb.DecreaseStockResponse{Success: true, Message: "Book stock decrease successfully"}, nil
}
func (handler *BookHandler) IncreaseStock(ctx context.Context, req *pb.IncreaseStockRequest) (*pb.IncreaseStockResponse, error) {
	err := handler.service.IncreaseStock(ctx, req.BookId)
	if err != nil {
		return &pb.IncreaseStockResponse{Success: false, Message: err.Message}, nil
	}
	return &pb.IncreaseStockResponse{Success: true, Message: "Book stock increase successfully"}, nil
}
