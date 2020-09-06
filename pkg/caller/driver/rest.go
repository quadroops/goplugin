package driver

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/quadroops/goplugin/pkg/caller"
	"github.com/quadroops/goplugin/pkg/errs"
)

const (
	// PathPing {self explained}
	PathPing = "/ping"

	// PathExec {self explained}
	PathExec = "/exec"

	// Timeout used to waiting client response and cancel the request after limit timeout reached
	Timeout = 5
)

// RESTOption used as main option data
type RESTOption struct {
	Timeout int
}

// JSONData used as main response
type JSONData struct {
	Response interface{} `json:"response"`
}

// JSONResponse following JSEND standard as api response
type JSONResponse struct {
	Status string   `json:"status"`
	Data   JSONData `json:"data"`
}

// JSONExecPayload used as main payload when sending exec request
type JSONExecPayload struct {
	Cmd     string `json:"command"`
	Payload string `json:"payload"`
}

type rest struct {
	address string
	option  *RESTOption
}

// NewREST used to create new instance and return Caller interface
func NewREST(addr string, o *RESTOption) caller.Caller {
	var opt *RESTOption
	opt = &RESTOption{
		Timeout: Timeout,
	}

	if o != nil {
		opt = o
	}

	// override timeout if less than 1s or not defined
	if opt.Timeout == 0 {
		opt.Timeout = Timeout
	}

	return &rest{addr, opt}
}

func (r *rest) request(method, endpoint string, payload *bytes.Buffer) (*http.Response, error) {
	timeout := time.Duration(r.option.Timeout) * time.Second
	client := http.Client{
		Timeout: timeout,
	}

	var req *http.Request
	var err error

	if payload != nil {
		req, err = http.NewRequest(method, endpoint, payload)
	} else {
		req, err = http.NewRequest(method, endpoint, nil)
	}

	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	return client.Do(req)
}

func (r *rest) Ping() (string, error) {
	endpoint := fmt.Sprintf("%s%s", r.address, PathPing)
	resp, err := r.request("GET", endpoint, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errs.ErrPluginPing
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("%w: %q", errs.ErrPluginCall, err)
	}

	var response JSONResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("%w: %q", errs.ErrPluginCall, err)
	}

	return fmt.Sprintf("%v", response.Data.Response), nil
}

func (r *rest) Exec(cmdName string, payload []byte) ([]byte, error) {
	endpoint := fmt.Sprintf("%s%s", r.address, PathExec)
	p := JSONExecPayload{
		Cmd:     cmdName,
		Payload: hex.EncodeToString(payload),
	}

	jsonBody, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	resp, err := r.request("POST", endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return nil, errs.ErrPluginExec
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response JSONResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("%w: %q", errs.ErrPluginCall, err)
	}

	respStr, ok := response.Data.Response.(string)
	if !ok {
		return nil, errs.ErrCastInterface
	}

	b, err := hex.DecodeString(respStr)
	if err != nil {
		return nil, fmt.Errorf("%w: %q", errs.ErrPluginExec, err)
	}

	return b, nil
}
