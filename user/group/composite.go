package nibbler_user_group

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/markdicksonjr/nibbler"
	"net/http"
)

type UserComposite struct {
	CurrentGroup       *nibbler.Group      `json:"currentGroup"`
	Groups             []nibbler.Group     `json:"groups"`
	RoleInCurrentGroup string              `json:"roleInCurrentGroup"`
}

// LoadUserCompositeRequestHandler gets the composite for the given user - you can either allow them to specify the ID
// as a path param, or just always make it return the composite for the caller.  If you allow the ID to be specified,
// protect this route with a check to see if the caller can ask for that user's composite info.
func (s *Extension) LoadUserCompositeRequestHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// if no ID was provided, assume the route did not have ID in the path, and the user wants to ask about the caller
	if id == "" {
		caller, err := s.SessionExtension.GetCaller(r)
		if err != nil {
			nibbler.Write500Json(w, err.Error())
			return
		}

		id = caller.ID
	}

	composite, err := s.LoadUserComposite(id)

	// fail on any error
	if err != nil {
		nibbler.Write500Json(w, err.Error())
		return
	}

	compositeJson, err := json.Marshal(composite)

	// fail on any error
	if err != nil {
		nibbler.Write500Json(w, err.Error())
		return
	}

	nibbler.Write200Json(w, string(compositeJson))
}

func (s *Extension) LoadUserComposite(userId string) (*UserComposite, error) {

	// prepare a composite instance to return
	composite := UserComposite{
		Groups:             make([]nibbler.Group, 0),
		CurrentGroup:       nil,
	}

	// load the user configuration
	u, err := s.UserExtension.GetUserById(userId)

	// fail on any error
	if err != nil {
		return nil, err
	}

	// check that we got a user back
	if u == nil {
		return nil, nil
	}

	// load the group memberships for the user
	memberships, err := s.GetGroupMembershipsForUser(userId)
	if err != nil {
		return nil, err
	}

	// load the groups for the set of memberships
	if len(memberships) > 0 {

		// get a list of group IDs from the group memberships
		var groupIds []string
		for _, m := range memberships {
			groupIds = append(groupIds, m.GroupID)
		}
		groups, err := s.GetGroups(groupIds)

		if err != nil {
			return nil, err
		}

		composite.Groups = groups

		// set the current group in the composite model we're returning
		if u.CurrentGroupID != nil {
			g, err := s.GetGroups([]string{ *u.CurrentGroupID })
			if err != nil {
				return nil, err
			}
			if len(g) > 0 { // TODO: log when empty?  bogus/deleted group ID is current.  Maybe wipe the current group from user?
				composite.CurrentGroup = &g[0]
			}
		}
	}

	// if there's no current group, we're done - return
	if composite.CurrentGroup == nil {
		return &composite, nil
	}

	// load the role for the current group
	var currentGroupMembership *nibbler.GroupMembership
	for _, c := range memberships {
		if c.GroupID == composite.CurrentGroup.ID {
			m := c
			currentGroupMembership = &m
		}
	}

	if currentGroupMembership == nil {
		return &composite, nil
	}

	composite.RoleInCurrentGroup = (currentGroupMembership).Role
	return &composite, nil
}
