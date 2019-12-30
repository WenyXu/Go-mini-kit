package session

import (
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"net/http"
	"strings"
	"time"
)

var (
	sessionIdNamePrefix = "session-id-"
	store               *sessions.CookieStore
)

func init() {
	// 随机生成32位加密的key，
	// 正式环境一定不要暴露，通过写到环境变量或其它安全方式
	store = sessions.NewCookieStore([]byte("OnNUU5RUr6Ii2HMI0d6E54bXTS52tCCL"))
}

// GetSession current session
func GetSession(w http.ResponseWriter, r *http.Request) *sessions.Session {
	// sessionId
	var sId string

	for _, c := range r.Cookies() {
		if strings.Index(c.Name, sessionIdNamePrefix) == 0 {
			sId = c.Name
			break
		}
	}

	// if sessionId is empty ,random a new one
	if sId == "" {
		sId = sessionIdNamePrefix + uuid.New().String()
	}

	ses, _ := store.Get(r, sId)

	if ses.ID == "" {
		// put sessionId into the cookie
		cookie := &http.Cookie{Name: sId, Value: sId, Path: "/", Expires: time.Now().Add(30 * time.Second), MaxAge: 0}
		http.SetCookie(w, cookie)

		// save the session
		ses.ID = sId
		_ = ses.Save(r, w)
	}
	return ses
}
