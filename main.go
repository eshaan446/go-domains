package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Go")
}
func formHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "ParseForm error %v", err)
		return
	}
	fmt.Fprintf(w, "Post request successful\n")
	domain := r.FormValue("domain")
	fmt.Fprintf(w, "Domain -> %v\n", domain)
	var hasMX, hasSPF, hasDMARC bool
	var spfRecord, dmarcRecord string
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	if len(mxRecords) > 0 {
		hasMX = true
	}
	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		log.Printf("Error:%v\n", err)
	}
	for _, record := range txtRecords {
		if strings.HasPrefix(record, "v=spf1") {
			hasSPF = true
			spfRecord = record
			break
		}
	}

	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		log.Printf("Error%v\n", err)
	}
	for _, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			hasDMARC = true
			dmarcRecord = record
			break
		}
	}
	fmt.Fprintf(w, "Domain: %s\nHasMX: %t\nHasSPF: %t\nSPFRecord: %s\nHasDMARC: %t\nDMARCRecord: %s\n", domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord)
	if hasMX {
		if hasSPF {
			if hasDMARC {
				fmt.Fprintf(w, "It's a valid domain\n")
			} else {
				fmt.Fprintf(w, "Missing DMARC record\n")
			}
		} else {
			fmt.Fprintf(w, "Missing SPF record\n")
		}
	} else {
		fmt.Fprintf(w, "Missing MX record\n")
	}
}

func main() {
	fileserver := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileserver)
	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/hello", helloHandler)
	fmt.Printf("Starting server at port 8080\n")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

// func checkDomain(domain string) {
// 	var hasMX, hasSPF, hasDMARC bool
// 	var spfRecord, dmarcRecord string
// 	mxRecords, err := net.LookupMX(domain)
// 	if err != nil {
// 		log.Printf("Error: %v\n", err)
// 	}
// 	if len(mxRecords) > 0 {
// 		hasMX = true
// 	}
// 	txtRecords, err := net.LookupTXT(domain)
// 	if err != nil {
// 		log.Printf("Error:%v\n", err)
// 	}
// 	for _, record := range txtRecords {
// 		if strings.HasPrefix(record, "v=spf1") {
// 			hasSPF = true
// 			spfRecord = record
// 			break
// 		}
// 	}

// 	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
// 	if err != nil {
// 		log.Printf("Error%v\n", err)
// 	}
// 	for _, record := range dmarcRecords {
// 		if strings.HasPrefix(record, "v=DMARC1") {
// 			hasDMARC = true
// 			dmarcRecord = record
// 			break
// 		}
// 	}
// 	fmt.Printf("%v,%v,%v,%v,%v,%v", domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord)

// }
