package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/miekg/dns"
	"github.com/miekg/unbound"
	"net/http"
	"time"
)

type Query struct {
	duration    time.Duration
	Versions    string
	Description string
	Server      string
}

type QuestionSection struct {
	Qtype string
	Qname string
}

type Answer struct {
	Query           *Query
	AnswerSection   []dns.RR
	ReturnCode      string
	QuestionSection *QuestionSection
}

func unboundToAnswer(u *unbound.Result) *Answer {
	a := new(Answer)
	a.Query = &Query{u.Rtt, ver, *loc, ""}
	a.AnswerSection = u.AnswerPacket.Answer
	a.ReturnCode = dns.RcodeToString[u.Rcode]
	a.QuestionSection = &QuestionSection{dns.TypeToString[u.Qtype], u.Qname}
	return a
}

func Html(w http.ResponseWriter, u *unbound.Result) {
	fmt.Fprintf(w, "%+v\n", unboundToAnswer(u))
}

func Json(w http.ResponseWriter, u *unbound.Result) {
	b, err := json.Marshal(unboundToAnswer(u))
	if err != nil { // must always work?
		fmt.Fprintf(w, "%s", err.Error())
		return
	}
	b1 := new(bytes.Buffer)
	json.Indent(b1, b, "", "  ")
	fmt.Fprintf(w, "%s\n", b1)
}

func Zone(w http.ResponseWriter, u *unbound.Result) {
	fmt.Fprintf(w, "%s\n", unboundToAnswer(u))
}

func Text(w http.ResponseWriter, u *unbound.Result) {
	fmt.Fprintf(w, "%s\n", unboundToAnswer(u))
}

func Xml(w http.ResponseWriter, u *unbound.Result) {
	fmt.Fprintf(w, "%s\n", unboundToAnswer(u))
}
