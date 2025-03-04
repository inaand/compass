package tenantfetcher

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	"github.com/kyma-incubator/compass/components/director/pkg/log"
	bndlErrors "github.com/pkg/errors"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	maxErrMessageLength = 50
)

// OAuth2Config missing godoc
type OAuth2Config struct {
	ClientID           string `envconfig:"APP_CLIENT_ID"`
	ClientSecret       string `envconfig:"APP_CLIENT_SECRET"`
	OAuthTokenEndpoint string `envconfig:"APP_OAUTH_TOKEN_ENDPOINT"`
	TokenPath          string `envconfig:"optional,APP_OAUTH_TOKEN_PATH"`
}

// APIConfig missing godoc
type APIConfig struct {
	EndpointTenantCreated     string `envconfig:"optional,APP_ENDPOINT_TENANT_CREATED"`
	EndpointTenantDeleted     string `envconfig:"optional,APP_ENDPOINT_TENANT_DELETED"`
	EndpointTenantUpdated     string `envconfig:"optional,APP_ENDPOINT_TENANT_UPDATED"`
	EndpointSubaccountCreated string `envconfig:"optional,APP_ENDPOINT_SUBACCOUNT_CREATED"`
	EndpointSubaccountDeleted string `envconfig:"optional,APP_ENDPOINT_SUBACCOUNT_DELETED"`
	EndpointSubaccountUpdated string `envconfig:"optional,APP_ENDPOINT_SUBACCOUNT_UPDATED"`
	EndpointSubaccountMoved   string `envconfig:"optional,APP_ENDPOINT_SUBACCOUNT_MOVED"`
}

func (c APIConfig) isUnassignedOptionalProperty(eventsType EventsType) bool {
	if eventsType == MovedSubaccountType && len(c.EndpointSubaccountMoved) == 0 {
		return true
	}
	return false
}

// MetricsPusher missing godoc
//go:generate mockery --name=MetricsPusher --output=automock --outpkg=automock --case=underscore
type MetricsPusher interface {
	RecordEventingRequest(method string, statusCode int, desc string)
}

// QueryParams describes the key and the corresponding value for query parameters when requesting the service
type QueryParams map[string]string

// Client implements the communication with the service
type Client struct {
	httpClient    *http.Client
	metricsPusher MetricsPusher

	apiConfig APIConfig
}

// NewClient missing godoc
func NewClient(oAuth2Config OAuth2Config, apiConfig APIConfig, timeout time.Duration) *Client {
	cfg := clientcredentials.Config{
		ClientID:     oAuth2Config.ClientID,
		ClientSecret: oAuth2Config.ClientSecret,
		TokenURL:     oAuth2Config.OAuthTokenEndpoint + oAuth2Config.TokenPath,
	}

	httpClient := cfg.Client(context.Background())
	httpClient.Timeout = timeout

	return &Client{
		httpClient: httpClient,
		apiConfig:  apiConfig,
	}
}

// SetMetricsPusher missing godoc
func (c *Client) SetMetricsPusher(metricsPusher MetricsPusher) {
	c.metricsPusher = metricsPusher
}

// FetchTenantEventsPage missing godoc
func (c *Client) FetchTenantEventsPage(eventsType EventsType, additionalQueryParams QueryParams) (TenantEventsResponse, error) {
	if c.apiConfig.isUnassignedOptionalProperty(eventsType) {
		log.D().Warnf("Optional property for event type %s was not set", eventsType)
		return nil, nil
	}

	endpoint, err := c.getEndpointForEventsType(eventsType)
	if endpoint == "" && err == nil {
		log.D().Warnf("Endpoint for event %s is not set", eventsType)
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	reqURL, err := c.buildRequestURL(endpoint, additionalQueryParams)
	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Get(reqURL)
	if err != nil {
		if c.metricsPusher != nil {
			desc := c.failedRequestDesc(err)
			c.metricsPusher.RecordEventingRequest(http.MethodGet, 0, desc)
		}
		return nil, bndlErrors.Wrap(err, "while sending get request")
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			log.D().Warnf("Unable to close response body. Cause: %v", err)
		}
	}()

	if c.metricsPusher != nil {
		c.metricsPusher.RecordEventingRequest(http.MethodGet, res.StatusCode, res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, bndlErrors.Wrap(err, "while reading response body")
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNoContent {
		return nil, fmt.Errorf("request to %q returned status code %d and body %q", reqURL, res.StatusCode, bytes)
	}

	return bytes, nil
}

func (c *Client) getEndpointForEventsType(eventsType EventsType) (string, error) {
	switch eventsType {
	case CreatedAccountType:
		return c.apiConfig.EndpointTenantCreated, nil
	case DeletedAccountType:
		return c.apiConfig.EndpointTenantDeleted, nil
	case UpdatedAccountType:
		return c.apiConfig.EndpointTenantUpdated, nil
	case CreatedSubaccountType:
		return c.apiConfig.EndpointSubaccountCreated, nil
	case DeletedSubaccountType:
		return c.apiConfig.EndpointSubaccountDeleted, nil
	case UpdatedSubaccountType:
		return c.apiConfig.EndpointSubaccountUpdated, nil
	case MovedSubaccountType:
		return c.apiConfig.EndpointSubaccountMoved, nil
	default:
		return "", apperrors.NewInternalError("unknown events type")
	}
}

func (c *Client) buildRequestURL(endpoint string, queryParams QueryParams) (string, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}

	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return "", err
	}

	for qKey, qValue := range queryParams {
		q.Add(qKey, qValue)
	}

	u.RawQuery = q.Encode()

	return u.String(), nil
}

func (c *Client) failedRequestDesc(err error) string {
	var e *net.OpError
	if errors.As(err, &e) && e.Err != nil {
		return e.Err.Error()
	}

	if len(err.Error()) > maxErrMessageLength {
		// not all errors are actually wrapped, sometimes the error message is just concatenated with ":"
		errParts := strings.Split(err.Error(), ":")
		return errParts[len(errParts)-1]
	}

	return err.Error()
}
