package wiz

import (
	"encoding/json"
	"fmt"
)

type methodParams map[string]any

type request struct {
	Id     int          `json:"id"`
	Method wizMethod    `json:"method"`
	Params methodParams `json:"params"`
}

func (w *request) marshal() ([]byte, error) {
	return json.Marshal(w)
}

type deviceErrorReport struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type baseResponse struct {
	Method string             `json:"method"`
	Env    string             `json:"env"`
	Id     int                `json:"id"`
	Error  *deviceErrorReport `json:"error"`
}

func (b *baseResponse) method() string {
	return b.Method
}

func (b *baseResponse) id() int {
	return b.Id
}

func (b *baseResponse) error() *deviceErrorReport {
	return b.Error
}

type PilotResult struct {
	Mac     string `json:"mac"`
	Rssi    int    `json:"rssi"`
	State   bool   `json:"state"`
	SceneId int    `json:"sceneId"`
	R       int    `json:"r"`
	G       int    `json:"g"`
	B       int    `json:"b"`
	C       int    `json:"c"`
	W       int    `json:"w"`
	Temp    int    `json:"temp"`
	Dimming int    `json:"dimming"`
	Speed   int    `json:"speed"`
}

func (w *PilotResult) String() string {
	return fmt.Sprintf("Mac: %s, Rssi: %d, State: %t, SceneId: %d, R: %d, G: %d, B: %d, C: %d, W: %d, Temp: %d, Dimming: %d, Speed: %d",
		w.Mac, w.Rssi, w.State, w.SceneId, w.R, w.G, w.B, w.C, w.W, w.Temp, w.Dimming, w.Speed)
}

// {"method":"getPilot","env":"pro","result":{"mac":"9877d5466ed4","rssi":-54,"state":true,"sceneId":0,"temp":5800,"dimming":100}}
// {"method":"getPilot","env":"pro","result":{"mac":"cc4085c5a8de","rssi":-52,"state":true,"sceneId":0,"r":48,"g":255,"b":0,"c":0,"w":0,"dimming":100}}
// {"method":"setPilot","id":1,"env":"pro","error":{"code":-32602,"message":"Invalid params"}}

type pilotResponse struct {
	baseResponse
	Result PilotResult `json:"result"`
}

func (w *pilotResponse) unmarshal(from []byte) error {
	return json.Unmarshal(from, w)
}

// {"method":"setState","id":1,"env":"pro","result":{"success":true}}
type SetStateResult struct {
	Success bool `json:"success"`
}

type setStateResponse struct {
	baseResponse
	Result SetStateResult `json:"result"`
}

func (w *setStateResponse) unmarshal(from []byte) error {
	return json.Unmarshal(from, w)
}

// {"method":"setPilot","id":1,"env":"pro","result":{"success":true}}
type SetPilotResult struct {
	baseResponse
	Success bool `json:"success"`
}

type setPilotResponse struct {
	baseResponse
	Result SetPilotResult `json:"result"`
}

func (w *setPilotResponse) unmarshal(from []byte) error {
	return json.Unmarshal(from, w)
}

type PulseResult struct {
	baseResponse
	Success bool `json:"success"`
}

type pulseResponse struct {
	baseResponse
	Result PulseResult `json:"result"`
}

func (w *pulseResponse) unmarshal(from []byte) error {
	return json.Unmarshal(from, w)
}

// {"devTotal":1,"headTotal":1,"swHead":0,"ps":2,"hasGradient":1,"nightLightOff":0,"minDimLevel":10,"devices":1,"devType":0,"lightType":1,"pwmFreq":1000,"pwmRes":13,"pwmRange":[0,100],"pwmRanges":[0,1000,0,1000,0,1000,0,1000,0,1000],"wcr":40,"nowc":2,"cctRange":[2200,2700,6500,6500],"renderFactor":[255,170,130,255,0,0,40,200,255,200],"wizc1":{"mode":[0,0,0,0,0,0,4000],"opts":{"dim":100}},"wizc2":{"mode":[11,0,0,0,0,0,0],"opts":{"dim":100}},"drvIface":4,"i2cDrv":[{"chip":"BP5758D","addr":255,"freq":200,"curr":[15,8,12,22,22],"output":[3,2,1,4,5]}]}
type ModelConfig struct {
	DevTotal      int   `json:"devTotal"`
	HeadTotal     int   `json:"headTotal"`
	SwHead        int   `json:"swHead"`
	Ps            int   `json:"ps"`
	HasGradient   int   `json:"hasGradient"`
	NightLightOff int   `json:"nightLightOff"`
	MinDimLevel   int   `json:"minDimLevel"`
	Devices       int   `json:"devices"`
	DevType       int   `json:"devType"`
	LightType     int   `json:"lightType"`
	PwmFreq       int   `json:"pwmFreq"`
	PwmRes        int   `json:"pwmRes"`
	PwmRange      []int `json:"pwmRange"`
	PwmRanges     []int `json:"pwmRanges"`
	Wcr           int   `json:"wcr"`
	Nowc          int   `json:"nowc"`
	CctRange      []int `json:"cctRange"`
	RenderFactor  []int `json:"renderFactor"`
}

func (c *ModelConfig) String() string {
	return fmt.Sprintf("DevTotal: %d, HeadTotal: %d, SwHead: %d, Ps: %d, HasGradient: %d, NightLightOff: %d, MinDimLevel: %d, Devices: %d, DevType: %d, LightType: %d, PwmFreq: %d, PwmRes: %d, PwmRange: %v, PwmRanges: %v, Wcr: %d, Nowc: %d, CctRange: %v, RenderFactor: %v",
		c.DevTotal, c.HeadTotal, c.SwHead, c.Ps, c.HasGradient, c.NightLightOff, c.MinDimLevel, c.Devices, c.DevType, c.LightType, c.PwmFreq, c.PwmRes, c.PwmRange, c.PwmRanges, c.Wcr, c.Nowc, c.CctRange, c.RenderFactor)
}

type modelConfigResponse struct {
	baseResponse
	Result ModelConfig `json:"result"`
}

func (w *modelConfigResponse) unmarshal(from []byte) error {
	return json.Unmarshal(from, w)
}

// {"fadeIn":500,"fadeOut":500,"dftDim":100,"opMode":0,"po":false,"minDimming":0,"tapSensor":1,"autoUpd":1,"devices":1,"dim2WarmPoints":[[1800,1],[1800,10],[2700,50],[4200,90],[4200,100]],"wizc1":{"mode":[0,0,0,0,0,0,4000],"opts":{"dim":100}},"wizc2":{"mode":[11,0,0,0,0,0,0],"opts":{"dim":100}},"apStkEn":false,"confTs":1}
type UserConfig struct {
	FadeIn         int     `json:"fadeIn"`
	FadeOut        int     `json:"fadeOut"`
	DftDim         int     `json:"dftDim"`
	OpMode         int     `json:"opMode"`
	Po             bool    `json:"po"`
	MinDimming     int     `json:"minDimming"`
	TapSensor      int     `json:"tapSensor"`
	AutoUpd        int     `json:"autoUpd"`
	Devices        int     `json:"devices"`
	Dim2WarmPoints [][]int `json:"dim2WarmPoints"`
	ApStkEn        bool    `json:"apStkEn"`
	ConfTs         int     `json:"confTs"`
}

func (c *UserConfig) String() string {
	return fmt.Sprintf("FadeIn: %d, FadeOut: %d, DftDim: %d, OpMode: %d, Po: %t, MinDimming: %d, TapSensor: %d, AutoUpd: %d, Devices: %d, Dim2WarmPoints: %v, ApStkEn: %t, ConfTs: %d",
		c.FadeIn, c.FadeOut, c.DftDim, c.OpMode, c.Po, c.MinDimming, c.TapSensor, c.AutoUpd, c.Devices, c.Dim2WarmPoints, c.ApStkEn, c.ConfTs)
}

type userConfigResponse struct {
	baseResponse
	Result UserConfig `json:"result"`
}

func (w *userConfigResponse) unmarshal(from []byte) error {
	return json.Unmarshal(from, w)
}

// {"mac":"9877d5a6b73a","homeId":4170338,"roomId":37587433,"rgn":"eu","moduleName":"ESP25_SHRGB_01","fwVersion":"1.38.0","groupId":0,"ping":0,"accUdpPropRate":100,"rdIdUidHash":"642f99e91c9dac57fd1276a5d0f28712916612aa0301a39cce6e5f65c1ac3c80"}
type SystemConfig struct {
	Mac            string `json:"mac"`
	HomeId         int    `json:"homeId"`
	RoomId         int    `json:"roomId"`
	Rgn            string `json:"rgn"`
	ModuleName     string `json:"moduleName"`
	FwVersion      string `json:"fwVersion"`
	GroupId        int    `json:"groupId"`
	Ping           int    `json:"ping"`
	AccUdpPropRate int    `json:"accUdpPropRate"`
	RdIdUidHash    string `json:"rdIdUidHash"`
}

func (c *SystemConfig) String() string {
	return fmt.Sprintf("Mac: %s, HomeId: %d, RoomId: %d, Rgn: %s, ModuleName: %s, FwVersion: %s, GroupId: %d, Ping: %d, AccUdpPropRate: %d, RdIdUidHash: %s",
		c.Mac, c.HomeId, c.RoomId, c.Rgn, c.ModuleName, c.FwVersion, c.GroupId, c.Ping, c.AccUdpPropRate, c.RdIdUidHash)
}

type systemConfigResponse struct {
	baseResponse
	Result SystemConfig `json:"result"`
}

func (w *systemConfigResponse) unmarshal(from []byte) error {
	return json.Unmarshal(from, w)
}
