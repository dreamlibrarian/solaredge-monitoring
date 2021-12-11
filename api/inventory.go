package api

type InventoryDocument struct {
	Inventory Inventory `json:"Inventory"`
}

type Inventory struct {
	Inverters           []Inverter  `json:"inverters"`
	ThirdPartyInverters []Inverter  `json:"thirdPartyInverters"`
	SMIDevices          []SMIDevice `json:"smiList"`
	Meters              []Meter     `json:"meters"`
	Sensors             []Sensor    `json:"sensors"`
	Gateways            []Gateway   `json:"gateways"`
	Batteries           []Battery   `json:"batteries"`
}

type Inverter struct {
	Name                string `json:"name"`
	Manufacturer        string `json:"manufacturer"`
	Model               string `json:"model"`
	CommunicationMethod string `json:"communicationMethod"`
	CPUVersion          string `json:"cpuVersion"`
	FirmwareVersion     string `json:"firmwareVersion"`
	SerialNumber        string `json:"SN"`
	ConnectedOptimizers int64  `json:"connectedOptimizers"`
}

type Meter struct {
	Name                       string `json:"name"`
	Manufacturer               string `json:"manufacturer"`
	Model                      string `json:"model"`
	SerialNumber               string `json:"SN"`
	Type                       string `json:"type"`
	FirmwareVersion            string `json:"firmwareVersion"`
	ConnectedTo                string `json:"connectedTo"`
	ConnectedSolaredgeDeviceSN string `json:"connectedSolaredgeDeviceSN"`
	Form                       string `json:"form"`
}

type Sensor struct {
	ID                         string `json:"id"`
	ConnectedTo                string `json:"connectedTo"`
	Category                   string `json:"category"`
	ConnectedSolaredgeDeviceSN string `json:"connectedSolaredgeDeviceSN"`
}

type Gateway struct {
	Name            string `json:"name"`
	SerialNumber    string `json:"serialNumber"`
	FirmwareVersion string `json:"firmwareVersion"`
}

type Battery struct {
	Name                       string `json:"name"`
	SerialNumber               string `json:"serialNumber"`
	Manufacturer               string `json:"manufacturer"`
	Model                      string `json:"model"`
	NameplateCapacity          string `json:"nameplateCapacity"`
	FirmwareVersion            string `json:"firmwareVersion"`
	ConnectedTo                string `json:"connectedTo"`
	ConnectedSolaredgeDeviceSN string `json:"connectedSolaredgeDeviceSN"`
}

type SMIDevice struct {
	Name                string `json:"name"`
	Manufacturer        string `json:"manufacturer"`
	Model               string `json:"model"`
	FirmwareVersion     string `json:"firmwareVersion"`
	CommunicationMethod string `json:"communicationMethod"`
	SerialNumber        string `json:"serialNumber"`
	ConnectedOptimizers string `json:"connectedOptimizers"`
}
