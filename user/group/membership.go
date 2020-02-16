package nibbler_user_group

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/markdicksonjr/nibbler"
	"io/ioutil"
	"net/http"
)

// SetGroupMembership upserts the group membership record for a given user and group
func (s *Extension) SetGroupMembership(groupId, userId string, role string) (nibbler.GroupMembership, error) {
	return s.PersistenceExtension.SetGroupMembership(groupId, userId, role)
}

// GetGroupMembershipsForUser lists the groups to which the user (with the provided ID) belongs
func (s *Extension) GetGroupMembershipsForUser(userId string) ([]nibbler.GroupMembership, error) {
	return s.PersistenceExtension.GetGroupMembershipsForUser(userId)
}

func (s *Extension) CreateGroupMembershipRequestHandler(w http.ResponseWriter, r *http.Request) {
	membership, err := getMembershipFromBody(r)
	if err != nil {
		nibbler.Write500Json(w, err.Error())
		return
	}

	if membership.MemberID == "" {
		nibbler.Write500Json(w, "{\"error\":\"no member ID provided\"")
		return
	}

	membership.GroupID = mux.Vars(r)["groupId"]

	// role is optional

	result, err := s.PersistenceExtension.SetGroupMembership(membership.GroupID, membership.MemberID, membership.Role)
	if err != nil {
		nibbler.Write500Json(w, err.Error())
		return
	}

	resultJson, err := json.Marshal(result)
	if err != nil {
		nibbler.Write500Json(w, err.Error())
		return
	}

	nibbler.Write200Json(w, string(resultJson))
}

func getMembershipFromBody(r *http.Request) (*nibbler.GroupMembership, error) {
	if r.Body == nil {
		return nil, errors.New("no body provided")
	}

	defer r.Body.Close()
	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	membership := nibbler.GroupMembership{}
	if err := json.Unmarshal(raw, &membership); err != nil {
		return nil, err
	}

	return &membership, nil
}
