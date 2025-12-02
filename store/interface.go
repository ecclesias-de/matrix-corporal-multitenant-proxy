package store

import (
	"devture-matrix-corporal/corporal/policy"
	"matrix-corporal-multitenant-proxy/mtp_policy"
)

type UpdateEvent struct {
	BasePolicies     *policy.Policy
	ExchangePolicies []*mtp_policy.ExchangePolicy
}

type Store interface {
	SetBasePolicy(policy *policy.Policy)
	GetBasePolicy() *policy.Policy
	SetExchangePolicy(policy *mtp_policy.ExchangePolicy, from string)
	GetExchangePolicies() []*mtp_policy.ExchangePolicy
	GetUpdateChannel() chan *UpdateEvent
}
