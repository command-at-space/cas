/* */

package main

import (
	online "casServer/online"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func checkFlags() {
	versionFlag := flag.Bool("v", false, "Show current version and exit")
	flag.Parse()
	switch {
	case *versionFlag:
		fmt.Printf("Version:\t: %s\n", version)
		fmt.Printf("Date   :\t: %s\n", releaseDate)
		os.Exit(0)
	}
}

func loadConfigJSON(a *app) {
	err := json.Unmarshal(getGlobalConfigJSON(), a)
	if err != nil {
		log.Fatal("Error parsing JSON config => \n", err)
	}
}

func createCustomInfoLogFile(f string) {
	infoLog, err := os.OpenFile(f, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("ERROR opening Info log file %s\n", err)
	}
	iLog = log.New(infoLog, "INFO :\t", log.Ldate|log.Ltime)
}

func createCustomErrorLogFile(f string) *os.File {
	mylog, err := os.OpenFile(f, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("ERROR opening Error log file %s\n", err)
	}
	log.SetOutput(mylog)
	return mylog
}

func showList(pjs *online.PJs) {
	for nick, _ := range pjs.Online {
		//fmt.Printf("%s\n", nick)
		fmt.Printf("%s ", nick)
	}
	fmt.Println("")
}

func showListInterval(pjs *online.PJs) {
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				showList(pjs)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
