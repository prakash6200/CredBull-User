package utils

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateOTP generates a 6-digit OTP
func GenerateOTP() string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano())) // Create a new random number generator
	otp := ""
	for i := 0; i < 6; i++ {
		otp += fmt.Sprintf("%d", rng.Intn(10)) // Generate a random digit (0-9) and append to OTP string
	}
	return otp
}
