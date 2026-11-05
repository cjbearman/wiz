package wiz

import (
	"fmt"
	"net"
	"sync/atomic"
	"time"

	"errors"

	probing "github.com/prometheus-community/pro-bing"
)

// Device represents a single device and the methods needed to operate it
type Device struct {
	ip      string
	name    string
	mac     string
	logger  Logger
	timeout time.Duration
	retries int
	nextId  atomic.Int32
}

// NewDevice creates a new device with the given IP address.
func NewDevice(ip string) (*Device, error) {

	if net.ParseIP(ip) == nil {
		return nil, errors.New("invalid IP address")
	}

	return &Device{
		ip:      ip,
		name:    ip,
		logger:  new(NopLogger),
		timeout: DefaultOperationTimeout,
		retries: DefaultRetries,
	}, nil
}

// MustDevice is same as NewDevice but panics on error (only error is invalid IP)
func MustDevice(ip string) *Device {
	d, err := NewDevice(ip)
	if err != nil {
		panic(err)
	}
	return d
}

// WithLogger attaches a logger to the device, primarily for debugging and development
func (d *Device) WithLogger(logger Logger) *Device {
	d.logger = logger
	return d
}

// IP returns the IP address of the device
func (d *Device) IP() string {
	return d.ip
}

// Name returns the device name, which is the IP address unless otherwise set
func (d *Device) Name() string {
	return d.name
}

// SetName sets the device name, which will override the default IP-based name
func (d *Device) SetName(name string) *Device {
	d.name = name
	return d
}

// MAC will return the device mac address for devices returned by discovery
// or if it's otherwise known, such as if a getPilot operation has been performed.
// Will return "" if the mac is not known
func (d *Device) MAC() string {
	return d.mac
}

// setMac for internal use to set the device MAC address
func (d *Device) setMac(mac string) *Device {
	d.mac = mac
	return d
}

// SetIP updates the IP address of the device
func (d *Device) SetIP(ip string) *Device {
	d.ip = ip
	return d
}

// WithTimeout sets the timeout for individual operations on the device
// actual end to end timeouts will be a multiple of this and the number of retries
func (d *Device) WithTimeout(timeout time.Duration) *Device {
	d.timeout = timeout
	return d
}

// WithRetries sets the number of retries for individual operations on the device
func (d *Device) WithRetries(retries int) *Device {
	d.retries = retries
	return d
}

// SetState will set the on/off state of the device
func (d *Device) SetState(state bool) (*SetStateResult, error) {
	request := request{
		Method: setState,
		Params: methodParams{
			pilotFieldState: state,
		},
	}
	response := &setStateResponse{}
	err := d.doWithRetry(request, response, d.timeout)
	if err != nil {
		d.logger.Errorf("SetState failure: %v", err)
		return nil, err
	}
	return &response.Result, nil
}

// SetStateAsync will set the on/off state of the device asynchronously
func (d *Device) SetStateAsync(state bool) *AsyncOperation[*SetStateResult] {
	return NewAsyncOperation(d, func() (*SetStateResult, error) {
		return d.SetState(state)
	})
}

// On will turn on the device
func (d *Device) On() (*SetStateResult, error) {
	return d.SetState(true)
}

// Off will turn off the device
func (d *Device) Off() (*SetStateResult, error) {
	return d.SetState(false)
}

// OnAsync will turn on the device asynchronously
func (d *Device) OnAsync() *AsyncOperation[*SetStateResult] {
	return d.SetStateAsync(true)
}

// OffAsync will turn off the device asynchronously
func (d *Device) OffAsync() *AsyncOperation[*SetStateResult] {
	return d.SetStateAsync(false)
}

// GetPilot retrieves the current operational status of the light bulb
func (d *Device) GetPilot() (*PilotResult, error) {
	request := request{
		Method: getPilot,
		Params: methodParams{},
	}
	response := &pilotResponse{}
	err := d.doWithRetry(request, response, d.timeout)
	if err != nil {
		d.logger.Errorf("Pilot failure: %v", err)
		return nil, err
	}
	d.mac = formatMac(response.Result.Mac)
	return &response.Result, nil
}

// GetPilotAsync will retrieve the current operational status of the light bulb asynchronously
func (d *Device) GetPilotAsync() *AsyncOperation[*PilotResult] {
	return NewAsyncOperation(d, func() (*PilotResult, error) {
		return d.GetPilot()
	})
}

// GetModelConfig retrieves the model configuration of the device
func (d *Device) GetModelConfig() (*ModelConfig, error) {
	request := request{
		Method: getModelConfig,
		Params: methodParams{},
	}
	response := &modelConfigResponse{}
	err := d.doWithRetry(request, response, d.timeout)
	if err != nil {
		d.logger.Errorf("GetModelConfig failure: %v", err)
		return nil, err
	}
	return &response.Result, nil
}

// GetModelConfigAsync retrieves the model configuration of the device asynchronously
func (d *Device) GetModelConfigAsync() *AsyncOperation[*ModelConfig] {
	return NewAsyncOperation(d, func() (*ModelConfig, error) {
		return d.GetModelConfig()
	})
}

// GetUserConfig retrieves the user configuration of the device
func (d *Device) GetUserConfig() (*UserConfig, error) {
	request := request{
		Method: getUserConfig,
		Params: methodParams{},
	}
	response := &userConfigResponse{}
	err := d.doWithRetry(request, response, d.timeout)
	if err != nil {
		d.logger.Errorf("GetUserConfig failure: %v", err)
		return nil, err
	}
	return &response.Result, nil
}

// GetUserConfigAsync retrieves the user configuration of the device asynchronously
func (d *Device) GetUserConfigAsync() *AsyncOperation[*UserConfig] {
	return NewAsyncOperation(d, func() (*UserConfig, error) {
		return d.GetUserConfig()
	})
}

// GetSystemConfig retrieves the system configuration of the device
func (d *Device) GetSystemConfig() (*SystemConfig, error) {
	request := request{
		Method: getSystemConfig,
		Params: methodParams{},
	}
	response := &systemConfigResponse{}
	err := d.doWithRetry(request, response, d.timeout)
	if err != nil {
		d.logger.Errorf("GetSystemConfig failure: %v", err)
		return nil, err
	}
	return &response.Result, nil
}

// GetSystemConfigAsync retrieves the system configuration of the device asynchronously
func (d *Device) GetSystemConfigAsync() *AsyncOperation[*SystemConfig] {
	return NewAsyncOperation(d, func() (*SystemConfig, error) {
		return d.GetSystemConfig()
	})
}

// SetPilot sets the operational status of the light bulb according to the provided settings
// The settings can be established with NewSettings or SettingsFromString
func (d *Device) SetPilot(settings *Settings) (*SetPilotResult, error) {
	request := request{
		Method: setPilot,
		Params: settings.p,
	}
	response := &setPilotResponse{}
	err := d.doWithRetry(request, response, d.timeout)
	if err != nil {
		d.logger.Errorf("SetPilot failure: %v", err)
		return nil, err
	}

	if !response.Result.Success {
		err = fmt.Errorf("Operation returned success=false")
	}

	return &response.Result, err
}

// SetPilotAsync sets the operational status of the light bulb according to the provided settings, asynchronously
func (d *Device) SetPilotAsync(settings *Settings) *AsyncOperation[*SetPilotResult] {
	return NewAsyncOperation(d, func() (*SetPilotResult, error) {
		return d.SetPilot(settings)
	})
}

// Pulse causes the light to do a single pulse.
// Delta -100 to 100 indicates the relative change in brightness
// durationMs is the millisecond duration of the pulse
func (d *Device) Pulse(delta, durationMs int) (*PulseResult, error) {
	request := request{
		Method: pulse,
		Params: methodParams{
			"delta":    delta,
			"duration": durationMs,
		},
	}
	response := &pulseResponse{}
	err := d.doWithRetry(request, response, d.timeout)
	if err != nil {
		d.logger.Errorf("Pulse failure: %v", err)
		return nil, err
	}

	if !response.Result.Success {
		err = fmt.Errorf("Operation returned success=false")
	}

	return &response.Result, err

}

// PulseAsync triggers a pulse effect on the light bulb asynchronously
func (d *Device) PulseAsync(delta, durationMs int) *AsyncOperation[*PulseResult] {
	return NewAsyncOperation(d, func() (*PulseResult, error) {
		return d.Pulse(delta, durationMs)
	})
}

// Reachable uses ICMP Echo (Ping) to check if the device is reachable
// The device is considered reachable if it responds to at least one of three pings
func (d *Device) Reachable() (bool, error) {
	pinger, err := probing.NewPinger(d.ip)
	if err != nil {
		return false, err
	}
	pinger.SetLogger(&probing.NoopLogger{})
	pinger.Timeout = pingTimeout
	pinger.Interval = pingInterval
	pinger.Count = pingCount
	err = pinger.Run()
	if err != nil {
		return false, err
	}
	return pinger.Statistics().PacketsRecv > 0, nil
}

// ReachableAsync tests device reachability asynchronously using ICMP Echo (Ping)
// The device is considered reachable if it responds to at least one of three pings
func (d *Device) ReachableAsync() *AsyncOperation[bool] {
	return NewAsyncOperation(d, func() (bool, error) {
		return d.Reachable()
	})
}

// doWithRetry performs the given request with retry logic based on the device's retry settings
func (d *Device) doWithRetry(request request, response jsonResponse, timeout time.Duration) error {
	if d.retries < 0 {
		return fmt.Errorf("invalid retries: %d", d.retries)
	}

	var err error
	for i := 0; i < d.retries; i++ {
		err, retryable := d.do(request, response, timeout)
		if err == nil || !retryable {
			return err
		}
		d.logger.Warnf("Retry %d/%d after error: %v", i+1, d.retries, err)
	}
	return err
}

// do performs the actual request to the device over UDP and returns any error along with a flag indicating if the error is retryable
func (d *Device) do(request request, response jsonResponse, timeout time.Duration) (err error, retryable bool) {
	bindAddr, err := net.ResolveUDPAddr("udp", ":0")
	if err != nil {
		return err, false
	}
	conn, err := net.ListenUDP("udp", bindAddr)
	if err != nil {
		return err, false
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(timeout))

	remoteAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", d.ip, wizPort))
	if err != nil {
		return err, false
	}

	request.Id = int(d.nextId.Add(1))

	if request.Params == nil {
		request.Params = make(methodParams)
	}

	rawRequest, err := request.marshal()
	if err != nil {
		return err, false
	}

	d.logger.Debugf("TX -> %s: %s", d.ip, string(rawRequest))
	_, err = conn.WriteToUDP(rawRequest, remoteAddr)
	if err != nil {
		return err, false
	}

	rawResponse := make([]byte, 1472)
	n, _, err := conn.ReadFromUDP(rawResponse)
	if err != nil {
		return err, true
	}

	if response != nil {
		// Passing a nil response is useful in development
		err = response.unmarshal(rawResponse[:n])
		if err != nil {
			return err, true
		}
	}

	if response.method() != string(request.Method) {
		return fmt.Errorf("unexpected response method: %s", response.method()), false
	}
	if response.id() != request.Id {
		return fmt.Errorf("unexpected response id: %d", response.id()), false
	}

	d.logger.Debugf("RX <- %s: %s", d.ip, string(rawResponse[:n]))

	if response.error() != nil {
		return fmt.Errorf("device reported failure, code: %d, message: %s", response.error().Code, response.error().Message), false
	}

	return nil, false

}

// String returns the string representation of the device, which is its IP address
func (d *Device) String() string {
	return d.ip
}
