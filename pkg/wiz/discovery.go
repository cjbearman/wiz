package wiz

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"
)

type DiscoveryCallback func(device *Device)

type Discovery struct {
	stop              chan struct{}
	endAfter          time.Duration
	broadcastInterval time.Duration
	err               error
	addr              *net.UDPAddr
	conn              net.PacketConn
	broadcastAddress  string
	logger            Logger
	devices           map[string]*Device
	callback          DiscoveryCallback
}

// Discover is used to create a discovery instance with the specified update status flag for discovered devices.
func Discover() *Discovery {
	d := &Discovery{
		stop:              make(chan struct{}),
		endAfter:          DefaultDiscoveryRunFor,
		broadcastInterval: DefaultDiscoveryBroadcastInterval,
		broadcastAddress:  "255.255.255.255",
		logger:            &NopLogger{},
		devices:           make(map[string]*Device),
	}
	return d
}

// WithCallback adds a function that will immediately be called on the discovery of a device
func (d *Discovery) WithCallback(callback DiscoveryCallback) *Discovery {
	d.callback = callback
	return d
}

// EndAfter sets the time after which discovery will automatically stop.
func (d *Discovery) EndAfter(endAfter time.Duration) *Discovery {
	d.endAfter = endAfter
	return d
}

// BroadcastInterval sets the interval at which discovery broadcast messages are sent.
func (d *Discovery) BroadcastInterval(interval time.Duration) *Discovery {
	// Set the broadcast interval duration
	// This is a placeholder for actual implementation
	d.broadcastInterval = interval
	return d
}

// BroadcastAddress sets the broadcast address used for discovery messages.
// This may be required on machines with multiple IP interfaces to set the
// appropriate interface to use for broadcast discovery
func (d *Discovery) BroadcastAddress(address string) *Discovery {
	d.broadcastAddress = address
	return d
}

// WithLogger sets the logger to be used for discovery messages.
func (d *Discovery) WithLogger(logger Logger) *Discovery {
	d.logger = logger
	return d
}

// GetDiscoveredDevices executes the discovery process and returns the list of discovered devices along with any error encountered.
func (d *Discovery) GetDiscoveredDevices() ([]*Device, error) {
	go d.run()
	<-d.stop

	// Turn our map into a slice
	devices := make([]*Device, 0, len(d.devices))
	for _, device := range d.devices {
		devices = append(devices, device)
	}

	return devices, d.err
}

// run runs the discovery process
func (d *Discovery) run() {
	defer close(d.stop)

	d.logger.Infof("Running discovery, for %v seconds, %v broadcast interval",
		d.endAfter.Seconds(), d.broadcastInterval)

	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", d.broadcastAddress, wizPort))
	if err != nil {
		d.logger.Errorf("Failed to resolve broadcast address: %v", err)
		d.err = err
		return
	}
	d.addr = addr

	conn, err := net.ListenPacket("udp4", ":0")
	if err != nil {
		d.logger.Errorf("Failed to listen on UDP port: %v", err)
		d.err = err
		return
	}
	d.conn = conn
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(d.endAfter))

	go d.broadcaster()

	for {
		buf := make([]byte, 1024)
		n, addr, err := d.conn.ReadFrom(buf)
		if err != nil {
			var netErr net.Error
			if errors.As(err, &netErr) && netErr.Timeout() {
				d.logger.Infof("Deadline reached, discovery ended")
				return
			}
			d.logger.Errorf("Error reading from UDP connection: %v", err)
			d.err = err
			return
		}

		discoveredDevice := d.handleReceivedPacket(buf[:n], addr)
		if discoveredDevice != nil {

			if _, found := d.devices[discoveredDevice.IP()]; found {
				d.logger.Infof("Duplicate discovery for device with IP %s", discoveredDevice.IP())
				continue
			}

			d.logger.Infof("New discovery for device with IP %s", discoveredDevice.IP())
			// Add to our map of discovered devices
			d.devices[discoveredDevice.IP()] = discoveredDevice
			if d.callback != nil {
				d.callback(discoveredDevice)
			}
		}
	}
}

type discoveryPacket struct {
	// {"method":"registration","env":"pro","result":{"mac":"9877d5466ed4","success":true}}
	Method string `json:"method"`
	Env    string `json:"env"`
	Result struct {
		Mac     string `json:"mac"`
		Success bool   `json:"success"`
	} `json:"result"`
}

func (d *Discovery) handleReceivedPacket(pkg []byte, addr net.Addr) *Device {
	d.logger.Debugf("Handling received packet, length = %d", len(pkg))
	var packet discoveryPacket
	if err := json.Unmarshal(pkg, &packet); err != nil {
		d.logger.Warnf("failed to unmarshal discovery packet: %v", err)
		return nil
	}

	if packet.Method == "registration" && packet.Result.Success {
		d.logger.Debugf("Valid registration packet received from %s, mac = %s", addr.String(), packet.Result.Mac)
		return MustDevice(ipFromAddr(addr)).WithLogger(d.logger).setMac(formatMac(packet.Result.Mac))
	} else {
		d.logger.Warnf("Erroneous registration packet received from %s", addr.String())
	}
	return nil
}

func (d *Discovery) broadcaster() {
	ticker := NewZeroableTicker(d.broadcastInterval)

	// IP is irrelevant since we are using register=false because we don't want updates
	discoveryPacket := []byte(`{"method":"registration","params":{"phoneMac":"AAAAAAAAAAAA","register":false,"phoneIp":"1.2.3.4","id":"1"}}`)
	for {
		// Broadcast a packet
		d.logger.Infof("Sending a discovery broadcast packet")
		_, err := d.conn.WriteTo(discoveryPacket, d.addr)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			d.err = fmt.Errorf("failed to write discovery packet: %v", err)
			return
		}

		// Wait for next tick, or stop signal
		select {
		case <-d.stop:
			d.logger.Infof("Stop signal received, stopping broadcaster")
			return
		case <-ticker.C:
			// Handle the broadcast interval tick
		}
	}
}

// Error returns the error encountered during the discovery process, if any.
func (d *Discovery) Error() error {
	return d.err
}
