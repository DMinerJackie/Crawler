package main

import (
	"github.com/kdar/factorlog"
	"os"
)

func Logger() {
	file, _ := os.OpenFile("logfile.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	log := factorlog.New(file, factorlog.NewStdFormatter("%{Date} %{Time}\tLOG\t%{Message}"))
	log.Println("Innerhalb Basic formatter")
}