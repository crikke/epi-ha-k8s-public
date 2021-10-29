package user

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Certificate struct {
	Id      uuid.UUID
	Name    string
	Issuer  string
	Issued  time.Time
	Expires time.Time
}

type User struct {
	Id              uuid.UUID
	Name            string
	Role            string
	YearsExperience int
	Keywords        []string
	Certificates    []Certificate
}

type UserRoute struct {
	repository Repository
}

func (u *UserRoute) AddRouter(router *gin.Engine) {
	router.GET("/user", u.GetUser)
	router.POST("/user", u.AddUser)
	router.PUT("/user", u.UpdateUser)
}

func (u *UserRoute) GetUser(c *gin.Context) {

	id, err := uuid.Parse(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	user := u.repository.GetUserById(c.Request.Context(), id)
	c.JSON(http.StatusOK, user)
}

func (u *UserRoute) AddUser(c *gin.Context) {
	user := &User{}

	err := c.BindJSON(user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u.repository.AddUser(c.Request.Context(), user)

	c.JSON(http.StatusCreated, user)
}

func (u *UserRoute) UpdateUser(c *gin.Context) {
	id, err := uuid.Parse(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	existing := u.repository.GetUserById(c.Request.Context(), id)

	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"id": id})
		return
	}

	data := &User{}

	err := json.Marshal(data)

	if err != nil {

	}
}
