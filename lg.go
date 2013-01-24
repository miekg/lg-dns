package main

// DNS Looking glass
// http://www.bortzmeyer.org/dns-lg.html

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/miekg/dns"
	"github.com/miekg/unbound"
	"log"
	"net/http"
	"strings"
)

func void(w http.ResponseWriter, r *http.Request) {}

// http://dns.bortzmeyer.org/{+domain}/{querytype}{?format,server,buffersize,dodnssec,tcp,reverse
func handler(w http.ResponseWriter, r *http.Request, typ string) {
	var (
		dnstype uint16
		ok      bool
	)
	log.Printf("request from %s %s\n", r.RemoteAddr, r.URL)
	if dnstype, ok = dns.StringToType[typ]; !ok {
		fmt.Fprintf(w, "Record type "+typ+" does not exist")
		return
	}
	values := r.URL.Query()
	domain := mux.Vars(r)["domain"]

	u := unbound.New()
	defer u.Destroy()
	fwd := false
	u.SetOption("edns-buffer-size:", "4096")
	for k, v := range values {
		switch k {
		case "tcp":
			if v[0] == "1" {
				u.SetOption("tcp-upstream:", "yes")
			}
		case "dodnssec":
			if v[0] == "1" {
				u.AddTaFile("Kroot.key")
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
			fwd = true
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
		}
	}
	if !fwd {
		u.ResolvConf("/etc/resolv.conf")
	}
	d, err := u.Resolve(domain, dnstype, dns.ClassINET)
	if err != nil {
		fmt.Fprintf(w, "error")
		return
	}
	if !d.HaveData {
		fmt.Fprintf(w, "Domain %s does not exist", domain)
		return
	}
	fmt.Fprintf(w, "%s\n", d.AnswerPacket)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/{domain}", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, "A")
	})
	router.HandleFunc("/{domain}/{type}", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, strings.ToUpper(mux.Vars(r)["type"]))
	})
	http.HandleFunc("/favicon.ico", void)
	http.Handle("/", router)

	e := http.ListenAndServe(":8080", nil)
	if e != nil {
		log.Fatal("ListenAndServe: ", e)
	}
}
