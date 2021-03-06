package git

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/svagner/go-gitlab/config"
)

/* {"name":string,"username":string,"id":int, "state":string,
"avatar_url":string,
"created_at":time.Time,
"is_admin":bool,"bio":null,"skype":string,"linkedin":string,"twitter":string,
"website_url":string,"email":string,"theme_id":int,
"color_scheme_id":int,"projects_limit":int,
"identities":nil,"can_create_group":bool,"can_create_project":bool} */

type UserInfo struct {
	Name             string `json:"name"`
	Username         string `json:"username"`
	Id               int    `json:"id"`
	State            string `json:"state"`
	Avatar_url       string `json:"avatar_url"`
	Is_admin         bool   `json:"is_admin"`
	Skype            string `json:"skype"`
	LinkedIn         string `json:"linkedin"`
	Twitter          string `json:"twitter"`
	Website          string `json:"website_url"`
	Email            string `json:"email"`
	ThremeId         int    `json:"threme_id"`
	ColourSchemeId   int    `json:"color_scheme_id"`
	ProjectsLimit    int    `json:"projects_limit"`
	CanCreateGroup   bool   `json:"can_create_group"`
	CanCreateProject bool   `json:"can_create_project"`
}

var (
	host   string
	scheme string
	user   string
	token  string
)

const (
	USERAPI_PREFIX  = "/api/v3/users/"
	USERAPI_POSTFIX = "/?private_token="
)

func InitGitLabApi(conf config.GitLab) {
	host = conf.Host
	scheme = conf.Scheme
	user = conf.User
	token = conf.Token
}

func GetUserInfo(userId int) (*UserInfo, error) {
	var prefix string

	if host == "" {
		return nil, errors.New("GitLab's host wasn't found")
	}
	if token == "" {
		return nil, errors.New("Token for gitlab wasn't defined")
	}

	if scheme == "" {
		prefix = "http://" // set prefix http by default
	} else {
		prefix = scheme + "://"
	}

	resp, err := http.Get(prefix + host + USERAPI_PREFIX + strconv.Itoa(userId) + USERAPI_POSTFIX + token)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("Get request fow get user info (GitLab API) returned wrong status: " + strconv.Itoa(resp.StatusCode))
	}
	data := make([]byte, resp.ContentLength)
	_, err = resp.Body.Read(data)
	if err != nil && err != io.EOF {
		return nil, err
	}

	user := new(UserInfo)
	err = json.NewDecoder(bytes.NewReader(data)).Decode(user)
	if user.Skype == "" {
		return nil, errors.New("User " + strconv.Itoa(userId) + " hasn't got any skype accounts")
	}
	return user, nil
}
