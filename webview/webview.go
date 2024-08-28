// webview.go
package webview

import webview "github.com/webview/webview_go"

func CreateWebview() {
	debug := true
	w := webview.New(debug)
	defer w.Destroy()
	w.SetTitle("Angular Application")
	w.SetSize(800, 500, webview.HintNone)
	w.Navigate("http://localhost:8080")
	w.Run()
}
