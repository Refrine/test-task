package database

import (
    "fmt"
    "subscription-service/internal/config"
    "subscription-service/internal/models"
    
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "github.com/sirupsen/logrus"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
        cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode)
    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        return nil, err
    }

    
    if err := db.AutoMigrate(&models.Subscription{}); err != nil {
        return nil, err
    }

    logrus.Info("Database connected and migrated")
    return db, nil
}