package userPorfileValidator

import (
	"fib/middleware"
	"regexp"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func isValidNumeric(input string) bool {
	_, err := strconv.Atoi(input)
	return err == nil
}

func isValidIFSC(ifsc string) bool {
	re := regexp.MustCompile(`^[A-Za-z]{4}0[A-Za-z0-9]{6}$`)
	return re.MatchString(ifsc)
}

func AddBankAccount() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Parse request body
		reqData := new(struct {
			BankName    string `json:"bankName"`
			AccountNo   string `json:"accountNo"`
			HolderName  string `json:"holderName"`
			IFSCCode    string `json:"ifscCode"`
			BranchName  string `json:"branchName"`
			AccountType string `json:"accountType"` // Optional
		})
		if err := c.BodyParser(reqData); err != nil {
			return middleware.JsonResponse(c, fiber.StatusBadRequest, false, "Invalid request body!", nil)
		}

		errors := make(map[string]string)

		// Validate Bank Name
		if len(strings.TrimSpace(reqData.BankName)) < 3 {
			errors["bankName"] = "Bank name must be at least 3 characters long!"
		}

		// Validate Account Number
		if len(strings.TrimSpace(reqData.AccountNo)) < 10 || len(reqData.AccountNo) > 18 {
			errors["accountNo"] = "Account number must be between 10 and 18 digits!"
		} else if !isValidNumeric(reqData.AccountNo) {
			errors["accountNo"] = "Account number must contain only numeric characters!"
		}

		// Validate Holder Name
		if len(strings.TrimSpace(reqData.HolderName)) < 3 {
			errors["holderName"] = "Holder name must be at least 3 characters long!"
		}

		// Validate IFSC Code
		if len(strings.TrimSpace(reqData.IFSCCode)) != 11 || !isValidIFSC(reqData.IFSCCode) {
			errors["ifscCode"] = "Invalid IFSC code! It must be 11 characters long and alphanumeric."
		}

		// Validate Branch Name (Optional but must not be empty if provided)
		if reqData.BranchName != "" && len(strings.TrimSpace(reqData.BranchName)) < 3 {
			errors["branchName"] = "Branch name must be at least 3 characters long if provided!"
		}

		// Validate Account Type (Optional but must match valid values)
		validAccountTypes := map[string]bool{"savings": true, "current": true}
		if reqData.AccountType != "" && !validAccountTypes[strings.ToLower(reqData.AccountType)] {
			errors["accountType"] = "Account type must be 'savings' or 'current'!"
		}

		// Respond with errors if any exist
		if len(errors) > 0 {
			return middleware.ValidationErrorResponse(c, errors)
		}

		// Pass validated bank details to the next middleware
		c.Locals("validatedBankDetails", reqData)
		return c.Next()
	}
}