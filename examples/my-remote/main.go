package main

import (
	"log"
	"os"

	"github.com/drgomesp/gitrmt"
)

var _ gitrmt.Handler = &MyRemoteHandler{}

type MyRemoteHandler struct {
}

func (m *MyRemoteHandler) Capabilities() string {
	return "push\nfetch\n"
}

func (m *MyRemoteHandler) List(forPush bool) ([]string, error) {
	log.Printf("List(forPush=%v)\n", forPush)

	return []string{
		"eee027b5728483d8089700e8fc3e7b9e14a3b5c4 refs/heads/main",
		"@refs/heads/main HEAD",
	}, nil
}

func (m *MyRemoteHandler) Push(localRef string, remoteRef string, force bool) (string, error) {
	log.Printf(
		"Push(local=%s, remote=%s, force=%v)\n",
		localRef,
		remoteRef,
		force,
	)

	return localRef, nil
}

func (m *MyRemoteHandler) Finish() error {
	return nil
}

func init() {
	log.SetFlags(0)
	log.SetPrefix("git-remote-my: ")
}

func main() {
	r := gitrmt.NewRemote(&MyRemoteHandler{})

	if err := r.Run(os.Stdin, os.Stdout); err != nil {
		log.Fatal(err)
	}
}
