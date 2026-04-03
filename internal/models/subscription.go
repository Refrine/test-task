package models

import (
    "time"
    "github.com/google/uuid"
)

type Subscription struct {
    ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
    ServiceName string     `gorm:"not null" json:"service_name"`
    Price       int        `gorm:"not null" json:"price"`
    UserID      uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
    StartDate   time.Time  `gorm:"not null" json:"start_date"`
    EndDate     *time.Time `json:"end_date,omitempty"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateSubscriptionRequest struct {
    ServiceName string     `json:"service_name" binding:"required"`
    Price       int        `json:"price" binding:"required,min=0"`
    UserID      string     `json:"user_id" binding:"required"`
    StartDate   string     `json:"start_date" binding:"required"`
    EndDate     *string    `json:"end_date,omitempty"`
}

type AggregateRequest struct {
    StartDate   string `form:"start_date" binding:"required"`
    EndDate     string `form:"end_date" binding:"required"`
    UserID      string `form:"user_id"`
    ServiceName string `form:"service_name"`
}