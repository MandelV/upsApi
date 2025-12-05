package app

import (
	"github.com/MandelV/raspberry-terraforming/upsmonitoring/models"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Battery metrics
	batteryCharge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "upsapi_battery_charge",
		Help: "Battery charge percentage",
	})
	batteryChargeLow = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "upsapi_battery_charge_low",
		Help: "Battery low charge threshold percentage",
	})
	batteryRuntime = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "upsapi_battery_runtime",
		Help: "Battery runtime in seconds",
	})

	// Input metrics
	inputTransferHigh = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "upsapi_input_transfer_high",
		Help: "Input transfer high voltage",
	})
	inputTransferLow = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "upsapi_input_transfer_low",
		Help: "Input transfer low voltage",
	})

	// Output metrics
	outputFreqNominal = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "upsapi_output_frequency_nominal",
		Help: "Output nominal frequency",
	})
	outputVoltage = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "upsapi_output_voltage",
		Help: "Output voltage",
	})
	outputVoltageNom = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "upsapi_output_voltage_nominal",
		Help: "Output nominal voltage",
	})

	// UPS metrics
	upsLoad = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "upsapi_ups_load",
		Help: "UPS load percentage",
	})
	upsPowerNominal = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "upsapi_ups_power_nominal",
		Help: "UPS nominal power in VA/W",
	})
	upsRealPower = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "upsapi_ups_real_power",
		Help: "UPS real power in Watts",
	})
	upsDelayShutdown = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "upsapi_ups_delay_shutdown",
		Help: "UPS delay shutdown in seconds",
	})
	upsDelayStart = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "upsapi_ups_delay_start",
		Help: "UPS delay start in seconds",
	})
	upsTimerShutdown = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "upsapi_ups_timer_shutdown",
		Help: "UPS timer shutdown in seconds",
	})
	upsTimerStart = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "upsapi_ups_timer_start",
		Help: "UPS timer start in seconds",
	})

	// Status info metrics (using labels)
	upsStatusInfo = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "upsapi_ups_info",
		Help: "UPS information with labels",
	}, []string{"status", "manufacturer", "model", "serial", "firmware", "type"})

	upsBeeperEnabled = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "upsapi_ups_beeper_enabled",
		Help: "UPS beeper status (1=enabled, 0=disabled)",
	})

	// Outlet metrics
	outlet1Status = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "upsapi_outlet_1_status",
		Help: "Outlet 1 status (1=on, 0=off)",
	})
	outlet1Switchable = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "upsapi_outlet_1_switchable",
		Help: "Outlet 1 switchable (1=yes, 0=no)",
	})
	outletSwitchable = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "upsapi_outlet_switchable",
		Help: "Outlet switchable (1=yes, 0=no)",
	})
)

// UpdateMetrics updates all Prometheus metrics from UPS status
func UpdateMetrics(status *models.UPSStatus) {
	// Battery metrics
	batteryCharge.Set(float64(status.BatteryCharge))
	batteryChargeLow.Set(float64(status.BatteryChargeLow))
	batteryRuntime.Set(float64(status.BatteryRuntime))

	// Input metrics
	inputTransferHigh.Set(status.InputTransferHigh)
	inputTransferLow.Set(status.InputTransferLow)

	// Output metrics
	outputFreqNominal.Set(status.OutputFreqNominal)
	outputVoltage.Set(status.OutputVoltage)
	outputVoltageNom.Set(status.OutputVoltageNom)

	// UPS metrics
	upsLoad.Set(float64(status.UPSLoad))
	upsPowerNominal.Set(float64(status.UPSPowerNominal))
	upsRealPower.Set(float64(status.UPSRealPower))
	upsDelayShutdown.Set(float64(status.UPSDelayShutdown))
	upsDelayStart.Set(float64(status.UPSDelayStart))
	upsTimerShutdown.Set(float64(status.UPSTimerShutdown))
	upsTimerStart.Set(float64(status.UPSTimerStart))

	// UPS status info (with labels)
	upsStatusInfo.WithLabelValues(
		status.UPSStatus,
		status.UPSMfr,
		status.UPSModel,
		status.UPSSerial,
		status.UPSFirmware,
		status.UPSType,
	).Set(1)

	// Beeper status
	if status.UPSBeeperStatus == "enabled" {
		upsBeeperEnabled.Set(1)
	} else {
		upsBeeperEnabled.Set(0)
	}

	// Outlet metrics
	if status.Outlet1Status == "on" {
		outlet1Status.Set(1)
	} else {
		outlet1Status.Set(0)
	}

	if status.Outlet1Switchable {
		outlet1Switchable.Set(1)
	} else {
		outlet1Switchable.Set(0)
	}

	if status.OutletSwitchable {
		outletSwitchable.Set(1)
	} else {
		outletSwitchable.Set(0)
	}
}
