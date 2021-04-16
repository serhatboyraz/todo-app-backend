package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	configuration "mytodoapp/config"
	"mytodoapp/handler"
	"mytodoapp/repository"
	"os"
)

func main() {
	e := echo.New()
	e.Use(middleware.CORS())

	config := new(configuration.MongoConfiguration).Init(getDbUri(), "mytodoapp")
	todoRepository := repository.NewTodoRepository(config.Database().Collection("todos"))
	todoHandler := handler.NewTodoHandler(todoRepository)
	e.GET("/todo", todoHandler.HandleGetAll)
	e.POST("/todo", todoHandler.HandleCreate)
	e.DELETE("/todo/:id", todoHandler.HandleDelete)

	e.Logger.Fatal(e.Start(":1323"))
}

func getDbUri() string {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		return "mongodb://localhost:27017"
	}
	return uri
}
