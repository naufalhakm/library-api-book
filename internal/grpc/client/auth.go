package client

import (
	"context"

	tkn "library-api-book/pkg/token"
	pb "library-api-book/proto/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	client pb.AuthServiceClient
	conn   *grpc.ClientConn
}

func NewAuthClient(addr string) (*AuthClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewAuthServiceClient(conn)
	return &AuthClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *AuthClient) ValidateToken(ctx context.Context, token string) (bool, *tkn.Token) {
	resp, err := c.client.ValidateToken(ctx, &pb.ValidateRequest{Token: token})
	if err != nil || !resp.Success {
		return false, nil
	}

	return true, &tkn.Token{
		AuthId: int(resp.AuthId),
		Role:   resp.Role,
	}
}

func (c *AuthClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
