# wiz - golang based library for wiz smart bulbs
A golang library for interacting with wiz connected devices. Supports discovery of wiz devices and interacting with them, such as programming colors.

In order to use this library you will need to enable local API control on your wiz bulbs using the wiz app.

## Discovery
Wiz bulbs permit discovery using IP broadcast, although as with most broadcast protocols your success rate may vary.

The following example will broadcast discovery attempts on 192.168.10.255 (should be your wiz device's network's broadcast address) at 1 second intervals for 3 seconds, returning any devices found.

```go
devices, err := wiz.Discover().
    BroadcastAddress("192.168.10.255").
    BroadcastInterval(1 * time.Second).
    EndAfter(3 * time.Second).
    GetDiscoveredDevices()
if err != nil {
    panic(err)
}
for _, dev := range devices {
    fmt.Printf("Discovered device: IP: %s, MAC: %s\n", dev.IP(), dev.MAC())
}
```

Note that it is not always necessary to set the broadcast address, you can often omit this step on machines that only have a single IP interface, or where the default IP interface connects to the network hosting the wiz devices.

You can also attach a callback function to receive devices immediately upon their discovery without waiting for the full discovery to complete:

```go
_, err := wiz.Discover().
BroadcastAddress("192.168.21.255").
BroadcastInterval(1 * time.Second).
EndAfter(3 * time.Second).
WithCallback(func(d *wiz.Device) {
	fmt.Printf("Discovered device: IP: %s, MAC: %s\n", d.IP(), d.MAC())
}).
GetDiscoveredDevices()
if err != nil {
	panic(err)
}
```

You can use the callback to initiate immediate asynchronous operations on the discovered devices.  The callback does block the discovery loop, so keep your actions quick and don't initiate any synchronous operations.

## Interacting with the devices
Once you have your wiz devices, a number of methods are provided to interact with them. These methods come in both synchronous and asynchronous flavors. The synchronous methods block until a response is received, the asynchronous ones do not block and you must separately wait for their responses.

The two most common operations are "GetPilot" and "SetPilot" which respectively retrieve and set the current color, brightness, on/off state and such for the device.

### Get Pilot (synchronous)
```go
device, err := wiz.NewDevice("192.168.20.163")
if err != nil {
	panic(err)
}

pilot, err := device.GetPilot()
if err != nil {
	panic(err)
}

fmt.Printf("Pilot: %+v\n", pilot)
// Outputs: 
//   Pilot: Mac: 9877d5469bd8, Rssi: -53, State: true, SceneId: 0, R: 0, G: 0, B: 0, C: 0, W: 0, Temp: 5800, Dimming: 100, Speed: 0
```
The internet provides plenty of explanations on the different values, so I won't go into detail here, but basically:
* State: power on state
* SceneId: The Wiz scene ID, if a scene is set, else 0
* R, G, B: RGB settings (0..255) if the light is RGB capable and in RGB mode
* C, W: Cool and Warm white LED settings if the light has that capability and is in that mode
* Temp: White temperature, if the light is programmed to output a specific white level
* Dimming: Brightness of the bulb (0..100%)
* Speed: Speed at which any scene in use is set to (or 0 if not set to a scene)

### Get Pilot (asynchronous)
```go
device, err := wiz.NewDevice("192.168.20.163")
if err != nil {
	panic(err)
}

operation := device.GetPilotAsync()

// Be doing other stuff whilst the operation happens

// Now wait for the operation to complete
pilot, err := operation.AwaitResult()
if err != nil {
	panic(err)
}

fmt.Printf("Pilot: %+v\n", pilot)
// Outputs: 
//   Pilot: Mac: 9877d5469bd8, Rssi: -53, State: true, SceneId: 0, R: 0, G: 0, B: 0, C: 0, W: 0, Temp: 5800, Dimming: 100, Speed: 0

```

### SetPilot
To set the state of the light, use SetPilot
```go
	device, err := wiz.NewDevice("192.168.20.163")
	if err != nil {
		panic(err)
	}

	_, err = device.SetPilot(
		wiz.NewSettings().
			On().
			Brightness(100).
			RGB(0xFF, 0x00, 0x00))

	if err != nil {
		fmt.Printf("Failed to set light: %v", err)
	}
```

As with all other operations, an async version (SetPilotAsync) is available

### Other operations
Various other operations are available:
* On - turns light on, alternative to using SetPilot
* Off - turns light off, alternative to using SetPilot
* GetModelConfig - returns information about the model
* GetSystemConfig - returns system settings
* GetUserConfig - returns user settings
* Pulse - pulses the light
* Reachable - uses ICMP echo (ping) to test the reachability of the IP address

More information about the responses and parameters for these operations can be found on the internet.

As with all other operations, async versions are available for all.

## Device groups
Sometimes, it is convenient to operate on a group of devices.  The DeviceGroup permits easy ways to do this with asynchronous operations (synchronous operations are not supported for device groups).

In this example, we discover all available devices and turn them green. 

```go
deviceGroup := wiz.NewDeviceGroup(
	wiz.MustDevice("192.168.20.163"),
	wiz.MustDevice("192.168.20.244"),
	wiz.MustDevice("192.168.20.163"))
// N.B. MustDevice is a convenience method for creating devices
// that panics instead of returning an error in the case of an invalid IP

results := deviceGroup.SetPilot(
	wiz.NewSettings().
		On().
		Brightness(100).
		RGB(0xFF, 0x00, 0x00)).
	AwaitAll()

for _, result := range results {
	if result.Err != nil {
		fmt.Printf("Failed to set light %q: %v", result.Device.IP(), result.Err)
	}
}
```

Device groups have corresponding (asynchronous) methods for all the operations listed above. For operations that return useful data (such as GetPilot), the result object will include that data.