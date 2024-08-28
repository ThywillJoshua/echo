// main.go
package main

import (
	"log"

	"github.com/thywilljoshua/echo/gitops"
	"github.com/thywilljoshua/echo/server"
	"github.com/thywilljoshua/echo/webview"
)

var message string

func main() {
	if err := gitops.StartCommit(&message); err != nil {
		log.Fatalf("Failed to start commit: %v", err)
	}
	go server.StartHTTPServer(&message)
	webview.CreateWebview()
}
