package api

import (
	"net/http"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/manage"
	oauthStore "gopkg.in/oauth2.v3/store"
	"gopkg.in/oauth2.v3/models"
	"github.com/jinzhu/configor"
	"poste/util"
	"poste/ticket"
	"github.com/go-oauth2/redis"
	"runtime"
	"path"
	"gopkg.in/oauth2.v3"
	"poste/mailman"
)

func handleRequest() {
	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		oauthSvr.HandleTokenRequest(w, r)
	})

	//For debugging
	http.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {
		ti, err := Verify(r)
		if err != nil {
			w.Write(Response{Err:err.Error()}.Marshal())
		} else {
			w.Write(Response{Data:ti}.Marshal())
		}
	})

	http.HandleFunc("/bind", func(w http.ResponseWriter, r *http.Request) {
		ti, err := Verify(r)
		if err != nil {
			w.Write(Response{Err:err.Error()}.Marshal())
			return
		}

		userId := r.URL.Query().Get("userId")
		t := ticket.GetTicket(userId, ti.GetClientID(), true)

		w.Write(Response{Data:map[string]string{"ticket":t}}.Marshal())
	})

	http.HandleFunc("/mailman", func(w http.ResponseWriter, r *http.Request) {
		ti, err := Verify(r)
		if err != nil {
			w.Write(Response{Err:err.Error()}.Marshal())
			return
		}

		t := r.URL.Query().Get("type")
		userId := r.URL.Query().Get("userId")
		uuid := ticket.UUID(userId, ti.GetClientID())
		if t == string(mailman.WsType) {
			node, ok := mailmenWsRing.GetNode(uuid)
			if !ok {
				w.Write(Response{Err: "get mailman node failed"}.Marshal())
			} else {
				w.Write(Response{Data:map[string]string{"node": node}}.Marshal())
			}
			return
		}
		if t == string(mailman.TcpType) {
			node, ok := mailmenTcpRing.GetNode(uuid)
			if !ok {
				w.Write(Response{Err:"get mailman node failed"}.Marshal())
			} else {
				w.Write(Response{Data:map[string]string{"node": node}}.Marshal())
			}
			return
		}

		w.Write(Response{Err:"param type invalid"}.Marshal())
		return
	})

}

func Verify(r *http.Request) (ti oauth2.TokenInfo, err error) {
	srv := buildSrv()

	ti, err = srv.ValidationBearerToken(r)
	if err != nil {
		util.LogError("access token verify failed : %s", err)
		return
	}
	if r.URL.Query().Get("scope") != ti.GetScope() {
		util.LogError("request out of scope. asking: %s, given: %s", r.URL.Query()["scope"], ti.GetScope())
		return
	}
	return
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

	_, filename, _, _ := runtime.Caller(1)
	configPath := path.Join(path.Dir(filename), "../config/redis.json")
	redisConf := redis.Config{}
	configor.Load(&redisConf, configPath)
	manager.MustTokenStorage(redis.NewTokenStore(&redisConf))

	clientStore := oauthStore.NewClientStore()
	clientsConf := []*models.Client{}
	configor.Load(&clientsConf, "config/oauth2.json")
	for _, clientConfig := range clientsConf {
		clientStore.Set(clientConfig.ID, clientConfig)
	}
	manager.MapClientStorage(clientStore)

	return
}