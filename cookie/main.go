package cookie

import (
	"net/http"
	"time"
)

var globalSessions *Manager

func init() {
	globalSessions, _ = NewManager("memory", "gosessionid", 3600)
	go globalSessions.GC()
}

func Main() {

}

func setC(w http.ResponseWriter) {
	expiration := time.Now()
	expiration = expiration.AddDate(1, 0, 0)
	cookie := http.Cookie{Name: "username", Value: "louis", Expires: expiration}
	http.SetCookie(w, &cookie)
}
