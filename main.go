package main

import (
	"fmt"
	"net"

	"github.com/MandelV/raspberry-terraforming/upsmonitoring/app"
	"github.com/MandelV/raspberry-terraforming/upsmonitoring/models"
	"github.com/MandelV/raspberry-terraforming/upsmonitoring/services"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	config := models.NewConfiguration()
	log.Logger = log.Level(zerolog.Level(config.LOG_LEVEL))

	if config.ENV != "dev" {
		gin.SetMode(gin.ReleaseMode)
	}

	log.Info().Msg("START UPS MONITORING API")

	if l, err := net.Listen("tcp", fmt.Sprintf(":%d", config.LISTEN_PORT)); err == nil {
		upsMonitor := app.NewUpsMonitorAPI(gin.Default(), l, services.NewUpscService(config.UPS_USER, config.UPS_HOST))

		if err := upsMonitor.ListenAndServe(); err != nil {
			log.Err(err).Msg("error at ListenAndServe")
		}

	} else {
		log.Fatal().Err(err).Msg("unable to listen")
	}

}
