package main

import (
	"encoding/json"
	"fmt"
	"github.com/miekg/unbound"
	"net/http"
)

func Html(w http.ResponseWriter, u *unbound.Result) {
	fmt.Fprintf(w, "%s\n", u.AnswerPacket)
}

func Json(w http.ResponseWriter, u *unbound.Result) {
	bytes, err := json.Marshal(u.AnswerPacket)
	if err != nil { // must always work?
		return
	}
	fmt.Fprintf(w, "%s\n", bytes)
}

func Zone(w http.ResponseWriter, u *unbound.Result) {
	fmt.Fprintf(w, "%s\n", u.AnswerPacket)
}

func Text(w http.ResponseWriter, u *unbound.Result) {
	fmt.Fprintf(w, "%s\n", u.AnswerPacket)
}

func Xml(w http.ResponseWriter, u *unbound.Result) {
	fmt.Fprintf(w, "%s\n", u.AnswerPacket)
}
