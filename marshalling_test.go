package nibbler

import (
	"reflect"
	"testing"
	"time"
)

type Score struct {
	ID						string		`json:"id" bson:"_id"`
	CreatedAt				time.Time	`json:"createdAt"`
	UpdatedAt				time.Time	`json:"updatedAt"`
	DeletedAt				*time.Time	`json:"deletedAt,omitempty" sql:"index"`
	Current					*int64		`json:"current"`
}

func TestToJson(t *testing.T) {
	var current int64 = 7
	score := Score{
		ID: "kj",
		Current: &current,
	}

	out, err := ToJson(score)

	if err != nil {
		t.Fail()
	}

	if out != "{\"id\":\"kj\",\"createdAt\":\"0001-01-01T00:00:00Z\",\"updatedAt\":\"0001-01-01T00:00:00Z\",\"current\":7}" {
		t.Fail()
	}
}

func TestFromJson(t *testing.T) {
	str, err := FromJson("{\"id\":\"kj\",\"createdAt\":\"2018-01-01T00:00:00Z\",\"updatedAt\":\"2018-02-01T00:00:00Z\",\"current\":7}", reflect.TypeOf(Score{}))

	if err != nil {
		t.Fail()
	}

	if str == nil {
		t.Fail()
	}

	score := str.(*Score)
	if score.ID != "kj" {
		t.Fail()
	}
}