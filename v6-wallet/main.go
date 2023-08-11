package main

import (
	"log"
	"os"
)

func main() {
	log.Default().SetFlags(log.Lshortfile | log.LstdFlags)

	f, err := os.OpenFile("output.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		panic("open log fail: " + err.Error())
	}
	defer f.Close()
	log.SetOutput(f)

	cli := NewCli()
	cli.Run()
}
