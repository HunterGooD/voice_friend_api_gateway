package user_service

import (
	"context"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase UserAuthUsecase
}

type UserAuthUsecase interface {
	RegisterUserUsecase(ctx context.Context, deviceID string) (any, error)
	LoginUserUsecase(ctx context.Context, deviceID string) (any, error)
	LogoutUserUsecase(ctx context.Context, refreshToken string) error
	UpdateRefreshTokenUsecase(ctx context.Context, refreshToken string) (any, error)
	UpdateAccessTokenUsecase(ctx context.Context, refreshToken string) (any, error)
}

func NewHandler(route *gin.Engine, userUsecase UserAuthUsecase) {
	userHandler := &UserHandler{userUsecase}

	userServiceGroup := route.Group("user")
	userServiceGroup.POST("/login", userHandler.LoginUserHandler)
	userServiceGroup.POST("/signup", userHandler.RegisterUserHandler)
	userServiceGroup.POST("/logout", userHandler.LogOutUserHandler)
	userServiceGroup.POST("/refresh/token/access", userHandler.UpdateAccessTokenUserHandler)
	userServiceGroup.POST("/refresh/token/refresh", userHandler.UpdateRefreshTokenUserHandler)

}

func (h *UserHandler) LoginUserHandler(c *gin.Context) {

}

func (h *UserHandler) RegisterUserHandler(c *gin.Context) {

}

func (h *UserHandler) LogOutUserHandler(c *gin.Context) {

}

func (h *UserHandler) UpdateAccessTokenUserHandler(c *gin.Context) {

}

func (h *UserHandler) UpdateRefreshTokenUserHandler(c *gin.Context) {

}
