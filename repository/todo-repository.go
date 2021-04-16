package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"mytodoapp/model"
)

type Repository interface {
	FindAll() interface{}
	Create(i interface{}) interface{}
	Delete(id string)
}

type TodoRepository struct {
	collection *mongo.Collection
}

func NewTodoRepository(collection *mongo.Collection) Repository {
	return TodoRepository{collection: collection}
}

func (repo TodoRepository) FindAll() interface{} {
	var todos = make([]model.TodoModel, 0)
	found, _ := repo.collection.Find(context.Background(), bson.D{})
	_ = found.All(context.Background(), &todos)
	return todos
}

func (repo TodoRepository) Create(i interface{}) interface{} {
	id, _ := repo.collection.InsertOne(context.Background(), i)
	objectId, _ := id.InsertedID.(primitive.ObjectID)
	filter := bson.M{"_id": objectId}
	one := repo.collection.FindOne(context.Background(), filter)
	var todo model.TodoModel
	_ = one.Decode(&todo)
	return todo
}

func (repo TodoRepository) Delete(id string) {
	objectId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objectId}
	_, _ = repo.collection.DeleteOne(context.Background(), filter)
}
