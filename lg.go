package main

// DNS Looking glass
// http://www.bortzmeyer.org/dns-lg.html

// TODO:
// * Add ?subnet option?

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/miekg/dns"
	"github.com/miekg/unbound"
	"html"
	"io"
	"log"
	"log/syslog"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	lg   *log.Logger
	mail *string
	loc  *string
	res  *string
	info string
	ver  string
)

func void(w http.ResponseWriter, r *http.Request) {}

// http://dns.bortzmeyer.org/{+domain}/{querytype}{?format,server,buffersize,dodnssec,tcp,reverse
func handler(w http.ResponseWriter, r *http.Request, typ string) {
	var (
		dnstype uint16
		ok      bool
	)
	lg.Printf("request from %s %s\n", r.RemoteAddr, r.URL)
	if dnstype, ok = dns.StringToType[typ]; !ok {
		fmt.Fprintf(w, "Record type "+typ+" does not exist")
		return
	}
	domain := mux.Vars(r)["domain"]
	domain = html.UnescapeString(domain)

	u := unbound.New()
	defer u.Destroy()
	forward := false
	format := "html"
	u.SetOption("module-config:", "iterator")
	for k, v := range r.URL.Query() {
		switch k {
		case "tcp":
			if v[0] == "1" {
				u.SetOption("tcp-upstream:", "yes")
			}
		case "dodnssec":
			if v[0] == "1" {
				u.SetOption("module-config:", "validator iterator")
				u.SetOption("edns-buffer-size:", "4096")
				u.AddTa(`;; ANSWER SECTION:
.                       168307 IN DNSKEY 257 3 8 (
                                AwEAAagAIKlVZrpC6Ia7gEzahOR+9W29euxhJhVVLOyQ
                                bSEW0O8gcCjFFVQUTf6v58fLjwBd0YI0EzrAcQqBGCzh
                                /RStIoO8g0NfnfL2MTJRkxoXbfDaUeVPQuYEhg37NZWA
                                JQ9VnMVDxP/VHL496M/QZxkjf5/Efucp2gaDX6RS6CXp
                                oY68LsvPVjR0ZSwzz1apAzvN9dlzEheX7ICJBBtuA6G3
                                LQpzW5hOA2hzCTMjJPJ8LbqF6dsV6DoBQzgul0sGIcGO
                                Yl7OyQdXfZ57relSQageu+ipAdTTJ25AsRTAoub8ONGc
                                LmqrAmRLKBP1dfwhYB4N7knNnulqQxA+Uk1ihz0=
                                ) ; key id = 19036`)
			}
		case "buffersize":
			if err := u.SetOption("edns-buffer-size:", v[0]); err != nil {
				fmt.Fprintf(w, "Not a valid buffer size: %s", v[0])
				return
			}
		case "server":
			if err := u.SetFwd(v[0]); err != nil {
				fmt.Fprintf(w, "Not a valid server `%s': %s", v[0], err.Error())
				return
			}
			forward = true
		case "reverse":
			if v[0] == "1" {
				var err error
				dnstype = dns.TypePTR
				domain, err = dns.ReverseAddr(domain)
				if err != nil {
					fmt.Fprintf(w, "Not a valid IP address: %s", v[0])
					return
				}
			}
		case "format": // unsupported format defaut to html
			for _, f := range []string{"html", "zone", "xml", "json", "text"} {
				if v[0] == f {
					format = f
				}
			}
		}
	}
	if !forward {
		u.ResolvConf("/etc/resolv.conf")
	}
	d, err := u.Resolve(domain, dnstype, dns.ClassINET)
	if err != nil {
		fmt.Fprintf(w, "error")
		return
	}
	if !d.HaveData {
		fmt.Fprintf(w, "Domain %s (type %s) does not exist", domain, dns.TypeToString[dnstype])
		return
	}
	switch format {
	case "json":
		Json(w, d)
	case "xml":
		Xml(w, d)
	case "html":
		fallthrough
	case "text":
		fallthrough
	case "zone":
		Zone(w, d)
	}
}

func indexhtml(w http.ResponseWriter, r *http.Request) {
	h, e := os.Open("README.html")
	if e != nil {
		fmt.Fprintf(w, "Documentation not found")
		return
	}
	defer h.Close()
	io.Copy(w, h)
}

func main() {
	port := flag.Int("port", 80, "port number to use")
	mail = flag.String("mail", "nobody@example.com", "email of service maintainer")
	loc = flag.String("loc", "COUNTRY, hosted at HOSTER, AS NNNN", "location of the server")
	ver = "DNS Looking Glass Go version"
	res = flag.String("res", "Unbound with DNSSEC validation", "resolver used")
	flag.Parse()

	info = "Service managed by " + *mail + " / Local resolver is " + *res +
		", the machine is in " + *loc + " / " + ver

	var err error
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		indexhtml(w, r)
	})
	router.HandleFunc("/index.html", func(w http.ResponseWriter, r *http.Request) {
		indexhtml(w, r)
	})
	router.HandleFunc("/{domain}", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, "A")
	})
	router.HandleFunc("/{domain}/{type}", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, strings.ToUpper(mux.Vars(r)["type"]))
	})
	http.HandleFunc("/favicon.ico", void)
	http.Handle("/", router)

	lg, err = syslog.NewLogger(syslog.LOG_INFO, log.LstdFlags)
	if err != nil {
		log.Fatal("NewLogger: ", err)
	}

	err = http.ListenAndServe(":"+strconv.Itoa(*port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
