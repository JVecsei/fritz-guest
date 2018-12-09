package guestmanager

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/JVecsei/fritz-guest/psk"

	"github.com/JVecsei/fritz-guest/session"
)

var (
	//ErrNoPSK indicates that there is no valid PSK given for configuration
	ErrNoPSK = errors.New("no valid psk was found for configuration")
	//ErrInvalidSession indicates that no valid session id was found
	ErrInvalidSession = errors.New("invalid session information")
	//ErrTurnOffFailed indicates that the request was successful but guest access is still enabled
	ErrTurnOffFailed = errors.New("could not disable guest access - still enabled")
)

//DataResponse holds the important `data` object as well as additional information such as the used SID
type DataResponse struct {
	Data Data   `json:"data"`
	PID  string `json:"pid"`
	Sid  string `json:"sid"`
}

//Data is a response for guest access from data.lua
type Data struct {
	Timestamp   int64       `json:"timestamp"`
	GuestAccess GuestAccess `json:"guestAccess"`
	Ok          *bool       `json:"ok"`
	Alert       string      `json:"alert"`
}

//GuestAccess holds all guest network information that is returned from the fritzbox
type GuestAccess struct {
	LPTxt               string      `json:"lPTxt"`
	Psk                 string      `json:"psk"`
	IsIPClient          bool        `json:"isIpClient"`
	SupportsRegulation  bool        `json:"supportsRegulation"`
	Notification        string      `json:"notification"`
	NotificationEnabled string      `json:"notificationEnabled"`
	AutoUpdate          string      `json:"autoUpdate"`
	ShowGuest           bool        `json:"showGuest"`
	Mode                string      `json:"mode"`
	LPEnabled           string      `json:"lPEnabled"`
	LPRedirect          string      `json:"lPRedirect"`
	ActiveNexusClient   bool        `json:"activeNexusClient"`
	HideRepAutoUpdate   bool        `json:"hideRepAutoUpdate"`
	Timeout             string      `json:"timeout"`
	IsMaster            bool        `json:"isMaster"`
	LPImg               string      `json:"lPImg"`
	TimeoutNoForcedOff  string      `json:"timeoutNoForcedOff"`
	LPReguire           string      `json:"lPReguire"`
	DefaultSSID         DefaultSSID `json:"defaultSsid"`
	GuestGroupAccess    string      `json:"guestGroupAccess"`
	Isolated            string      `json:"isolated"`
	BoxType             string      `json:"boxType"`
	SSID                string      `json:"ssid"`
	IsEnabled           string      `json:"isEnabled"`
	LPRedirectURL       string      `json:"lPRedirectUrl"`
	WpsActive           bool        `json:"wpsActive"`
	IsTimeoutActive     string      `json:"isTimeoutActive"`
}

//DefaultSSID represents the default fritzbox guest ssid
type DefaultSSID struct {
	Private string `json:"private"`
	Public  string `json:"public"`
}

//GuestManager manages access to the guest network
type GuestManager struct {
	session *session.Session
}

//NewGuestManager returns a new guest manager
func NewGuestManager(s *session.Session) (*GuestManager, error) {
	if s.SID == session.DefaultSID || s.SID == "" {
		return nil, ErrInvalidSession
	}
	return &GuestManager{
		s,
	}, nil
}

//TurnOn turns the guest network on without changing its settings
func (g *GuestManager) TurnOn() error {
	return g.TurnOnWithPsk(psk.Noop())
}

//TurnOnWithPsk turns guest access on with given PSK. If PSK is empty it uses the currently configured PSK. If there is no PSK set and no PSK given this method will return an error (errNoPSK).
func (g *GuestManager) TurnOnWithPsk(p psk.Psk) error {
	dataURL := fmt.Sprintf("%s/data.lua", g.session.URL)
	currentConfigReq, err := http.PostForm(dataURL, url.Values{
		"sid":   {g.session.SID},
		"xhr":   {"1"},
		"page":  {"wGuest"},
		"lang":  {"de"},
		"xhrId": {"all"},
	})
	if err != nil {
		return err
	}
	dataRes := &DataResponse{}
	decoder := json.NewDecoder(currentConfigReq.Body)
	err = decoder.Decode(dataRes)

	if err != nil {
		return err
	}

	if p == "" && dataRes.Data.GuestAccess.Psk != "" {
		p = psk.FromString(dataRes.Data.GuestAccess.Psk)
	} else if p == "" {
		return ErrNoPSK
	}

	res, err := http.PostForm(dataURL, url.Values{
		"isEnabled":       {"1"},
		"guestAccessType": {"1"},
		"ssid":            {dataRes.Data.GuestAccess.SSID},
		"psk":             {p.String()},
		"sid":             {g.session.SID},
		"xhr":             {"1"},
		"page":            {"wGuest"},
		"lang":            {"de"},
		"timeout":         {"30"},
		"showGuest":       {"true"},
		"apply":           {""},
	})
	if err != nil {
		return err
	}
	turnOnDataResponse := &DataResponse{}
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(turnOnDataResponse)

	if err != nil {
		return err
	}

	if turnOnDataResponse.Data.Ok != nil {
		return fmt.Errorf("fb: %s", turnOnDataResponse.Data.Alert)
	}

	return nil
}

//TurnOff turns the guest network off
func (g *GuestManager) TurnOff() error {
	dataURL := fmt.Sprintf("%s/data.lua", g.session.URL)
	res, err := http.PostForm(dataURL, url.Values{
		"isEnabled":       {"0"},
		"guestAccessType": {"1"},
		"sid":             {g.session.SID},
		"xhr":             {"1"},
		"page":            {"wGuest"},
		"lang":            {"de"},
		"timeout":         {"30"},
		"showGuest":       {"true"},
		"apply":           {""},
	})

	if err != nil {
		return err
	}
	dataRes := &DataResponse{}
	dec := json.NewDecoder(res.Body)
	dec.Decode(dataRes)

	if dataRes.Data.Ok != nil {
		return fmt.Errorf("fb: %s", dataRes.Data.Alert)
	}

	if dataRes.Data.GuestAccess.IsEnabled != "false" {
		return ErrTurnOffFailed
	}
	return nil
}
