package connector

import (
	"github.com/bborbe/teamvault-utils/model"
)

type Connector interface {
	Password(key model.TeamvaultKey) (model.TeamvaultPassword, error)
	User(key model.TeamvaultKey) (model.TeamvaultUser, error)
	Url(key model.TeamvaultKey) (model.TeamvaultUrl, error)
	File(key model.TeamvaultKey) (model.TeamvaultFile, error)
}
