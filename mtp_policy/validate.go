package mtp_policy

import (
	"fmt"
)

func Validate(policy MultiTenantPolicy) error {
	if policy.User != nil {
		for i, user := range policy.User {
			if user.Id == "" {
				return fmt.Errorf(`User at index %d misses Id.`, i)
			}

			//todo make optional
			if user.JoinedRooms != nil {
				for i, room := range user.JoinedRooms {
					if room.RoomId == "" {
						return fmt.Errorf(`RoomMembership info at index %d on user %s misses id.`, i, user.Id)
					}
				}
			}
		}
	}

	//todo make optional?
	if policy.RemoteUsers != nil {
		for i, user := range policy.RemoteUsers {
			if user.Id == "" {
				return fmt.Errorf(`User at index %d misses Id.`, i)
			}

			if user.JoinedRooms == nil {
				return fmt.Errorf(`User %s misses JoinedRooms`, user.Id)
			}

			for i, room := range user.JoinedRooms {
				if room.RoomId == "" {
					return fmt.Errorf(`RoomMembership info at index %d on user %s misses id.`, i, user.Id)
				}
			}
		}
	} else {
		policy.RemoteUsers = make([]*RemoteUserPolicy, 0)
	}

	return nil
}
