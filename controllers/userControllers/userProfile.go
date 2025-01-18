package userProfileController

import (
	"fib/database"
	"fib/middleware"
	"fib/models"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AddBankAccount(c *fiber.Ctx) error {
	// Retrieve the userId from the JWT token (added by JWTMiddleware)
	userId := c.Locals("userId").(uint)
	fmt.Println(userId)
	// Parse the request body to get the bank details
	reqData := new(struct {
		BankName    string `json:"bankName"`
		AccountNo   string `json:"accountNo"`
		HolderName  string `json:"holderName"`
		IFSCCode    string `json:"ifscCode"`
		BranchName  string `json:"branchName"`
		AccountType string `json:"accountType"` // Optional, default to "savings"
	})

	if err := c.BodyParser(reqData); err != nil {
		return middleware.JsonResponse(c, fiber.StatusBadRequest, false, "Failed to parse request body!", nil)
	}

	var user models.User
	if err := database.Database.Db.Where("id = ? AND is_deleted = ?", userId, false).First(&user).Error; err != nil {
		return middleware.JsonResponse(c, fiber.StatusInternalServerError, false, "Failed to fetch user!", nil)
	}

	// Check if the user already has a bank account
	if user.BankDetails != 0 {
		return middleware.JsonResponse(c, fiber.StatusConflict, false, "You already have a bank account!", nil)
	}

	// Check if the bank account already exists
	var existingBankDetails models.BankDetails
	result := database.Database.Db.Where("account_no = ?", reqData.AccountNo).First(&existingBankDetails)

	if result.RowsAffected > 0 {
		return middleware.JsonResponse(c, fiber.StatusConflict, false, "Bank account already exists!", nil)
	}

	// Create a new BankDetails object
	newBankDetails := models.BankDetails{
		BankName:    reqData.BankName,
		AccountNo:   reqData.AccountNo,
		HolderName:  reqData.HolderName,
		IFSCCode:    reqData.IFSCCode,
		BranchName:  reqData.BranchName,
		AccountType: reqData.AccountType,
		UserID:      userId,
	}

	// Save the new bank account to the database
	if err := database.Database.Db.Create(&newBankDetails).Error; err != nil {
		return middleware.JsonResponse(c, fiber.StatusInternalServerError, false, "Failed to add bank account!", nil)
	}

	// If user exists, update their bank details field with the new bank account ID
	user.BankDetails = newBankDetails.ID

	// Save the updated user with the new bank details
	if err := database.Database.Db.Save(&user).Error; err != nil {
		return middleware.JsonResponse(c, fiber.StatusInternalServerError, false, "Failed to update user with bank details!", nil)
	}

	// Respond with success message
	return middleware.JsonResponse(c, fiber.StatusOK, true, "Bank account added successfully.", newBankDetails)
}

func SendAdharOtp(c *fiber.Ctx) error {
	// Retrieve the userId from the JWT token (added by JWTMiddleware)
	userId := c.Locals("userId").(uint)

	reqData := new(struct {
		AadharNumber string `json:"aadharNumber"`
	})

	if err := c.BodyParser(reqData); err != nil {
		return middleware.JsonResponse(c, fiber.StatusBadRequest, false, "Failed to parse request body!", nil)
	}

	// Check if the user exists
	var user models.User
	if err := database.Database.Db.Where("id = ? AND is_deleted = ?", userId, false).First(&user).Error; err != nil {
		return middleware.JsonResponse(c, fiber.StatusNotFound, false, "User not found!", nil)
	}

	// Check if a KYC record already exists for the user
	var existingKYC models.UserKYC
	if err := database.Database.Db.Where("user_id = ?", userId).First(&existingKYC).Error; err == nil {
		return middleware.JsonResponse(c, fiber.StatusConflict, false, "Your KYC record already exists!", nil)
	}

	// Check if the Aadhar number already exists
	var existingAadhar models.AadharDetails
	if err := database.Database.Db.Where("aadhar_number = ?", reqData.AadharNumber).First(&existingAadhar).Error; err == nil {
		return middleware.JsonResponse(c, fiber.StatusConflict, false, "Aadhar number already exists!", nil)
	}

	// Send Aadhaar OTP using the sandbox API
	url := "https://api.sandbox.co.in/kyc/aadhaar/okyc/otp"

	// Construct payload
	payload := fmt.Sprintf(`{"aadhaar_number":"%s"}`, reqData.AadharNumber)

	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		return middleware.JsonResponse(c, fiber.StatusInternalServerError, false, "Failed to create request for Aadhaar OTP!", nil)
	}

	// Add headers
	req.Header.Add("accept", "application/json")
	req.Header.Add("x-api-version", "2.0")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", c.Get("Authorization"))
	req.Header.Add("x-api-key", "key_live_HZYsCB58PuDIMsyhCW2Uvxq576V6Pr6n")       // Replace with your actual API key
	req.Header.Add("x-api-secret", "secret_live_6GBggEXGr5OCxbVXpuwESvKcHXFcQ7MZ") // Replace with your actual API secret

	// Make the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return middleware.JsonResponse(c, fiber.StatusInternalServerError, false, "Failed to send Aadhaar OTP!", nil)
	}
	defer res.Body.Close()

	// Read and parse the response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return middleware.JsonResponse(c, fiber.StatusInternalServerError, false, "Failed to read Aadhaar OTP response!", nil)
	}

	// Check for successful response
	if res.StatusCode != http.StatusOK {
		return middleware.JsonResponse(c, fiber.StatusInternalServerError, false, "Failed to send Aadhaar OTP: "+string(body), nil)
	}

	// Respond with success message
	return middleware.JsonResponse(c, fiber.StatusOK, true, "Adhar OTP sent successfully.", nil)
}
