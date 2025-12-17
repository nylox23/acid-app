package handler

import (
	"web_service_auth/internal/app/config"
	"web_service_auth/internal/app/redis"
	"web_service_auth/internal/app/repository"
	"web_service_auth/internal/app/role"
	"web_service_auth/internal/app/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Repository   *repository.Repository
	MinioService *service.MinioService
	Config       *config.Config
	Redis        *redis.Client
}

func NewHandler(r *repository.Repository, minioService *service.MinioService, cfg *config.Config, redisClient *redis.Client) *Handler {
	return &Handler{
		Repository:   r,
		MinioService: minioService,
		Config:       cfg,
		Redis:        redisClient,
	}
}

// Регистрация эндпоинтов API
func (h *Handler) RegisterAPIHandler(router *gin.Engine) {

	api := router.Group("/api")
	{
		// Домен кислот
		acids := api.Group("/acids")
		{
			acids.GET("", h.GetAcidsAPI)
			acids.GET("/:id", h.GetAcidAPI)
			acids.Use(h.WithAuthCheck(role.User, role.Admin)).POST("/:id/toCarbonate", h.AddAcidToCarbonateAPI)
			acids.Use(h.WithAuthCheck(role.Admin)).POST("", h.CreateAcidAPI)
			acids.Use(h.WithAuthCheck(role.Admin)).PUT("/:id", h.UpdateAcidAPI)
			acids.Use(h.WithAuthCheck(role.Admin)).DELETE("/:id", h.DeleteAcidAPI)
			acids.Use(h.WithAuthCheck(role.Admin)).POST("/:id/image", h.AddAcidImageAPI)
		}

		// Домен заявок
		carbonates := api.Group("/carbonates")
		{
			carbonates.Use(h.WithAuthCheck(role.User, role.Admin)).GET("/current", h.GetCurrentCarbonateAPI)
			carbonates.Use(h.WithAuthCheck(role.User, role.Admin)).GET("", h.GetCarbonatesAPI)
			carbonates.Use(h.WithAuthCheck(role.User, role.Admin)).GET("/:id", h.GetCarbonateAPI)
			carbonates.Use(h.WithAuthCheck(role.User, role.Admin)).PUT("", h.UpdateCarbonateAPI)
			carbonates.Use(h.WithAuthCheck(role.User, role.Admin)).PUT("/form", h.FormCarbonateAPI)
			carbonates.Use(h.WithAuthCheck(role.Admin)).PUT("/:id/status", h.SetCarbonateStatusAPI)
			carbonates.Use(h.WithAuthCheck(role.User, role.Admin)).DELETE("/:id", h.DeleteCarbonateAPI)
		}

		// Домен м-м
		carbonateAcids := api.Group("/carbonate-acids")
		{
			carbonateAcids.Use(h.WithAuthCheck(role.User, role.Admin)).PUT("/:id", h.UpdateCarbonateAcidAPI)
			carbonateAcids.Use(h.WithAuthCheck(role.User, role.Admin)).DELETE("/:id", h.DeleteCarbonateAcidAPI)
		}

		// Домен пользователей
		users := api.Group("/users")
		{
			users.POST("/register", h.RegisterUserAPI)
			users.POST("/login", h.LoginUserAPI)
			users.Use(h.WithAuthCheck(role.User, role.Admin)).POST("/logout", h.LogoutUserAPI)
			users.Use(h.WithAuthCheck(role.User, role.Admin)).GET("/profile", h.GetUserProfileAPI)
			users.Use(h.WithAuthCheck(role.User, role.Admin)).PUT("/profile", h.UpdateUserProfileAPI)
		}
	}
}
