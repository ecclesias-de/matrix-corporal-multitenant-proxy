package test

import (
	"io"
	"matrix-corporal-multitenant-proxy/config"
	"matrix-corporal-multitenant-proxy/container"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func putPolicy(t *testing.T, addr, tenant, policy string) {
	httpClient := http.Client{}

	req, err := http.NewRequest("PUT", "http://"+addr+"/_matrix/corporal/policy/"+tenant, strings.NewReader(policy))

	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Authorization", "Bearer "+"well-known-test-token-proxy-"+tenant)
	req.Header.Set("Content-Type", "application/json")

	res, err := httpClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf(`Status code should be 200 Ok, but is %d`, res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	if string(body) != "{}" {
		t.Errorf(`Body should be {}, but is: %s`, string(body))
	}
}

func TestE2ESuccessSimple(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	aAddr, aToken, aPolicyPath, aClose := CreateMatrixCorporalInstance(t, "a.test")
	defer aClose()
	bAddr, bToken, bPolicyPath, bClose := CreateMatrixCorporalInstance(t, "b.test")
	defer bClose()

	config := config.Config{
		Instances: []config.Instance{
			{
				HomeServerName:      "a.test",
				ApiAccessToken:      "well-known-test-token-proxy-a.test",
				CorporalApiUrl:      "http://" + aAddr,
				CorporalAccessToken: aToken,
			},
			{
				HomeServerName:      "b.test",
				ApiAccessToken:      "well-known-test-token-proxy-b.test",
				CorporalApiUrl:      "http://" + bAddr,
				CorporalAccessToken: bToken,
			},
		},
		ReconcileRetryInterval: 30,
		ListenAddress:          GetListenerAddr(t),
	}
	services := container.Bootstrap(config)
	for _, service := range services {
		err := service.Start()
		if err != nil {
			t.Fatal(err)
		}
	}

	defer func() {
		for _, service := range services {
			err := service.Start()
			if err != nil {
				t.Error(err)
			}
		}
	}()

	putPolicy(t, config.ListenAddress, "a.test", `{
		"schemaVersion": 2,
		"managedRoomIds": [
			"!1:a.test"
		],
		"users": [
			{
				"id":"@a:a.test",
				"authType": "plain",
				"joinedRooms": [
					{"roomId":"!1:a.test", "powerLevel": 10}
				]
			}
		],
		"remoteUsers": [
			{
				"id": "@a:b.test",
				"joinedRooms":[
					{"roomId": "!1:a.test", "powerLevel": 20}
				]
			}
		]
	}`)
	time.Sleep(time.Microsecond * 2000)

	TestFileMachetes(t, aPolicyPath, `{"schemaVersion":2,"identificationStamp":null,"flags":{"allowCustomUserDisplayNames":false,"allowCustomUserAvatars":false,"allowCustomPassthroughUserPasswords":false,"allowUnauthenticatedPasswordResets":false,"forbidRoomCreation":false,"forbidEncryptedRoomCreation":false,"forbidUnencryptedRoomCreation":false,"allow3pidLogin":false},"hooks":null,"managedRoomIds":["!1:a.test"],"users":[{"id":"@a:a.test","active":false,"authType":"plain","authCredential":"","displayName":"","avatarUri":"","joinedRooms":[{"roomId":"!1:a.test","powerLevel":10}],"forbidRoomCreation":null,"forbidEncryptedRoomCreation":null,"forbidUnencryptedRoomCreation":null}]}`)
	TestFileMachetes(t, bPolicyPath, ``)

	putPolicy(t, config.ListenAddress, "b.test", `{
		"schemaVersion": 2,
		"managedRoomIds": [
			"!1:b.test"
		],
		"users": [
			{
				"id":"@a:b.test",
				"authType": "plain",
				"joinedRooms": [
					{"roomId":"!1:b.test", "powerLevel": 10}
				]
			}
		],
		"remoteUsers": [
			{
				"id": "@a:a.test",
				"joinedRooms":[
					{"roomId": "!1:b.test", "powerLevel": 20}
				]
			}
		]
	}`)
	time.Sleep(time.Microsecond * 2000)

	TestFileMachetes(t, aPolicyPath, `{"schemaVersion":2,"identificationStamp":null,"flags":{"allowCustomUserDisplayNames":false,"allowCustomUserAvatars":false,"allowCustomPassthroughUserPasswords":false,"allowUnauthenticatedPasswordResets":false,"forbidRoomCreation":false,"forbidEncryptedRoomCreation":false,"forbidUnencryptedRoomCreation":false,"allow3pidLogin":false},"hooks":null,"managedRoomIds":["!1:a.test","!1:b.test"],"users":[{"id":"@a:a.test","active":false,"authType":"plain","authCredential":"","displayName":"","avatarUri":"","joinedRooms":[{"roomId":"!1:a.test","powerLevel":10},{"roomId":"!1:b.test","powerLevel":20}],"forbidRoomCreation":null,"forbidEncryptedRoomCreation":null,"forbidUnencryptedRoomCreation":null}]}`)
	TestFileMachetes(t, bPolicyPath, `{"schemaVersion":2,"identificationStamp":null,"flags":{"allowCustomUserDisplayNames":false,"allowCustomUserAvatars":false,"allowCustomPassthroughUserPasswords":false,"allowUnauthenticatedPasswordResets":false,"forbidRoomCreation":false,"forbidEncryptedRoomCreation":false,"forbidUnencryptedRoomCreation":false,"allow3pidLogin":false},"hooks":null,"managedRoomIds":["!1:b.test","!1:a.test"],"users":[{"id":"@a:b.test","active":false,"authType":"plain","authCredential":"","displayName":"","avatarUri":"","joinedRooms":[{"roomId":"!1:b.test","powerLevel":10},{"roomId":"!1:a.test","powerLevel":20}],"forbidRoomCreation":null,"forbidEncryptedRoomCreation":null,"forbidUnencryptedRoomCreation":null}]}`)
}
