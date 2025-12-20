package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	db "maildefender/engine/internal/database"
	"maildefender/engine/internal/models"
)

func GetRules(c *gin.Context) {
	tx := db.Instance().Gorm
	r, err := models.GetAllRules(tx)

	if err != nil {
		c.JSON(http.StatusInternalServerError, newError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, r)
}

func GetRuleByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, newError(err.Error()))
		return
	}

	tx := db.Instance().Gorm
	r, err := models.GetRuleByID(tx, uint(id))

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, newError("rule not found"))
			return
		}

		c.JSON(http.StatusInternalServerError, newError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, r)
}

func DeleteRule(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, newError(err.Error()))
		return
	}

	tx := db.Instance().Gorm
	err = models.DeleteRuleByID(tx, uint(id))

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, newError("rule not found"))
			return
		}

		c.JSON(http.StatusInternalServerError, newError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, nil)
}
