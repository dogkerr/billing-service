package mutations

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type mutationsHandler struct {
	mutationService Service
}

func NewMutationsHandler(mutationService Service) *mutationsHandler {
	return &mutationsHandler{mutationService}
}

func (h *mutationsHandler) GetMutationsByUserID(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userIDString, ok := userID.(string)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	mutations, err := h.mutationService.GetMutationsByUserID(userIDString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mutations)
}

func (h *mutationsHandler) GetMutationByID(c *gin.Context) {
	mutationID := c.Param("id")

	mutation, err := h.mutationService.GetMutationByID(mutationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mutation)
}
