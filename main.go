package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gliderlabs/ssh"
)

type handler struct {
	guard  chan struct{}
	logger *log.Logger
}

func main() {

	var lpath string
	var concurrency int
	var port int
	var wait int

	flag.StringVar(&lpath, "l", "attempts.log", "Log file path to write ssh logins attempt to")
	flag.IntVar(&concurrency, "c", 20, "max concurrent attempts allowed")
	flag.IntVar(&port, "p", 22, "Port to use")
	flag.IntVar(&wait, "w", 15, "Wait duration before returning response(in seconds)")

	flag.Parse()

	if concurrency < 1 {
		log.Fatal("Chosen concurrency limit is too low")
	}
	if wait < 5 {
		log.Fatal("Wait duration is too low, please choose at least 5 seconds")
	}

	f, err := os.OpenFile(lpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	v := handler{guard: make(chan struct{}, concurrency),
		logger: log.New(f, "Login Attempt: ", log.LstdFlags)}
	ssh.Handle(handle)

	p := strconv.Itoa(port)
	log.Printf("starting ssh server on port %s...\n", p)
	log.Fatal(ssh.ListenAndServe(":"+p, nil, ssh.PasswordAuth(v.passHandler)))
}

func handle(s ssh.Session) {
	io.WriteString(s, fmt.Sprintf("%s $\n", s.User()))
}

func (g *handler) passHandler(ctx ssh.Context, password string) bool {
	g.guard <- struct{}{}
	defer func() { <-g.guard }()
	time.Sleep(15 * time.Second)
	g.logger.Printf("User %s, ip: %s\n", ctx.User(), ctx.RemoteAddr())
	return false
}
