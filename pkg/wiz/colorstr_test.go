package wiz

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBasicRGBStrings(t *testing.T) {

	assert := assert.New(t)
	require := require.New(t)

	s, err := SettingsFromString("#fe993c")
	require.NoError(err)
	assert.EqualValues(0xfe, s.p[pilotFieldR])
	assert.EqualValues(0x99, s.p[pilotFieldG])
	assert.EqualValues(0x3c, s.p[pilotFieldB])
	assert.Nil(s.p[pilotFieldDimming])
	assert.Nil(s.p[pilotFieldC])
	assert.Nil(s.p[pilotFieldW])

	s, err = SettingsFromString("#fe993c/99")
	require.NoError(err)
	assert.EqualValues(0xfe, s.p[pilotFieldR])
	assert.EqualValues(0x99, s.p[pilotFieldG])
	assert.EqualValues(0x3c, s.p[pilotFieldB])
	assert.EqualValues(99, s.p[pilotFieldDimming])
	assert.Nil(s.p[pilotFieldC])
	assert.Nil(s.p[pilotFieldW])

	s, err = SettingsFromString("#fe993c/123")
	require.Error(err)

	s, err = SettingsFromString("#fe99x2")
	require.Error(err)

}

func TestBasicRGBSCWtrings(t *testing.T) {

	assert := assert.New(t)
	require := require.New(t)

	s, err := SettingsFromString("#fe993caabb")
	require.NoError(err)
	assert.EqualValues(0xfe, s.p[pilotFieldR])
	assert.EqualValues(0x99, s.p[pilotFieldG])
	assert.EqualValues(0x3c, s.p[pilotFieldB])
	assert.EqualValues(0xaa, s.p[pilotFieldC])
	assert.EqualValues(0xbb, s.p[pilotFieldW])
	assert.Nil(s.p[pilotFieldDimming])

	s, err = SettingsFromString("#fe993caabb/99")
	require.NoError(err)
	assert.EqualValues(0xfe, s.p[pilotFieldR])
	assert.EqualValues(0x99, s.p[pilotFieldG])
	assert.EqualValues(0x3c, s.p[pilotFieldB])
	assert.EqualValues(0xaa, s.p[pilotFieldC])
	assert.EqualValues(0xbb, s.p[pilotFieldW])
	assert.EqualValues(99, s.p[pilotFieldDimming])

	s, err = SettingsFromString("#fe993caabb/123")
	require.Error(err)

	s, err = SettingsFromString("#fe99x299e3")
	require.Error(err)

}

func TestBasicCWStrings(t *testing.T) {

	assert := assert.New(t)
	require := require.New(t)

	s, err := SettingsFromString("#abcd")
	require.NoError(err)
	assert.EqualValues(0xab, s.p[pilotFieldC])
	assert.EqualValues(0xcd, s.p[pilotFieldW])
	assert.Nil(s.p[pilotFieldR])
	assert.Nil(s.p[pilotFieldG])
	assert.Nil(s.p[pilotFieldB])
	assert.Nil(s.p[pilotFieldDimming])

	s, err = SettingsFromString("#abcd/99")
	require.NoError(err)
	assert.EqualValues(0xab, s.p[pilotFieldC])
	assert.EqualValues(0xcd, s.p[pilotFieldW])
	assert.Nil(s.p[pilotFieldR])
	assert.Nil(s.p[pilotFieldG])
	assert.Nil(s.p[pilotFieldB])
	assert.EqualValues(99, s.p[pilotFieldDimming])

	s, err = SettingsFromString("#1234/123")
	require.Error(err)

	s, err = SettingsFromString("#fexa")
	require.Error(err)

}

func TestTempStrings(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	s, er := SettingsFromString("W1234")
	require.NoError(er)
	assert.EqualValues(1234, s.p[pilotFieldTemp])
	assert.Nil(s.p[pilotFieldDimming])

	s, err := SettingsFromString("w1234/99")
	require.NoError(err)
	assert.EqualValues(1234, s.p[pilotFieldTemp])
	assert.EqualValues(99, s.p[pilotFieldDimming])

}

func TestSceneStrings(t *testing.T) {

	assert := assert.New(t)
	require := require.New(t)

	s, err := SettingsFromString("steampunk")
	require.NoError(err)
	assert.EqualValues(Steampunk, s.p[pilotFieldSceneId])
	assert.Nil(s.p[pilotFieldSpeed])
	assert.Nil(s.p[pilotFieldDimming])
	assert.Nil(s.p[pilotFieldR])
	assert.Nil(s.p[pilotFieldG])
	assert.Nil(s.p[pilotFieldB])
	assert.Nil(s.p[pilotFieldC])
	assert.Nil(s.p[pilotFieldW])

	s, err = SettingsFromString("steampunk/100")
	require.NoError(err)
	assert.EqualValues(Steampunk, s.p[pilotFieldSceneId])
	assert.EqualValues(100, s.p[pilotFieldDimming])
	assert.Nil(s.p[pilotFieldSpeed])
	assert.Nil(s.p[pilotFieldR])
	assert.Nil(s.p[pilotFieldG])
	assert.Nil(s.p[pilotFieldB])
	assert.Nil(s.p[pilotFieldC])
	assert.Nil(s.p[pilotFieldW])

	s, err = SettingsFromString("steampunk/100/50")
	require.NoError(err)
	assert.EqualValues(Steampunk, s.p[pilotFieldSceneId])
	assert.EqualValues(100, s.p[pilotFieldDimming])
	assert.EqualValues(50, s.p[pilotFieldSpeed])
	assert.Nil(s.p[pilotFieldR])
	assert.Nil(s.p[pilotFieldG])
	assert.Nil(s.p[pilotFieldB])
	assert.Nil(s.p[pilotFieldC])
	assert.Nil(s.p[pilotFieldW])

	s, err = SettingsFromString("junk")
	require.Error(err)

}

func TestOff(t *testing.T) {
	require := require.New(t)
	s, err := SettingsFromString("off")
	require.NoError(err)
	require.NotNil(s.p[pilotFieldState])
	require.False(s.p[pilotFieldState].(bool))
}

func TestRationalization(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		in, out string
	}{
		{"#112233", "#112233/100"},
		{"#112233/99", "#112233/99"},
		{"#1122334455", "#1122334455/100"},
		{"#1122334455/99", "#1122334455/99"},
		{"#1122", "#1122/100"},
		{"#1122/99", "#1122/99"},
		{"W1234", "W1234/100"},
		{"W1234/99", "W1234/99"},
		{"steampunk", "steampunk/100/200"},
		{"steampunk/100", "steampunk/100/200"},
		{"steampunk/100/50", "steampunk/100/50"},
	}

	t.Run("Rationalization", func(t *testing.T) {
		for _, tt := range tests {
			out, err := RationalizeColor(tt.in)
			assert.NoError(err)
			assert.Equal(tt.out, out)
		}
	})

}
