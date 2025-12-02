package reconcile

import (
	"devture-matrix-corporal/corporal/policy"
	"encoding/json"
	"matrix-corporal-multitenant-proxy/instance"
	"matrix-corporal-multitenant-proxy/mtp_policy"
	"matrix-corporal-multitenant-proxy/store"
	"testing"
)

var serverAIncomingPolicy = mtp_policy.MultiTenantPolicy{
	Policy: policy.Policy{
		ManagedRoomIds: []string{"!1:a.test", "!2:a.test", "!3:a.test"},
		User: []*policy.UserPolicy{
			{
				Id:          "@a:a.test",
				Active:      true,
				DisplayName: "a",
				JoinedRooms: []*policy.RoomState{
					{RoomId: "!1:a.test", PowerLevel: 10},
					{RoomId: "!2:a.test", PowerLevel: 20},
				},
			},
			{
				Id:          "@b:a.test",
				Active:      true,
				DisplayName: "b",
				JoinedRooms: []*policy.RoomState{
					{RoomId: "!2:a.test", PowerLevel: 10},
				},
			},
			{
				Id:          "@c:a.test",
				Active:      true,
				DisplayName: "b",
				JoinedRooms: []*policy.RoomState{
					{RoomId: "!3:a.test", PowerLevel: 10},
				},
			},
		},
	},
	RemoteUsers: []*mtp_policy.RemoteUserPolicy{
		{
			Id: "@a:b.test",
			JoinedRooms: []*policy.RoomState{
				{RoomId: "!1:a.test"},
				{RoomId: "!2:a.test"},
			},
		},
		{
			Id: "@b:b.test",
			JoinedRooms: []*policy.RoomState{
				{RoomId: "!2:a.test"},
				{RoomId: "!3:a.test"},
			},
		},
		{
			Id: "@a:c.test",
			JoinedRooms: []*policy.RoomState{
				{RoomId: "!1:a.test"},
				{RoomId: "!3:a.test"},
			},
		},
	},
}

var serverABasePolicy = policy.Policy{
	ManagedRoomIds: []string{"!1:a.test", "!2:a.test", "!3:a.test"},
	User: []*policy.UserPolicy{
		{
			Id:          "@a:a.test",
			Active:      true,
			DisplayName: "a",
			JoinedRooms: []*policy.RoomState{
				{RoomId: "!1:a.test", PowerLevel: 10},
				{RoomId: "!2:a.test", PowerLevel: 20},
			},
		},
		{
			Id:          "@b:a.test",
			Active:      true,
			DisplayName: "b",
			JoinedRooms: []*policy.RoomState{
				{RoomId: "!2:a.test", PowerLevel: 10},
			},
		},
		{
			Id:          "@c:a.test",
			Active:      true,
			DisplayName: "b",
			JoinedRooms: []*policy.RoomState{
				{RoomId: "!3:a.test", PowerLevel: 10},
			},
		},
	},
}

var serverABExchangePolicy = mtp_policy.ExchangePolicy{
	User: []*mtp_policy.RemoteUserPolicy{
		{
			Id: "@a:a.test",
			JoinedRooms: []*policy.RoomState{
				{RoomId: "!3:a.test", PowerLevel: 10},
				{RoomId: "!4:a.test", PowerLevel: 10},
			},
		},
		{
			Id: "@b:a.test",
			JoinedRooms: []*policy.RoomState{
				{RoomId: "!2:a.test", PowerLevel: 30},
				{RoomId: "!3:a.test", PowerLevel: 10},
			},
		},
		{
			Id: "@d:a.test",
			JoinedRooms: []*policy.RoomState{
				{RoomId: "!2:a.test", PowerLevel: 10},
			},
		},
	},
}

func TestBuildOutgoingPolicy(t *testing.T) {
	serverABasePolicyBeforeJson, _ := json.Marshal(serverABasePolicy)

	updateEvent := &store.UpdateEvent{
		BasePolicies: &serverABasePolicy,
		ExchangePolicies: []*mtp_policy.ExchangePolicy{
			&serverABExchangePolicy,
		},
	}

	policy := BuildOutgoingPolicy(updateEvent)

	if len(policy.ManagedRoomIds) != 4 {
		t.Errorf(`Expected ManagedRoomIds to be length 4, but it is of length %d.`, len(policy.ManagedRoomIds))
	}

	if len(policy.User) != 3 {
		t.Errorf(`Expected User to be length 3, but it is of length %d.`, len(policy.User))
	}

	for _, expected := range []struct {
		id    string
		count int
	}{
		{"@a:a.test", 4},
		{"@b:a.test", 3},
		{"@c:a.test", 1},
	} {
		_, user := FindUser(&policy.User, expected.id)
		if user == nil {
			t.Errorf(`Expected policy to have user %s, but it dose not.`, expected.id)
			continue
		}

		if len(user.JoinedRooms) != expected.count {
			t.Errorf(`Expected user %s to be joined %d rooms, but they are joined %d rooms.`, expected.id, expected.count, len(user.JoinedRooms))
		}
	}

	for _, expected := range []struct {
		id    string
		room  string
		level int
	}{
		{"@a:a.test", "!1:a.test", 10},
		{"@a:a.test", "!2:a.test", 20},
		{"@a:a.test", "!3:a.test", 10},
		{"@a:a.test", "!4:a.test", 10},
		{"@b:a.test", "!2:a.test", 10},
		{"@b:a.test", "!3:a.test", 10},
		{"@c:a.test", "!3:a.test", 10},
	} {
		_, user := FindUser(&policy.User, expected.id)
		if user == nil {
			t.Errorf(`Expected policy to have user %s, but it dose not.`, expected.id)
			continue
		}

		room := FindRoom(user.JoinedRooms, expected.room)
		if room == nil {
			t.Errorf(`Expected user %s to be in room %s, but they are not.`, expected.id, expected.room)
			continue
		}

		if room.PowerLevel != expected.level {
			t.Errorf(`Expected user %s to have power leven %d in room %s, but they have power level %d.`, expected.id, expected.level, expected.room, room.PowerLevel)
		}
	}

	serverABasePolicyAfterJson, _ := json.Marshal(serverABasePolicy)
	if string(serverABasePolicyBeforeJson) != string(serverABasePolicyAfterJson) {
		t.Errorf(`Server a base policy changes. Before: %s After: %s`, serverABasePolicyBeforeJson, serverABasePolicyAfterJson)
	}
}

func findRemoteUser(remoteUsers []*mtp_policy.RemoteUserPolicy, id string) *mtp_policy.RemoteUserPolicy {
	for _, user := range remoteUsers {
		if user.Id == id {
			return user
		}
	}
	return nil
}

func TestUpdateIncomingPolicy(t *testing.T) {
	instanceRepository := instance.NewInMemoryInstanceRepository()

	instanceA := &instance.Instance{Store: store.NewInMemoryStore()}
	instanceB := &instance.Instance{Store: store.NewInMemoryStore()}
	instanceC := &instance.Instance{Store: store.NewInMemoryStore()}
	instanceRepository.AddInstance("a.test", instanceA)
	instanceRepository.AddInstance("b.test", instanceB)
	instanceRepository.AddInstance("c.test", instanceC)

	UpdateIncomingPolicy(instanceRepository, &serverAIncomingPolicy, "a.test")

	aBasePolicy := instanceA.Store.GetBasePolicy()
	if len(aBasePolicy.User) != 3 {
		t.Errorf(`Expected base policy of instance %s to have exactly 3 users.`, "a.test")
	}

	for _, expected := range []struct {
		instance string
		userId   string
		roomId   string
	}{
		{"b.test", "@a:b.test", "!1:a.test"},
		{"b.test", "@a:b.test", "!2:a.test"},
		{"b.test", "@b:b.test", "!2:a.test"},
		{"b.test", "@b:b.test", "!3:a.test"},
		{"c.test", "@a:c.test", "!1:a.test"},
		{"c.test", "@a:c.test", "!3:a.test"},
	} {

		bExchangePolicies := instanceRepository.GetInstance(expected.instance).Store.GetExchangePolicies()
		user := findRemoteUser(bExchangePolicies[0].User, expected.userId)
		if user == nil {
			t.Errorf(`Expected user %s to be in exchange policy for instance %s, but it is not.`, expected.userId, expected.instance)
			continue
		}

		if FindRoom(user.JoinedRooms, expected.roomId) == nil {
			t.Errorf(`Expected JoinedRooms for user %s to list room %s in exchange policy for instance %s`, expected.userId, expected.roomId, expected.instance)
		}
	}
}
