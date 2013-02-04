package main

import (
	"fmt"
	"github.com/miekg/unbound"
	"net/http"
)

func Html(w http.ResponseWriter, u *unbound.Result) {
	s, _ := unboundToHTML(u)
	fmt.Fprintf(w, "%s\n", s)
}

func Xml(w http.ResponseWriter, u *unbound.Result) {
	s, _ := unboundToXML(u)
	fmt.Fprintf(w, "%s\n", s)
}

func Json(w http.ResponseWriter, u *unbound.Result) {
	s, _ := unboundToJson(u)
	fmt.Fprintf(w, "%s\n", s)
}

func Zone(w http.ResponseWriter, u *unbound.Result) {
	s, _ := unboundToZone(u)
	fmt.Fprintf(w, "%s\n", s)
}
