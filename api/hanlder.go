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
	"gopkg.in/oauth2.v3"
	"poste/consul"
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
		uuid := ticket.UUID(userId, ti.GetClientID())

		node, ok := mailmenRing.GetNode(uuid)

		if !ok {
			w.Write(Response{Err: "get mailman node failed"}.Marshal())
		} else {
			w.Write(Response{Data:map[string]string{"node": node, "ticket":t}}.Marshal())
		}
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

	services := consul.Get(consul.Redis)
	if len(services) <= 0 {
		util.LogPanic("redis is not currently available, try `poste init` to initialize the services from config.")
	}
	service := services[0]
	redisConf := redis.Config{
		Addr: util.ToAddr(service.Host, service.Port),
	}
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