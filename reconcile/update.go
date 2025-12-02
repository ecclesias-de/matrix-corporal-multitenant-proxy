package reconcile

import (
	"devture-matrix-corporal/corporal/policy"
	"matrix-corporal-multitenant-proxy/instance"
	"matrix-corporal-multitenant-proxy/mtp_policy"
	"matrix-corporal-multitenant-proxy/store"
	"slices"
	"strings"

	"github.com/sirupsen/logrus"
)

func appendJoinedRooms(managedRooms []string, joinedRooms []*policy.RoomState) []string {
	for _, roomState := range joinedRooms {
		if slices.Contains(managedRooms, roomState.RoomId) {
			continue
		}

		managedRooms = append(managedRooms, roomState.RoomId)
	}

	return managedRooms
}

func BuildOutgoingPolicy(updateEvent *store.UpdateEvent) *policy.Policy {
	basePolicy := updateEvent.BasePolicies
	if basePolicy == nil {
		logrus.Warnf(`Updated OutgoingPolicy, but base policy is not set. Skipping tenant`)
		return nil
	}
	outgoingPolicy := *basePolicy

	outgoingPolicy.User = make([]*policy.UserPolicy, len(basePolicy.User))
	copy(outgoingPolicy.User, basePolicy.User)
	for _, val := range updateEvent.ExchangePolicies {
		for _, remoteUser := range val.User {
			i, user_pt := FindUser(&outgoingPolicy.User, remoteUser.Id)
			if user_pt == nil {
				// print warning / info
				continue
			}

			user := *user_pt
			//todo fix
			user.JoinedRooms = append(user.JoinedRooms, remoteUser.JoinedRooms...)
			outgoingPolicy.User[i] = &user

			outgoingPolicy.ManagedRoomIds = appendJoinedRooms(outgoingPolicy.ManagedRoomIds, remoteUser.JoinedRooms)
		}
	}

	return &outgoingPolicy
}

func UpdateIncomingPolicy(instances instance.InstanceRepository, incomingPolicy *mtp_policy.MultiTenantPolicy, tenant string) {
	logrus.Debugf(`Update Incoming Policy for %s`, tenant)

	tenantInstance := instances.GetInstance(tenant)
	if tenantInstance == nil {
		logrus.Warnf(`Unknown tenant %s. Skipping it.`, tenant)
		return
	}

	exchangePolicies := make(map[string]*mtp_policy.ExchangePolicy, 0)
	for _, user := range incomingPolicy.RemoteUsers {
		idParts := strings.Split(user.Id, ":")
		if len(idParts) != 2 {
			logrus.Infof(`Invalid matrix id. Ids should contain exactly one :, but splitting on : did not produce 2 elements. Id: %s\n`, user.Id)
			continue
		}

		if exchangePolicies[idParts[1]] == nil {
			exchangePolicies[idParts[1]] = &mtp_policy.ExchangePolicy{}
		}

		exchangePolicies[idParts[1]].User = append(exchangePolicies[idParts[1]].User, &mtp_policy.RemoteUserPolicy{
			Id:          user.Id,
			JoinedRooms: user.JoinedRooms,
		})
	}

	tenantInstance.Store.SetBasePolicy(&incomingPolicy.Policy)

	for exchangeTenant, exchangePolicy := range exchangePolicies {
		instance := instances.GetInstance(exchangeTenant)
		if instance == nil {
			logrus.Infof(`Unknown exchange tenant %s. Skipping it.`, exchangeTenant)
		}
		instance.Store.SetExchangePolicy(exchangePolicy, tenant)
	}
}
