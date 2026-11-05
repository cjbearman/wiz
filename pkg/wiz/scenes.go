package wiz

import "strings"

// Scene represents a predefined lighting scene for the device.
type Scene int

const (
	Ocean        Scene = 1
	Romance      Scene = 2
	Sunset       Scene = 3
	Party        Scene = 4
	Fireplace    Scene = 5
	Cozy         Scene = 6
	Forest       Scene = 7
	PastelColors Scene = 8
	WakeUp       Scene = 9
	Bedtime      Scene = 10
	WarmWhite    Scene = 11
	Daylight     Scene = 12
	CoolWhite    Scene = 13
	NightLight   Scene = 14
	Focus        Scene = 15
	Relax        Scene = 16
	TrueColors   Scene = 17
	TVTime       Scene = 18
	PlantGrowth  Scene = 19
	Spring       Scene = 20
	Summer       Scene = 21
	Fall         Scene = 22
	DeepDive     Scene = 23
	Jungle       Scene = 24
	Mojito       Scene = 25
	Club         Scene = 26
	Christmas    Scene = 27
	Halloween    Scene = 28
	Candlelight  Scene = 29
	GoldenWhite  Scene = 30
	Pulse        Scene = 31
	Steampunk    Scene = 32
	Diwali       Scene = 33
	LightAlarm   Scene = 35
	Rhythm       Scene = 1000
)

var sceneMap = map[string]Scene{
	"Ocean":        Ocean,
	"Romance":      Romance,
	"Sunset":       Sunset,
	"Party":        Party,
	"Fireplace":    Fireplace,
	"Cozy":         Cozy,
	"Forest":       Forest,
	"PastelColors": PastelColors,
	"WakeUp":       WakeUp,
	"Bedtime":      Bedtime,
	"WarmWhite":    WarmWhite,
	"Daylight":     Daylight,
	"CoolWhite":    CoolWhite,
	"NightLight":   NightLight,
	"Focus":        Focus,
	"Relax":        Relax,
	"TrueColors":   TrueColors,
	"TVTime":       TVTime,
	"PlantGrowth":  PlantGrowth,
	"Spring":       Spring,
	"Summer":       Summer,
	"Fall":         Fall,
	"DeepDive":     DeepDive,
	"Jungle":       Jungle,
	"Mojito":       Mojito,
	"Club":         Club,
	"Christmas":    Christmas,
	"Halloween":    Halloween,
	"Candlelight":  Candlelight,
	"GoldenWhite":  GoldenWhite,
	"Pulse":        Pulse,
	"Steampunk":    Steampunk,
	"Diwali":       Diwali,
	"LightAlarm":   LightAlarm,
	"Rhythm":       Rhythm,
}

// String returns the name of the scene as a string. If the scene is not recognized, it returns an empty string.
func (s Scene) String() string {
	for name, scene := range sceneMap {
		if scene == s {
			return name
		}
	}
	return ""
}

// GetSceneByName retrieves the Scene corresponding to the given name.
// The comparison is case-insensitive.
// Returns the Scene and true if found, otherwise returns 0 and false.
func GetSceneByName(name string) (Scene, bool) {
	nameLower := strings.ToLower(name)
	for mapName, scene := range sceneMap {
		if strings.ToLower(mapName) == nameLower {
			return scene, true
		}
	}
	return 0, false
}

// AllSceneNames returns a slice containing the names of all available scenes.
func AllSceneNames() []string {
	names := make([]string, 0, len(sceneMap))
	for name := range sceneMap {
		names = append(names, name)
	}
	return names
}
