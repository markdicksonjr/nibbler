package nibbler_user_group

import (
	"github.com/markdicksonjr/nibbler"
	"net/http"
)

// EnforceHasPrivilege will use HasPrivilege to produce a result for the caller - it will return a 500 if something
// went wrong, a 401 if no user is authenticated, a 404 if there is no access.  It will pass through to the routerFunc
// if the caller has access
func (s *Extension) EnforceHasPrivilege(action string, routerFunc func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		caller, err := s.SessionExtension.GetCaller(r)
		if err != nil {
			nibbler.Write500Json(w, err.Error())
			return
		}

		if caller == nil {
			nibbler.Write401Json(w)
			return
		}

		if has, err := s.HasPrivilege(caller.ID, action); err != nil {
			nibbler.Write500Json(w, err.Error())
		} else if !has {
			nibbler.Write404Json(w)
		} else {
			routerFunc(w, r)
		}
	}
}

// EnforceHasPrivilegeOnResource will use HasPrivilegeOnResource to produce a result for the caller - it will return
// a 500 if something went wrong, a 401 if no user is authenticated, a 404 if there is no access.  It will pass through
// to the routerFunc if the caller has access
func (s *Extension) EnforceHasPrivilegeOnResource(action string, getResourceIdFn func(r *http.Request) (string, error), routerFunc func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		caller, err := s.SessionExtension.GetCaller(r)
		if err != nil {
			nibbler.Write500Json(w, err.Error())
			return
		}

		if caller == nil {
			nibbler.Write401Json(w)
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

		// if the user does not have the privilege on the target resource, fall back to the "global" version of that privilege
		if has, err := s.HasPrivilegeOnResource(caller.ID, targetGroup, action); err != nil {
			nibbler.Write500Json(w, err.Error())
		} else if !has {

			// this is the check against the global privilege for this action (e.g. admins, etc)
			if has, err := s.HasPrivilege(caller.ID, action); err != nil {
				nibbler.Write500Json(w, err.Error())
			} else if !has {
				nibbler.Write404Json(w)
			} else {
				routerFunc(w, r)
			}
		} else {
			routerFunc(w, r)
		}
	}
}
