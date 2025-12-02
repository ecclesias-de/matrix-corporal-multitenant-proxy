package client

import (
	"devture-matrix-corporal/corporal/policy"
	"matrix-corporal-multitenant-proxy/config"
	"matrix-corporal-multitenant-proxy/test"
	"testing"
)

var serverAPolicy = policy.Policy{
	SchemaVersion:  2,
	ManagedRoomIds: []string{"!1:a.test"},
	User: []*policy.UserPolicy{
		{
			Id:          "@a:a.test",
			Active:      true,
			DisplayName: "a",
			AuthType:    "plain",
			JoinedRooms: []*policy.RoomState{
				{RoomId: "!1:a.test", PowerLevel: 10},
			},
		},
	},
}

func TestClient(t *testing.T) {
	addr, token, _, stop := test.CreateMatrixCorporalInstance(t, "a.test")

	defer stop()

	client := NewClient(config.Instance{
		CorporalApiUrl:      "http://" + addr,
		CorporalAccessToken: token,
	})

	if err := client.UpdatePolicy(&serverAPolicy); err != nil {
		t.Error(err)
	}
}
