// UPSStatus représente l'ensemble des données retournées par la commande `upsc`.
// Les tags JSON correspondent aux clés exactes de la sortie de upsc (ex: "battery.charge").
// Cette struct est remplie par ParseUPSOuput dans le service upsc, puis sérialisée
// en JSON sur l'endpoint /api/read et utilisée pour alimenter les métriques Prometheus.
type UPSStatus struct {
	// --- Batterie ---
	BatteryCharge    int    `json:"battery.charge"`     // Charge actuelle en %
	BatteryChargeLow int    `json:"battery.charge.low"` // Seuil de charge basse en %
	BatteryRuntime   int    `json:"battery.runtime"`    // Autonomie restante en secondes
	BatteryType      string `json:"battery.type"`       // Type de batterie (ex: "PbAc")

	// --- Appareil ---
	DeviceMfr    string `json:"device.mfr"`    // Fabricant
	DeviceModel  string `json:"device.model"`  // Modèle
	DeviceSerial string `json:"device.serial"` // Numéro de série
	DeviceType   string `json:"device.type"`   // Type d'appareil (ex: "ups")

	// --- Driver NUT ---
	DriverDebug           int    `json:"driver.debug"`
	DriverAllowKillPower  bool   `json:"driver.flag.allow_killpower"` // Autorise le kill power (0/1 -> bool)
	DriverName            string `json:"driver.name"`                 // Nom du driver NUT (ex: "usbhid-ups")
	DriverPollFreq        int    `json:"driver.parameter.pollfreq"`   // Fréquence de polling en secondes
	DriverPollInterval    int    `json:"driver.parameter.pollinterval"`
	DriverPort            string `json:"driver.parameter.port"`       // Port du device (ex: "auto")
	DriverSynchronous     string `json:"driver.parameter.synchronous"`
	DriverVendorID        string `json:"driver.parameter.vendorid"` // Vendor ID USB
	DriverState           string `json:"driver.state"`
	DriverVersion         string `json:"driver.version"`
	DriverVersionData     string `json:"driver.version.data"`
	DriverVersionInternal string `json:"driver.version.internal"`
	DriverVersionUSB      string `json:"driver.version.usb"`

	// --- Entrée secteur ---
	InputTransferHigh float64 `json:"input.transfer.high"` // Tension haute de transfert (V)
	InputTransferLow  float64 `json:"input.transfer.low"`  // Tension basse de transfert (V)

	// --- Prise commandée (outlet 1) ---
	Outlet1Desc       string `json:"outlet.1.desc"`
	Outlet1ID         int    `json:"outlet.1.id"`
	Outlet1Status     string `json:"outlet.1.status"`     // "on" / "off"
	Outlet1Switchable bool   `json:"outlet.1.switchable"` // "yes"/"no" -> bool

	// --- Prise principale ---
	OutletDesc       string `json:"outlet.desc"`
	OutletID         int    `json:"outlet.id"`
	OutletSwitchable bool   `json:"outlet.switchable"`

	// --- Sortie ---
	OutputFreqNominal float64 `json:"output.frequency.nominal"` // Fréquence nominale (Hz)
	OutputVoltage     float64 `json:"output.voltage"`           // Tension de sortie actuelle (V)
	OutputVoltageNom  float64 `json:"output.voltage.nominal"`   // Tension nominale (V)

	// --- UPS global ---
	UPSBeeperStatus  string `json:"ups.beeper.status"`  // "enabled" / "disabled"
	UPSDelayShutdown int    `json:"ups.delay.shutdown"` // Délai avant shutdown (s)
	UPSDelayStart    int    `json:"ups.delay.start"`    // Délai avant redémarrage (s)
	UPSFirmware      string `json:"ups.firmware"`
	UPSLoad          int    `json:"ups.load"`          // Charge de l'UPS en %
	UPSMfr           string `json:"ups.mfr"`           // Fabricant
	UPSModel         string `json:"ups.model"`         // Modèle
	UPSPowerNominal  int    `json:"ups.power.nominal"` // Puissance nominale (VA ou W selon le modèle)
	UPSProductID     string `json:"ups.productid"`     // Product ID USB
	UPSRealPower     int    `json:"ups.realpower"`     // Puissance réelle consommée (W)
	UPSSerial        string `json:"ups.serial"`
	UPSStatus        string `json:"ups.status"` // Statut NUT : OL=online, OB=on battery, LB=low battery
	UPSTimerShutdown int    `json:"ups.timer.shutdown"`
	UPSTimerStart    int    `json:"ups.timer.start"`
	UPSType          string `json:"ups.type"`
	UPSVendorID      string `json:"ups.vendorid"`
}
