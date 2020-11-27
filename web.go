package hitbtc

import (
	"net/http"
	"text/template"
)

var (
	homeTemplate    template.Template
	symbolsTemplate template.Template
)

func (m *MessageRouter) Serve() {
	template.ParseFiles()

	m.server.HandleFunc("/", m.home)
	m.server.HandleFunc("/symbols", m.symbols)

}

func (m *MessageRouter) home(w http.ResponseWriter, r *http.Request) {

}

func (m *MessageRouter) symbols(w http.ResponseWriter, r *http.Request) {

}
