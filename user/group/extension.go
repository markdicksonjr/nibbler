package nibbler_user_group

import (
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
	GetGroupsById(ids []string) ([]nibbler.Group, error)
	AddPrivilegeToGroups(groupIdList []string, targetGroupId string, action string) error
	GetPrivilegesForAction(groupId string, resourceId *string, action string) ([]nibbler.GroupPrivilege, error)
}

type Extension struct {
	nibbler.NoOpExtension
	PersistenceExtension PersistenceExtension
	SessionExtension     *session.Extension
	UserExtension        *user.Extension
}

func (s *Extension) GetName() string {
	return "user-group"
}

func (s *Extension) EnforceHasPrivilege(action string, routerFunc func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		caller, err := s.SessionExtension.GetCaller(r)
		if err != nil {
			nibbler.Write500Json(w, err.Error())
			return
		}

		if caller == nil {
			nibbler.Write404Json(w)
			return
		}

		userFromDb, err := s.UserExtension.GetUserById(caller.ID)
		if err != nil {
			nibbler.Write500Json(w, err.Error())
			return
		}

		if userFromDb == nil || userFromDb.CurrentGroupID == nil {
			nibbler.Write404Json(w)
			return
		}

		privileges, err := s.PersistenceExtension.GetPrivilegesForAction(*userFromDb.CurrentGroupID, nil, action)
		if err != nil {
			nibbler.Write500Json(w, err.Error())
			return
		}
		if len(privileges) == 0 {
			nibbler.Write404Json(w)
			return
		}

		routerFunc(w, r)
	}
}

func (s *Extension) EnforceHasPrivilegeOnResource(action string, getResourceIdFn func(r *http.Request) (string, error), routerFunc func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		caller, err := s.SessionExtension.GetCaller(r)
		if err != nil {
			nibbler.Write500Json(w, err.Error())
			return
		}

		if caller == nil {
			nibbler.Write404Json(w)
			return
		}

		userFromDb, err := s.UserExtension.GetUserById(caller.ID)
		if err != nil {
			nibbler.Write500Json(w, err.Error())
			return
		}

		if userFromDb == nil || userFromDb.CurrentGroupID == nil {
			nibbler.Write404Json(w)
			return
		}

		targetGroup, err := getResourceIdFn(r)
		if err != nil {
			nibbler.Write500Json(w, err.Error())
			return
		}

		if targetGroup == "" {
			nibbler.Write404Json(w)
			return
		}

		privileges, err := s.PersistenceExtension.GetPrivilegesForAction(*userFromDb.CurrentGroupID, &targetGroup, action)
		if err != nil {
			nibbler.Write500Json(w, err.Error())
			return
		}
		if len(privileges) == 0 {
			nibbler.Write404Json(w)
			return
		}

		routerFunc(w, r)
	}
}

func (s *Extension) PostInit(context *nibbler.Application) error {
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
