package utils

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
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

func SendOTPToMobile(mobile string, otp string) error {
	// Construct the SMS message
	smsMsg := fmt.Sprintf("OTP for Credbull App Registration is %s. Do not share it with anyone.", otp)

	// Construct the TextLocal API URL
	apiURL := "https://api.textlocal.in/send/"
	data := url.Values{}
	data.Set("apikey", "NDM0NDZjNzE2ZDQ1Njg3ODM4NTE0ZDU2NmQ2ZjUxNmI=") // Replace with your actual API key
	data.Set("numbers", mobile)
	data.Set("sender", "CRDBUL")
	data.Set("message", smsMsg)

	// Make the API request
	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		log.Printf("Error while sending OTP: %v", err)
		return err
	}
	defer resp.Body.Close()

	// Check if the response status code is not OK
	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to send OTP, response code: %d", resp.StatusCode)
		return fmt.Errorf("failed to send OTP")
	}

	log.Println("OTP sent successfully to", mobile)
	return nil
}
