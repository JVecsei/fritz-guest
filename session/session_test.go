package session

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// todo
func TestSession_NewSessionByPassword_success(t *testing.T) {
	getCalled := false
	postCalled := false
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() == "/login_sid.lua" && req.Method == "GET" {
			rw.Write([]byte(`<SessionInfo>
			<SID>0000000000000000</SID>
			<Challenge>321</Challenge>
			<BlockTime>0</BlockTime>
		</SessionInfo>`))
			getCalled = true
		} else if req.URL.String() == "/login_sid.lua" && req.Method == "POST" {
			if req.FormValue("response") != "321-5056f210b5dd4c2a92a0d0bb9aced2fa" {
				t.Error("invalid challenge response sent")
			}
			rw.Write([]byte(`<SessionInfo>
			<SID>123</SID>
			<Challenge>321</Challenge>
			<BlockTime>0</BlockTime>
		</SessionInfo>`))
			postCalled = true
		}
	}))
	defer server.Close()
	s, err := NewSessionByPassword(server.URL, "test")

	if err != nil {
		t.Error("error while retreiving SID", err)
	}

	if s.SID == DefaultSID || s.SID == "" {
		t.Errorf("no valid SID received '%v'", s.SID)
	}

	if !getCalled {
		t.Error("did not get challenge request")
	}

	if !postCalled {
		t.Error("did not get response request")
	}

}

func TestSession_NewSessionByPassword_invalidDefaultSID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() == "/login_sid.lua" && req.Method == "GET" {
			rw.Write([]byte(`<SessionInfo>
			<SID>00000000000000001</SID>
			<Challenge>321</Challenge>
			<BlockTime>0</BlockTime>
		</SessionInfo>`))
		}
	}))
	defer server.Close()
	s, err := NewSessionByPassword(server.URL, "test")

	if err == nil {
		t.Error("err should not be nil", err)
	}

	if s != nil {
		t.Errorf("session should be nil but was %v", s)
	}

}

func TestSession_NewSessionByPassword_wrongURL(t *testing.T) {
	s, err := NewSessionByPassword("invalid_addr", "test")
	if err == nil {
		t.Error("err should not be nil", err)
	}

	if s != nil {
		t.Errorf("session should be nil but was %v", s)
	}

}

func TestSession_NewSessionPassword_failURLParse(t *testing.T) {
	s, err := NewSessionByPassword(" http://invalid_url_leading_whitespace", "")
	if err == nil {
		t.Error("err should not be nil", err)
	}

	if s != nil {
		t.Errorf("session should be nil but was %v", s)
	}

}
