package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/thejasmeetsingh/go-ecommerce/product_service/api"
	"github.com/thejasmeetsingh/go-ecommerce/product_service/internal/database"
)

func main() {
	godotenv.Load()
	engine := gin.Default()

	// Load env varriables
	port := os.Getenv("PORT")
	mode := os.Getenv("GIN_MODE")
	dbURL := os.Getenv("DB_URL")

	if port == "" {
		panic("Port is not configured")
	}

	if mode == "" || mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	if dbURL == "" {
		panic("DB URL is not configured")
	}

	// Get DB connection
	dbConn := api.GetDBConn(dbURL)
	defer dbConn.Close()

	// Get Redis connection
	redisClient := api.GetRedisClient()
	defer redisClient.Close()

	apiCfg := api.APIConfig{
		DB:      dbConn,
		Queries: database.New(dbConn),
		Cache:   redisClient,
	}

	// Initialize prometheus
	httpRequestsTotal := api.GetPromRequestTotal()
	httpRequestDuration := api.GetPromRequestDuration()

	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)

	// Load API routes
	api.GetRoutes(engine, &apiCfg)

	log.Infoln("Product services started!")
	engine.Run(":" + port)
}
