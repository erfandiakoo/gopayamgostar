package gopayamgostar_test

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/erfandiakoo/gopayamgostar"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"
)

type configAdmin struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type Config struct {
	HostName string      `json:"hostname"`
	Proxy    string      `json:"proxy,omitempty"`
	Admin    configAdmin `json:"admin"`
}

var (
	config     *Config
	configOnce sync.Once
	setupOnce  sync.Once
	testUserID string
)

type RestyLogWriter struct {
	io.Writer
	t testing.TB
}

func (w *RestyLogWriter) Errorf(format string, v ...interface{}) {
	w.write("[ERROR] "+format, v...)
}

func (w *RestyLogWriter) Warnf(format string, v ...interface{}) {
	w.write("[WARN] "+format, v...)
}

func (w *RestyLogWriter) Debugf(format string, v ...interface{}) {
	w.write("[DEBUG] "+format, v...)
}

func (w *RestyLogWriter) write(format string, v ...interface{}) {
	w.t.Logf(format, v...)
}

func GetConfig(t testing.TB) *Config {
	configOnce.Do(func() {
		rand.Seed(uint64(time.Now().UTC().UnixNano()))
		configFileName, ok := os.LookupEnv("GOPAYAMGOSTAR_TEST_CONFIG")
		if !ok {
			configFileName = filepath.Join("testdata", "config.json")
		}
		configFile, err := os.Open(configFileName)
		require.NoError(t, err, "cannot open config.json")
		defer func() {
			err := configFile.Close()
			require.NoError(t, err, "cannot close config file")
		}()
		data, err := ioutil.ReadAll(configFile)
		require.NoError(t, err, "cannot read config.json")
		config = &Config{}
		err = json.Unmarshal(data, config)
		require.NoError(t, err, "cannot parse config.json")
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		if len(config.Proxy) != 0 {
			proxy, err := url.Parse(config.Proxy)
			require.NoError(t, err, "incorrect proxy url: "+config.Proxy)
			http.DefaultTransport.(*http.Transport).Proxy = http.ProxyURL(proxy)
		}
	})
	return config
}

func NewClientWithDebug(t testing.TB) *gopayamgostar.GoPayamgostar {
	cfg := GetConfig(t)
	client := gopayamgostar.NewClient(cfg.HostName)
	cond := func(resp *resty.Response, err error) bool {
		if resp != nil && resp.IsError() {
			if e, ok := resp.Error().(*gopayamgostar.HTTPErrorResponse); ok {
				msg := e.String()
				return strings.Contains(msg, "Cached clientScope not found") || strings.Contains(msg, "unknown_error")
			}
		}
		return false
	}

	restyClient := client.RestyClient()

	// restyClient.AddRetryCondition(
	// 	func(r *resty.Response, err error) bool {
	// 		if err != nil || r.RawResponse.StatusCode == 500 || r.RawResponse.StatusCode == 502 {
	// 			return true
	// 		}

	// 		return false
	// 	},
	// ).SetRetryCount(5).SetRetryWaitTime(10 * time.Millisecond)

	restyClient.
		// SetDebug(true).
		SetLogger(&RestyLogWriter{
			t: t,
		}).
		SetRetryCount(10).
		SetRetryWaitTime(2 * time.Second).
		AddRetryCondition(cond)

	return client
}

func Test_PostAuth(t *testing.T) {
	t.Parallel()
	cfg := GetConfig(t)
	client := NewClientWithDebug(t)
	newToken, err := client.PostAuth(
		context.Background(),
		cfg.Admin.UserName,
		cfg.Admin.Password,
	)
	require.NoError(t, err, "Login failed")
	t.Logf("New token: %+v", *newToken)
	require.Equal(t, newToken.ExpiresAt, 0, "Got a refresh token instead of offline")
	require.NotEmpty(t, newToken.AccessToken, "Got an empty if token")
}
