package controllers

import (
	"github.com/dharmasastra/learn-echolabstack-twitter/app/models"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (h *Handler) SignUp(c echo.Context) (err error) {
	//bind
	u := &models.User{ID: bson.NewObjectId()}
	if err = c.Bind(u); err != nil {
		return
	}

	//validate
	if u.Email == "" || u.Password == "" {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Invalid email or password"}
	}

	//save user
	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB("twitter").C("users").Insert(u); err != nil {
		return
	}
	u.Password = "" //don't send password
	return c.JSON(http.StatusCreated, u)
}

func (h *Handler) Login(c echo.Context) (err error) {
	//bind
	u := new(models.User)
	if err = c.Bind(u); err != nil {
		return
	}

	//find user
	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB("twitter").C("users").
		Find(bson.M{"email": u.Email, "password": u.Password}).One(u); err != nil {
		if err == mgo.ErrNotFound {
			return &echo.HTTPError{Code: http.StatusUnauthorized, Message: "Invalid email or password"}
		}
		return
	}

	//JWT

	//create token
	token := jwt.New(jwt.SigningMethodHS256)

	//set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = u.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	//generate encoded token and send it as response
	u.Token, err = token.SignedString([]byte(Key))
	if err != nil {
		return err
	}

	u.Password = "" //don't send password
	return c.JSON(http.StatusOK, u)
}

func (h *Handler) Follow(c echo.Context) (err error) {
	userID := userIDFromToken(c)
	id := c.Param("id")

	// add  a follower to user
	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB("twitter").C("users").
		UpdateId(bson.ObjectIdHex(id), bson.M{"$addToSet": bson.M{"followers": userID}}); err != nil {
		if err == mgo.ErrNotFound {
			return echo.ErrNotFound
		}
	}
	return
}

func userIDFromToken(c echo.Context) string {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims["id"].(string)
}
