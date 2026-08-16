package httptransport

import (
	"context"
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/wyw14/cry005/internal/application"
)

var _ = context.Background

type Handler struct {
	service  *application.Service
	validate *validator.Validate
}

func New(service *application.Service) *Handler {
	return &Handler{service: service, validate: validator.New()}
}

func (h *Handler) Router() http.Handler {
	r := gin.New()
	api := r.Group("/api/v1")
	api.GET("/items", h.list)
	api.POST("/items", h.create)
	api.GET("/events", h.replay)
	api.POST("/workflow", h.workflow)
	r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
	r.GET("/readyz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ready"}) })
	return r
}

func (h *Handler) list(c *gin.Context) {
	scope := "*"

	items, err := h.service.List(c.Request.Context(), scope)
	if err != nil {
		c.JSON(500, gin.H{"code": "LIST_FAILED", "message": err.Error()})
		return
	}
	c.JSON(200, items)
}

type createRequest struct {
	Scope   string `json:"scope" validate:"required"`
	Actor   string `json:"actor" validate:"required"`
	Payload string `json:"payload" validate:"required"`
}

func (h *Handler) create(c *gin.Context) {
	var req createRequest
	if c.ShouldBindJSON(&req) != nil || h.validate.Struct(req) != nil {
		c.JSON(422, gin.H{"code": "VALIDATION_FAILED"})
		return
	}
	key := c.GetHeader("Idempotency-Key")

	item, err := h.service.Create(c.Request.Context(), req.Scope, req.Actor, key, req.Payload)
	if err != nil {
		c.JSON(409, gin.H{"code": "CREATE_FAILED", "message": err.Error()})
		return
	}
	c.JSON(201, item)
}

func (h *Handler) replay(c *gin.Context) {
	after, _ := strconv.ParseInt(c.GetHeader("Last-Event-ID"), 10, 64)

	events, err := h.service.Replay(c.Request.Context(), after)
	if err != nil {
		c.JSON(500, gin.H{"code": "REPLAY_FAILED"})
		return
	}
	c.JSON(200, events)
}

func (h *Handler) workflow(c *gin.Context) {
	ctx := c.Request.Context()

	failAt := c.Query("fail_at")

	if err := h.service.ApplyAtomic(ctx, "completed", failAt); err != nil {
		c.JSON(409, gin.H{"code": "WORKFLOW_FAILED"})
		return
	}
	c.Status(204)
}
