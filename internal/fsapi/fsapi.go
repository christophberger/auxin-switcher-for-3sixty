package fsapi

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/christophberger/3sixty/internal/xml"
)

type fsapi struct {
	url    string
	pin    string
	sid    string
	client *http.Client
}

const (
	responseTag = "fsapiResponse"
	statusOK    = "FS_OK"
)

func New(url, pin string) *fsapi {
	return &fsapi{
		url: url,
		pin: pin,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (f *fsapi) CreateSession() (err error) {
	query := fmt.Sprintf("CREATE_SESSION?pin=%s", f.pin)
	f.sid, err = f.get(query, "sessionId")
	if err != nil {
		return fmt.Errorf("CreateSession: cannot get SID: %w", err)
	}
	return nil
}

func (f *fsapi) SetMode(mode string) (err error) {
	query := fmt.Sprintf("SET/netRemote.sys.mode?pin=%s&sid=%s&value=%s", f.pin, f.sid, mode)
	_, err = f.get(query, "status")
	if err != nil {
		return fmt.Errorf("SetMode: cannot set mode: %w", err)
	}
	return nil
}

// get receives a query endpoint (minus the base URL) and a
// query pqth to the desired value in the XML response.
// It returns the value as a string, or an error if the query fails.
func (f fsapi) get(query, resPath string) (string, error) {
	endpoint := fmt.Sprintf("%s/%s", f.url, query)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("call: creating request failed:", err)
	}

	res, err := f.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("call: running request failed:", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("call: cannot read body:", err)
	}

	status, err := xml.Get(body, ".fsapiResponse.status")
	if err != nil {
		return "", fmt.Errorf("get: cannot get status: %w", err)
	}
	if status != statusOK {
		return "", fmt.Errorf("get: status is %s", status)
	}
	val, err := xml.Get(body, "fsapiResponse."+resPath)
	if err != nil {
		return "", err
	}
	return val, err
}

func (f fsapi) Sid() string {
	return f.sid
}
