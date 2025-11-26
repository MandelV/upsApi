package models

type UPSStatus struct {
	BatteryCharge    int    `json:"battery.charge"`     // %
	BatteryChargeLow int    `json:"battery.charge.low"` // %
	BatteryRuntime   int    `json:"battery.runtime"`    // secondes
	BatteryType      string `json:"battery.type"`

	DeviceMfr    string `json:"device.mfr"`
	DeviceModel  string `json:"device.model"`
	DeviceSerial string `json:"device.serial"`
	DeviceType   string `json:"device.type"`

	DriverDebug           int    `json:"driver.debug"`
	DriverAllowKillPower  bool   `json:"driver.flag.allow_killpower"`
	DriverName            string `json:"driver.name"`
	DriverPollFreq        int    `json:"driver.parameter.pollfreq"`
	DriverPollInterval    int    `json:"driver.parameter.pollinterval"`
	DriverPort            string `json:"driver.parameter.port"`
	DriverSynchronous     string `json:"driver.parameter.synchronous"`
	DriverVendorID        string `json:"driver.parameter.vendorid"`
	DriverState           string `json:"driver.state"`
	DriverVersion         string `json:"driver.version"`
	DriverVersionData     string `json:"driver.version.data"`
	DriverVersionInternal string `json:"driver.version.internal"`
	DriverVersionUSB      string `json:"driver.version.usb"`

	InputTransferHigh float64 `json:"input.transfer.high"`
	InputTransferLow  float64 `json:"input.transfer.low"`

	Outlet1Desc       string `json:"outlet.1.desc"`
	Outlet1ID         int    `json:"outlet.1.id"`
	Outlet1Status     string `json:"outlet.1.status"`     // "on"/"off"
	Outlet1Switchable bool   `json:"outlet.1.switchable"` // "yes"/"no" → bool

	OutletDesc       string `json:"outlet.desc"`
	OutletID         int    `json:"outlet.id"`
	OutletSwitchable bool   `json:"outlet.switchable"`

	OutputFreqNominal float64 `json:"output.frequency.nominal"`
	OutputVoltage     float64 `json:"output.voltage"`
	OutputVoltageNom  float64 `json:"output.voltage.nominal"`

	UPSBeeperStatus  string `json:"ups.beeper.status"` // "enabled"/"disabled"
	UPSDelayShutdown int    `json:"ups.delay.shutdown"`
	UPSDelayStart    int    `json:"ups.delay.start"`
	UPSFirmware      string `json:"ups.firmware"`
	UPSLoad          int    `json:"ups.load"` // %
	UPSMfr           string `json:"ups.mfr"`
	UPSModel         string `json:"ups.model"`
	UPSPowerNominal  int    `json:"ups.power.nominal"` // VA/W (dépend du modèle)
	UPSProductID     string `json:"ups.productid"`
	UPSRealPower     int    `json:"ups.realpower"` // W
	UPSSerial        string `json:"ups.serial"`
	UPSStatus        string `json:"ups.status"` // OL, OB, LB, etc.
	UPSTimerShutdown int    `json:"ups.timer.shutdown"`
	UPSTimerStart    int    `json:"ups.timer.start"`
	UPSType          string `json:"ups.type"`
	UPSVendorID      string `json:"ups.vendorid"`
}
