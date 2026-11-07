package controller

import (
	"github.com/gofiber/fiber/v3"

	"github.com/rescuebyte/rescuebyte-be/domain/contracts"
	"github.com/rescuebyte/rescuebyte-be/internal/middlewares"
)

type commentController struct {
	commentService contracts.CommentService
}

func InitCommentController(router fiber.Router, commentService contracts.CommentService, middleware *middlewares.Middleware) {
	ctr := &commentController{commentService: commentService}

	storeComments := router.Group("/stores/:store_id/comments")
	storeComments.Get("", ctr.GetStoreComments)
	storeComments.Post("", middleware.RequireAuth(), ctr.CreateStoreComment)

	productComments := router.Group("/products/:product_id/comments")
	productComments.Get("", ctr.GetProductComments)
	productComments.Post("", middleware.RequireAuth(), ctr.CreateProductComment)
}
