package api

import (
	"net/http"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/store"
	"gopkg.in/oauth2.v3/models"
	"github.com/jinzhu/configor"
	"poste/util"
	"github.com/go-oauth2/redis"
)

func handleRequest() {
	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		oauthSvr.HandleTokenRequest(w, r)
	})

	http.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {
		ok := Verify(r)
		if !ok {
			w.Write([]byte("invalid token"))
		} else {
			w.Write([]byte("success"))
		}
	})
}

func Verify(r *http.Request) bool {
	srv := buildSrv()

	ti, err := srv.ValidationBearerToken(r)
	if err != nil {
		util.LogError("access token verify failed : %s", err)
		return false
	}
	if r.URL.Query().Get("scope") != ti.GetScope() {
		util.LogError("request out of scope. asking: %s, given: %s", r.URL.Query()["scope"], ti.GetScope())
		return false
	}

	return true
}

func buildSrv() (srv *server.Server) {
	manager := buildManager()

	srv = server.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)

	return
}

func buildManager() (manager *manage.Manager) {
	manager = manage.NewDefaultManager()

	redisConf := redis.Config{}
	configor.Load(&redisConf, "config/redis.json")
	manager.MustTokenStorage(redis.NewTokenStore(&redisConf))

	clientStore := store.NewClientStore()
	clientsConf := []*models.Client{}
	configor.Load(&clientsConf, "config/oauth2.json")
	for _, clientConfig := range clientsConf {
		clientStore.Set(clientConfig.ID, clientConfig)
	}
	manager.MapClientStorage(clientStore)

	return
}