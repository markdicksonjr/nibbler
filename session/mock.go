package session

import (
	"github.com/gorilla/sessions"
	"net/http"
)

type MockStoreConnector struct {
	Store sessions.Store
}

func (m *MockStoreConnector) Connect() (error, sessions.Store) {
	return nil, m.Store
}
func (m *MockStoreConnector) MaxAge() int {
	return 360000
}

type MockStore struct {
	Session *sessions.Session
}

func (m *MockStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return m.Session, nil
}

func (m *MockStore) New(r *http.Request, name string) (*sessions.Session, error) {
	return m.Session, nil
}

func (m *MockStore) Save(r *http.Request, w http.ResponseWriter, s *sessions.Session) error {
	return nil
}

type MockStoreFailedGet struct {
	MockStore
	ErrOnGet error
}

func (m *MockStoreFailedGet) Get(r *http.Request, name string) (*sessions.Session, error) {
	return m.Session, m.ErrOnGet
}
