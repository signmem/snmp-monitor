package http


import (
	"net/http"
	"encoding/json"
	"github.com/signmem/snmp-monitor/g"
)

func init() {
	healthCheck()
}

type Dto struct {
	Msg		string			`json:"msg"`
	Data    interface{}     `json:"data"`
}


func RenderJson(w http.ResponseWriter, v interface{}) {
	bs, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(bs)
}

func RenderDataJson(w http.ResponseWriter, data interface{}) {
	RenderJson(w, Dto{Msg: "success", Data: data})
}

func RenderMsgJson(w http.ResponseWriter, msg string) {
	RenderJson(w, map[string]string{"msg": msg})
}

func AutoRender(w http.ResponseWriter, data interface{}, err error) {
	if err != nil {
		RenderMsgJson(w, err.Error())
		return
	}

	RenderDataJson(w, data)
}

func Start() {

	listen := g.Config().Listen

	s := &http.Server{
		Addr:           listen,
		MaxHeaderBytes: 1 << 30,
	}

	g.Logger.Infof("listening %s", listen)
	g.Logger.Fatalln(s.ListenAndServe())
}

func healthCheck() {
	http.HandleFunc("/_health_check",
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})
}
