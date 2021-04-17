package handler

import (
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mytodoapp/model"
	"mytodoapp/repository"
)

type TodoHandler struct {
	repository repository.Repository
}

type Handler interface {
	HandleGetAll(ctx echo.Context) error
	HandleCreate(ctx echo.Context) error
	HandleDelete(ctx echo.Context) error
}

func NewTodoHandler(repository repository.Repository) Handler {
	return TodoHandler{repository}
}

func (handler TodoHandler) HandleGetAll(ctx echo.Context) error {
	return ctx.JSON(200, handler.repository.FindAll())
}

func (handler TodoHandler) HandleCreate(ctx echo.Context) error {
	var request model.TodoModel
	if request.Id == primitive.NilObjectID {
		request.Id = primitive.NewObjectID()
	}
	_ = ctx.Bind(&request)

	response := handler.repository.Create(request)
	return ctx.JSON(201, response)
}

func (handler TodoHandler) HandleDelete(ctx echo.Context) error {
	id := ctx.Param("id")
	handler.repository.Delete(id)
	return ctx.NoContent(204)
}
