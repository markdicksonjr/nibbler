package nibbler_user_group

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/markdicksonjr/nibbler"
	"net/http"
)

func (s *Extension) CreateGroupRequestHandler(w http.ResponseWriter, r *http.Request) {

	// grab and validate group name
	groupName := r.FormValue("name")
	if groupName == "" {
		nibbler.Write500Json(w, "group name is a required field")
		return
	}

	// get the current user from the session
	caller, err := s.SessionExtension.GetCaller(r)
	if err != nil {
		nibbler.Write500Json(w, err.Error())
		return
	}

	ext, err := s.PersistenceExtension.StartTransaction()
	if err != nil {
		nibbler.Write500Json(w, err.Error())
		return
	}

	group := nibbler.Group{
		ID:   uuid.New().String(),
		Name: groupName,
	}

	// create the group
	err = ext.CreateGroup(group)
	if err != nil {
		ext.RollbackTransaction()
		nibbler.Write500Json(w, err.Error())
		return
	}

	// add the user as a manager of the new group - this call will also change the current group for the current user
	// to the new group, as well as create a score or any other model that's made once a user is added to a group
	_, err = ext.SetGroupMembership(group.ID, caller.ID, "admin")
	if err != nil {
		ext.RollbackTransaction()
		nibbler.Write500Json(w, err.Error())
		return
	}

	// commit the transaction
	err = ext.CommitTransaction()
	if err != nil {
		nibbler.Write500Json(w, err.Error())
		return
	}

	// stringify the group in order to return it
	groupJson, err := json.Marshal(&group)
	if err != nil {
		nibbler.Write500Json(w, err.Error())
		return
	}

	nibbler.Write200Json(w, string(groupJson))
}

// QueryGroupsRequestHandler lists groups - it does not yet support queries
func (s *Extension) QueryGroupsRequestHandler(w http.ResponseWriter, r *http.Request) {
	params := nibbler.SearchParameters{
		Query:        nil,
		Offset:       nil,
		Size:         nil,
		IncludeTotal: true,
	}

	// TODO: looks wrong - this should be a query param here
	includePrivs := mux.Vars(r)["includePrivs"] == "true"

	if g, err := s.PersistenceExtension.SearchGroups(params, includePrivs); err != nil {
		nibbler.Write500Json(w, err.Error())
		return
	} else {
		r, err := json.Marshal(g)
		if err != nil {
			nibbler.Write500Json(w, err.Error())
			return
		}
		nibbler.Write200Json(w, string(r))
	}
}

func (s *Extension) CreateGroup(name string) (nibbler.Group, error) {
	group := nibbler.Group{
		ID:   uuid.New().String(),
		Name: name,
	}
	err := s.PersistenceExtension.CreateGroup(group)
	return group, err
}

func (s *Extension) DeleteGroupRequestHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: allow query param for hard delete
	if err := s.PersistenceExtension.DeleteGroup(mux.Vars(r)["groupId"], false); err != nil {
		nibbler.Write500Json(w, err.Error())
		return
	}

	nibbler.Write200Json(w, "{\"result\":\"ok\"")
}

func (s *Extension) GetGroups(groupIds []string, includePrivileges bool) ([]nibbler.Group, error) {

	// load the groups for the set of memberships
	if len(groupIds) == 0 {
		return []nibbler.Group{}, nil
	}

	return s.PersistenceExtension.GetGroupsById(groupIds, includePrivileges)
}
