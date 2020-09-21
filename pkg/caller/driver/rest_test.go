package driver_test

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/quadroops/goplugin/pkg/caller/driver"
	"github.com/quadroops/goplugin/pkg/errs"
	"github.com/stretchr/testify/assert"
	"github.com/vmihailenco/msgpack/v5"
)

func createServerPing(resp interface{}, status int) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := json.Marshal(resp)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		fmt.Fprintf(w, string(b))
	}))

	return server
}

func gethostport(addr string) (string, int) {
	u, _ := url.Parse(addr)
	log.Printf("Schema: %s", u.Scheme)
	log.Printf("Host: %s", u.Host)
	log.Printf("Port: %s", u.Port())

	port, _ := strconv.Atoi(u.Port())
	host := fmt.Sprintf("%s://%s", u.Scheme, u.Hostname())
	return host, port
}

func TestRESTPingSuccess(t *testing.T) {
	server := createServerPing(func() driver.JSONResponse {
		data := driver.JSONData{
			Response: "pong",
		}

		resp := driver.JSONResponse{
			Status: "success",
			Data:   data,
		}

		return resp
	}(), http.StatusOK)
	defer server.Close()

	host, port := gethostport(server.URL)
	rest := driver.NewREST(&driver.RESTOptions{
		Addr: host,
		Port: port,
	})

	resp, err := rest.Ping()
	assert.NoError(t, err)
	assert.Equal(t, "pong", resp)
	server.Close()
}

func TestPingFailed(t *testing.T) {
	server := createServerPing(func() interface{} {
		return "test"
	}(), http.StatusOK)
	defer server.Close()

	host, port := gethostport(server.URL)
	rest := driver.NewREST(&driver.RESTOptions{
		Addr: host,
		Port: port,
	})

	_, err := rest.Ping()
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrPluginCall))
}

func TestPingUnknownResponse(t *testing.T) {
	server := createServerPing(func() interface{} {
		return nil
	}(), http.StatusInternalServerError)
	defer server.Close()

	host, port := gethostport(server.URL)
	rest := driver.NewREST(&driver.RESTOptions{
		Addr: host,
		Port: port,
	})

	_, err := rest.Ping()
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrPluginPing))
}

func TestExecSuccess(t *testing.T) {
	server := createServerPing(func() interface{} {
		content := "test"
		b, _ := msgpack.Marshal(content)
		data := driver.JSONData{
			Response: hex.EncodeToString(b),
		}

		resp := driver.JSONResponse{
			Status: "success",
			Data:   data,
		}

		return resp
	}(), http.StatusAccepted)

	host, port := gethostport(server.URL)
	rest := driver.NewREST(&driver.RESTOptions{
		Addr: host,
		Port: port,
	})

	b, err := rest.Exec("rest.testing", []byte("test"))
	assert.NoError(t, err)

	var s string
	err = msgpack.Unmarshal(b, &s)
	assert.NoError(t, err)
	assert.Equal(t, "test", s)
}

func TestExecErrorContent(t *testing.T) {
	server := createServerPing(func() interface{} {
		return "test"
	}(), http.StatusAccepted)

	host, port := gethostport(server.URL)
	rest := driver.NewREST(&driver.RESTOptions{
		Addr: host,
		Port: port,
	})

	_, err := rest.Exec("rest.testing", []byte("test"))
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrPluginCall))
}

func TestExecErrorResponseInvalid(t *testing.T) {
	server := createServerPing(func() interface{} {
		data := driver.JSONData{
			Response: 10,
		}

		resp := driver.JSONResponse{
			Status: "success",
			Data:   data,
		}

		return resp
	}(), http.StatusAccepted)

	host, port := gethostport(server.URL)
	rest := driver.NewREST(&driver.RESTOptions{
		Addr: host,
		Port: port,
	})

	_, err := rest.Exec("rest.testing", []byte("test"))
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrCastInterface))
}

func TestExecErrorStatusCode(t *testing.T) {
	server := createServerPing(func() interface{} {
		data := driver.JSONData{
			Response: 10,
		}

		resp := driver.JSONResponse{
			Status: "success",
			Data:   data,
		}

		return resp
	}(), http.StatusInternalServerError)

	host, port := gethostport(server.URL)
	rest := driver.NewREST(&driver.RESTOptions{
		Addr: host,
		Port: port,
	})

	_, err := rest.Exec("rest.testing", []byte("test"))
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrPluginExec))
}
