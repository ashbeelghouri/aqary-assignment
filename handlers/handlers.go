package handlers

import (
	"context"
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

func NewUserHandler(db *pgx.Conn, store *database.Queries, ctx context.Context) *UserHandler {
	return &UserHandler{
		db:    db,
		store: store,
		ctx:   ctx,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var newUser *database.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}
}
