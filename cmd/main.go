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
	"sync"
	"time"
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
	if err := IsValidate(expression.Expression); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid format of expression"})
		return
	}

	result, err := EvaluateExpression(expression.Expression, expression.AddTime, expression.SubTime, expression.MultiplyTime, expression.DivideTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
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

func EvaluateExpression(expression string, AddTime, SubTime, MultiplyTime, DivideTime int64) (float64, error) {
	var operands []float64
	var operators []rune
	numStr := ""

	for _, char := range expression {
		switch char {
		case '+', '-', '*', '/':
			num, err := strconv.ParseFloat(numStr, 64)
			if err != nil {
				return 0, err
			}
			operands = append(operands, num)
			operators = append(operators, char)
			numStr = ""
		default:
			numStr += string(char)
		}
	}

	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, err
	}
	operands = append(operands, num)

	resultCh := make(chan float64, len(operators))
	errCh := make(chan error)

	var wg sync.WaitGroup

	for i := 0; i < len(operators); {
		switch operators[i] {
		case '*':
			wg.Add(1)
			go MultiplyAsync(operands, i, MultiplyTime, resultCh, errCh, &wg)
		case '/':
			wg.Add(1)
			go DivideAsync(operands, i, DivideTime, resultCh, errCh, &wg)
		case '+':
			wg.Add(1)
			go AddAsync(operands, i, AddTime, resultCh, errCh, &wg)
		case '-':
			wg.Add(1)
			go SubtractAsync(operands, i, SubTime, resultCh, errCh, &wg)
		}
		i++
	}

	wg.Wait()
	var result float64
	for range operators {
		select {
		case res := <-resultCh:
			result = res
		case err := <-errCh:
			return 0, err
		}
	}

	return result, nil
}

func AddAsync(operands []float64, i int, Timer int64, resultCh chan<- float64, errCh chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Duration(Timer) * time.Second)
	operands[i] += operands[i+1]
	resultCh <- operands[i]
}

func SubtractAsync(operands []float64, i int, Timer int64, resultCh chan<- float64, errCh chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Duration(Timer) * time.Second)
	operands[i] -= operands[i+1]
	resultCh <- operands[i]
}

func MultiplyAsync(operands []float64, i int, Timer int64, resultCh chan<- float64, errCh chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Duration(Timer) * time.Second)
	operands[i] *= operands[i+1]
	resultCh <- operands[i]
}

func DivideAsync(operands []float64, i int, Timer int64, resultCh chan<- float64, errCh chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Duration(Timer) * time.Second)
	if operands[i+1] == 0 {
		errCh <- errors.New("division by zero")
		return
	}
	operands[i] /= operands[i+1]
	resultCh <- operands[i]
}
