package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"os"
	"time"
)

type (
	Todo struct {
		Id       primitive.ObjectID `bson:"_id" json:"id"`
		Title    string             `bson:"title" json:"title" validate:"required"`
		Priority string             `bson:"priority" json:"priority" validate:"required"`
	}
)

func main() {
	// Echo instance
	e := echo.New()
	e.Use(middleware.CORS())

	e.GET("/todo", getToDoList)
	e.PUT("/todo", addToDo)
	e.DELETE("/todo", deleteToDo)

	e.Logger.Fatal(e.Start(":1323"))
}

func getToDoList(c echo.Context) error {

	var ctx, _ = getContext()
	var client, _ = getMongoConnection(ctx)
	collection := client.Database("todoApp").Collection("todos")

	var todoList []Todo
	var cursor, _ = collection.Find(ctx, bson.M{})
	cursor.All(ctx, &todoList)

	return c.JSON(200, todoList)
}

func addToDo(c echo.Context) (err error) {
	todo := new(Todo)
	if err = c.Bind(todo); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var ctx, _ = getContext()
	var client, _ = getMongoConnection(ctx)

	collection := client.Database("todoApp").Collection("todos")
	todo.Id = primitive.NewObjectID()
	collection.InsertOne(ctx, todo)
	return c.JSON(http.StatusOK, todo)
}

func deleteToDo(c echo.Context) (err error) {
	todoId := c.QueryParam("todoId")
	idPrimitive, err := primitive.ObjectIDFromHex(todoId)

	var ctx, _ = getContext()
	var client, _ = getMongoConnection(ctx)

	collection := client.Database("todoApp").Collection("todos")

	var todo Todo
	if err = collection.FindOne(ctx, bson.M{"_id": idPrimitive}).Decode(&todo); err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	deleteResult, _ := collection.DeleteOne(ctx, bson.M{"_id": idPrimitive})

	if deleteResult.DeletedCount > 0 {
		return c.NoContent(http.StatusOK)
	} else {
		return c.NoContent(http.StatusInternalServerError)
	}
}

func getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

func getMongoConnection(ctx context.Context) (*mongo.Client, error) {
	return mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
}
