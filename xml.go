package main

// http://tools.ietf.org/html/draft-mohan-dns-query-xml-00
// DNS message in XML
type Answer []ResourceRecord
type Authority []ResourceRecord
type Additional []ResourceRecord

type Response struct {
	Id uint16
	Aa bool
	Ad bool
	Cd bool
	Rcode string
	Anscount int
	Answers Answer
	Nscount int
	Authorities Authority
	Arcount int
	Additionals Additional
}

type ResourceRecord struct {
	Name string
	Type string
	Rdlength int
	Rdata string
}
