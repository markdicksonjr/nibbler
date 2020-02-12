package nibbler_user_group

import (
	"github.com/markdicksonjr/nibbler"
)

// SetGroupMembership upserts the group membership record for a given user and group
func (s *Extension) SetGroupMembership(groupId, userId string, role string) (nibbler.GroupMembership, error) {
	return s.PersistenceExtension.SetGroupMembership(groupId, userId, role)
}

func (s *Extension) GetGroupMembershipsForUser(userId string) ([]nibbler.GroupMembership, error) {
	return s.PersistenceExtension.GetGroupMembershipsForUser(userId)
}
