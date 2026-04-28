package handler

import (
	"crypto/subtle"
	"devture-matrix-corporal/corporal/httpapi/handler"
	"devture-matrix-corporal/corporal/httphelp"
	"devture-matrix-corporal/corporal/matrix"
	"encoding/json"
	"matrix-corporal-multitenant-proxy/instance"
	"matrix-corporal-multitenant-proxy/mtp_policy"
	"matrix-corporal-multitenant-proxy/reconcile"
	"strings"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type PolicyHandler struct {
	instances instance.InstanceRepository
}

func NewPolicyHandler(instances instance.InstanceRepository) *PolicyHandler {
	return &PolicyHandler{
		instances: instances,
	}
}

func (me *PolicyHandler) RegisterRoutesWithRouter(router *mux.Router) {
	router.HandleFunc("/_matrix/corporal/policy/{tenant}", me.policyPutTenantInUrl).Methods("PUT")
	router.HandleFunc("/_matrix/corporal/policy", me.policyPutTenantIsHost).Methods("PUT")
}

func (me *PolicyHandler) policyPutTenantInUrl(w http.ResponseWriter, r *http.Request) {
	tenant := mux.Vars(r)["tenant"]
	me.policyPut(w, r, tenant)
}

func (me *PolicyHandler) policyPutTenantIsHost(w http.ResponseWriter, r *http.Request) {
	tenant := strings.Split(r.Host, ":")[0]
	me.policyPut(w, r, tenant)
}

func (me *PolicyHandler) policyPut(w http.ResponseWriter, r *http.Request, tenant string) {
	access_token := httphelp.GetAccessTokenFromRequest(r)
	if access_token == "" {
		logrus.Info(`Policy Put: missing access token`)
		handler.Respond(w, http.StatusUnauthorized, handler.ApiResponseError{
			ErrorCode:    handler.ErrorCodeMissingToken,
			ErrorMessage: "Missing access token",
		})
		return
	}

	instance := me.instances.GetInstance(tenant)
	if instance == nil {
		logrus.Infof(`Policy Put: tenant not found: %s`, tenant)
		handler.Respond(w, http.StatusNotFound, handler.ApiResponseError{
			ErrorCode:    matrix.ErrorNotFound,
			ErrorMessage: "Tenant not found",
		})
		return
	}

	if subtle.ConstantTimeCompare([]byte(access_token), []byte(instance.Config.ApiAccessToken)) != 1 {
		logrus.Infof(`Policy Put: Wrong access token: tenant: %s`, tenant)
		handler.Respond(w, http.StatusUnauthorized, handler.ApiResponseError{
			ErrorCode:    handler.ErrorCodeUnknownToken,
			ErrorMessage: "Bad access token",
		})
		return
	}

	bodyBytes, err := httphelp.GetRequestBody(r)
	if err != nil {
		logrus.Infof(`Policy Put: Failed reading request body: %s`, err.Error())
		handler.Respond(w, http.StatusBadRequest, handler.ApiResponseError{
			ErrorCode:    handler.ErrorCodeBadJson,
			ErrorMessage: "Bad body payload",
		})
		return
	}

	logrus.Debugf(`request body: %s`, string(bodyBytes))

	var policy mtp_policy.MultiTenantPolicy
	err = json.Unmarshal(bodyBytes, &policy)
	if err != nil {
		logrus.Infof(`Policy Put: Could not decode json: %s`, err.Error())
		handler.Respond(w, http.StatusBadRequest, handler.ApiResponseError{
			ErrorCode:    handler.ErrorCodeBadJson,
			ErrorMessage: "Bad body payload",
		})
		return
	}

	if err := mtp_policy.Validate(policy); err != nil {
		logrus.Warnf(`Policy Put: Policy invalid: %s`, err.Error())
		handler.Respond(w, http.StatusBadRequest, handler.ApiResponseError{
			ErrorCode:    handler.ErrorCodeBadJson,
			ErrorMessage: err.Error(),
		})
		return
	}

	reconcile.UpdateIncomingPolicy(me.instances, &policy, tenant)

	handler.Respond(w, http.StatusOK, map[string]interface{}{})
}
