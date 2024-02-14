package main

import (
	"github.com/RakhimovAns/Calculus/Initializers"
	"github.com/RakhimovAns/Calculus/pkg/handler"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
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
	r := gin.Default()
	r.StaticFile("/expression", "./static/index.html")
	r.StaticFile("/calculate", "./static/calculate.html")
	r.StaticFile("/result", "./static/result.html")
	r.POST("/expression", handler.PostExpression)
	r.POST("/calculate/:id", handler.StartCount)
	r.GET("/result/:id", handler.GetStatus)
	r.Run()
}
