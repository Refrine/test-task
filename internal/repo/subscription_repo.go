package repository

import (
    "time"
    "subscription-service/internal/models"
    
    "gorm.io/gorm"
    "github.com/google/uuid"
    "github.com/sirupsen/logrus"
)

type SubscriptionRepository struct {
    db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository{
	return &SubscriptionRepository{db: db}
}


func (r *SubscriptionRepository) Create(sub *models.Subscription) error {
    logrus.WithField("service", sub.ServiceName).Info("Create sub")
    return r.db.Create(sub).Error
}

func (r *SubscriptionRepository) GetByID(id uuid.UUID) (*models.Subscription, error) {
    var sub models.Subscription
    err := r.db.First(&sub, "id = ?", id).Error
    return &sub, err
}

func (r *SubscriptionRepository) GetAll() ([]models.Subscription, error) {
    var subs []models.Subscription
    err := r.db.Find(&subs).Error
    return subs, err
}

func (r *SubscriptionRepository) Update(id uuid.UUID, sub *models.Subscription) error {
    logrus.WithField("id", id).Info("Update sub")
    return r.db.Model(&models.Subscription{}).Where("id = ?", id).Updates(sub).Error
}

func (r *SubscriptionRepository) Delete(id uuid.UUID) error {
    logrus.WithField("id", id).Info("Delete sub")
    return r.db.Delete(&models.Subscription{}, "id = ?", id).Error
}

func (r *SubscriptionRepository) Aggregate(startDate, endDate time.Time, userID, serviceName string) (int64, error) {
    var total int64
    
    query := r.db.Model(&models.Subscription{}).
        Where("start_date >= ? AND start_date <= ?", startDate, endDate)
    
    if userID != "" {
        query = query.Where("user_id = ?", userID)
    }
    
    if serviceName != "" {
        query = query.Where("service_name = ?", serviceName)
    }
    
    err := query.Select("COALESCE(SUM(price), 0)").Scan(&total).Error
    logrus.WithFields(logrus.Fields{
        "total": total,
        "user_id": userID,
        "service": serviceName,
    }).Info("Aggregated total price")
    
    return total, err
}