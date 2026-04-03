package config

import (
    "github.com/spf13/viper"
    "github.com/sirupsen/logrus"
)

type Config struct {
    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string
    DBSSLMode  string
    ServerPort string
    LogLevel   string
}

func LoadConfig() (*Config, error) {
    viper.SetConfigFile(".env")
    viper.AutomaticEnv()

    if err := viper.ReadInConfig(); err != nil {
        logrus.Warn("no .env")
    }

    config := &Config{
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnv("DB_PORT", "5432"),
        DBUser:     getEnv("DB_USER", "postgres"),
        DBPassword: getEnv("DB_PASSWORD", "postgres"),
        DBName:     getEnv("DB_NAME", "subscriptions"),
        DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
        ServerPort: getEnv("SERVER_PORT", "8080"),
        LogLevel:   getEnv("LOG_LEVEL", "info"),
    }

    return config, nil
}

func getEnv(key, defaultValue string) string {
    if viper.IsSet(key) {
        return viper.GetString(key)
    }
    return defaultValue
}