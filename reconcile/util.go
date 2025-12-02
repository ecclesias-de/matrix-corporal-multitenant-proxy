package reconcile

import "devture-matrix-corporal/corporal/policy"

func FindUser(users *[]*policy.UserPolicy, id string) (int, *policy.UserPolicy) {
	for i, user := range *users {
		if user.Id == id {
			return i, user
		}
	}
	return 0, nil
}

func FindRoom(memberships []*policy.RoomState, roomId string) *policy.RoomState {
	for _, membership := range memberships {
		if membership.RoomId == roomId {
			return membership
		}
	}

	return nil
}
