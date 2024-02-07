package main

import (
	"errors"
	"fmt"
	"github.com/RakhimovAns/Calculus/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	fmt.Println(port)
	r := gin.Default()
	r.POST("/expression", PostExpression)
	r.Run()
}

func PostExpression(c *gin.Context) {
	var expression models.Expression
	if err := c.ShouldBindJSON(&expression); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err := IsValidate(expression.Expression)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid format of expression"})
		return
	}

	result, err := EvaluateExpression(expression.Expression)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

func IsValidate(expression string) error {
	size := len(expression)
	if !(expression[0] >= '0' && expression[0] <= '9') {
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

func EvaluateExpression(expression string) (float64, error) {
	var operands []float64
	var operators []rune

	for _, char := range expression {
		if char == '+' || char == '-' || char == '*' || char == '/' {
			operators = append(operators, char)
		} else {
			num, _ := strconv.ParseFloat(string(char), 64)
			if len(operands) == 0 {
				operands = append(operands, num)
			} else {
				// Check for previous operator
				prevOperator := operators[len(operators)-1]
				if prevOperator == '*' {
					operands[len(operands)-1] *= num
					operators = operators[:len(operators)-1] // Pop the previous operator
				} else if prevOperator == '/' {
					if num == 0 {
						return 0, errors.New("division by zero")
					}
					operands[len(operands)-1] /= num
					operators = operators[:len(operators)-1] // Pop the previous operator
				} else {
					operands = append(operands, num)
				}
			}
		}
	}

	result := operands[0]
	for i, op := range operators {
		if op == '+' {
			result += operands[i+1]
		} else if op == '-' {
			result -= operands[i+1]
		}
	}

	return result, nil
}
