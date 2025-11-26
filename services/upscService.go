package services

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/MandelV/raspberry-terraforming/upsmonitoring/models"
)

type IUpscService interface {
	Read() (*models.UPSStatus, error)
}

type upscService struct {
	user string
	host string
}

func NewUpscService(user string, host string) IUpscService {
	return &upscService{
		user: user,
		host: host,
	}
}

// Read
func (u *upscService) Read() (*models.UPSStatus, error) {
	app := "upsc"
	arg0 := fmt.Sprintf("%s@%s", u.user, u.host)
	cmd := exec.Command(app, arg0)
	stdout, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	//Parse & output
	return ParseUPSOuput(stdout)
}

func ParseUPSOuput(b []byte) (*models.UPSStatus, error) {
	status := &models.UPSStatus{}
	scanner := bufio.NewScanner(bytes.NewReader(b))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// On coupe uniquement sur le premier ":", pour gérer les valeurs qui en contiennent
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		switch key {
		// --- battery ---
		case "battery.charge":
			status.BatteryCharge = atoi(val)
		case "battery.charge.low":
			status.BatteryChargeLow = atoi(val)
		case "battery.runtime":
			status.BatteryRuntime = atoi(val)
		case "battery.type":
			status.BatteryType = val

		// --- device ---
		case "device.mfr":
			status.DeviceMfr = val
		case "device.model":
			status.DeviceModel = val
		case "device.serial":
			status.DeviceSerial = val
		case "device.type":
			status.DeviceType = val

		// --- driver ---
		case "driver.debug":
			status.DriverDebug = atoi(val)
		case "driver.flag.allow_killpower":
			status.DriverAllowKillPower = parseBool01(val)
		case "driver.name":
			status.DriverName = val
		case "driver.parameter.pollfreq":
			status.DriverPollFreq = atoi(val)
		case "driver.parameter.pollinterval":
			status.DriverPollInterval = atoi(val)
		case "driver.parameter.port":
			status.DriverPort = val
		case "driver.parameter.synchronous":
			status.DriverSynchronous = val
		case "driver.parameter.vendorid":
			status.DriverVendorID = val
		case "driver.state":
			status.DriverState = val
		case "driver.version":
			status.DriverVersion = val
		case "driver.version.data":
			status.DriverVersionData = val
		case "driver.version.internal":
			status.DriverVersionInternal = val
		case "driver.version.usb":
			status.DriverVersionUSB = val

		// --- input ---
		case "input.transfer.high":
			status.InputTransferHigh = atof(val)
		case "input.transfer.low":
			status.InputTransferLow = atof(val)

		// --- outlet 1 ---
		case "outlet.1.desc":
			status.Outlet1Desc = val
		case "outlet.1.id":
			status.Outlet1ID = atoi(val)
		case "outlet.1.status":
			status.Outlet1Status = val
		case "outlet.1.switchable":
			status.Outlet1Switchable = parseBoolYesNo(val)

		// --- outlet global ---
		case "outlet.desc":
			status.OutletDesc = val
		case "outlet.id":
			status.OutletID = atoi(val)
		case "outlet.switchable":
			status.OutletSwitchable = parseBoolYesNo(val)

		// --- output ---
		case "output.frequency.nominal":
			status.OutputFreqNominal = atof(val)
		case "output.voltage":
			status.OutputVoltage = atof(val)
		case "output.voltage.nominal":
			status.OutputVoltageNom = atof(val)

		// --- ups ---
		case "ups.beeper.status":
			status.UPSBeeperStatus = val
		case "ups.delay.shutdown":
			status.UPSDelayShutdown = atoi(val)
		case "ups.delay.start":
			status.UPSDelayStart = atoi(val)
		case "ups.firmware":
			status.UPSFirmware = val
		case "ups.load":
			status.UPSLoad = atoi(val)
		case "ups.mfr":
			status.UPSMfr = val
		case "ups.model":
			status.UPSModel = val
		case "ups.power.nominal":
			status.UPSPowerNominal = atoi(val)
		case "ups.productid":
			status.UPSProductID = val
		case "ups.realpower":
			status.UPSRealPower = atoi(val)
		case "ups.serial":
			status.UPSSerial = val
		case "ups.status":
			status.UPSStatus = val
		case "ups.timer.shutdown":
			status.UPSTimerShutdown = atoi(val)
		case "ups.timer.start":
			status.UPSTimerStart = atoi(val)
		case "ups.type":
			status.UPSType = val
		case "ups.vendorid":
			status.UPSVendorID = val
		}
	}

	return status, scanner.Err()
}

func atoi(s string) int {
	s = strings.TrimSpace(s)
	i, err := strconv.Atoi(s)
	if err != nil {
		return -1
	}
	return i
}

func atof(s string) float64 {
	s = strings.TrimSpace(s)
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return -1.0
	}
	return f
}

// "yes"/"no" -> bool
func parseBoolYesNo(s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))
	return s == "yes" || s == "enabled" || s == "on"
}

// "0"/"1" -> bool (driver.flag.allow_killpower)
func parseBool01(s string) bool {
	s = strings.TrimSpace(s)
	return s == "1" || s == "true" || s == "yes"
}
