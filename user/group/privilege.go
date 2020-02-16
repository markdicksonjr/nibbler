package nibbler_user_group

// some hopefully reusable or useful privileges - at some point, the group extension
// will likely provide a default implementation of privileges around group editing
const AddMemberToGroupPrivilege = "add-member-to-group"
const CreateGroupPrivilege = "create-group"
const DeleteGroupPrivilege = "delete-group"

// allows all groups in the groupIdList to perform the provided action
// on the targetGroupId.  If targetGroupId is blank, it means "all groups"
func (s *Extension) AddPrivilegeToGroups(
	groupIdList []string,
	targetGroupId string,
	action string,
) error {
	return s.PersistenceExtension.AddPrivilegeToGroups(groupIdList, targetGroupId, action)
}
