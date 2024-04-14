package deposits

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type depositsHandler struct {
	depositService Service
}

func NewDepositsHandler(depositService Service) *depositsHandler {
	return &depositsHandler{depositService}
}

func (h *depositsHandler) InitiateDeposit(c *gin.Context) {
	var input DepositInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	snapResponse, err := h.depositService.InitiateDeposit(input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": snapResponse})
}

func (h *depositsHandler) AddDeposit(c *gin.Context) {
	var input DepositInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	deposit, err := h.depositService.AddDeposit(input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": deposit})
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

func (h *depositsHandler) UpdateDeposit(c *gin.Context) {
	var input DepositUpdateInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	deposit, err := h.depositService.UpdateDeposit(input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": deposit})
}
