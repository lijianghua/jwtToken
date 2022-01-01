package tokenService

import (
	"fmt"
	"jwtToken/util"
	"net/http"
	"strings"
	"time"
)

type Token struct {
	AccessID        string        `json:"-"`
	Access          string        `json:"access-token"`
	AccessCreateAt  time.Time     `json:"access-token-create-at"`
	AccessExpiresIn time.Duration `json:"access-token-expires-in"`

	RefreshID        string        `json:"-"`
	Refresh          string        `json:"refresh-token"`
	RefreshCreateAt  time.Time     `json:"refresh-token-create-at"`
	RefreshExpiresIn time.Duration `json:"refresh-token-expires-in"`
}

type AccessDetails struct {
	AccessID string
	UserID   string
}

//ExtractToken: Extract access token from Bearer Authorization header
func ExtractToken(r *http.Request) (string, bool) {
	auth := r.Header.Get("Authorization")
	prefix := "Bearer "
	token := ""

	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	}

	return token, token != ""
}

func (m *Token) Response(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")

	resp := util.NewRespMsg(0, "OK", struct {
		Access           string `json:"access_token"`
		TokenType        string `json:"token_type"`
		AccessExpiresIn  string `json:"access_expires_in"`
		Refresh          string `json:"refresh_token"`
		RefreshExpiresIn string `json:"refresh_expires_in"`
	}{
		m.Access,
		"Bearer",
		fmt.Sprintf("%.0f", m.AccessExpiresIn.Seconds()),
		m.Refresh,
		fmt.Sprintf("%.0f", m.RefreshExpiresIn.Seconds()),
	})

	resp.WriteTo(w)
}
