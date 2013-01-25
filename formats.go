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

/*
Query for: www.miek.nl., type SOA
Flags: 
Canonical name: a.miek.nl.
TTL: 39624
Resolver queried: ::1
Query done at: 2013-01-25 09:31:22Z
Query duration: 0:00:00.001288
Service description: / Local resolver is Unbound with DNSSEC validation, the machine is in the USA, hosted at 6sync, AS 46636.
DNS Looking Glass 2013012101, DNSpython version 1.10.0, Python version CPython 2.7.3 on Linux
*/
func (a *Answer) String() string {
	s := "Query for: " + a.QuestionSection.Qname + ", type " + a.QuestionSection.Qtype + "\n"
	s += "Flags:"
	return s
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
