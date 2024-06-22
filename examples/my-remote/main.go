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

func main() {
	log.Printf("args: %+v\n", os.Args)

	r := gitrmt.NewRemote(os.Stdin, os.Stdout, &MyRemoteHandler{})

	if err := r.Run(); err != nil {
		log.Fatal(err)
	}
}
