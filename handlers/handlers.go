package handlers

import (
	"context"
	"log"
	"net/http"
	"sort"

	"github.com/ashbeelghouri/aqary-assignment/internal/database"
	"github.com/ashbeelghouri/aqary-assignment/utilities"
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

type VerifyOTPInout struct {
	PhoneNumber string `json:"phone_number"`
	Otp         string `json:"otp"`
}

type RearrangeString struct {
	Str string `json:"s"`
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
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  "Invalid request data",
		})
		return
	}
	tx, err := h.db.Begin(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	_, err = h.store.GetUserByPhone(context.Background(), string(newUser.PhoneNumber))

	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  "phone number already exists",
		})
		return
	}

	userCreated, err := h.store.CreateUser(h.ctx, database.CreateUserParams{
		Name:        newUser.Name,
		PhoneNumber: newUser.PhoneNumber,
	})

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  "unexpected error while creating user",
		})
		tx.Rollback(context.Background())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "user is created successfully",
		"data":    userCreated,
	})
}

func (h *UserHandler) GenerateOTP(c *gin.Context) {
	var request *GenerateOTPInput

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": false,
			"error":  "Invalid request data",
		})
		return
	}
	tx, err := h.db.Begin(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	_, err = h.store.GetUserByPhone(context.Background(), string(request.PhoneNumber))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  "phone number does not exists",
		})
		return
	}
	otp := utilities.GenerateOTP()

	_, err = h.store.UpdateUserOTP(h.ctx, database.UpdateUserOTPParams{
		PhoneNumber: request.PhoneNumber,
		Otp:         pgtype.Text{String: otp, Valid: true},
	})

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"error":  "Unexpected error occurred while creating OTP",
		})
		tx.Rollback(context.Background())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "otp is generated",
		"data":    otp,
	})
}

func (h *UserHandler) VerifyOTP(c *gin.Context) {
	var request *VerifyOTPInout
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": false,
			"error":  "Invalid request data",
		})
		return
	}
	userAccount, err := h.store.GetUserByPhone(context.Background(), request.PhoneNumber)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  "phone number does not exists",
		})
		return
	}

	otpProvided := pgtype.Text{String: request.Otp, Valid: true}

	if userAccount.Otp != otpProvided {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  "invalid OTP",
		})
		return
	}

	if utilities.IsOtpExpired(userAccount.OtpExpirationTime.Time) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  "Your OTP has been expired",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "your OTP is valid",
	})
}

func ReArrangeString(c *gin.Context) {
	var request *RearrangeString
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": false,
			"error":  "Invalid request data",
		})
		return
	}

	freq := make(map[rune]int)
	for _, char := range request.Str {
		freq[char]++
	}

	var sortedChars []rune

	for char := range freq {
		sortedChars = append(sortedChars, char)
	}

	sort.Slice(sortedChars, func(i, j int) bool {
		return freq[sortedChars[i]] > freq[sortedChars[j]] || (freq[sortedChars[i]] == freq[sortedChars[j]] && sortedChars[i] < sortedChars[j])
	})

	if freq[sortedChars[0]] > (len(request.Str)+1)/2 {
		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"data":    "",
			"message": "re-arrangement not possible",
		})
		return
	}

	result := make([]rune, len(request.Str))

	idx := 0
	for _, char := range sortedChars {
		count := freq[char]
		for count > 0 {
			result[idx] = char
			idx += 2
			if idx >= len(request.Str) {
				idx = 1
			}
			count--
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"data":    string(result),
		"message": "re-arrangement done",
	})
}
