package tasmota

import (
	"context"
	"encoding/json"
	"time"
)

// StatusResponse represents the complete status response from a Tasmota device.
type StatusResponse struct {
	Status    *StatusInfo     `json:"Status,omitempty"`
	StatusPRM *StatusParam    `json:"StatusPRM,omitempty"`
	StatusFWR *StatusFirmware `json:"StatusFWR,omitempty"`
	StatusLOG *StatusLog      `json:"StatusLOG,omitempty"`
	StatusMEM *StatusMemory   `json:"StatusMEM,omitempty"`
	StatusNET *StatusNetwork  `json:"StatusNET,omitempty"`
	StatusMQT *StatusMQTT     `json:"StatusMQT,omitempty"`
	StatusTIM *StatusTime     `json:"StatusTIM,omitempty"`
	StatusSNS *StatusSensor   `json:"StatusSNS,omitempty"`
	StatusSTS *StatusState    `json:"StatusSTS,omitempty"`
	StatusPTH *StatusPower    `json:"StatusPTH,omitempty"`
}

// StatusInfo contains basic device information (Status 0, Status 1).
type StatusInfo struct {
	Module       int      `json:"Module"`
	DeviceName   string   `json:"DeviceName"`
	FriendlyName []string `json:"FriendlyName"`
	Topic        string   `json:"Topic"`
	ButtonTopic  string   `json:"ButtonTopic"`
	Power        int      `json:"Power"`
	PowerOnState int      `json:"PowerOnState"`
	LedState     int      `json:"LedState"`
	LedMask      string   `json:"LedMask"`
	SaveData     int      `json:"SaveData"`
	SaveState    int      `json:"SaveState"`
	SwitchTopic  string   `json:"SwitchTopic"`
	SwitchMode   []int    `json:"SwitchMode"`
	ButtonRetain int      `json:"ButtonRetain"`
	SwitchRetain int      `json:"SwitchRetain"`
	SensorRetain int      `json:"SensorRetain"`
	PowerRetain  int      `json:"PowerRetain"`
}

// EthernetInfo contains ethernet interface information.
type EthernetInfo struct {
	Hostname   string   `json:"Hostname"`
	IPAddress  []IPAddr `json:"IPAddress"`
	Gateway    IPAddr   `json:"Gateway"`
	Subnetmask IPAddr   `json:"Subnetmask"`
	DNSServer  IPAddr   `json:"DNSServer"`
	Mac        MACAddr  `json:"Mac"`
}

// StatusParam contains device parameters (Status 2).
type StatusParam struct {
	Hostname    string        `json:"Hostname"`
	IPAddress   []IPAddr      `json:"IPAddress"`
	Gateway     IPAddr        `json:"Gateway"`
	Subnetmask  IPAddr        `json:"Subnetmask"`
	DNSServer   IPAddr        `json:"DNSServer"`
	Mac         MACAddr       `json:"Mac"`
	Ethernet    *EthernetInfo `json:"Ethernet,omitempty"`
	WebServer   int           `json:"WebServer"`
	WebPassword int           `json:"WebPassword"`
	Sleep       int           `json:"Sleep"`
	BootCount   int           `json:"BootCount"`
	BCResetTime string        `json:"BCResetTime"`
	SaveCount   int           `json:"SaveCount"`
	SaveAddress string        `json:"SaveAddress"`
}

// StatusFirmware contains firmware information (Status 2).
type StatusFirmware struct {
	Version       string `json:"Version"`
	BuildDateTime string `json:"BuildDateTime"`
	Boot          int    `json:"Boot"`
	Core          string `json:"Core"`
	SDK           string `json:"SDK"`
	CpuFrequency  int    `json:"CpuFrequency"` //nolint:revive // Tasmota API field name
	Hardware      string `json:"Hardware"`
	CR            string `json:"CR"`
}

// StatusLog contains logging information (Status 3).
type StatusLog struct {
	SerialLog  int      `json:"SerialLog"`
	WebLog     int      `json:"WebLog"`
	MqttLog    int      `json:"MqttLog"`
	SysLog     int      `json:"SysLog"`
	LogHost    string   `json:"LogHost"`
	LogPort    int      `json:"LogPort"`
	SSId       []string `json:"SSId"`
	TelePeriod int      `json:"TelePeriod"`
	Resolution string   `json:"Resolution"`
	SetOption  []string `json:"SetOption"`
}

// StatusMemory contains memory information (Status 4).
type StatusMemory struct {
	ProgramSize      int      `json:"ProgramSize"`
	Free             int      `json:"Free"`
	Heap             int      `json:"Heap"`
	ProgramFlashSize int      `json:"ProgramFlashSize"`
	FlashSize        int      `json:"FlashSize"`
	FlashChipId      string   `json:"FlashChipId"` //nolint:revive // Tasmota API field name
	FlashFrequency   int      `json:"FlashFrequency"`
	FlashMode        int      `json:"FlashMode"`
	Features         []string `json:"Features"`
	Drivers          string   `json:"Drivers"`
	Sensors          string   `json:"Sensors"`
	DisplayWidth     int      `json:"DisplayWidth,omitempty"`
	DisplayHeight    int      `json:"DisplayHeight,omitempty"`
	DisplayMode      int      `json:"DisplayMode,omitempty"`
	DisplayRotate    int      `json:"DisplayRotate,omitempty"`
	DisplayCols      int      `json:"DisplayCols,omitempty"`
	DisplayRows      int      `json:"DisplayRows,omitempty"`
	DisplayFont      int      `json:"DisplayFont,omitempty"`
}

// StatusNetwork contains network information (Status 5).
type StatusNetwork struct {
	Hostname   string  `json:"Hostname"`
	IPAddress  IPAddr  `json:"IPAddress"`
	Gateway    IPAddr  `json:"Gateway"`
	Subnetmask IPAddr  `json:"Subnetmask"`
	DNSServer  IPAddr  `json:"DNSServer"`
	DNSServer2 IPAddr  `json:"DNSServer2"`
	Mac        MACAddr `json:"Mac"`
	Webserver  int     `json:"Webserver"`
	HTTPAPI    int     `json:"HTTP_API"` //nolint:revive // Tasmota API field name
	WifiConfig int     `json:"WifiConfig"`
	WifiPower  float64 `json:"WifiPower"`
}

// StatusMQTT contains MQTT configuration (Status 6).
type StatusMQTT struct {
	MqttHost        string `json:"MqttHost"`
	MqttPort        int    `json:"MqttPort"`
	MqttClientMask  string `json:"MqttClientMask"`
	MqttClient      string `json:"MqttClient"`
	MqttUser        string `json:"MqttUser"`
	MqttCount       int    `json:"MqttCount"`
	MaxPacketSize   int    `json:"MAX_PACKET_SIZE"`   //nolint:revive // Tasmota API field name
	Keepalive       int    `json:"KEEPALIVE"`         //nolint:revive // Tasmota API field name
	SocketTimeout   int    `json:"SOCKET_TIMEOUT"`    //nolint:revive // Tasmota API field name
}

// StatusTime contains time information (Status 7).
type StatusTime struct {
	UTC       string `json:"UTC"`
	Local     string `json:"Local"`
	StartDST  string `json:"StartDST"`
	EndDST    string `json:"EndDST"`
	Timezone  int    `json:"Timezone"`
	Sunrise   string `json:"Sunrise"`
	Sunset    string `json:"Sunset"`
}

// StatusSensor contains sensor data (Status 8, Status 10).
type StatusSensor struct {
	Time   string              `json:"Time"`
	Switch []string            `json:"Switch,omitempty"`
	Energy *EnergyData         `json:"ENERGY,omitempty"`
	// Add more sensor types as needed
	Raw map[string]interface{} `json:"-"` // Catch-all for unknown sensors
}

// EnergyData contains power monitoring information.
type EnergyData struct {
	TotalStartTime string  `json:"TotalStartTime"`
	Total          float64 `json:"Total"`
	Yesterday      float64 `json:"Yesterday"`
	Today          float64 `json:"Today"`
	Period         float64 `json:"Period"`
	Power          float64 `json:"Power"`
	ApparentPower  float64 `json:"ApparentPower"`
	ReactivePower  float64 `json:"ReactivePower"`
	Factor         float64 `json:"Factor"`
	Voltage        float64 `json:"Voltage"`
	Current        float64 `json:"Current"`
}

// StatusState contains current device state (Status 11).
type StatusState struct {
	Time     string   `json:"Time"`
	Uptime   string   `json:"Uptime"`
	UptimeSec int     `json:"UptimeSec"`
	Heap     int      `json:"Heap"`
	SleepMode string  `json:"SleepMode"`
	Sleep    int      `json:"Sleep"`
	LoadAvg  int      `json:"LoadAvg"`
	MqttCount int     `json:"MqttCount"`
	POWER    string   `json:"POWER,omitempty"`
	POWER1   string   `json:"POWER1,omitempty"`
	POWER2   string   `json:"POWER2,omitempty"`
	POWER3   string   `json:"POWER3,omitempty"`
	POWER4   string   `json:"POWER4,omitempty"`
	POWER5   string   `json:"POWER5,omitempty"`
	POWER6   string   `json:"POWER6,omitempty"`
	POWER7   string   `json:"POWER7,omitempty"`
	POWER8   string   `json:"POWER8,omitempty"`
	Wifi     *WifiInfo `json:"Wifi,omitempty"`
}

// WifiInfo contains WiFi connection information.
type WifiInfo struct {
	AP            int    `json:"AP"`
	SSId          string `json:"SSId"`
	BSSId         string `json:"BSSId"`
	Channel       int    `json:"Channel"`
	Mode          string `json:"Mode"`
	RSSI          int    `json:"RSSI"`
	Signal        int    `json:"Signal"`
	LinkCount     int    `json:"LinkCount"`
	Downtime      string `json:"Downtime"`
}

// StatusPower contains power usage details.
type StatusPower struct {
	Power     float64 `json:"Power"`
	Energy    float64 `json:"Energy"`
	Voltage   float64 `json:"Voltage"`
	Current   float64 `json:"Current"`
}

// Status queries device status information.
// category can be 0 (all) or 1-11 for specific status types:
//   0 = All status information
//   1 = Device parameters
//   2 = Firmware version
//   3 = Logging information
//   4 = Memory information
//   5 = Network information
//   6 = MQTT information
//   7 = Time information
//   8 = Sensor information
//   9 = Power threshold
//   10 = Sensor information
//   11 = State information
func (c *Client) Status(ctx context.Context, category int) (*StatusResponse, error) {
	if category < 0 || category > 11 {
		return nil, NewError(ErrorTypeCommand, "status category must be between 0 and 11", nil)
	}

	cmd := "Status"
	if category > 0 {
		cmd = "Status " + string(rune('0'+category))
	}

	raw, err := c.ExecuteCommand(ctx, cmd)
	if err != nil {
		return nil, err
	}

	var resp StatusResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, NewError(ErrorTypeParse, "failed to parse status response", err)
	}

	return &resp, nil
}

// GetDeviceInfo retrieves basic device information (Status 1).
func (c *Client) GetDeviceInfo(ctx context.Context) (*StatusInfo, error) {
	resp, err := c.Status(ctx, 1)
	if err != nil {
		return nil, err
	}
	if resp.Status == nil {
		return nil, NewError(ErrorTypeParse, "status response missing Status field", nil)
	}
	return resp.Status, nil
}

// GetFirmwareInfo retrieves firmware version information (Status 2).
func (c *Client) GetFirmwareInfo(ctx context.Context) (*StatusFirmware, error) {
	resp, err := c.Status(ctx, 2)
	if err != nil {
		return nil, err
	}
	if resp.StatusFWR == nil {
		return nil, NewError(ErrorTypeParse, "status response missing StatusFWR field", nil)
	}
	return resp.StatusFWR, nil
}

// GetNetworkInfo retrieves network configuration (Status 5).
func (c *Client) GetNetworkInfo(ctx context.Context) (*StatusNetwork, error) {
	resp, err := c.Status(ctx, 5)
	if err != nil {
		return nil, err
	}
	if resp.StatusNET == nil {
		return nil, NewError(ErrorTypeParse, "status response missing StatusNET field", nil)
	}
	return resp.StatusNET, nil
}

// GetMQTTInfo retrieves MQTT configuration (Status 6).
func (c *Client) GetMQTTInfo(ctx context.Context) (*StatusMQTT, error) {
	resp, err := c.Status(ctx, 6)
	if err != nil {
		return nil, err
	}
	if resp.StatusMQT == nil {
		return nil, NewError(ErrorTypeParse, "status response missing StatusMQT field", nil)
	}
	return resp.StatusMQT, nil
}

// GetSensorData retrieves sensor readings (Status 10).
func (c *Client) GetSensorData(ctx context.Context) (*StatusSensor, error) {
	resp, err := c.Status(ctx, 10)
	if err != nil {
		return nil, err
	}
	if resp.StatusSNS == nil {
		return nil, NewError(ErrorTypeParse, "status response missing StatusSNS field", nil)
	}
	return resp.StatusSNS, nil
}

// GetState retrieves current device state (Status 11).
func (c *Client) GetState(ctx context.Context) (*StatusState, error) {
	resp, err := c.Status(ctx, 11)
	if err != nil {
		return nil, err
	}
	if resp.StatusSTS == nil {
		return nil, NewError(ErrorTypeParse, "status response missing StatusSTS field", nil)
	}
	return resp.StatusSTS, nil
}

// GetUptime retrieves device uptime.
func (c *Client) GetUptime(ctx context.Context) (time.Duration, error) {
	state, err := c.GetState(ctx)
	if err != nil {
		return 0, err
	}
	return time.Duration(state.UptimeSec) * time.Second, nil
}

// GetWiFiSignal retrieves WiFi signal strength (RSSI).
func (c *Client) GetWiFiSignal(ctx context.Context) (int, error) {
	state, err := c.GetState(ctx)
	if err != nil {
		return 0, err
	}
	if state.Wifi == nil {
		return 0, NewError(ErrorTypeDevice, "WiFi information not available", nil)
	}
	return state.Wifi.RSSI, nil
}
