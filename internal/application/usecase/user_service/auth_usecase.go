package user_service

type AuthUserUsecase struct {
	grpcClient InfraGRPCClient
}

type InfraGRPCClient interface {
}

func NewAuthUserUsecase(grpcClient InfraGRPCClient) *AuthUserUsecase {
	return &AuthUserUsecase{grpcClient: grpcClient}
}
