package wiz

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var rgbRegexp = regexp.MustCompile(`^#([0-9a-f]{2})([0-9a-f]{2})([0-9a-f]{2})(\/(\d+))?$`)
var rgbcwRegexp = regexp.MustCompile(`^#([0-9a-f]{2})([0-9a-f]{2})([0-9a-f]{2})([0-9a-f]{2})([0-9a-f]{2})(\/(\d+))?$`)
var cwRegexp = regexp.MustCompile(`^#([0-9a-f]{2})([0-9a-f]{2})(\/(\d+))?$`)
var sceneRegexp = regexp.MustCompile(`^([a-z]+)(\/(\d+))?(\/(\d+))?$`)
var tempRegexp = regexp.MustCompile(`^w(\d+)(\/(\d+))?$`)

const ColorOff = "off"

// SettingsFromString will generate RGB CW Scene Speed and Brightness settings
// from simple strings
// Accepted are:
// #RRGGBB (hex)
// #RRGGBB/brightness (hex with optional brightness in decimal)
// #RRGGBBCCDDEE (hex with optional CW values)
// #RRGGBBCCDDEE/brightness (hex with optional CW values and brightness)
// #CCWW (hex)
// #CCWW/brightness (hex with optional brightness in decimal)
// scene/brightness/speed (scene name with optional brightness and speed)
// See WizScene for available scene names
func SettingsFromString(str string) (*Settings, error) {

	lower := strings.ToLower(strings.TrimSpace(str))
	if lower == ColorOff {
		return NewSettings().Off(), nil
	}
	settings := NewSettings()

	matches := rgbRegexp.FindStringSubmatch(lower)
	if matches != nil {
		// RGB
		r, _ := strconv.ParseUint(matches[1], 16, 32)
		g, _ := strconv.ParseUint(matches[2], 16, 32)
		b, _ := strconv.ParseUint(matches[3], 16, 32)
		settings.RGB(uint8(r), uint8(g), uint8(b))

		if matches[5] != "" {
			bright, _ := strconv.ParseUint(matches[5], 10, 32)
			if bright < 1 || bright > 100 {
				return nil, fmt.Errorf("brightness must be between 1 and 100")
			}
			settings.Brightness(int(bright))
		}
		return settings, nil
	}

	matches = rgbcwRegexp.FindStringSubmatch(lower)
	if matches != nil {
		// RGBCW
		r, _ := strconv.ParseUint(matches[1], 16, 32)
		g, _ := strconv.ParseUint(matches[2], 16, 32)
		b, _ := strconv.ParseUint(matches[3], 16, 32)
		c, _ := strconv.ParseUint(matches[4], 16, 32)
		w, _ := strconv.ParseUint(matches[5], 16, 32)
		settings.RGB(uint8(r), uint8(g), uint8(b)).CW(uint8(c)).WW(uint8(w))

		if matches[7] != "" {
			bright, _ := strconv.ParseUint(matches[7], 10, 32)

			if bright < 1 || bright > 100 {
				return nil, fmt.Errorf("brightness must be between 1 and 100")
			}
			settings.Brightness(int(bright))
		}
		return settings, nil
	}

	matches = cwRegexp.FindStringSubmatch(lower)
	if matches != nil {
		// CW
		c, _ := strconv.ParseUint(matches[1], 16, 32)
		w, _ := strconv.ParseUint(matches[2], 16, 32)
		settings.CW(uint8(c)).WW(uint8(w))

		if matches[4] != "" {
			bright, _ := strconv.ParseUint(matches[4], 10, 32)

			if bright < 1 || bright > 100 {
				return nil, fmt.Errorf("brightness must be between 1 and 100")
			}
			settings.Brightness(int(bright))
		}
		return settings, nil
	}

	matches = tempRegexp.FindStringSubmatch(lower)
	if matches != nil {
		// Temperature
		temp, _ := strconv.ParseUint(matches[1], 10, 32)
		settings.Temperature(int(temp))

		if matches[3] != "" {
			bright, _ := strconv.ParseUint(matches[3], 10, 32)
			if bright < 1 || bright > 100 {
				return nil, fmt.Errorf("brightness must be between 1 and 100")
			}
			settings.Brightness(int(bright))
		}

		return settings, nil
	}

	matches = sceneRegexp.FindStringSubmatch(lower)
	if matches != nil {
		// Scene
		scene := matches[1]

		sceneId, found := GetSceneByName(scene)
		if !found {
			return nil, fmt.Errorf("%q is not a valid color string or wiz scene name", scene)
		}

		settings.Scene(sceneId)

		if matches[3] != "" {
			bright, _ := strconv.ParseUint(matches[3], 10, 32)
			if bright < 1 || bright > 100 {
				return nil, fmt.Errorf("brightness must be between 1 and 100")
			}
			settings.Brightness(int(bright))
		}

		if matches[5] != "" {
			speed, _ := strconv.ParseUint(matches[5], 10, 32)
			if speed < 10 || speed > 200 {
				return nil, fmt.Errorf("speed must be between 10 and 200")
			}
			settings.Speed(int(speed))
		}
		return settings, nil
	}

	return nil, fmt.Errorf("invalid color string")
}

// MustSettingsFromString will produce settings in the same way as SettingsFromString, but will panic on error
func MustSettingsFromString(s string) *Settings {
	settings, err := SettingsFromString(s)
	if err != nil {
		panic(err)
	}
	return settings
}

// ColorStringFromPilot will return a color string that coresponds to a given pilot result
func ColorStringFromPilot(p *PilotResult) (string, error) {

	if !p.State {
		return ColorOff, nil
	}

	if p.R != 0 || p.G != 0 || p.B != 0 {
		// RGB or RGBCW
		if p.C != 0 || p.W != 0 {
			// RGBCW
			return fmt.Sprintf("#%02x%02x%02x%02x%02x/%d", p.R, p.G, p.B, p.C, p.W, p.Dimming), nil
		} else {
			// RGB
			return fmt.Sprintf("#%02x%02x%02x/%d", p.R, p.G, p.B, p.Dimming), nil
		}
	} else if p.C != 0 || p.W != 0 {
		// CW
		return fmt.Sprintf("#%02x%02x/%d", p.C, p.W, p.Dimming), nil
	} else if p.Temp != 0 {
		// Temperature
		return fmt.Sprintf("W%d/%d", p.Temp, p.Dimming), nil
	} else if p.SceneId != 0 {
		// Scene
		sceneName := Scene(p.SceneId).String()
		if sceneName == "" {
			return "", fmt.Errorf("unknown scene ID %d", p.SceneId)
		}
		// Always return scene with dimming and speed
		return fmt.Sprintf("%s/%d/%d", sceneName, p.Dimming, p.Speed), nil
	}
	return "", fmt.Errorf("invalid pilot result")
}

// ColorStringRationalize will take a basic color that may not include brightness (or speed)
// and if those are missing, will add them at default values
func RationalizeColor(col string) (string, error) {

	if _, err := SettingsFromString(col); err != nil {
		return "", err
	}

	lower := strings.ToLower(strings.TrimSpace(col))
	if lower == ColorOff {
		return ColorOff, nil
	}
	matches := rgbRegexp.FindStringSubmatch(lower)
	if matches != nil {
		// RGB
		if matches[5] == "" {
			matches[5] = "100"
		}
		return fmt.Sprintf("#%s%s%s/%s", matches[1], matches[2], matches[3], matches[5]), nil
	}

	matches = rgbcwRegexp.FindStringSubmatch(lower)
	if matches != nil {
		if matches[7] == "" {
			matches[7] = "100"
		}
		return fmt.Sprintf("#%s%s%s%s%s/%s", matches[1], matches[2], matches[3], matches[4], matches[5], matches[7]), nil
		// RGBCW
	}

	matches = cwRegexp.FindStringSubmatch(lower)
	if matches != nil {
		// CW
		if matches[4] == "" {
			matches[4] = "100"
		}
		return fmt.Sprintf("#%s%s/%s", matches[1], matches[2], matches[4]), nil
	}

	matches = tempRegexp.FindStringSubmatch(lower)
	if matches != nil {
		// Temperature
		if matches[3] == "" {
			matches[3] = "100"
		}
		return fmt.Sprintf("W%s/%s", matches[1], matches[3]), nil
	}

	matches = sceneRegexp.FindStringSubmatch(lower)
	if matches != nil {
		// Scene
		if matches[3] == "" {
			matches[3] = "100"
		}
		if matches[5] == "" {
			matches[5] = "200"
		}
		return fmt.Sprintf("%s/%s/%s", matches[1], matches[3], matches[5]), nil
	}
	return "", fmt.Errorf("invalid color string")
}
