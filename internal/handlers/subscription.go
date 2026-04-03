package handlers

import (
    "net/http"
    "time"
    "subscription-service/internal/models"
    "subscription-service/internal/repo"
    
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "github.com/sirupsen/logrus"
)

type SubscriptionHandler struct {
    repo *repository.SubscriptionRepository
}

func NewSubscriptionHandler(repo *repository.SubscriptionRepository) *SubscriptionHandler {
    return &SubscriptionHandler{repo: repo}
}

// CreateSubscription godoc
// @Summary Create new subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Success 201 {object} models.Subscription
// @Router /subscriptions [post]
func (h *SubscriptionHandler) Create(c *gin.Context) {
    var req models.CreateSubscriptionRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        logrus.WithError(err).Error("invalid request body")
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    startDate, err := time.Parse("01-2006", req.StartDate)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format"})
        return
    }

    userUUID, err := uuid.Parse(req.UserID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id UUID"})
        return
    }

    sub := &models.Subscription{
        ServiceName: req.ServiceName,
        Price:       req.Price,
        UserID:      userUUID,
        StartDate:   startDate,
    }

    if req.EndDate != nil && *req.EndDate != "" {
        endDate, err := time.Parse("01-2006", *req.EndDate)
        if err == nil {
            sub.EndDate = &endDate
        }
    }

    if err := h.repo.Create(sub); err != nil {
        logrus.WithError(err).Error("fail create sub")
        c.JSON(http.StatusInternalServerError, gin.H{"error": "fail sub"})
        return
    }

    c.JSON(http.StatusCreated, sub)
}

// GetSubscriptions godoc
// @Summary Get all subscriptions
// @Tags subscriptions
// @Produce json
// @Success 200 {array} models.Subscription
// @Router /subscriptions [get]
func (h *SubscriptionHandler) GetAll(c *gin.Context) {
    subs, err := h.repo.GetAll()
    if err != nil {
        logrus.WithError(err).Error("fail to sub")
        c.JSON(http.StatusInternalServerError, gin.H{"error": "fail to sub"})
        return
    }
    c.JSON(http.StatusOK, subs)
}

// GetSubscription godoc
// @Summary Get subscription by ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} models.Subscription
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetByID(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
        return
    }

    sub, err := h.repo.GetByID(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "sub not found"})
        return
    }
    c.JSON(http.StatusOK, sub)
}

// UpdateSubscription godoc
// @Summary Update subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Success 200 {object} models.Subscription
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) Update(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
        return
    }

    var req models.CreateSubscriptionRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    startDate, _ := time.Parse("01-2006", req.StartDate)
    userUUID, _ := uuid.Parse(req.UserID)

    sub := &models.Subscription{
        ServiceName: req.ServiceName,
        Price:       req.Price,
        UserID:      userUUID,
        StartDate:   startDate,
    }

    if err := h.repo.Update(id, sub); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "fail  update"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "updated ok"})
}

// DeleteSubscription godoc
// @Summary Delete subscription
// @Tags subscriptions
// @Produce json
// @Success 204 "No Content"
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) Delete(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
        return
    }

    if err := h.repo.Delete(id); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail delete"})
        return
    }
    c.Status(http.StatusNoContent)
}

// AggregateTotal godoc
// @Summary Get total price of subscriptions in period
// @Tags aggregate
// @Produce json
// @Success 200 {object} map[string]int64
// @Router /subscriptions/aggregate [get]
func (h *SubscriptionHandler) Aggregate(c *gin.Context) {
    var req models.AggregateRequest
    if err := c.ShouldBindQuery(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    startDate, err := time.Parse("01-2006", req.StartDate)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format"})
        return
    }

    endDate, err := time.Parse("01-2006", req.EndDate)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format"})
        return
    }

    total, err := h.repo.Aggregate(startDate, endDate, req.UserID, req.ServiceName)
    if err != nil {
        logrus.WithError(err).Error("failed to aggregate")
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to total"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"total_price": total})
}