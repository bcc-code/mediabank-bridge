package proto

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ansel1/merry/v2"
	"github.com/davecgh/go-spew/spew"
	"github.com/rs/zerolog"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const audience = "media.bcc.mediabanken"
const auth0URL = "https://login.bcc.no/oauth/token"

type clientCredentialResponse struct {
	Token     string `json:"access_token"`
	ExpiresIn uint64 `json:"expires_in"`
	Type      string `json:"token_type"`
}

type credentialsProvider struct {
	clientID     string
	clientSecret string
	token        string
	expires      time.Time
	log          zerolog.Logger
}

// NewClient returns fully functional client, including all needed authentication etc providers
func NewClient(
	ctx context.Context,
	endpoint, clientID, clientSecret string,
	logger *zerolog.Logger,
	insecure bool, // Used for local debug purposes
) (MediabankBridgeClient, error) {
	log := logger.With().Str("package", "mediabank-bridge/proto").Logger()
	provider, err := newCredentialsProvider(clientID, clientSecret, log)
	if err != nil {
		return nil, err
	}

	conn, err := grpc.DialContext(ctx, endpoint, grpc.WithPerRPCCredentials(provider), grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	if err != nil {
		return nil, err
	}

	client := NewMediabankBridgeClient(conn)
	return client, nil
}

// NewCredentialsProvider for this api
func newCredentialsProvider(clientID, clientSecret string, log zerolog.Logger) (*credentialsProvider, error) {
	log.Debug().Msg("called NewCredentialsProvider")
	p := &credentialsProvider{
		clientID:     clientID,
		clientSecret: clientSecret,
		log:          log,
	}

	debugToken := os.Getenv("DEBUG_CLIENT_TOKEN")
	if debugToken != "" {
		log.Debug().Msg("Using DEBUG_CLIENT_TOKEN")
		p.token = debugToken
		p.expires = time.Now().Add(24 * time.Hour)
	}

	return p, p.refreshToken()
}

func (c credentialsProvider) RequireTransportSecurity() bool {
	return false
}

func (c *credentialsProvider) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	if err := c.refreshToken(); err != nil {
		return nil, err
	}

	return map[string]string{"token": c.token}, nil
}

func (c *credentialsProvider) refreshToken() error {
	if c.expires.After(time.Now()) {
		// We have a non-expired token. NOOP
		return nil
	}

	payload := strings.NewReader(
		fmt.Sprintf(
			"grant_type=client_credentials&client_id=%s&client_secret=%s&audience=%s",
			c.clientID, c.clientSecret, audience,
		),
	)

	req, _ := http.NewRequest("POST", auth0URL, payload)
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return merry.Wrap(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	response := &clientCredentialResponse{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return err
	}

	c.token = response.Token

	expiresTime, _ := time.ParseDuration(fmt.Sprintf("%ds", response.ExpiresIn))
	spew.Dump(expiresTime)
	c.expires = time.Now().Add(expiresTime)

	return nil
}
