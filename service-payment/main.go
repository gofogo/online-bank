package main

import (
	"service-payment/service"
	"github.com/gin-gonic/gin"
)

const (
	PRODUCER_URL string = "kafka:9092"
	KAFKA_TOPIC  string = "paym"
)

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(CORSMiddleware())
	v1 := r.Group("/")
	service.ServiceCheckRegister(v1.Group("/system"))
	service.PaymentRegister(v1.Group("/payment"))
	return r
}

type Extension struct {

}
type Blueprint struct {
        extension *Extension
	router *gin.RouterGroup
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}
