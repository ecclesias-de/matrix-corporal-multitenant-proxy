package handler

import (
	"devture-matrix-corporal/corporal/httpapi/handler"
	"encoding/json"
	"io"
	"matrix-corporal-multitenant-proxy/config"
	"matrix-corporal-multitenant-proxy/instance"
	"matrix-corporal-multitenant-proxy/store"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestPolicyPut(t *testing.T) {
	instances := instance.NewInMemoryInstanceRepository()
	instances.AddInstance("a.test", &instance.Instance{
		Store: store.NewInMemoryStore(),
		Config: config.Instance{
			ApiAccessToken: "well-known-test-token",
		},
	})
	instances.AddInstance("b.test", &instance.Instance{
		Store: store.NewInMemoryStore(),
		Config: config.Instance{
			ApiAccessToken: "other-well-known-test-token",
		},
	})
	policyHandler := NewPolicyHandler(instances)

	router := mux.NewRouter()
	policyHandler.RegisterRoutesWithRouter(router)

	correctPolicy := `{"managedRoomIds": ["!1:a.test"], "users": [{"id":"@a:a.test", "joinedRooms": [{"roomId":"!1:a.test", "powerLevel": 10}]}], "remoteUsers": [{"id": "@a:b.test", "joinedRooms": [{"roomId": "!1:a.test", "powerLevel": 20}]}]}`

	for i, expected := range []struct {
		path          string
		authorization string
		body          string
		statusCode    int
		errorMessage  string
	}{
		{
			path:         "/_matrix/corporal/policy/a.test",
			body:         correctPolicy,
			statusCode:   http.StatusUnauthorized,
			errorMessage: "Missing access token",
		}, {
			path:          "/_matrix/corporal/policy/noExistentTenant.test",
			body:          correctPolicy,
			authorization: "Bearer well-known-test-token",
			statusCode:    http.StatusNotFound,
			errorMessage:  "Tenant not found",
		}, {
			path:          "/_matrix/corporal/policy/a.test",
			body:          correctPolicy,
			authorization: "Bearer wrong-token",
			statusCode:    http.StatusUnauthorized,
			errorMessage:  "Bad access token",
		}, {
			path:          "/_matrix/corporal/policy/a.test",
			body:          "notValidJson",
			authorization: "Bearer well-known-test-token",
			statusCode:    http.StatusBadRequest,
			errorMessage:  "Bad body payload",
		}, {
			path:          "/_matrix/corporal/policy/a.test",
			body:          correctPolicy,
			authorization: "Bearer well-known-test-token",
			statusCode:    http.StatusOK,
			errorMessage:  "",
		},
		// todo: move following test cases UpdateIncomingPolicyTest? They that function. And the validated function ...
		// no remote users
		{
			path:          "/_matrix/corporal/policy/a.test",
			body:          `{"managedRoomIds": ["!1:a.test"], "users": [{"id":"@a:a.test", "joinedRooms": [{"roomId":"!1:a.test", "powerLevel": 10}]}]}`,
			authorization: "Bearer well-known-test-token",
			statusCode:    http.StatusOK,
			errorMessage:  "",
		},
		// no users
		{
			path:          "/_matrix/corporal/policy/a.test",
			body:          `{"managedRoomIds": ["!1:a.test"], "remoteUsers": [{"id": "@a:b.test", "joinedRooms": [{"roomId": "!1:a.test", "powerLevel": 20}]}]}`,
			authorization: "Bearer well-known-test-token",
			statusCode:    http.StatusOK,
			errorMessage:  "",
		},
		// no joined rooms + no managed room ids
		{
			path:          "/_matrix/corporal/policy/a.test",
			body:          `{"users": [{"id":"@a:a.test"}], "remoteUsers": [{"id": "@a:b.test", "joinedRooms": [{"roomId": "!1:a.test", "powerLevel": 20}]}]}`,
			authorization: "Bearer well-known-test-token",
			statusCode:    http.StatusOK,
			errorMessage:  "",
		},
	} {
		reader := strings.NewReader(expected.body)
		req, err := http.NewRequest("PUT", expected.path, reader)
		if err != nil {
			t.Error(err)
			return
		}

		if expected.authorization != "" {
			req.Header.Add("Authorization", expected.authorization)
		}

		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		if res.Code != expected.statusCode {
			t.Errorf("%d: Returned wrong status code: got %v want %v", i, res.Code, expected.statusCode)
		}

		var apiError handler.ApiResponseError
		bytes, err := io.ReadAll(res.Body)
		if err != nil {
			t.Error(err)
			return
		}

		if err := json.Unmarshal(bytes, &apiError); err != nil {
			t.Error(err)
			return
		}

		if apiError.ErrorMessage != expected.errorMessage {
			t.Errorf(`%d: "Returned wrong error message: got %s want %s"`, i, apiError.ErrorMessage, expected.errorMessage)
		}
	}
}
