package main

import (
	"encoding/xml"
	"encoding/json"
	"github.com/miekg/dns"
	"github.com/miekg/unbound"
)

// http://tools.ietf.org/html/draft-mohan-dns-query-xml-00
// DNS message in XML

type Response struct {
	Id          uint16           `xml:"id,json:"id"`
	Aa          int              `xml:"aa"`
	Ad          int              `xml:"ad"`
	Cd          int              `xml:"cd"`
	Rcode       string           `xml:"rcode"`
	Anscount    int              `xml:"anscount"`
	Answers     []ResourceRecord `xml:"answers>answer,json:"answers>answer"`
	Nscount     int              `xml:"nscount"`
	Authorities []ResourceRecord `xml:"authorities>authority,json:authorities>authority"`
	Arcount     int              `xml:"arcount"`
	Additionals []ResourceRecord `xml:"additionals>additional,json:additionals>additional"`
}

type ResourceRecord struct {
	Name     string `xml:"name"`
	Type     string `xml:"type"`
	Class    string `xml:"class"`
	Ttl      uint32 `xml:"ttl"`
	Rdlength uint16 `xml:"rdlength"`
	Rdata    string `xml:"rdata"`
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func unboundToXML(u *unbound.Result) (string, error) {
	output, err := xml.MarshalIndent(toResponse(u), "  ", "    ")
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func unboundToJson(u *unbound.Result) (string, error) {
	output, err := json.MarshalIndent(toResponse(u), "  ", "    ")
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func toResponse(u *unbound.Result) *Response {
	r := &Response{Id: u.AnswerPacket.Id,
		Aa:       boolToInt(u.AnswerPacket.Authoritative),
		Ad:       boolToInt(u.AnswerPacket.AuthenticatedData),
		Cd:       boolToInt(u.AnswerPacket.CheckingDisabled),
		Rcode:    dns.RcodeToString[u.AnswerPacket.Rcode],
		Anscount: len(u.AnswerPacket.Answer),
		Nscount:  len(u.AnswerPacket.Ns),
		Arcount:  len(u.AnswerPacket.Extra)}

	r.Answers = sectionToResourceRecords(u.AnswerPacket.Answer)
	r.Authorities = sectionToResourceRecords(u.AnswerPacket.Ns)
	r.Additionals = sectionToResourceRecords(u.AnswerPacket.Extra)
	return r
}

func sectionToResourceRecords(section []dns.RR) []ResourceRecord {
	var a []ResourceRecord
	for _, r := range section {
		x := new(dns.RFC3597)
		x.ToRFC3597(r)
		a = append(a, ResourceRecord{r.Header().Name,
			dns.TypeToString[r.Header().Rrtype],
			dns.ClassToString[r.Header().Class],
			r.Header().Ttl,
			r.Header().Rdlength,
			x.Rdata})
	}
	return a
}
