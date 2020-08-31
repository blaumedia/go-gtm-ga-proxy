package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// Dispatcher saves all dispatch-functions in it.
var Dispatcher map[string][]func(*http.ResponseWriter, *http.Request, *int, *[]byte)

func BeforeGaJs(w *http.ResponseWriter, r *http.Request, req *http.Request) {
	fmt.Println(req.URL)
	fmt.Println(req.Method)
	fmt.Println(`BeforeGaJs()`)
}

func AfterGaJs(w *http.ResponseWriter, r *http.Request, statusCode *int, src *[]byte) {
	fmt.Println(r.RemoteAddr)
	fmt.Println(r.UserAgent())
	fmt.Println(string(*src))

	fmt.Println(`AfterGaJs()`)
}

func BeforeGaCollect(w *http.ResponseWriter, r *http.Request, req *http.Request) {
	fmt.Println(req.URL)
	fmt.Println(req.Method)

	fmt.Println(req.Body)

	// Payload Manipulation

	newPayload := bytes.NewBuffer([]byte(`v=1&_v=j83&a=917697895&t=event&ni=0&_s=1&dl=https%3A%2F%2Fblaumedia.com%2Fblog%2F&ul=de-de&de=UTF-8&dt=blaumedia%20%C2%BB%20Technik%2C%20Internet%20und%20BTS&sd=24-bit&sr=2560x1440&vp=1420x1329&je=0&ec=site_interaction&ea=scroll_depth&el=1000%25&_u=aADAAEABC~&jid=&gjid=&cid=815424026.1594577735&tid=UA-142129676-1&_gid=49508200.1598904687&gtm=2wg8j256VXSSC&cd1=GTM-56VXSSC&cd2=815424026.1594577735&z=307776779`))

	req.ContentLength = int64(newPayload.Len())
	buf := newPayload.Bytes()
	req.Body = ioutil.NopCloser(newPayload)
	req.GetBody = func() (io.ReadCloser, error) {
		r := bytes.NewReader(buf)
		return ioutil.NopCloser(r), nil
	}

	// Payload Manipulation

	fmt.Println(`BeforeGaCollect()`)
}

func AfterGaCollect(w *http.ResponseWriter, r *http.Request, statusCode *int, src *[]byte) {
	fmt.Println(`AfterGaCollect()`)
}

func BeforeGtmJs(w *http.ResponseWriter, r *http.Request, req *http.Request) {
	fmt.Println(`BeforeGtmJs()`)
}

func AfterGtmJs(w *http.ResponseWriter, r *http.Request, statusCode *int, src *[]byte) {
	fmt.Println(`AfterGtmJs()`)
}

// Main function integrates the trigger and functions into the Dispatcher map.
func Main() {
	fmt.Println("PLUGIN: Injected!")
}

// To silent the debugger we give her what she wants.
func main() {}
