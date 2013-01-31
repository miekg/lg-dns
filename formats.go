package main

import (
	"fmt"
	"github.com/miekg/unbound"
	"net/http"
)

func Html(w http.ResponseWriter, u *unbound.Result) {
	for _, r := range u.Rr {
		fmt.Fprintf(w, "%s\n", r.String())
	}
}

func Json(w http.ResponseWriter, u *unbound.Result) {
	for _, r := range u.Rr {
		fmt.Fprintf(w, "%s\n", r.String())
	}
}

func Zone(w http.ResponseWriter, u *unbound.Result) {
	for _, r := range u.Rr {
		fmt.Fprintf(w, "%s\n", r.String())
	}
}

func Text(w http.ResponseWriter, u *unbound.Result) {
	for _, r := range u.Rr {
		fmt.Fprintf(w, "%s\n", r.String())
	}
}

func Xml(w http.ResponseWriter, u *unbound.Result) {
	s, _ := unboundToXML(u)
	fmt.Fprintf(w, "%s\n", s)
}
