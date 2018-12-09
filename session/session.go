package session

import (
	"crypto/md5"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

const (
	//DefaultSID is the default session id provided by fritzbox
	DefaultSID = "0000000000000000"

	//ExpirationTime indicates the default expiration time for a fritzbox session
	ExpirationTime = 10 * time.Minute
)

// Session holds the fritzbox session information
type Session struct {
	URL       string
	SID       string
	Challenge string
	BlockTime int8
	Username  string
}

func (s *Session) buildChallengeResponse(password string) string {
	hasher := md5.New()
	transformer := transform.NewWriter(hasher, unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder())
	transformer.Write([]byte(s.Challenge + "-" + password))
	hashed := hasher.Sum(nil)
	return fmt.Sprintf("%s-%x", s.Challenge, hashed)
}

//NewSessionByUsernamePassword authenticates to given fritzbox url and returns the session information
func NewSessionByUsernamePassword(fbURL, fbUsername, fbPassword string) (*Session, error) {
	s := &Session{}
	parsedURL, err := url.Parse(fbURL)
	if err != nil {
		return nil, err
	}
	s.URL = fmt.Sprintf("http://%s:%s", parsedURL.Hostname(), parsedURL.Port())
	s.Username = fbUsername

	loginURL := fmt.Sprintf("%s/login_sid.lua", s.URL)
	challengeCall, err := http.Get(loginURL)

	if err != nil {
		return nil, fmt.Errorf("invalid fbUrl passed: %v", err)
	}

	challengeRes := xml.NewDecoder(challengeCall.Body)
	challengeRes.Decode(&s)

	if s.SID == DefaultSID {
		loginCall, err := http.PostForm(loginURL, url.Values{
			"username": {fbUsername},
			"response": {s.buildChallengeResponse(fbPassword)},
		})

		if err != nil {
			return nil, err
		}

		loginRes := xml.NewDecoder(loginCall.Body)
		loginRes.Decode(s)

		return s, nil
	}
	return nil, errors.New("did not retrieve a valid challenge response")
}

//NewSessionByPassword authenticates to given fritzbox url and returns the session information
func NewSessionByPassword(fbURL, fbPassword string) (*Session, error) {
	return NewSessionByUsernamePassword(fbURL, "", fbPassword)
}
