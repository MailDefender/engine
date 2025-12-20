package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	db "maildefender/engine/internal/database"
	engineErrors "maildefender/engine/internal/errors"
	"maildefender/engine/internal/models"
	"maildefender/engine/internal/validation"
)

type validateTokenIn struct {
	Token string `uri:"token" binding:"required"`
}

func ValidateToken(c *gin.Context) {
	var in validateTokenIn
	if err := c.ShouldBindUri(&in); err != nil {
		c.JSON(http.StatusBadRequest, newError(err.Error()))
		return
	}

	tx := db.Instance().Gorm
	token, err := models.GetValidationTokenByToken(tx, in.Token)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	err = validation.Validate(tx, token)

	if err == nil {
		c.JSON(http.StatusOK, nil)
		return
	}

	if errors.Is(err, engineErrors.ErrExpiredToken) || errors.Is(err, engineErrors.ErrAlreadyValidatedToken) {
		c.JSON(http.StatusConflict, newError(err.Error()))
		return
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	c.JSON(http.StatusInternalServerError, newError(err.Error()))
}
