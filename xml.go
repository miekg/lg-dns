package main

import (
	"encoding/json"
	"encoding/xml"
	"github.com/miekg/dns"
	"github.com/miekg/unbound"
	"time"
)

// http://tools.ietf.org/html/draft-mohan-dns-query-xml-00
// DNS message in XML

type Query struct {
	When        time.Time
	Duration    time.Duration
	Version     string
	Description string
	Server      string
}

type Response struct {
	Id       uint16           `xml:"-" json:"-"`
	Aa       int              `xml:"aa" json:"aa"`
	Ad       int              `xml:"ad" json:"ad"`
	Cd       int              `xml:"cd" json:"cd"`
	Rcode    string           `xml:"rcode" json:"rcode"`
	Anscount int              `xml:"anscount" json:"anscount"`
	Answers  []ResourceRecord `xml:"answers>answer" json:"answers"`
	//	Nscount     int              `xml:"nscount" json:"nscount"`
	//	Authorities []ResourceRecord `xml:"authorities>authority" json:"authorities"`
	//	Arcount     int              `xml:"arcount" json:"arcount"`
	//	Additionals []ResourceRecord `xml:"additionals>additional" json:"additionals"`
}

type LookingGlass struct {
	XMLName    xml.Name `xml:"Result" json:"-"`
	Query      Query    `xml:"Query" json:"Query"`
	Response   Response `xml:"Response" json:"Response"`
	Validation string   `xml:"Validation" json:"Validation"`
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

func boolToString(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

func unboundToXML(u *unbound.Result) (string, error) {
	output, err := xml.MarshalIndent(toLookingGlass(u), "  ", "    ")
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func unboundToJson(u *unbound.Result) (string, error) {
	output, err := json.MarshalIndent(toLookingGlass(u), "  ", "    ")
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func unboundToZone(u *unbound.Result) (string, error) {
	lg := toLookingGlass(u)
	output := "; When: " + lg.Query.When.String() + "\n"
	output += "; Query duration: " + lg.Query.Duration.String() + "\n"
	output += "; Version: " + lg.Query.Version + "\n"
	output += "; Description: " + lg.Query.Description + "\n"
	output += "; Server: " + lg.Query.Server + "\n"
	output += "; Flags: Aa: " + boolToString(u.AnswerPacket.Authoritative)
	output += ", Ad: " + boolToString(u.AnswerPacket.AuthenticatedData)
	output += ", Cd: " + boolToString(u.AnswerPacket.CheckingDisabled)
	output += ", Rcode: " + dns.RcodeToString[u.AnswerPacket.Rcode] + "\n"
	output += "\n; Answer Section:\n"
	for _, r := range u.AnswerPacket.Answer {
		output += r.String() + "\n"
	}
	if u.Secure {
		output += "\n; Validation: Secure\n"
	}
	if u.Bogus {
		output += "\n; Validation: " + u.WhyBogus + "\n"
	}
	return output, nil
}

func toLookingGlass(u *unbound.Result) *LookingGlass {
	l := &LookingGlass{Query: Query{time.Now(), u.Rtt, ver, "Managed by " + *mail + ", " + *loc, *res},
		Response: Response{Id: u.AnswerPacket.Id,
			Aa:       boolToInt(u.AnswerPacket.Authoritative),
			Ad:       boolToInt(u.AnswerPacket.AuthenticatedData),
			Cd:       boolToInt(u.AnswerPacket.CheckingDisabled),
			Rcode:    dns.RcodeToString[u.AnswerPacket.Rcode],
			Anscount: len(u.AnswerPacket.Answer)}}
	l.Response.Answers = sectionToResourceRecords(u.AnswerPacket.Answer)
	if u.Secure {
		l.Validation = "Secure"
	}
	if u.Bogus {
		l.Validation = u.WhyBogus
	}
	return l
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
