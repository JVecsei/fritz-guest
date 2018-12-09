package guestmanager

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JVecsei/fritz-guest/psk"

	"github.com/JVecsei/fritz-guest/session"
)

func TestNewGuestManager(t *testing.T) {
	g, err := NewGuestManager(&session.Session{
		SID: "",
	})

	if g != nil {
		t.Errorf("g should be nil but was %v", g)
	}

	if err != ErrInvalidSession {
		t.Errorf("err should be ErrInvalidSession but was %v", err)
	}

}

func TestGuestManager_TurnOn(t *testing.T) {
	getConfig := false
	setConfig := false
	d := DataResponse{
		Data: Data{
			GuestAccess: GuestAccess{
				Psk:  "private",
				SSID: "public",
			},
		},
	}
	res, _ := json.Marshal(d)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() == "/data.lua" && req.Method == "POST" && req.FormValue("xhrId") == "all" {
			rw.Write(res)
			getConfig = true
		} else if req.URL.String() == "/data.lua" && req.Method == "POST" && req.FormValue("psk") != "" {
			if req.FormValue("psk") != "private" || req.FormValue("ssid") != "public" {
				t.Errorf("unexpected values '%s' should be 'private' and '%s' should be 'public'", req.FormValue("psk"), req.FormValue("ssid"))
			}
			rw.Write(res)
			setConfig = true
		}
	}))

	s := &session.Session{
		BlockTime: 0,
		Challenge: "321",
		SID:       "123",
		URL:       server.URL,
		Username:  "",
	}
	g, err := NewGuestManager(s)

	if err != nil {
		t.Errorf("err should be nil but was %v", err)
	}
	err = g.TurnOn()

	if err != nil {
		t.Errorf("err should be nil but was %v", err)
	}

	if !getConfig || !setConfig {
		t.Errorf("%v should be true and %v should be true", getConfig, setConfig)
	}

}

func TestGuestManager_TurnOnWithPsk(t *testing.T) {
	getConfig := false
	setConfig := false
	d := DataResponse{
		Data: Data{
			GuestAccess: GuestAccess{
				Psk:  "private",
				SSID: "public",
			},
		},
	}
	res, _ := json.Marshal(d)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() == "/data.lua" && req.Method == "POST" && req.FormValue("xhrId") == "all" {
			rw.Write(res)
			getConfig = true
		} else if req.URL.String() == "/data.lua" && req.Method == "POST" && req.FormValue("psk") != "" {
			if req.FormValue("psk") != "privateNew" || req.FormValue("ssid") != "public" {
				t.Errorf("unexpected values '%s' should be 'private' and '%s' should be 'public'", req.FormValue("psk"), req.FormValue("ssid"))
			}
			rw.Write(res)
			setConfig = true
		}
	}))

	s := &session.Session{
		BlockTime: 0,
		Challenge: "321",
		SID:       "123",
		URL:       server.URL,
		Username:  "",
	}
	g, err := NewGuestManager(s)

	if err != nil {
		t.Errorf("err should be nil but was %v", err)
	}
	err = g.TurnOnWithPsk(psk.FromString("privateNew"))

	if err != nil {
		t.Errorf("err should be nil but was %v", err)
	}

	if !getConfig || !setConfig {
		t.Errorf("%v should be true and %v should be true", getConfig, setConfig)
	}

}

func TestGuestManager_TurnOff(t *testing.T) {
	setConfig := false
	d := DataResponse{
		Data: Data{
			GuestAccess: GuestAccess{
				Psk:       "private",
				SSID:      "public",
				IsEnabled: "false",
			},
		},
	}
	res, _ := json.Marshal(d)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() == "/data.lua" && req.Method == "POST" && req.FormValue("isEnabled") == "0" {
			rw.Write(res)
			setConfig = true
		} else {
			t.Errorf("did not receive isEnabled = 0 request")
		}
	}))

	s := &session.Session{
		BlockTime: 0,
		Challenge: "321",
		SID:       "123",
		URL:       server.URL,
		Username:  "",
	}
	g, err := NewGuestManager(s)

	if err != nil {
		t.Errorf("err should be nil but was %v", err)
	}
	err = g.TurnOff()

	if err != nil {
		t.Errorf("err should be nil but was %v", err)
	}

	if !setConfig {
		t.Errorf("%v should be true", setConfig)
	}
}

func TestGuestManager_TurnOnWithRandomPsk(t *testing.T) {
	getConfig := false
	setConfig := false
	psk := psk.Random(20)
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() == "/data.lua" && req.Method == "POST" && req.FormValue("xhrId") == "all" {
			d := DataResponse{
				Data: Data{
					GuestAccess: GuestAccess{
						Psk:  "private",
						SSID: "public",
					},
				},
			}
			res, _ := json.Marshal(d)
			rw.Write(res)
			getConfig = true
		} else if req.URL.String() == "/data.lua" && req.Method == "POST" && req.FormValue("psk") != "" {
			if req.FormValue("psk") != psk.String() || req.FormValue("ssid") != "public" {
				t.Errorf("unexpected values '%s' should be '%s' and '%s' should be 'public'", req.FormValue("psk"), psk, req.FormValue("ssid"))
			}
			d := DataResponse{
				Data: Data{
					GuestAccess: GuestAccess{
						Psk:  "private",
						SSID: "public",
					},
				},
			}
			res, _ := json.Marshal(d)
			rw.Write(res)
			setConfig = true
		}
	}))

	s := &session.Session{
		BlockTime: 0,
		Challenge: "321",
		SID:       "123",
		URL:       server.URL,
		Username:  "",
	}
	g, err := NewGuestManager(s)

	if err != nil {
		t.Errorf("err should be nil but was %v", err)
	}
	err = g.TurnOnWithPsk(psk)

	if err != nil {
		t.Errorf("err should be nil but was %v", err)
	}

	if !getConfig || !setConfig {
		t.Errorf("%v should be true and %v should be true", getConfig, setConfig)
	}

}
