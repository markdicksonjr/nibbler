package nibbler_user_group

import (
	"github.com/gorilla/mux"
	"github.com/markdicksonjr/nibbler"
	"github.com/markdicksonjr/nibbler/session"
	"github.com/markdicksonjr/nibbler/user"
	"net/http"
)

type PersistenceExtension interface {
	StartTransaction() (PersistenceExtension, error)
	RollbackTransaction() error
	CommitTransaction() error
	GetGroupMembershipsForUser(id string) ([]nibbler.GroupMembership, error)
	SetGroupMembership(groupId string, userId string, role string) (nibbler.GroupMembership, error)
	CreateGroup(group nibbler.Group) error
	DeleteGroup(groupId string, hardDelete bool) error
	SearchGroups(query nibbler.SearchParameters, includePrivileges bool) (nibbler.SearchResults, error)
	GetGroupsById(ids []string, includePrivileges bool) ([]nibbler.Group, error)
	AddPrivilegeToGroups(groupIdList []string, resourceId string, action string) error
	GetPrivilegesForAction(groupId string, resourceId *string, action string) ([]nibbler.GroupPrivilege, error)
	DeletePrivilege(id string) error
}

type Extension struct {
	nibbler.NoOpExtension
	PersistenceExtension PersistenceExtension
	SessionExtension     *session.Extension
	UserExtension        *user.Extension
	DisableDefaultRoutes bool
}

func (s *Extension) GetName() string {
	return "user-group"
}

func GetParamValueFromRequest(paramName string) func(r *http.Request) (s string, err error) {
	return func(r *http.Request) (s string, err error) {
		return mux.Vars(r)[paramName], nil
	}
}

func (s *Extension) PostInit(app *nibbler.Application) error {
	if !s.DisableDefaultRoutes {
		app.Router.HandleFunc(app.Config.ApiPrefix+"/group/composite", s.SessionExtension.EnforceLoggedIn(s.GetUserCompositeRequestHandler)).Methods("GET")
		app.Router.HandleFunc(app.Config.ApiPrefix+"/group", s.EnforceHasPrivilege(ListGroupsAction, s.QueryGroupsRequestHandler)).Methods("GET")
		app.Router.HandleFunc(app.Config.ApiPrefix+"/group", s.EnforceHasPrivilege(CreateGroupAction, s.CreateGroupRequestHandler)).Methods("PUT")
		app.Router.HandleFunc(app.Config.ApiPrefix+"/group/:groupId/privilege", s.EnforceHasPrivilegeOnResource(DeleteGroupPrivilegeAction, GetParamValueFromRequest("groupId"), s.DeleteGroupPrivilegeRequestHandler)).Methods("DELETE")
		app.Router.HandleFunc(app.Config.ApiPrefix+"/group/:groupId/privilege", s.EnforceHasPrivilegeOnResource(CreateGroupPrivilegeAction, GetParamValueFromRequest("groupId"), s.CreateGroupPrivilegeRequestHandler)).Methods("PUT")
		app.Router.HandleFunc(app.Config.ApiPrefix+"/group/:groupId", s.EnforceHasPrivilegeOnResource(DeleteGroupAction, GetParamValueFromRequest("groupId"), s.DeleteGroupRequestHandler)).Methods("DELETE")
		app.Router.HandleFunc(app.Config.ApiPrefix+"/group/:groupId/membership", s.EnforceHasPrivilegeOnResource(CreateGroupMembershipAction, GetParamValueFromRequest("groupId"), s.CreateGroupMembershipRequestHandler)).Methods("PUT")
		app.Router.HandleFunc(app.Config.ApiPrefix+"/group/:groupId/membership", s.EnforceHasPrivilegeOnResource(RemoveMemberFromGroupAction, GetParamValueFromRequest("groupId"), s.CreateGroupMembershipRequestHandler)).Methods("DELETE")
	}
	return nil
}

func GetModels() []interface{} {
	var models []interface{}
	models = append(models, nibbler.Group{})
	models = append(models, nibbler.GroupPrivilege{})
	models = append(models, nibbler.User{})
	models = append(models, nibbler.GroupMembership{})

	return models
}
