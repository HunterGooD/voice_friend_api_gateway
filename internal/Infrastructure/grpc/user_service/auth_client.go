package user_service

import (
	"context"

	"github.com/HunterGooD/voice_friend_api_gateway/internal/domain/user_service/entity"
	pd "github.com/HunterGooD/voice_friend_contracts/gen/go/user_service"
	"google.golang.org/grpc"
)

type GrpcConnProvider interface {
	GetConn() (*grpc.ClientConn, error)
}

type Client struct {
	grpcConnPool GrpcConnProvider
}

func NewAuthClient(provider GrpcConnProvider) *Client {
	return &Client{provider}
}

func (c *Client) Register(ctx context.Context) (*entity.User, error) {
	client, err := c.getConnection(ctx)
	if err != nil {
		return nil, err
	}

	client.Register(ctx, nil)
	return nil, nil
}

func (c *Client) Login(ctx context.Context) (*entity.User, error) {
	return nil, nil
}

func (c *Client) LogOut(ctx context.Context) (*entity.User, error) {
	return nil, nil
}

func (c *Client) UpdateAccessToken(ctx context.Context) (*entity.User, error) {
	return nil, nil
}

func (c *Client) UpdateRefreshToken(ctx context.Context) (*entity.User, error) {
	return nil, nil
}

func (c *Client) getConnection(ctx context.Context) (pd.AuthClient, error) {
	conn, err := c.grpcConnPool.GetConn()
	if err != nil {
		return nil, err
	}

	return pd.NewAuthClient(conn), nil
}
