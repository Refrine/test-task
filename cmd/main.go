package main

import (
    "fmt"
    "subscription-service/internal/config"
    "subscription-service/internal/database"
    "subscription-service/internal/handlers"
    "subscription-service/internal/middleware"
    "subscription-service/internal/repo"

	
    
    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
)

// @title Subscription service API
// @description REST API
// @host localhost:8080
// @BasePath /
func main() {
    // load cfg
    cfg, err := config.LoadConfig()
    if err != nil {
        logrus.Fatal("failed to load config:", err)
    }
    
    // setup logging
    level, _ := logrus.ParseLevel(cfg.LogLevel)
    logrus.SetLevel(level)
    logrus.SetFormatter(&logrus.JSONFormatter{})
    
    // connect 
    db, err := database.InitDB(cfg)
    if err != nil {
        logrus.Fatal("failed to connect", err)
    }
    
    // setup repo && handlers
    repo := repository.NewSubscriptionRepository(db)
    handler := handlers.NewSubscriptionHandler(repo)
    
    // setup 
    router := gin.Default()
    router.Use(middleware.Logger())
    
    // routes
    api := router.Group("/api/v1")
    {
        api.POST("/subscriptions", handler.Create)
        api.GET("/subscriptions", handler.GetAll)
        api.GET("/subscriptions/:id", handler.GetByID)
        api.PUT("/subscriptions/:id", handler.Update)
        api.DELETE("/subscriptions/:id", handler.Delete)
        api.GET("/subscriptions/aggregate/total", handler.Aggregate)
    }
    
    // start server
    addr := fmt.Sprintf(":%s", cfg.ServerPort)
    logrus.WithField("port", cfg.ServerPort).Info("start server")
    router.Run(addr)
}