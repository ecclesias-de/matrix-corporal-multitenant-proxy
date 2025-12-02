package store

import (
	"devture-matrix-corporal/corporal/policy"
	"encoding/json"
	"matrix-corporal-multitenant-proxy/mtp_policy"
	"matrix-corporal-multitenant-proxy/test"
	"testing"
	"time"
)

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

func TestPersistentStorage(t *testing.T) {
	basePath := t.TempDir()
	store := BootstrapPersistentStore(basePath, "a.test")

	store.SetBasePolicy(&serverABasePolicy)
	store.SetExchangePolicy(&serverABExchangePolicy, "b.test")

	time.Sleep(time.Microsecond * 2000)

	test.TestFileMachetes(t, basePath+"/a.test.json", `{"basePolicy":{"schemaVersion":0,"identificationStamp":null,"flags":{"allowCustomUserDisplayNames":false,"allowCustomUserAvatars":false,"allowCustomPassthroughUserPasswords":false,"allowUnauthenticatedPasswordResets":false,"forbidRoomCreation":false,"forbidEncryptedRoomCreation":false,"forbidUnencryptedRoomCreation":false,"allow3pidLogin":false},"hooks":null,"managedRoomIds":["!1:a.test","!2:a.test","!3:a.test"],"users":[{"id":"@a:a.test","active":true,"authType":"","authCredential":"","displayName":"a","avatarUri":"","joinedRooms":[{"roomId":"!1:a.test","powerLevel":10},{"roomId":"!2:a.test","powerLevel":20}],"forbidRoomCreation":null,"forbidEncryptedRoomCreation":null,"forbidUnencryptedRoomCreation":null},{"id":"@b:a.test","active":true,"authType":"","authCredential":"","displayName":"b","avatarUri":"","joinedRooms":[{"roomId":"!2:a.test","powerLevel":10}],"forbidRoomCreation":null,"forbidEncryptedRoomCreation":null,"forbidUnencryptedRoomCreation":null},{"id":"@c:a.test","active":true,"authType":"","authCredential":"","displayName":"b","avatarUri":"","joinedRooms":[{"roomId":"!3:a.test","powerLevel":10}],"forbidRoomCreation":null,"forbidEncryptedRoomCreation":null,"forbidUnencryptedRoomCreation":null}]},"exchangePolicies":{"b.test":{"remoteUsers":[{"id":"@a:a.test","joinedRooms":[{"roomId":"!3:a.test","powerLevel":10},{"roomId":"!4:a.test","powerLevel":10}]},{"id":"@b:a.test","joinedRooms":[{"roomId":"!2:a.test","powerLevel":30},{"roomId":"!3:a.test","powerLevel":10}]},{"id":"@d:a.test","joinedRooms":[{"roomId":"!2:a.test","powerLevel":10}]}]}}}`)

	storeLoaded := BootstrapPersistentStore(basePath, "a.test")
	basePolicy, err := json.Marshal(store.GetBasePolicy())
	if err != nil {
		t.Error(err)
	} else {
		basePolicyLoaded, err := json.Marshal(storeLoaded.GetBasePolicy())
		if err != nil {
			t.Error(err)
		} else {
			if string(basePolicy) != string(basePolicyLoaded) {
				t.Errorf(`Loaded basePolicy dose not match basePolicy: got %s expected %s`, string(basePolicyLoaded), string(basePolicy))
			}
		}
	}

	exchangePolicies, err := json.Marshal(store.GetExchangePolicies())
	if err != nil {
		t.Error(err)
	} else {
		exchangePoliciesLoaded, err := json.Marshal(storeLoaded.GetExchangePolicies())
		if err != nil {
			t.Error(err)
		} else {
			if string(exchangePolicies) != string(exchangePoliciesLoaded) {
				t.Errorf(`Loaded basePolicy dose not match basePolicy: got %s expected %s`, string(exchangePolicies), string(exchangePoliciesLoaded))
			}
		}
	}
}
