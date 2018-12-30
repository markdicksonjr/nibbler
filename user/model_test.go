package user

import (
	"testing"
)

func TestToJson(t *testing.T) {
	user := User{
		ID: "654",
	}

	out, err := ToJson(&user)

	if err != nil {
		t.Fail()
	}

	if out != "{\"id\":\"654\",\"createdAt\":\"0001-01-01T00:00:00Z\",\"updatedAt\":\"0001-01-01T00:00:00Z\"}" {
		t.Fail()
	}
}

func TestFromJson(t *testing.T) {
	user, err := FromJson("{\"id\":\"123\"}")

	if err != nil {
		t.Fail()
	}

	if user == nil {
		t.Fail()
	}

	if user.ID != "123" {
		t.Fail()
	}
}
