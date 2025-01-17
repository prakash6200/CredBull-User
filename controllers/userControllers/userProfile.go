package userProfileController

import (
	"fib/database"
	"fib/middleware"
	"fib/models"

	"github.com/gofiber/fiber/v2"
)

func AddBankAccount(c *fiber.Ctx) error {
	// Retrieve the userId from the JWT token (added by JWTMiddleware)
	userId := c.Locals("userId").(uint)

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
