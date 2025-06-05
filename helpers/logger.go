package helpers

import (
	"log"
	"os"
)

var Logger = log.New(os.Stdout, "[debug] ",log.LstdFlags)

func Init() {
	f, err := os.OpenFile("debug.log", os.O_APPEND| os.O_CREATE| os.O_WRONLY,0644)
	if err != nil {
		panic(err)
	}
	Logger.SetOutput(f)
}