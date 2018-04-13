package connector

import (
	"testing"

	. "github.com/bborbe/assert"
	"github.com/bborbe/teamvault_utils/model"
)

func TestDummyUser(t *testing.T) {
	key := model.TeamvaultKey("key123")
	du := NewDummy()
	user, err := du.User(key)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(user, Is(model.TeamvaultUser("key123"))); err != nil {
		t.Fatal(err)
	}
}

func TestDummyPassword(t *testing.T) {
	key := model.TeamvaultKey("key123")
	du := NewDummy()
	password, err := du.Password(key)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(password, Is(model.TeamvaultPassword("LgIWz7BC2r68P9WTtVJdfFOYrpT2tv_yw95BzhzECiU="))); err != nil {
		t.Fatal(err)
	}
}

func TestDummyURL(t *testing.T) {
	key := model.TeamvaultKey("key123")
	du := NewDummy()
	url, err := du.Url(key)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(url, Is(model.TeamvaultUrl("dk9kTUjDqGcvPlvF0ZOovq3sBE-0_-Y62i8mlTX_g1M="))); err != nil {
		t.Fatal(err)
	}
}
