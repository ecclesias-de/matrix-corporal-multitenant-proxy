package store

import (
	"devture-matrix-corporal/corporal/policy"
	"matrix-corporal-multitenant-proxy/mtp_policy"
)

type InMemoryStore struct {
	basePolicy       *policy.Policy
	exchangePolicies map[string]*mtp_policy.ExchangePolicy
	channels         []chan *UpdateEvent
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		exchangePolicies: make(map[string]*mtp_policy.ExchangePolicy),
		channels:         make([]chan *UpdateEvent, 0),
	}
}

func (me *InMemoryStore) SetBasePolicy(policy *policy.Policy) {
	me.basePolicy = policy

	me.update()
}

func (me *InMemoryStore) GetBasePolicy() *policy.Policy {
	return me.basePolicy
}

func (me *InMemoryStore) SetExchangePolicy(policy *mtp_policy.ExchangePolicy, from string) {
	me.exchangePolicies[from] = policy

	me.update()
}

func (me *InMemoryStore) GetExchangePolicies() []*mtp_policy.ExchangePolicy {
	policies := []*mtp_policy.ExchangePolicy{}

	for _, policy := range me.exchangePolicies {
		policies = append(policies, policy)
	}

	return policies
}

func (me *InMemoryStore) GetUpdateChannel() chan *UpdateEvent {
	channel := make(chan *UpdateEvent)

	me.channels = append(me.channels, channel)

	return channel
}

func (me *InMemoryStore) update() {
	if me.basePolicy == nil {
		return
	}

	updateEvent := &UpdateEvent{
		BasePolicies:     me.basePolicy,
		ExchangePolicies: me.GetExchangePolicies(),
	}

	for _, channel := range me.channels {
		channel <- updateEvent
	}
}
