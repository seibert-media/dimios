package connector

import (
	"fmt"
	"net/http"

	"net/url"

	http_header "github.com/bborbe/http/header"
	"github.com/bborbe/http/rest"
	"github.com/bborbe/teamvault-utils/model"
)

type teamvaultPasswordProvider struct {
	url  model.TeamvaultUrl
	user model.TeamvaultUser
	pass model.TeamvaultPassword
	rest rest.Rest
}

func New(
	executeRequest func(req *http.Request) (resp *http.Response, err error),
	url model.TeamvaultUrl,
	user model.TeamvaultUser,
	pass model.TeamvaultPassword,
) *teamvaultPasswordProvider {
	t := new(teamvaultPasswordProvider)
	t.rest = rest.New(executeRequest)
	t.url = url
	t.user = user
	t.pass = pass
	return t
}

func (t *teamvaultPasswordProvider) Password(key model.TeamvaultKey) (model.TeamvaultPassword, error) {
	currentRevision, err := t.CurrentRevision(key)
	if err != nil {
		return "", err
	}
	var response struct {
		Password model.TeamvaultPassword `json:"password"`
	}
	if err := t.rest.Call(fmt.Sprintf("%sdata", currentRevision.String()), nil, http.MethodGet, nil, &response, t.createHeader()); err != nil {
		return "", err
	}
	return response.Password, nil
}

func (t *teamvaultPasswordProvider) User(key model.TeamvaultKey) (model.TeamvaultUser, error) {
	var response struct {
		User model.TeamvaultUser `json:"username"`
	}
	if err := t.rest.Call(fmt.Sprintf("%s/api/secrets/%s/", t.url.String(), key.String()), nil, http.MethodGet, nil, &response, t.createHeader()); err != nil {
		return "", err
	}
	return response.User, nil
}

func (t *teamvaultPasswordProvider) Url(key model.TeamvaultKey) (model.TeamvaultUrl, error) {
	var response struct {
		Url model.TeamvaultUrl `json:"url"`
	}
	if err := t.rest.Call(fmt.Sprintf("%s/api/secrets/%s/", t.url.String(), key.String()), nil, http.MethodGet, nil, &response, t.createHeader()); err != nil {
		return "", err
	}
	return response.Url, nil
}

func (t *teamvaultPasswordProvider) CurrentRevision(key model.TeamvaultKey) (model.TeamvaultCurrentRevision, error) {
	var response struct {
		CurrentRevision model.TeamvaultCurrentRevision `json:"current_revision"`
	}
	if err := t.rest.Call(fmt.Sprintf("%s/api/secrets/%s/", t.url.String(), key.String()), nil, http.MethodGet, nil, &response, t.createHeader()); err != nil {
		return "", err
	}
	return response.CurrentRevision, nil
}

func (t *teamvaultPasswordProvider) File(key model.TeamvaultKey) (model.TeamvaultFile, error) {
	rev, err := t.CurrentRevision(key)
	if err != nil {
		return "", fmt.Errorf("get current revision failed: %v", err)
	}
	var response struct {
		File model.TeamvaultFile `json:"file"`
	}
	if err := t.rest.Call(fmt.Sprintf("%sdata", rev.String()), nil, http.MethodGet, nil, &response, t.createHeader()); err != nil {
		return "", err
	}
	return response.File, nil
}

func (t *teamvaultPasswordProvider) createHeader() http.Header {
	header := make(http.Header)
	header.Add("Authorization", fmt.Sprintf("Basic %s", http_header.CreateAuthorizationToken(t.user.String(), t.pass.String())))
	header.Add("Content-Type", "application/json")
	return header
}

func (t *teamvaultPasswordProvider) Search(search string) ([]model.TeamvaultKey, error) {
	var response struct {
		Results []struct {
			ApiUrl model.TeamvaultApiUrl `json:"api_url"`
		} `json:"results"`
	}
	values := url.Values{}
	values.Add("search", search)
	if err := t.rest.Call(fmt.Sprintf("%s/api/secrets/", t.url.String()), values, http.MethodGet, nil, &response, t.createHeader()); err != nil {
		return nil, err
	}
	var result []model.TeamvaultKey
	for _, re := range response.Results {
		key, err := re.ApiUrl.Key()
		if err != nil {
			return nil, err
		}
		result = append(result, key)
	}
	return result, nil
}
