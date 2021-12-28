package fsapi

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/christophberger/3sixty/internal/xml"
)

type fsapi struct {
	url string
	pin string
	sid string
	xq  *xml.Query
}

const (
	responseTag = "fsapiResponse"
	statusOK    = "FS_OK"
)

func New(url, pin string) *fsapi {
	return &fsapi{
		url: url,
		pin: pin,
		xq:  xml.New()}
}

func (f *fsapi) CreateSession() error {

	url := "http://k--che.fritz.box/fsapi/CREATE_SESSION?pin=1812"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return fmt.Errorf("CreateSession: creating request failed:", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("CreateSession: running request failed:", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("CreateSession: cannot read body:", err)
	}

	f.sid, err = f.get(body, "sessionId")
	if err != nil {
		return fmt.Errorf("CreateSession: cannot get SID: %w", err)
	}
	return nil
}

func (f fsapi) get(body []byte, path string) (string, error) {
	status, err := f.xq.Get(body, "fsapiResponse.status")
	if err != nil {
		return "", fmt.Errorf("get: cannot get status: %w", err)
	}
	if status != statusOK {
		return "", fmt.Errorf("get: status is %s", status)
	}
	val, err := f.xq.Get(body, "fsapiResponse."+path)
	if err != nil {
		return "", err
	}
	return val, err
}

func (f fsapi) Sid() string {
	return f.sid
}
