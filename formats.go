package main

import (
	"github.com/miekg/unbound"
	"net/http"
)

func Html(w http.ResponseWriter, u *unbound.Result) {
}

func Json(w http.ResponseWriter, u *unbound.Result) {
}

func Zone(w http.ResponseWriter, u *unbound.Result) {
}

func Text(w http.ResponseWriter, u *unbound.Result) {
}

func Xml(w http.ResponseWriter, u *unbound.Result) {
}
