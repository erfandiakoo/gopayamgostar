package gopayamgostar

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type GoPayamgostar struct {
	basePath    string
	restyClient *resty.Client
	Config      struct {
		AuthEndpoint           string
		GetFormEndpoint        string
		CreateFormEndpoint     string
		FindFormEndpoint       string
		UpdateFormEndpoint     string
		GetPersonEndpoint      string
		FindPersonEndpoint     string
		CreatePurchaseEndpoint string
		DeletePurchaseEndpoint string
	}
}

const (
	adminClientID string = "admin-cli"
	urlSeparator  string = "/"
)

func makeURL(path ...string) string {
	return strings.Join(path, urlSeparator)
}

// GetRequest returns a request for calling endpoints.
func (g *GoPayamgostar) GetRequest(ctx context.Context) *resty.Request {
	var err HTTPErrorResponse
	return injectTracingHeaders(
		ctx, g.restyClient.R().
			SetContext(ctx).
			SetError(&err),
	)
}

func getID(resp *resty.Response) (string, error) {
	// Define a struct to match the expected response structure
	var result struct {
		CrmId string `json:"crmId"`
	}

	// Unmarshal the response body into the result struct
	err := json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Return the CrmId if available
	return result.CrmId, nil
}

func injectTracingHeaders(ctx context.Context, req *resty.Request) *resty.Request {
	// look for span in context, do nothing if span is not found
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return req
	}

	// look for tracer in context, use global tracer if not found
	tracer, ok := ctx.Value(tracerContextKey).(opentracing.Tracer)
	if !ok || tracer == nil {
		tracer = opentracing.GlobalTracer()
	}

	// inject tracing header into request
	err := tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	if err != nil {
		return req
	}

	return req
}

// GetRequestWithBearerAuthNoCache returns a JSON base request configured with an auth token and no-cache header.
func (g *GoPayamgostar) GetRequestWithBearerAuthNoCache(ctx context.Context, token string) *resty.Request {
	return g.GetRequest(ctx).
		SetAuthToken(token).
		SetHeader("Content-Type", "application/json").
		SetHeader("Cache-Control", "no-cache")
}

// GetRequestWithBearerAuth returns a JSON base request configured with an auth token.
func (g *GoPayamgostar) GetRequestWithBearerAuth(ctx context.Context, token string) *resty.Request {
	return g.GetRequest(ctx).
		SetAuthToken(token).
		SetHeader("Content-Type", "application/json")
}

func NewClient(basePath string, options ...func(*GoPayamgostar)) *GoPayamgostar {
	c := GoPayamgostar{
		basePath:    strings.TrimRight(basePath, urlSeparator),
		restyClient: resty.New(),
	}

	c.Config.AuthEndpoint = makeURL("api", "v2", "auth", "login")
	c.Config.GetFormEndpoint = makeURL("api", "v2", "crmobject", "form", "get")
	c.Config.CreateFormEndpoint = makeURL("api", "v2", "crmobject", "form", "create")
	c.Config.UpdateFormEndpoint = makeURL("api", "v2", "crmobject", "form", "update")
	c.Config.FindFormEndpoint = makeURL("api", "v2", "crmobject", "form", "find")
	c.Config.GetPersonEndpoint = makeURL("api", "v2", "crmobject", "person", "get")
	c.Config.FindPersonEndpoint = makeURL("api", "v2", "crmobject", "person", "find")
	c.Config.CreatePurchaseEndpoint = makeURL("api", "v2", "crmobject", "invoice", "purchase", "create")
	c.Config.DeletePurchaseEndpoint = makeURL("api", "v2", "crmobject", "invoice", "purchase", "delete")

	for _, option := range options {
		option(&c)
	}

	return &c
}

// RestyClient returns the internal resty g.
// This can be used to configure the g.
func (g *GoPayamgostar) RestyClient() *resty.Client {
	return g.restyClient
}

// SetRestyClient overwrites the internal resty g.
func (g *GoPayamgostar) SetRestyClient(restyClient *resty.Client) {
	g.restyClient = restyClient
}

func checkForError(resp *resty.Response, err error, errMessage string) error {
	if err != nil {
		return &APIError{
			Code:    0,
			Message: errors.Wrap(err, errMessage).Error(),
			Type:    ParseAPIErrType(err),
		}
	}

	if resp == nil {
		return &APIError{
			Message: "empty response",
			Type:    ParseAPIErrType(err),
		}
	}

	if resp.IsError() {
		var msg string

		if e, ok := resp.Error().(*HTTPErrorResponse); ok && e.NotEmpty() {
			msg = fmt.Sprintf("%s: %s", resp.Status(), e)
		} else {
			msg = resp.Status()
		}

		return &APIError{
			Code:    resp.StatusCode(),
			Message: msg,
			Type:    ParseAPIErrType(err),
		}
	}

	return nil
}

func (g *GoPayamgostar) getFullEndpointURL(path ...string) string {
	path = append([]string{g.basePath, g.Config.AuthEndpoint}, path...)
	return makeURL(path...)
}

func (g *GoPayamgostar) AdminAuthenticate(ctx context.Context, username string, password string) (*JWT, error) {
	const errMessage = "could not get token"

	var token JWT
	var req *resty.Request

	// Initialize the request here
	req = g.GetRequest(ctx)

	model := AuthRequest{
		Username:     username,
		Password:     password,
		PlatformType: 1,
		DeviceId:     uuid.NewString(),
	}
	resp, err := req.SetBody(model).
		SetResult(&token).
		Post(g.basePath + "/" + g.Config.AuthEndpoint)

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &token, nil
}

func (g *GoPayamgostar) UserAuthenticate(ctx context.Context, username string, password string) (*JWT, error) {
	const errMessage = "could not get token(customer)"

	var token JWT
	var req *resty.Request

	// Initialize the request here
	req = g.GetRequest(ctx)

	model := AuthRequest{
		Username:     username,
		Password:     password,
		PlatformType: 1,
		DeviceId:     uuid.NewString(),
	}
	resp, err := req.SetBody(model).
		SetResult(&token).
		Post(g.basePath + "/" + g.Config.AuthEndpoint)

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &token, nil
}

func (g *GoPayamgostar) GetPersonInfoById(ctx context.Context, accessToken, crmId string) (*PersonInfo, error) {
	const errMessage = "could not get user info"

	var result PersonInfo

	model := GetRequest{
		ID:                   crmId,
		ShowPreviews:         *BoolP(false),
		ShowExtendedPreviews: *BoolP(true),
	}

	resp, err := g.GetRequestWithBearerAuth(ctx, accessToken).
		SetBody(model).
		SetResult(&result).
		Post(g.basePath + "/" + g.Config.GetPersonEndpoint)

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

func (g *GoPayamgostar) GetFormInfoById(ctx context.Context, accessToken, crmId string) (*FormInfo, error) {
	const errMessage = "could not get form info"

	var result FormInfo

	model := GetRequest{
		ID:                   crmId,
		ShowPreviews:         *BoolP(true),
		ShowExtendedPreviews: *BoolP(true),
	}

	resp, err := g.GetRequestWithBearerAuth(ctx, accessToken).
		SetBody(model).
		SetResult(&result).
		Post(g.basePath + "/" + g.Config.GetFormEndpoint)

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	return &result, nil
}

func (g *GoPayamgostar) CreatePurchase(ctx context.Context, accessToken string, purchase CreatePurchase) (string, error) {
	const errMessage = "could not create purchase"

	resp, err := g.GetRequestWithBearerAuth(ctx, accessToken).
		SetBody(purchase).
		Post(g.basePath + "/" + g.Config.CreatePurchaseEndpoint)

	if err := checkForError(resp, err, errMessage); err != nil {
		return "", err
	}

	crmid, err := getID(resp)
	if err != nil {
		return "", err
	}

	return crmid, nil
}

func (g *GoPayamgostar) DeletePurchase(ctx context.Context, accessToken string, purchaseID string) error {
	const errMessage = "could not delete purchase"

	request := DeleteRequest{
		Id:     purchaseID,
		Option: 1,
	}

	resp, err := g.GetRequestWithBearerAuth(ctx, accessToken).
		SetBody(request).
		Post(g.basePath + "/" + g.Config.DeletePurchaseEndpoint)

	return checkForError(resp, err, errMessage)
}

func (g *GoPayamgostar) FindPersonByName(ctx context.Context, accessToken string, typeKey string, firstName string, lastName string) (*FindResponse, error) {
	const errMessage = "could find person"

	var result FindResponse

	request := FindRequest{
		TypeKey: typeKey,
		Queries: []Query{
			{
				LogicalOperator: 0,
				Operator:        0,
				Field:           "FirstName",
				Value:           firstName,
			},
			{
				LogicalOperator: 0,
				Operator:        0,
				Field:           "LastName",
				Value:           lastName,
			},
		},
		PageNumber: 1,
		PageSize:   10,
	}

	resp, err := g.GetRequestWithBearerAuthNoCache(ctx, accessToken).
		SetBody(request).
		Post(g.basePath + "/" + g.Config.FindPersonEndpoint)

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	// Unmarshal response into the result struct
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("%s: %w", errMessage, err)
	}

	// Return the result
	return &result, nil
}

func (g *GoPayamgostar) FindForm(ctx context.Context, accessToken string, typeKey string, queries []Query) (*FindFormResponse, error) {
	const errMessage = "could find form"

	var result FindFormResponse

	request := FindRequest{
		TypeKey:    *StringP(typeKey),
		Queries:    queries,
		PageNumber: *Int64P(1),
		PageSize:   *Int64P(10),
	}

	resp, err := g.GetRequestWithBearerAuthNoCache(ctx, accessToken).
		SetBody(request).
		Post(g.basePath + "/" + g.Config.FindFormEndpoint)

	if err := checkForError(resp, err, errMessage); err != nil {
		return nil, err
	}

	// Unmarshal response into the result struct
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("%s: %w", errMessage, err)
	}

	// Return the result
	return &result, nil
}

func (g *GoPayamgostar) UpdateForm(ctx context.Context, accessToken string, request UpdateFormRequest) (string, error) {
	const errMessage = "could not update form"

	resp, err := g.GetRequestWithBearerAuthNoCache(ctx, accessToken).
		SetBody(request).
		Post(g.basePath + "/" + g.Config.UpdateFormEndpoint)

	if err := checkForError(resp, err, errMessage); err != nil {
		return "", err
	}

	crmid, err := getID(resp)
	if err != nil {
		return "", err
	}

	return crmid, nil
}
