package main

import (
	"log"
	"os"

	"github.com/drgomesp/gitrmt"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("git-remote-my: ")
}

var _ gitrmt.RemoteHandler = &MyRemoteHandler{}

type MyRemoteHandler struct {
}

func (m *MyRemoteHandler) Capabilities() string {
	return "push\nfetch\n"
}

func (m *MyRemoteHandler) List(forPush bool) ([]string, error) {
	log.Printf("List(forPush=%v)\n", forPush)

	return []string{
		"a43fb037b04c8d592a71dd32b24748ca7f2f2b7a refs/heads/main",
		"@refs/heads/main HEAD",
	}, nil
}

func (m *MyRemoteHandler) Push(localRef string, remoteRef string, force bool) (string, error) {
	log.Printf(
		"Push(localRef=%s, remoteRef=%s, force=%v)\n",
		localRef,
		remoteRef,
		force,
	)

	return localRef, nil
}

func (m *MyRemoteHandler) Finish() error {
	return nil
}

func main() {
	log.Printf("args: %+v\n", os.Args)

	r := gitrmt.NewRemote(os.Stdin, os.Stdout, &MyRemoteHandler{})

	if err := r.Run(); err != nil {
		log.Fatal(err)
	}
}
