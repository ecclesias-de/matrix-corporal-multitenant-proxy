package mtp_policy

import "devture-matrix-corporal/corporal/policy"

type MultiTenantPolicy struct {
	policy.Policy
	RemoteUsers []*RemoteUserPolicy `json:"remoteUsers"`
}

type RemoteUserPolicy struct {
	Id          string              `json:"id"`
	JoinedRooms []*policy.RoomState `json:"joinedRooms"`
}

type ExchangePolicy struct {
	User []*RemoteUserPolicy `json:"remoteUsers"`
}
