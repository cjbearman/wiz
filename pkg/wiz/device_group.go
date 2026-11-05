package wiz

// DeviceGroup represents a group of devices that can be controlled together.
// all operations on a group are performed concurrently and asynchronously
type DeviceGroup struct {
	devices []*Device
}

// NewDeviceGroup creates a new DeviceGroup with the given devices.
func NewDeviceGroup(devices ...*Device) *DeviceGroup {
	return &DeviceGroup{
		devices: devices,
	}
}

func (dg *DeviceGroup) AddDevice(device *Device) {
	dg.devices = append(dg.devices, device)
}

// GetPilot retrieves the pilot information for all devices in the group asynchronously.
func (dg *DeviceGroup) GetPilot() AsyncOperations[*PilotResult] {
	if len(dg.devices) == 0 {
		return nil
	}

	ops := make(AsyncOperations[*PilotResult], len(dg.devices))
	for i, device := range dg.devices {
		ops[i] = device.GetPilotAsync()
	}

	return ops
}

// GetModelConfig retrieves the model configuration for all devices in the group asynchronously.
func (dg *DeviceGroup) GetModelConfig() AsyncOperations[*ModelConfig] {
	if len(dg.devices) == 0 {
		return nil
	}

	ops := make(AsyncOperations[*ModelConfig], len(dg.devices))
	for i, device := range dg.devices {
		ops[i] = device.GetModelConfigAsync()
	}

	return ops
}

// GetUserConfig retrieves the user configuration for all devices in the group asynchronously.
func (dg *DeviceGroup) GetUserConfig() AsyncOperations[*UserConfig] {
	if len(dg.devices) == 0 {
		return nil
	}

	ops := make(AsyncOperations[*UserConfig], len(dg.devices))
	for i, device := range dg.devices {
		ops[i] = device.GetUserConfigAsync()
	}

	return ops
}

// GetSystemConfig retrieves the system configuration for all devices in the group asynchronously.
func (dg *DeviceGroup) GetSystemConfig() AsyncOperations[*SystemConfig] {
	if len(dg.devices) == 0 {
		return nil
	}

	ops := make(AsyncOperations[*SystemConfig], len(dg.devices))
	for i, device := range dg.devices {
		ops[i] = device.GetSystemConfigAsync()
	}

	return ops
}

// SetPilot sets the pilot settings for all devices in the group asynchronously.
func (dg *DeviceGroup) SetPilot(settings *Settings) AsyncOperations[*SetPilotResult] {
	if len(dg.devices) == 0 {
		return nil
	}

	ops := make(AsyncOperations[*SetPilotResult], len(dg.devices))
	for i, device := range dg.devices {
		ops[i] = device.SetPilotAsync(settings)
	}

	return ops
}

// SetState sets the state (on/off) for all devices in the group asynchronously.
func (dg *DeviceGroup) SetState(state bool) AsyncOperations[*SetStateResult] {
	if len(dg.devices) == 0 {
		return nil
	}

	ops := make(AsyncOperations[*SetStateResult], len(dg.devices))
	for i, device := range dg.devices {
		ops[i] = device.SetStateAsync(state)
	}

	return ops
}

// On turns on all devices in the group asynchronously.
func (dg *DeviceGroup) On() AsyncOperations[*SetStateResult] {
	return dg.SetState(true)
}

// Off turns off all devices in the group asynchronously.
func (dg *DeviceGroup) Off() AsyncOperations[*SetStateResult] {
	return dg.SetState(false)
}

func (dg *DeviceGroup) Pulse(delta, durationMs int) AsyncOperations[*PulseResult] {
	if len(dg.devices) == 0 {
		return nil
	}

	ops := make(AsyncOperations[*PulseResult], len(dg.devices))
	for i, device := range dg.devices {
		ops[i] = device.PulseAsync(delta, durationMs)
	}

	return ops
}

// Reachable checks if all devices in the group are reachable asynchronously.
// ICMP Echo (Ping) is used to test the device.
func (dg *DeviceGroup) Reachable() AsyncOperations[bool] {
	if len(dg.devices) == 0 {
		return nil
	}

	ops := make(AsyncOperations[bool], len(dg.devices))
	for i, device := range dg.devices {
		ops[i] = device.ReachableAsync()
	}

	return ops
}
