package connector

import (
	"crypto/sha256"
	"encoding/base64"

	"github.com/bborbe/teamvault-utils/model"
)

type dummyPasswordProvider struct {
}

func NewDummy() *dummyPasswordProvider {
	t := new(dummyPasswordProvider)
	return t
}

func (t *dummyPasswordProvider) Password(key model.TeamvaultKey) (model.TeamvaultPassword, error) {
	h := sha256.New()
	h.Write([]byte(key + "-password"))
	result := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return model.TeamvaultPassword(result), nil
}

func (t *dummyPasswordProvider) User(key model.TeamvaultKey) (model.TeamvaultUser, error) {
	return model.TeamvaultUser(key.String()), nil
}

func (t *dummyPasswordProvider) Url(key model.TeamvaultKey) (model.TeamvaultUrl, error) {
	h := sha256.New()
	h.Write([]byte(key + "-url"))
	result := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return model.TeamvaultUrl(result), nil
}

func (t *dummyPasswordProvider) File(key model.TeamvaultKey) (model.TeamvaultFile, error) {
	result := base64.URLEncoding.EncodeToString([]byte(key + "-file"))
	return model.TeamvaultFile(result), nil
}
