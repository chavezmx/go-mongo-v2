package controllers

import (
	"context"
	"github.com/labstack/gommon/log"
	"github.com/newrelic/go-agent/v3/integrations/nrecho-v4"
	"go-mongo-v2/config"
	"go-mongo-v2/models"
	"go-mongo-v2/responses"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var userCollection *mongo.Collection = config.GetCollection(config.DB, "users")
var validate = validator.New()

func CreateUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()


	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest,
			Message: "Error", Data: &echo.Map{"data": err.Error()}})
	}

	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest,
			Message: "Error", Data: &echo.Map{"data": validationErr.Error()}})
	}

	newUser := models.User{
		Id:			primitive.NewObjectID(),
		Name:		user.Name,
		Location:	user.Location,
		Title:		user.Title,
	}

	tnx := nrecho.FromContext(c)
	tnx.AddAttribute("NewUser", user.Id)

	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError,
			Message: "Error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusCreated, responses.UserResponse{Status: http.StatusCreated, Message: "Success",
		Data: &echo.Map{"data": result}})
}

func GetUser( c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Param("userId")
	tnx := nrecho.FromContext(c)
	tnx.AddAttribute("userId", userId)

	var user models.User
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)
	log.Info(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError,
			Message: "Error User", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "Success",
		Data: &echo.Map{"data": user}})

	//if err := c.Bind(&user); err != nil {
	//	return c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest,
	//		Message: "Error", Data: &echo.Map{"data": err.Error()}})
	//}
	//
	//if validationErr := validate.Struct(&user); validationErr != nil {
	//}

}

func GetAllUsers(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []models.User
	defer cancel()

	results, err := userCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleUser models.User
		if err = results.Decode(&singleUser); err != nil {
			return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
		}

		users = append(users, singleUser)
	}

	return c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": users}})
}
