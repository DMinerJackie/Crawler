package main

import (
	"fmt"
	"github.com/kdar/factorlog"
	"io/ioutil"
	"net/http"
	"strings"
	"os"
)

var ip = "http://127.0.0.1:8080"
var file, _ = os.OpenFile("logfile.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
var logg = factorlog.New(os.Stdout, factorlog.NewStdFormatter("%{Message}"))
var loggFile = factorlog.New(file, factorlog.NewStdFormatter("%{Message}"))


func Phantom(link string) {

	buf := strings.NewReader(link)
	resp, err := http.Post(ip, "text/plain", buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	bodyString := string(bodyBytes)
	logg.Println(bodyString + " # " + link)
	loggFile.Println(bodyString + " # " + link)
	//Info.Println("TEST" + bodyString + " : " + link)

	return
}
