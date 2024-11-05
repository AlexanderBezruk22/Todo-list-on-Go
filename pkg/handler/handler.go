package handler

import (
	"awesomeProject/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (s *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("sign-up", s.signUp)
		auth.POST("sign-in", s.signIn)
	}
	api := router.Group("/api", s.userIdentity)
	{
		lists := api.Group("/lists")
		{
			lists.POST("/", s.createList)
			lists.GET("/", s.getAllLists)
			lists.PUT("/:id", s.updateList)
			lists.GET("/:id", s.getListById)
			lists.DELETE("/:id", s.deleteList)

			items := lists.Group(":id/items")
			{
				items.POST("/", s.createItem)
				items.GET("/", s.getAllItem)
			}
		}
		items := api.Group("/items")
		{
			items.PUT("/:id", s.updateItem)
			items.GET("/:id", s.getItemById)
			items.DELETE("/:id", s.deleteItem)
		}

	}
	return router
}
