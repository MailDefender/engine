package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	db "maildefender/engine/internal/database"
	"maildefender/engine/internal/models"
)

type searchReputationIn struct {
	Email  string `form:"email"`
	Status string `form:"status"`
}

func GetReputations(c *gin.Context) {
	tx := db.Instance().Gorm

	reps, err := models.GetAllReputations(tx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, newError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, reps)
}

func SearchReputation(c *gin.Context) {
	var in searchReputationIn
	if err := c.ShouldBindQuery(&in); err != nil {
		c.JSON(http.StatusBadRequest, newError(err.Error()))
		return
	}

	tx := db.Instance().Gorm

	reps, err := models.SearchReputation(tx, models.SearchReputationIn{
		Email:  in.Email,
		Status: models.ReputationStatus(in.Status),
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, newError("rule not found"))
			return
		}

		c.JSON(http.StatusInternalServerError, newError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, reps)
}
