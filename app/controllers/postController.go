package controllers

import (
	"github.com/dharmasastra/learn-echolabstack-twitter/app/models"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strconv"
)

func (h *Handler) CreatePost(c echo.Context) (err error) {
	u := &models.User{
		ID: bson.ObjectIdHex(userIDFromToken(c)),
	}
	p := &models.Post{
		ID:   bson.NewObjectId(),
		From: u.ID.Hex(),
	}
	if err = c.Bind(p); err != nil {
		return
	}

	//validation
	if p.To == "" || p.Message == "" {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid to or message fields"}
	}

	//find user from database
	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB("twitter").C("users").FindId(u.ID).One(u); err != nil {
		if err == mgo.ErrNotFound {
			return echo.ErrNotFound
		}
		return
	}

	//save post in database
	if err = db.DB("twitter").C("posts").Insert(p); err != nil {
		return
	}
	return c.JSON(http.StatusCreated, p)
}

func (h *Handler) FetchPost(c echo.Context) (err error) {
	userId := userIDFromToken(c)
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	//default
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 100
	}

	//retrieve posts from database
	var posts []*models.Post
	db := h.DB.Clone()
	if err = db.DB("twitter").C("posts").
		Find(bson.M{"to": userId}).
		Skip((page - 1) * limit).
		Limit(limit).
		All(&posts); err != nil {
		return
	}
	defer db.Close()

	return c.JSON(http.StatusOK, posts)
}
