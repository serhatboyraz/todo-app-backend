package handler

import (
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mytodoapp/model"
	"mytodoapp/repository"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_ShouldReturnEmptyTodoList(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	todos := make([]model.TodoModel, 0)
	repo := repository.NewMockRepository(controller)
	repo.EXPECT().FindAll().Return(todos).Times(1)

	handler := TodoHandler{repo}
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	req.Header.Set(echo.HeaderAccept, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	expected := `[]
`
	if assert.NoError(t, handler.HandleGetAll(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, rec.Body.String())
	}
}

func Test_ShouldFindAllTodos(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	var todos []model.TodoModel
	primId, _ := primitive.ObjectIDFromHex("60798540288942e18e737a34")
	todos = append(todos, model.TodoModel{
		Id:    primId,
		Title: "example my todo",
	})

	repo := repository.NewMockRepository(controller)
	repo.EXPECT().FindAll().Return(todos).Times(1)
	handler := TodoHandler{repo}

	expectVal := `[{"id":"60798540288942e18e737a34","title":"example my todo"}]
`

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/todo", nil)
	req.Header.Set(echo.HeaderAccept, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	if assert.NoError(t, handler.HandleGetAll(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expectVal, rec.Body.String())
	}
}

func Test_ShouldCreateTodo(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	primId, _ := primitive.ObjectIDFromHex("60798540288942e18e737a34")
	request := model.TodoModel{Id: primId, Title: "example my todo"}

	requestAsString := `{"id":"60798540288942e18e737a34","title":"example my todo"}`
	repo := repository.NewMockRepository(controller)
	handler := TodoHandler{repo}
	repo.EXPECT().Create(gomock.Eq(request)).Return(requestAsString).Times(1)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/todo", strings.NewReader(requestAsString))
	req.Header.Set(echo.HeaderAccept, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	if assert.NoError(t, handler.HandleCreate(ctx)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}
func Test_ShouldDeleteTodo(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	repo := repository.NewMockRepository(controller)
	repo.EXPECT().Delete(gomock.Eq("60798540288942e18e737a34"))
	handler := TodoHandler{repo}

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/todos/60798540288942e18e737a34", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues("60798540288942e18e737a34")

	if assert.NoError(t, handler.HandleDelete(ctx)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
	}
}
