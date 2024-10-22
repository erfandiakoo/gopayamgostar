package gopayamgostar_test

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
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

	"github.com/erfandiakoo/gopayamgostar/v2"
	"github.com/erfandiakoo/gopayamgostar/v2/shared/enums"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"
)

type configAdmin struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type configUser struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type Config struct {
	HostName string      `json:"hostname"`
	Proxy    string      `json:"proxy,omitempty"`
	Admin    configAdmin `json:"admin"`
	User     configUser  `json:"user"`
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

// FailRequest fails requests and returns an error
//
//	err - returned error or nil to return the default error
//	failN - number of requests to be failed
//	skipN = number of requests to be executed and not failed by this function
func FailRequest(client *gopayamgostar.GoPayamgostar, err error, failN, skipN int) *gopayamgostar.GoPayamgostar {
	client.RestyClient().OnBeforeRequest(
		func(c *resty.Client, r *resty.Request) error {
			if skipN > 0 {
				skipN--
				return nil
			}
			if failN == 0 {
				return nil
			}
			failN--
			if err == nil {
				err = fmt.Errorf("an error for request: %+v", r)
			}
			return err
		},
	)
	return client
}

func GetToken(t testing.TB, client *gopayamgostar.GoPayamgostar) *gopayamgostar.JWT {
	cfg := GetConfig(t)
	token, err := client.AdminAuthenticate(
		context.Background(),
		cfg.Admin.UserName,
		cfg.Admin.Password,
	)
	require.NoError(t, err, "Login failed")
	return token
}

// ---------
// API tests
// ---------

func Test_AdminAuthenticate(t *testing.T) {
	t.Parallel()
	cfg := GetConfig(t)
	client := NewClientWithDebug(t)
	newToken, err := client.AdminAuthenticate(
		context.Background(),
		cfg.Admin.UserName,
		cfg.Admin.Password,
	)
	require.NoError(t, err, "Login failed")
	t.Logf("New token: %+v", *newToken)
	//require.Equal(t, newToken.ExpiresAt, 0, "Got a refresh token instead of offline")
	require.NotEmpty(t, newToken.AccessToken, "Got an empty if token")
}

func Test_UserAuthenticate(t *testing.T) {
	t.Parallel()
	cfg := GetConfig(t)
	client := NewClientWithDebug(t)
	newToken, err := client.UserAuthenticate(
		context.Background(),
		cfg.User.UserName,
		cfg.User.Password,
	)
	require.NoError(t, err, "User Login failed")
	t.Logf("New token: %+v", *newToken)
	//require.Equal(t, newToken.ExpiresAt, 0, "Got a refresh token instead of offline")
	require.NotEmpty(t, newToken.AccessToken, "Got an empty if token")
}

func GetUserInfo(t *testing.T) {
	t.Parallel()
	client := NewClientWithDebug(t)
	token := GetToken(t, client)
	userInfo, err := client.GetPersonInfoById(
		context.Background(),
		token.AccessToken,
		"f845cf77-fec4-4631-b106-7f3d8580321b",
	)
	require.NoError(t, err, "Failed to fetch userinfo")
	t.Log(userInfo)
	FailRequest(client, nil, 1, 0)
	_, err = client.GetPersonInfoById(
		context.Background(),
		token.AccessToken,
		"f845cf77-fec4-4631-b106-7f3d8580321b")
	require.Error(t, err, "")
}

func Test_GetUserInfo(t *testing.T) {
	t.Parallel()
	client := NewClientWithDebug(t)
	token := GetToken(t, client)
	userInfo, err := client.GetPersonInfoById(
		context.Background(),
		token.AccessToken,
		"f845cf77-fec4-4631-b106-7f3d8580321b",
	)
	require.NoError(t, err, "Failed to fetch userinfo")
	t.Log(userInfo)
	FailRequest(client, nil, 1, 0)
	_, err = client.GetPersonInfoById(
		context.Background(),
		token.AccessToken,
		"f845cf77-fec4-4631-b106-7f3d8580321b")
	require.Error(t, err, "")
}

func Test_GetFormInfo(t *testing.T) {
	t.Parallel()
	client := NewClientWithDebug(t)
	token := GetToken(t, client)
	formInfo, err := client.GetPersonInfoById(
		context.Background(),
		token.AccessToken,
		"d81d07dd-cdc2-479a-99d5-0270a1f8f07d",
	)
	require.NoError(t, err, "Failed to fetch forminfo")
	t.Log(formInfo)
	FailRequest(client, nil, 1, 0)
	_, err = client.GetPersonInfoById(
		context.Background(),
		token.AccessToken,
		"d81d07dd-cdc2-479a-99d5-0270a1f8f07d")
	require.Error(t, err, "")
}

func CreatePurchase(t *testing.T, client *gopayamgostar.GoPayamgostar) (func(), string) {
	token := GetToken(t, client)

	purchase := gopayamgostar.CreatePurchase{
		CRMObjectTypeCode: "PurchaseInvoice",
		Details: []gopayamgostar.Detail{
			{
				IsService:      *gopayamgostar.BoolP(true),
				BaseUnitPrice:  *gopayamgostar.Int64P(1000),
				FinalUnitPrice: *gopayamgostar.Int64P(1000),
				TotalUnitPrice: *gopayamgostar.Int64P(1000),
				Count:          1,
				ReturnedCount:  1,
				TotalVat:       0,
				TotalToll:      0,
				TotalDiscount:  0,
				ProductCode:    "product-2",
			},
		},
		FinalValue: *gopayamgostar.Int64P(1000),
		TotalValue: *gopayamgostar.Int64P(1000),
		IdentityID: "f845cf77-fec4-4631-b106-7f3d8580321b",
		ColorID:    1,
		Discount:   0,
		Vat:        0,
		Toll:       0,
	}
	purchaseID, err := client.CreatePurchase(
		context.Background(),
		token.AccessToken,
		purchase,
	)

	require.NoError(t, err, "CreatePurchase failed")
	purchase.CrmId = purchaseID

	t.Logf("Created Purchase: %+v", purchase)
	// tearDown := func() {
	// 	err := client.DeletePurchase(
	// 		context.Background(),
	// 		token.AccessToken,
	// 		purchaseID)
	// 	require.NoError(t, err, "Delete Purchase")
	// }

	return nil, purchase.CRMObjectTypeCode
}

func Test_CreatePurchase(t *testing.T) {
	t.Parallel()
	client := NewClientWithDebug(t)

	tearDown, _ := CreatePurchase(t, client)
	defer tearDown()
}

func Test_FindPersonInfoByName(t *testing.T) {
	t.Parallel()
	client := NewClientWithDebug(t)
	token := GetToken(t, client)
	personInfo, err := client.FindPersonByName(
		context.Background(),
		token.AccessToken,
		*gopayamgostar.StringP("Kanon01"),
		*gopayamgostar.StringP("عرفان"),
		*gopayamgostar.StringP("دیاکونژاد"),
	)
	require.NoError(t, err, "Failed to fetch personInfo")
	t.Log(personInfo)
	FailRequest(client, nil, 1, 0)
	_, err = client.FindPersonByName(
		context.Background(),
		token.AccessToken,
		*gopayamgostar.StringP("Kanon01"),
		*gopayamgostar.StringP("عرفان"),
		*gopayamgostar.StringP("دیاکونژاد"),
	)
	require.Error(t, err, "")
}

func FindPersonInfoByName(t *testing.T) {
	t.Parallel()
	client := NewClientWithDebug(t)
	token := GetToken(t, client)
	userInfo, err := client.FindPersonByName(
		context.Background(),
		token.AccessToken,
		*gopayamgostar.StringP("Kanon01"),
		*gopayamgostar.StringP("عرفان"),
		*gopayamgostar.StringP("دیاکونژاد"),
	)
	require.NoError(t, err, "Failed to fetch personInfo")
	t.Log(userInfo)
	FailRequest(client, nil, 1, 0)
	_, err = client.FindPersonByName(
		context.Background(),
		token.AccessToken,
		*gopayamgostar.StringP("Kanon01"),
		*gopayamgostar.StringP("عرفان"),
		*gopayamgostar.StringP("دیاکونژاد"),
	)
	require.Error(t, err, "")
}

func Test_FindForm(t *testing.T) {
	t.Parallel()
	client := NewClientWithDebug(t)
	token := GetToken(t, client)

	// Define the queries for filtering BankAccount forms
	queries := []gopayamgostar.Query{
		{
			LogicalOperator: int(enums.And),
			Field:           *gopayamgostar.StringP("TrackingNumber"),
			Value:           *gopayamgostar.StringP("778756"),
		},
		{
			LogicalOperator: int(enums.And),
			Field:           *gopayamgostar.StringP("DepositAmount"),
			Value:           *gopayamgostar.StringP("625000000"),
		},
	}

	// Fetch the form information
	formInfo, err := client.FindForm(
		context.Background(),
		token.AccessToken,
		"BankAccount",
		queries,
	)
	require.NoError(t, err, "Failed to fetch formInfo")
	t.Log(formInfo)
	// Ensure that subsequent request results in an error (simulate failure)
	_, err = client.FindForm(
		context.Background(),
		token.AccessToken,
		"BankAccount",
		queries,
	)
	require.Error(t, err, "Expected an error after failing request")
}

func Test_UpdateForm(t *testing.T) {
	t.Parallel()
	client := NewClientWithDebug(t)
	token := GetToken(t, client)

	updateRequest := gopayamgostar.UpdateFormRequest{
		CrmId:              "d81d07dd-cdc2-479a-99d5-0270a1f8f07d",
		ParentCrmObjectId:  nil,
		ExtendedProperties: nil,
		Tags: []string{
			"تایید کارشناس",
		},
		StageId:    nil,
		ColorId:    1,
		IdentityId: "f845cf77-fec4-4631-b106-7f3d8580321b",
	}

	// Test successful request
	crmid, err := client.UpdateForm(
		context.Background(),
		token.AccessToken,
		updateRequest,
	)
	require.NoError(t, err, "Failed to update form")
	t.Log("CRMId:", crmid)

	// Test failure case (simulate request failure)
	FailRequest(client, nil, 1, 0)

	_, err = client.UpdateForm(
		context.Background(),
		token.AccessToken,
		updateRequest,
	)
	require.Error(t, err, "Expected error but got nil")
}

func Test_CreateForm(t *testing.T) {
	t.Parallel()
	client := NewClientWithDebug(t)
	token := GetToken(t, client)

	createRequest := gopayamgostar.CreateFormRequest{
		CRMObjectTypeCode: "SettlementRequest",
		ParentCRMObjectID: nil,
		ExtendedProperties: []gopayamgostar.ExtendedProperty{
			{
				UserKey: "DepositDate",
				Value:   *gopayamgostar.StringP("1403/12/12"),
			},
			{
				UserKey: "DepositAmount",
				Value:   *gopayamgostar.StringP("1000"),
			},
			{
				UserKey: "TrackingNumber",
				Value:   *gopayamgostar.StringP("112233"),
			},
			{
				UserKey: "CenterDetails",
				Value:   *gopayamgostar.StringP("CITY CENTER TEST"),
			},
			{
				UserKey: "fishType",
				Value:   *gopayamgostar.StringP("فردی"),
			},
			{
				UserKey: "fishSubType",
				Value:   *gopayamgostar.StringP("فردی"),
			},
			{
				UserKey: "Users",
				Value:   *gopayamgostar.StringP("1"),
			},
			{
				UserKey: "fishnum2",
				Value:   *gopayamgostar.StringP("131213"),
			},
			{
				UserKey: "Erjanum",
				Value:   *gopayamgostar.StringP("131213"),
			},
			{
				UserKey: "FurtherDescription",
				Value:   *gopayamgostar.StringP("test"),
			},
		},
		IdentityID:         "",
		Tags:               nil,
		RefID:              nil,
		ColorID:            1,
		AssignedToUserName: nil,
		StageID:            nil,
		Subject:            nil,
	}

	// Test successful request
	crmid, err := client.CreateForm(
		context.Background(),
		token.AccessToken,
		createRequest,
	)
	require.NoError(t, err, "Failed to create form")
	t.Log("CRMId:", crmid)

	// Test failure case (simulate request failure)
	FailRequest(client, nil, 1, 0)

	_, err = client.CreateForm(
		context.Background(),
		token.AccessToken,
		createRequest,
	)
	require.Error(t, err, "Expected error but got nil")
}
