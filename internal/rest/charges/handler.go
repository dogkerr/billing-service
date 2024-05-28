package charges

import (
	"dogker/andrenk/billing-service/internal/rest/mutations"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type chargeHandler struct {
	depositService  Service
	mutationService mutations.Service
}

func NewChargeHandler(depositService Service, mutationService mutations.Service) *chargeHandler {
	return &chargeHandler{depositService, mutationService}
}

func (h *chargeHandler) InitiateCharge(c *gin.Context) {
	var input ChargeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the userID from the context JWT
	userID, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// Convert the userID to a string
	userIDString, ok := userID.(string)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	input.UserID = userIDString

	// Create charge
	charge, err := h.depositService.CreateCharge(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create mutation
	mutationInput := mutations.MutationInput{
		ID:       uuid.NewString(),
		UserID:   input.UserID,
		Mutation: input.TotalCost,
		Type:     "charge",
		ChargeID: input.ID,
	}

	if _, err := h.mutationService.CreateMutation(mutationInput); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Charge initiated", "data": charge})
}

func (h *chargeHandler) GetChargeByID(c *gin.Context) {
	chargeID := c.Param("id")

	charge, err := h.depositService.GetChargeByID(chargeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Charge fetched!", "data": charge})
}
