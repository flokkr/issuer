package main

import (
	"flag"
	"fmt"
	"github.com/progrium/go-basher"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

var bash, _ = basher.NewContext("/bin/bash", false)
var keyDir = "/data/"

func source(script string) {
	err := bash.Source("bash/"+script, Asset)
	if err != nil {
		panic(err)
	}
}
func init() {
	source("keytab.sh")
	source("root.sh")
	source("issue.sh")
	source("check.sh")
}

func checkTrustFile() (string) {
	trustFile := path.Join(keyDir, "trust.keystore")
	if _, err := os.Stat(trustFile); os.IsNotExist(err) {
		_, err := bash.Run("root", []string{keyDir})
		if err != nil {
			log.Printf(err.Error())
		}
	}
	return trustFile
}

func checkCertFile(name string) (string) {
	certFile := path.Join(keyDir, name + ".keystore")
	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		_, err := bash.Run("issue", []string{keyDir, name})
		if err != nil {
			log.Printf(err.Error())
		}
	}
	return certFile
}

func checkKeytab(host string, serviceName string) string {
	principal := serviceName + "/" + host + "@EXAMPLE.COM"
	keytabFile := path.Join(keyDir, serviceName + "." + host + ".keytab")
	if _, err := os.Stat(keytabFile); os.IsNotExist(err) {
		args := []string{principal, keytabFile}
		_, err := bash.Run("keytab", args)
		if err != nil {
			log.Printf(err.Error())
		}

	}
	return keytabFile

}

func trustStoreGenerator(w http.ResponseWriter, r *http.Request) {
	trustFile := checkTrustFile()
	content, _ := ioutil.ReadFile(trustFile)
	w.Write(content)
}

func certGenerator(w http.ResponseWriter, r *http.Request) {
	checkTrustFile()
	var segments = strings.Split(r.URL.String(), "/")
	certFile := checkCertFile(segments[2])
	content, _ := ioutil.ReadFile(certFile)
	w.Write(content)
}

func keytabGenerator(w http.ResponseWriter, r *http.Request) {
	var segments = strings.Split(r.URL.String(), "/")
	host := segments[2]
	serviceName := segments[3]
	keytabFile := checkKeytab(host, serviceName)
	content, _ := ioutil.ReadFile(keytabFile)
	w.Write(content)
}

//wait until kds is available
func waitForKdc() {
	for {
		resp, err := bash.Run("check", []string{keyDir})
		if err != nil || resp != 0 {
			if resp != 0 {
				fmt.Println("KDC is not yet available.  Shell return code is " + strconv.Itoa(resp))
			} else {
				fmt.Println("KDC is not yet available " + err.Error())
			}
			time.Sleep(1 * time.Second)
		} else {
			return
		}
	}
}

func main() {
	port := flag.String("port", "8081", "Http port to listen on")
	flag.Parse()

	waitForKdc()
	http.HandleFunc("/keystore/", certGenerator)
	http.HandleFunc("/truststore", trustStoreGenerator)
	http.HandleFunc("/keytab/", keytabGenerator)
	print("Issuer is listening on : " + *port)
	if err := http.ListenAndServe(":" + *port, nil); err != nil {
		log.Fatal(err)
	}

}
