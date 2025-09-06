package transport

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Health struct {
	DB *sqlx.DB
}

// Register добавляет эндпоинты для проверки здоровья сервиса /healthz и /readyz
// @Summary Liveness probe
// @Description Проверяет, что сервис работает и может отвечать на запросы
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string{"status":"ok"}
// @Failure 503 {object} map[string]string{"status":"unhealthy","error":"<error message>"}
// @Router /healthz [get]
//
// @Summary Readiness probe
// @Description Проверяет, что сервис готов принимать трафик
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string{"status":"ready"}
// @Failure 503 {object} map[string]string{"status":"unready","error":"<error message>"}
// @Router /readyz [get]
func (h Health) Register(r *gin.Engine) {
	r.GET("/healthz", func(c *gin.Context) {
		if err := h.DB.Ping(); err != nil {
			c.JSON(503, gin.H{"status": "unhealthy", "error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "ok"})
	})
	r.GET("/readyz", func(c *gin.Context) {
		if err := h.DB.Ping(); err != nil {
			c.JSON(503, gin.H{"status": "unready", "error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "ready"})
	})
}
