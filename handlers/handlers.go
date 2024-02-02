package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/ashbeelghouri/aqary-assignment/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type UserHandler struct {
	db    *pgx.Conn
	store *database.Queries
	ctx   context.Context
}

type CreateUserInput struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}

func NewUserHandler(db *pgx.Conn, store *database.Queries, ctx context.Context) *UserHandler {
	return &UserHandler{
		db:    db,
		store: store,
		ctx:   ctx,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var newUser *CreateUserInput

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	_, err := h.store.GetUserByPhone(context.Background(), string(newUser.PhoneNumber))

	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phone number already exists"})
		return
	}

	userCreated, err := h.store.CreateUser(h.ctx, database.CreateUserParams{
		Name:        newUser.Name,
		PhoneNumber: newUser.PhoneNumber,
	})

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "unexpected error while creating user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user is created successfully",
		"data":    userCreated,
	})
}
