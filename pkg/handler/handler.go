package handler

import (
	"github.com/RakhimovAns/Calculus/Initializers"
	"github.com/RakhimovAns/Calculus/models"
	"github.com/RakhimovAns/Calculus/pkg/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func PostExpression(c *gin.Context) {
	var expression models.Expression
	if err := c.ShouldBindJSON(&expression); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := service.IsValidate(expression.Expression); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid format of expression"})
		return
	}
	id := Initializers.CreateModel(expression)
	c.JSON(http.StatusOK, gin.H{"Your ID": id})
}
func StartCount(c *gin.Context) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	expression := Initializers.GetByID(int64(ID))
	if len(expression.Expression) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "non-existent ID"})
		return
	}
	result, err := service.CountExpression(expression)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	Initializers.SetResult(int64(ID), result)
	c.JSON(http.StatusOK, gin.H{"Your result": result})
}

func GetStatus(c *gin.Context) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	expression := Initializers.GetByID(int64(ID))
	if expression.IsCounted == true {
		c.JSON(http.StatusOK, gin.H{"Your result": expression.Result})
	} else {
		c.JSON(http.StatusOK, gin.H{"Your result:": "counting"})
	}
}
