package main

import (
	"errors"
	"fmt"
	"github.com/RakhimovAns/Calculus/Initializers"
	"github.com/RakhimovAns/Calculus/govaluate"
	"github.com/RakhimovAns/Calculus/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func init() {
	Initializers.ConnectToDB()
	Initializers.CreateTable()
}
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	fmt.Println(port)
	r := gin.Default()
	r.POST("/expression", PostExpression)
	r.POST("/expression/:id", StartCount)
	r.GET("/expression/:id", GetStatus)
	r.Run()
}

func PostExpression(c *gin.Context) {
	var expression models.Expression
	if err := c.ShouldBindJSON(&expression); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := IsValidate(expression.Expression); err != nil {
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
	result, err := CountExpression(expression)
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
func IsValidate(expression string) error {
	size := len(expression)
	if size == 0 || !(expression[0] >= '0' && expression[0] <= '9') {
		return errors.New("invalid format of expression")
	}
	if !(expression[size-1] >= '0' && expression[size-1] <= '9') {
		return errors.New("invalid format of expression")
	}
	for i := 1; i < size; i++ {
		if !(expression[i] >= '0' && expression[i] <= '9') && !(expression[i-1] >= '0' && expression[i-1] <= '9') {
			return errors.New("invalid format of expression")
		}
	}
	return nil
}

func CountExpression(expression models.Expression) (interface{}, error) {
	expr, err := govaluate.NewEvaluableExpression(expression.Expression)
	if err != nil {
		return -1, errors.New("error with parsing")
	}
	result, err := expr.Evaluate(nil)
	if err != nil {
		return -1, errors.New("error evaluating expression")
	}
	Time := FindTime(expression)
	time.Sleep(time.Duration(Time) * time.Second)
	return result, nil
}

func FindTime(expression models.Expression) int64 {
	result := int64(0)
	for _, char := range expression.Expression {
		switch char {
		case '+':
			result += expression.AddTime
		case '-':
			result += expression.SubTime
		case '/':
			result += expression.DivideTime
		case '*':
			result += expression.MultiplyTime
		}
	}
	return result
}
