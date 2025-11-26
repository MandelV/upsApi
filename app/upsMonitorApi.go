package app

import (
	"net"
	"net/http"

	"github.com/MandelV/raspberry-terraforming/upsmonitoring/models"
	"github.com/MandelV/raspberry-terraforming/upsmonitoring/services"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

type UpsMonitorAPI struct {
	vroom       *gin.Engine
	api         *gin.RouterGroup
	listener    net.Listener
	upscService services.IUpscService
}

func NewUpsMonitorAPI(engine *gin.Engine, listener net.Listener, upscService services.IUpscService) *UpsMonitorAPI {

	obj := &UpsMonitorAPI{
		vroom:       engine,
		listener:    listener,
		upscService: upscService,
	}

	obj.api = obj.vroom.Group("/api")
	obj.BuildRoute()
	return obj
}

func (s *UpsMonitorAPI) BuildRoute() {
	//s.vroom.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	s.api.GET("/read", s.endpoint_get_upsc_status)
	s.vroom.GET("/metrics", gin.WrapH(promhttp.Handler()))
	health := s.vroom.Group("/health")
	health.GET("/liveness", s.endpoint_liveness)
	health.GET("/readiness", s.endpoint_readiness)
}

func (s *UpsMonitorAPI) ListenAndServe() error {
	return s.vroom.RunListener(s.listener)
}

func (s *UpsMonitorAPI) endpoint_get_upsc_status(ctx *gin.Context) {

	ctx.Header("Cache-Control", "no-store")

	if status, err := s.upscService.Read(); err == nil {
		ctx.JSONP(http.StatusOK, status)
	} else {
		log.Err(err).Msg("error while reading upsc status")
		ctx.JSONP(http.StatusInternalServerError, models.HttpError{Code: http.StatusInternalServerError, Message: "Error while reading the status"})
	}
}

// endpoint_liveness godoc
// @Summary      Liveness probe
// @Description  Endpoint utilisé pour vérifier que le process est vivant.
// @Tags         health
// @Produce      plain
// @Success      200  {string}  string  "OK"
// @Router       /health/liveness [get]
func (s *UpsMonitorAPI) endpoint_liveness(ctx *gin.Context) {
	ctx.Header("Cache-Control", "no-store")
	ctx.IndentedJSON(http.StatusOK, "OK")
}

// endpoint_readiness godoc
// @Summary      Readiness probe
// @Description  Vérifie que tous les contrôleurs sont prêts. Retourne 503 si l'un d'eux échoue.
// @Tags         health
// @Produce      plain
// @Success      200  {string}  string  "OK"
// @Failure      503  {string}  string  "Service Unavailable"
// @Router       /health/readiness [get]
func (s *UpsMonitorAPI) endpoint_readiness(ctx *gin.Context) {
	ctx.Header("Cache-Control", "no-store")
	ctx.IndentedJSON(http.StatusOK, "OK")

}
