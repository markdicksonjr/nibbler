package nibbler_user_group

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/markdicksonjr/nibbler"
	"io/ioutil"
	"net/http"
)

// some (hopefully) reusable or useful privileges
const CreateGroupMembershipAction = "add-member-to-group"
const CreateGroupAction = "create-group"
const DeleteGroupAction = "delete-group"
const CreatePrivilegeAction = "create-privilege"
const DeletePrivilegeAction = "delete-privilege"
const CreateGroupPrivilegeAction = "create-group-privilege"
const DeleteGroupPrivilegeAction = "delete-group-privilege"
const ListGroupsAction = "list-groups"
const RemoveMemberFromGroupAction = "remove-member-from-group"

// AddPrivilegeToGroups adds a specific privilege definition to save to multiple groups.  It allows all groups in the
// groupIdList to perform the provided action on the targetGroupId.  If targetGroupId is blank, it means
// "all resources/groups"
func (s *Extension) AddPrivilegeToGroups(
	groupIdList []string,
	targetGroupId string,
	action string,
) error {
	return s.PersistenceExtension.AddPrivilegeToGroups(groupIdList, targetGroupId, action)
}

// HasPrivilege returns whether the caller has a privilege for a resource-agnostic action.  This is suitable
// for something like "create-admin" or any other "global"-type privilege
func (s *Extension) HasPrivilege(userId, action string) (bool, error) {
	userFromDb, err := s.UserExtension.GetUserById(userId)
	if err != nil {
		return false, err
	}

	if userFromDb == nil || userFromDb.CurrentGroupID == nil {
		return false, nil
	}

	privileges, err := s.PersistenceExtension.GetPrivilegesForAction(*userFromDb.CurrentGroupID, nil, action)
	if err != nil {
		return false, err
	}

	return len(privileges) > 0, nil
}

// HasPrivilegeOnResource will state whether the caller can perform an action on a specific resource.  If there is no
// resource-specific privilege, it will check to see if the caller has the global privilege for that action.  For
// example, some users may have "create-user" privileges for a specific group, but an admin may have a resource-agnostic
// "create-user" privilege.  This function will check both.
func (s *Extension) HasPrivilegeOnResource(userId, resourceId, action string) (bool, error) {
	userFromDb, err := s.UserExtension.GetUserById(userId)
	if err != nil {
		return false, err
	}

	if userFromDb == nil || userFromDb.CurrentGroupID == nil {
		return false, nil
	}

	privileges, err := s.PersistenceExtension.GetPrivilegesForAction(*userFromDb.CurrentGroupID, &resourceId, action)
	if err != nil {
		return false, err
	}

	return len(privileges) > 0, nil
}

// DeleteGroupPrivilegeRequestHandler handles an http request with a privilege in its body and "groupId" in the path params
func (s *Extension) DeleteGroupPrivilegeRequestHandler(w http.ResponseWriter, r *http.Request) {
	priv, err := getPrivilegeFromBody(r)
	if err != nil {
		nibbler.Write500Json(w, err.Error())
		return
	}

	priv.GroupID = mux.Vars(r)["groupId"]

	// if a user passes no resource for the privilege, this is a little dangerous - we need to be sure they are allowed
	// to allocate privileges that are independent of a resource (global privileges)
	if priv.ResourceID == "" {
		caller, err := s.SessionExtension.GetCaller(r)
		if err != nil {
			nibbler.Write500Json(w, err.Error())
			return
		}

		if caller == nil {
			nibbler.Write401Json(w)
			return
		}

		// if the user does not have the right to create such privileges, stop them here
		if has, err := s.HasPrivilege(caller.ID, DeletePrivilegeAction); err != nil {
			nibbler.Write500Json(w, err.Error())
			return
		} else if !has {
			nibbler.Write404Json(w)
			return
		}
	}

	// get group/resource/action from request and delete match(es)
	privileges, err := s.PersistenceExtension.GetPrivilegesForAction(priv.GroupID, &priv.ResourceID, priv.Action)
	if err != nil {
		nibbler.Write500Json(w, err.Error())
		return
	}

	if len(privileges) == 0 {
		nibbler.Write404Json(w)
		return
	}

	for _, p := range privileges {

		// TODO: hard del query param
		if err := s.PersistenceExtension.DeletePrivilege(p.ID, false); err != nil {
			nibbler.Write500Json(w, err.Error())
			return
		}
	}

	nibbler.Write200Json(w, "{\"result\":\"ok\"")
}

// CreateGroupPrivilegeRequestHandler handles an http request with a path param of groupId and body that is a Privilege
func (s *Extension) CreateGroupPrivilegeRequestHandler(w http.ResponseWriter, r *http.Request) {
	priv, err := getPrivilegeFromBody(r)
	if err != nil {
		nibbler.Write500Json(w, err.Error())
		return
	}

	// override the group ID in the request body/contents with the one from the path (TODO: perhaps warn if different - could identify hack attempts)
	priv.GroupID = mux.Vars(r)["groupId"]

	// if a user passes no resource for the privilege, this is a little dangerous - we need to be sure they are allowed
	// to allocate privileges that are independent of a resource (global privileges)
	if priv.ResourceID == "" {
		caller, err := s.SessionExtension.GetCaller(r)
		if err != nil {
			nibbler.Write500Json(w, err.Error())
			return
		}

		if caller == nil {
			nibbler.Write401Json(w)
			return
		}

		// if the user does not have the right to create such privileges, stop them here
		if has, err := s.HasPrivilege(caller.ID, CreatePrivilegeAction); err != nil {
			nibbler.Write500Json(w, err.Error())
			return
		} else if !has {
			nibbler.Write404Json(w)
			return
		}
	}

	if err := s.PersistenceExtension.AddPrivilegeToGroups([]string{priv.GroupID}, priv.ResourceID, priv.Action); err != nil {
		nibbler.Write500Json(w, err.Error())
		return
	}

	nibbler.Write200Json(w, "{\"result\":\"ok\"")
}

// getPrivilegeFromBody parses the request body into a GroupPrivilege struct
func getPrivilegeFromBody(r *http.Request) (*nibbler.GroupPrivilege, error) {
	if r.Body == nil {
		return nil, errors.New("no body provided")
	}

	defer r.Body.Close()
	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	priv := nibbler.GroupPrivilege{}
	if err := json.Unmarshal(raw, &priv); err != nil {
		return nil, err
	}

	return &priv, nil
}
