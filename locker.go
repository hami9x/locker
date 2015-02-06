package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

const (
	KeyComment = " #LOCKED BY LOCKER"
	Secs       = 42
	EntD       = 30
)

func stop(hostsFile *os.File) {
	hostsdata, err := ioutil.ReadAll(hostsFile)
	if err != nil {
		panic(err)
	}

	shd := string(hostsdata)

	lines := strings.Split(shd, "\n")
	result := []string{}
	for _, l := range lines {
		if !strings.HasSuffix(l, KeyComment) {
			result = append(result, l)
		}
	}

	sr := strings.Join(result, "\n")
	if sr == shd {
		return
	}

	err = hostsFile.Truncate(0)
	if err != nil {
		panic(err)
	}

	_, err = hostsFile.Seek(0, 0)
	if err != nil {
		panic(err)
	}

	_, err = hostsFile.WriteString(sr)
	if err != nil {
		panic(err)
	}
}

func start(conf []byte, hostsFile *os.File) {
	confLines := strings.Split(string(conf), "\n")
	for i, _ := range confLines {
		if strings.TrimSpace(confLines[i]) == "" {
			continue
		}

		confLines[i] = "0.0.0.0 " + confLines[i] + KeyComment
	}

	_, err := hostsFile.Seek(0, 2)
	if err != nil {
		panic(err)
	}

	_, err = hostsFile.WriteString("\n" + strings.Join(confLines, "\n"))
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()
	hostsFile, err := os.OpenFile("/etc/hosts", os.O_RDWR, 0644)
	defer hostsFile.Close()
	if err != nil {
		panic(err)
	}

	conf, err := ioutil.ReadFile(path.Join("/etc/locker_hosts"))
	if err != nil {
		panic(err)
	}

	switch flag.Arg(0) {
	case "iwannafckngentertainpleaseletmesir":
		stop(hostsFile)
		fmt.Printf("You got %v Minutes\n", EntD)
		time.Sleep(EntD * time.Minute)
		start(conf, hostsFile)
	case "tempoaccs":
		stop(hostsFile)
		fmt.Printf("You got %v Seconds\n", Secs)
		time.Sleep(Secs * time.Second)
		start(conf, hostsFile)
	case "":
		stop(hostsFile)
		start(conf, hostsFile)
	default:
		fmt.Printf("What are you smoking?\n")
	}
}
