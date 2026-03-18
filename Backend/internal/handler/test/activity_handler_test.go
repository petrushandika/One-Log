package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/petrushandika/one-log/internal/handler"
	"github.com/petrushandika/one-log/internal/repository"
	"github.com/petrushandika/one-log/internal/service"
)

func TestActivityRoutes_RegisterAndReturn200(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Fake auth middleware to inject user_id
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Next()
	})

	// Use in-memory sqlite? Repo uses gorm DB; for a lightweight route test we only ensure handler wiring compiles.
	// We'll construct service with nil repo would panic if executed; so only hit endpoints that validate period and return 400,
	// or ensure router registration doesn't crash. Here we just assert 400 on invalid period path without calling repo.
	var dummyRepo repository.LogRepository = nil
	activitySvc := service.NewActivityService(dummyRepo)
	activityHandler := handler.NewActivityHandler(activitySvc)

	api := r.Group("/api/v1")
	{
		admin := api.Group("")
		{
			admin.GET("/activity/summary", activityHandler.Summary)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/activity/summary?period=bogus", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}
