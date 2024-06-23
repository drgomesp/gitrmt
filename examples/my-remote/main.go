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
	r := gitrmt.NewRemote(&MyRemoteHandler{})

	if err := r.Run(os.Stdin, os.Stdout); err != nil {
		log.Fatal(err)
	}
}
