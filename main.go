package main

import (
	"mobility-server/app/controller"
	"mobility-server/app/database"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDatabase()
	router := gin.Default()

	api := router.Group("/api")

	authRoute := api.Group("/auth")
	{
		authRoute.GET("/get-otp", controller.GetOtp)
		authRoute.GET("/verify-otp", controller.VerifyOTP)
	}

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	serveErr := server.ListenAndServe()
	if serveErr != nil {
		panic(serveErr)
	}
}
