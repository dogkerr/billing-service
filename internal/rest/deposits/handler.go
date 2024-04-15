package deposits

import (
	"dogker/andrenk/billing-service/internal/rest/auth"
	"dogker/andrenk/billing-service/internal/rest/mutations"
	"dogker/andrenk/billing-service/internal/rest/payment"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
)

type depositsHandler struct {
	depositService  Service
	paymentService  payment.Service
	authService     auth.Service
	mutationService mutations.Service
}

func NewDepositsHandler(depositService Service, paymentService payment.Service, authService auth.Service, mutationService mutations.Service) *depositsHandler {
	return &depositsHandler{depositService, paymentService, authService, mutationService}
}

func (h *depositsHandler) InitiateDeposit(c *gin.Context) {
	// Bind the input to the TransactionInput struct
	var input payment.TransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the userID from the context JWT
	userID, ok := c.Get("userID")
	if !ok {
		fmt.Println("6")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userIDString, ok := userID.(string)
	if !ok {
		fmt.Println("7")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Get the user from the authService
	user, err := h.authService.GetUserById(userIDString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info("User: ", user)

	fmt.Println("Calling paymentService.GetPayment")

	// Call the paymentService.GetPayment method
	transaction, err := h.paymentService.GetPayment(input, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("Transaction: ", transaction)

	// Create a deposit (PENDING)
	depositInput := DepositInput{
		ID:     transaction.ID,
		UserID: transaction.UserID,
		Amount: float64(transaction.Amount),
		Status: "pending",
	}

	if _, err := h.depositService.AddDeposit(depositInput); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction created successfully",
		"data":    transaction,
	})
}

func (h *depositsHandler) HandleNotificationPayment(c *gin.Context) {
	var input payment.TransactionNotificationInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//process notification
	if !h.paymentService.HandleNotificationPayment(input) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Payment rejected"})
		return
	}

	//update deposit status
	deposit, err := h.depositService.UpdateDepositStatus(input.OrderID, input.TransactionStatus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//create mutation
	mutationInput := mutations.MutationInput{
		ID:        uuid.NewString(),
		UserID:    deposit.UserID,
		Mutation:  deposit.Amount,
		Type:      "deposit",
		DepositID: deposit.ID,
	}
	if _, err := h.mutationService.CreateMutation(mutationInput); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction is paid and deposit is billed successfully", "data": deposit})
}

func (h *depositsHandler) GetDepositByID(c *gin.Context) {
	id := c.Param("id")

	deposit, err := h.depositService.GetDepositByID(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": deposit})
}
