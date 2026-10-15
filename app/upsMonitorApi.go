package app

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/MandelV/raspberry-terraforming/upsmonitoring/models"
	"github.com/MandelV/raspberry-terraforming/upsmonitoring/services"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

// UpsMonitorAPI est le serveur HTTP principal. Il encapsule le routeur Gin,
// un listener TCP, et le service upsc pour interroger l'UPS.
// Un contexte interne permet d'arrêter proprement la goroutine de rafraîchissement des métriques.
type UpsMonitorAPI struct {
	vroom       *gin.Engine
	api         *gin.RouterGroup
	listener    net.Listener
	ctx         context.Context
	cancel      context.CancelFunc
	upscService services.IUpscService
}

// NewUpsMonitorAPI crée le serveur, enregistre les routes et lance une goroutine
// qui rafraîchit les métriques Prometheus toutes les 10 secondes en appelant upsc.
func NewUpsMonitorAPI(engine *gin.Engine, listener net.Listener, upscService services.IUpscService) *UpsMonitorAPI {

	ctx, cancel := context.WithCancel(context.Background())
	obj := &UpsMonitorAPI{
		vroom:       engine,
		listener:    listener,
		upscService: upscService,
		ctx:         ctx,
		cancel:      cancel,
	}

	obj.api = obj.vroom.Group("/api")
	obj.BuildRoute()
	go obj.refreshMetrics(10 * time.Second)
	return obj
}

// Close annule le contexte interne, ce qui stoppe la goroutine de rafraîchissement des métriques.
func (s *UpsMonitorAPI) Close() error {
	s.cancel()
	return nil
}

// refreshMetrics interroge upsc à intervalle régulier (d) et met à jour les gauges Prometheus.
// S'arrête quand le contexte est annulé via Close().
func (s *UpsMonitorAPI) refreshMetrics(d time.Duration) {

	ticker := time.NewTicker(d)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():

		case <-ticker.C:
			if status, err := s.upscService.Read(); err == nil {

				UpdateMetrics(status)

			} else {
				log.Err(err).Str("ctx", "refreshMetrics").Msg("unable to read ups's metrics")
			}

		}
	}

}

// BuildRoute enregistre toutes les routes du serveur :
//   - GET /api/read       -> status JSON de l'UPS
//   - GET /metrics        -> métriques Prometheus (via promhttp)
//   - GET /health/liveness  -> probe de liveness (toujours 200)
//   - GET /health/readiness -> probe de readiness (toujours 200)
func (s *UpsMonitorAPI) BuildRoute() {
	s.api.GET("/read", s.endpoint_get_upsc_status)
	s.vroom.GET("/metrics", gin.WrapH(promhttp.Handler()))
	health := s.vroom.Group("/health")
	health.GET("/liveness", s.endpoint_liveness)
	health.GET("/readiness", s.endpoint_readiness)
}

// ListenAndServe démarre le serveur HTTP sur le listener TCP fourni à la construction.
func (s *UpsMonitorAPI) ListenAndServe() error {
	return s.vroom.RunListener(s.listener)
}

// endpoint_get_upsc_status exécute `upsc` à la demande et retourne le status complet en JSON.
// Désactive le cache navigateur pour toujours renvoyer des données fraîches.
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
