package main

// @title Wallet API
// @version 1.0
// @description A simple wallet API
// @host localhost:9000
// @BasePath /

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "wallet-api/docs"
	"wallet-api/handlers"
	"wallet-api/store"
	custommw "wallet-api/middleware"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	godotenv.Load()
	store.InitDB()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e := echo.New()

	e.Use(custommw.Logger)
	e.Use(middleware.CORS())
	e.Use(custommw.RateLimiter)
	
	

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.POST("/wallets", handlers.CreateWallet)
	e.GET("/wallets/:id", handlers.GetWallet)
	e.POST("/wallets/:id/transactions", handlers.AddTransaction)
	e.GET("/wallets/:id/transactions", handlers.GetTransactions)

	fmt.Println("Server running on port", port)
	e.Start(":" + port)
}
