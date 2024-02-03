package utilities

import (
	"fmt"
	"math/rand"
	"regexp"
	"time"
)

func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%04d", rand.Intn(10000))
}

func IsOtpExpired(usersOTPExpiry time.Time) bool {
	return usersOTPExpiry.Before(time.Now())
}

func IsValidPhoneNumber(phone string) bool {
	return regexp.MustCompile(`^\+[1-9]\d{1,14}$`).MatchString(phone)
}
