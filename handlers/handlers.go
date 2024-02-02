package handlers

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/ashbeelghouri/aqary-assignment/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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

type GenerateOTPInput struct {
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

func (h *UserHandler) GenerateOTP(c *gin.Context) {
	var request *GenerateOTPInput

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid request data"})
		return
	}

	_, err := h.store.GetUserByPhone(context.Background(), string(request.PhoneNumber))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phone number does not exists"})
		return
	}

	rand.Seed(time.Now().UnixNano())
	otp := fmt.Sprintf("%04d", rand.Intn(10000))

	_, err = h.store.UpdateUserOTP(h.ctx, database.UpdateUserOTPParams{
		PhoneNumber: request.PhoneNumber,
		Otp:         pgtype.Text{String: otp, Valid: true},
	})

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unexpected error occurred while creating OTP",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "otp is generated",
		"data":    otp,
	})
}
