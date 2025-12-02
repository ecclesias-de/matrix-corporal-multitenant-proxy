package test

import (
	"devture-matrix-corporal/corporal/configuration"
	"devture-matrix-corporal/corporal/httpapi"
	"devture-matrix-corporal/corporal/httpapi/handler"
	"devture-matrix-corporal/corporal/httphelp"
	"devture-matrix-corporal/corporal/policy"
	"devture-matrix-corporal/corporal/policy/provider"
	"io"
	"net"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func GetListenerAddr(t *testing.T) string {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}

	defer listener.Close()

	return listener.Addr().String()
}

// bootstrapping matrix corporal https api, store and last seen provider
func CreateMatrixCorporalInstance(t *testing.T, tenant string) (string, string, string, func()) {
	logger := logrus.New()

	store := policy.NewStore(logger, policy.NewValidator(tenant))
	cachePath := t.TempDir() + "/" + tenant + "-policy.json"
	provider, err := provider.NewLastSeenStorePolicyProvider(configuration.PolicyProvider{
		"CachePath": cachePath,
	}, store, logger)
	if err != nil {
		t.Fatal(err)
	}
	handler := handler.NewPolicyApiHandlerRegistrator(store, provider)

	configuration := configuration.HttpApi{
		Enabled:                  true,
		ListenAddress:            GetListenerAddr(t),
		AuthorizationBearerToken: "well-known-test-token-corporal-" + tenant,
	}
	server := httpapi.NewServer(logger, configuration, []httphelp.HandlerRegistrator{handler}, time.Second*5)

	if err := provider.Start(); err != nil {
		t.Fatal(err)
	}

	if err := server.Start(); err != nil {
		t.Fatal(err)
	}

	return configuration.ListenAddress, configuration.AuthorizationBearerToken, cachePath, func() {
		if err := server.Stop(); err != nil {
			t.Log(err)
		}

		provider.Stop()
	}
}

func TestFileMachetes(t *testing.T, path string, expected string) {
	file, err := os.Open(path)
	if err != nil {
		t.Log(err)

		if expected == "" {
			// if expected is "" and there is an os error => "e.g. file dose not exits" this test is treaded as successful
			return
		}

		t.Error(err)
	}
	content, err := io.ReadAll(file)
	if err != nil {
		t.Error(err)
	}

	if string(content) != expected {
		t.Errorf(`Strings do not match. got: %s expected: %s`, string(content), expected)
	}
}
