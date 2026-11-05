package wiz

const (
	pilotFieldR       = "r"
	pilotFieldG       = "g"
	pilotFieldB       = "b"
	pilotFieldC       = "c"
	pilotFieldW       = "w"
	pilotFieldTemp    = "temp"
	pilotFieldDimming = "dimming"
	pilotFieldSpeed   = "speed"
	pilotFieldRatio   = "ratio"
	pilotFieldSceneId = "sceneId"
	pilotFieldState   = "state"
)

// Settings provides a fluent interface for configuring LED settings
type Settings struct {
	p methodParams
}

// NewSettings returns a new settings instance
func NewSettings() *Settings {
	return &Settings{
		p: make(methodParams),
	}
}

// State sets the on/off state of the LEDs
func (s *Settings) State(state bool) *Settings {
	s.p[pilotFieldState] = state
	return s
}

// On turns the LEDs on (equivalent to State(true))
func (p *Settings) On() *Settings {
	return p.State(true)
}

// Off turns the LEDs off (equivalent to State(false))
func (p *Settings) Off() *Settings {
	return p.State(false)
}

// RGB sets the values of the RGB LEDs
// Can be set in combination with the CW and WW LEDs if desired and supported
func (p *Settings) RGB(r, g, b uint8) *Settings {
	p.p[pilotFieldR] = r
	p.p[pilotFieldG] = g
	p.p[pilotFieldB] = b
	return p
}

// CW sets the value of the CW LED.
// Can be set in combination with RGB and WW LEDs if desired and supported
func (p *Settings) CW(cw uint8) *Settings {
	p.p[pilotFieldC] = cw
	return p
}

// WW sets the value of the WW LED.
// Can be set in combination with RGB and CW LEDs if desired and supported
func (p *Settings) WW(w uint8) *Settings {
	p.p[pilotFieldW] = w
	return p
}

// Brightness sets the brightness level of the LED.
// 1..100 are usually supported, straying outside this range may cause
// automatica adjustment or operation failure
func (p *Settings) Brightness(brightness int) *Settings {
	p.p[pilotFieldDimming] = brightness
	return p
}

// Temperature sets the color temperature of the LED.
// 2200-6500 are usually supported, straying outside this range may cause
// automatic adjustment or operation failure
func (p *Settings) Temperature(temp int) *Settings {
	p.p[pilotFieldTemp] = temp
	return p
}

// Speed sets the speed of the effect.
// 10-200 are usually supported, straying outside this range may cause
// automatic adjustment or operation failure
func (p *Settings) Speed(speed int) *Settings {
	p.p[pilotFieldSpeed] = speed
	return p
}

// Ratio sets the ratio of the effect.
// TODO: Figure out exactly what is going on with this one
func (p *Settings) Ratio(ratio int) *Settings {
	p.p[pilotFieldRatio] = ratio
	return p
}

// Scene sets the scene ID for the LED.
// The exact range of valid scene IDs depends on the device and firmware.
func (p *Settings) Scene(sceneId Scene) *Settings {
	p.p[pilotFieldSceneId] = sceneId
	return p
}
