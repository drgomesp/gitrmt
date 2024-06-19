package main

import "log"

func init() {
	log.SetFlags(0)
	log.SetPrefix("foo: ")
}

func main() {
	log.Println("hey!")
}
