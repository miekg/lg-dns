package main

import (
	"fmt"
	"github.com/miekg/unbound"
	"net/http"
)

func Xml(w http.ResponseWriter, u *unbound.Result) {
	s, _ := unboundToXML(u)
	w.Header().Set("Content-Type", "application/xml")
	fmt.Fprintf(w, "%s\n", s)
}

func Json(w http.ResponseWriter, u *unbound.Result) {
	s, _ := unboundToJson(u)
	w.Header().Set("Content-Type" ,"application/json")
	fmt.Fprintf(w, "%s\n", s)
}

func Zone(w http.ResponseWriter, u *unbound.Result) {
	s, _ := unboundToZone(u)
	fmt.Fprintf(w, "%s\n", s)
}
