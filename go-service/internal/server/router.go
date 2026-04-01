package server

import (
	"errors"
	"net/http"
	"strings"

	"laboratorywork10/go-service/internal/auth"
	"laboratorywork10/go-service/internal/models"

	"github.com/gin-gonic/gin"
)

const (
	demoUsername = "student"
	demoPassword = "securepass123"
	demoRole     = "integration-client"
)

type Router struct {
	authService *auth.Service
}

func NewRouter(authService *auth.Service) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := &Router{authService: authService}
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	engine.GET("/health", router.health)
	engine.POST("/auth/token", router.issueToken)

	protected := engine.Group("/api", router.authMiddleware())
	protected.POST("/process", router.processPayload)

	return engine
}

func (r *Router) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (r *Router) issueToken(c *gin.Context) {
	var request models.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.Username != demoUsername || request.Password != demoPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := r.authService.GenerateToken(request.Username, demoRole)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"type":  "Bearer",
	})
}

func (r *Router) processPayload(c *gin.Context) {
	claimsValue, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing claims"})
		return
	}

	claims, ok := claimsValue.(*auth.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
		return
	}

	var request models.ProcessRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var total float64
	for _, item := range request.Items {
		total += float64(item.Quantity) * item.Price
	}

	response := models.ProcessResponse{
		RequestID:   request.RequestID,
		ApprovedBy:  claims.Username,
		ItemsCount:  len(request.Items),
		TotalAmount: total,
		Tags:        request.Metadata.Tags,
		Status:      "accepted",
	}

	c.JSON(http.StatusOK, response)
}

func (r *Router) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			return
		}

		token, err := extractBearerToken(header)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		claims, err := r.authService.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token validation failed"})
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}

func extractBearerToken(header string) (string, error) {
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
		return "", errors.New("invalid authorization header")
	}

	return strings.TrimSpace(parts[1]), nil
}
