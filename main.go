package main

import (
	"es-project/src/mqtt"
	"es-project/src/webserver"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		mqtt.Boot()
	}() // we use goroutine to run the subscription function
	wg.Add(1)
	go func() {
		webserver.WebServer()
	}()
	wg.Wait()
}
