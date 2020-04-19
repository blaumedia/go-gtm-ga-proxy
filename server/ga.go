package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func googleAnalyticsJsHandle(w http.ResponseWriter, r *http.Request, debug bool) {
	client := &http.Client{}

	var req *http.Request
	var err error
	if debug == false {
		req, err = http.NewRequest(`GET`, `https://www.google-analytics.com/analytics.js`, nil)
	} else {
		req, err = http.NewRequest(`GET`, `https://www.google-analytics.com/analytics_debug.js`, nil)
	}

	if err != nil {
		fmt.Println(`Experienced problems on requesting analytics.js from google. Aborting.`)

		return
	}

	req.Header.Set(`User-Agent`, `GoGtmGaProxy 1.0.0; github.com/blaumedia/go-gtm-ga-proxy`)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(`Experienced problems on requesting analytics.js from google. Aborting.`)

		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(`Experienced problems on requesting analytics.js from google. Aborting.`)

		return
	}

	JsSubdirectoryWithoutLeadingSlash := []rune(JsSubdirectory)
	JsSubdirectoryWithoutLeadingSlash = JsSubdirectoryWithoutLeadingSlash[1:]

	re := regexp.MustCompile(`googletagmanager.com`)
	body = re.ReplaceAll([]byte(body), []byte(EndpointURL))

	re = regexp.MustCompile(`\/gtm.js`)
	body = re.ReplaceAll([]byte(body), []byte(string(JsSubdirectoryWithoutLeadingSlash)+GtmFilename))

	re = regexp.MustCompile(`www.google-analytics.com`)
	body = re.ReplaceAll([]byte(body), []byte(EndpointURL))

	re = regexp.MustCompile(`analytics.js`)
	body = re.ReplaceAll([]byte(body), []byte(string(JsSubdirectoryWithoutLeadingSlash)+GaFilename))

	re = regexp.MustCompile(`u\/analytics_debug.js`)
	body = re.ReplaceAll([]byte(body), []byte(string(JsSubdirectoryWithoutLeadingSlash)+GaDebugFilename))

	re = regexp.MustCompile(`\"/r\/collect`)
	body = re.ReplaceAll([]byte(body), []byte(`"`+GaCollectEndpointRedirect))

	re = regexp.MustCompile(`\"/j\/collect`)
	body = re.ReplaceAll([]byte(body), []byte(`"`+GaCollectEndpointJ))

	re = regexp.MustCompile(`\"/collect`)
	body = re.ReplaceAll([]byte(body), []byte(`"`+GaCollectEndpoint))

	setResponseHeaders(w, resp.Header)
	w.WriteHeader(resp.StatusCode)

	w.Write([]byte(body))
}

func googleAnalyticsCollectHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != `GET` && r.Method != `POST` {
		fmt.Println(`ERROR: Connection to collect endpoint through ` + r.Method + ` method. Aborting.`)
		return
	}

	client := &http.Client{}
	clientURL := ``

	switch r.URL.Path {
	case GaCollectEndpointRedirect:
		clientURL = `https://www.google-analytics.com/r/collect`
	case GaCollectEndpointJ:
		clientURL = `https://www.google-analytics.com/j/collect`
	case GaCollectEndpoint:
		fallthrough
	default:
		clientURL = `https://www.google-analytics.com/collect`
	}

	var req *http.Request
	var err error
	var bodyPayload = make(map[string]string)

	switch r.Method {
	case `GET`:
		for k, v := range r.URL.Query() {
			bodyPayload[url.QueryEscape(k)] = url.QueryEscape(v[0])
		}
	case `POST`:
		postPayloadRaw, _ := ioutil.ReadAll(r.Body)

		fmt.Println(string(postPayloadRaw))

		postPayload := strings.Split(string(postPayloadRaw), `&`)

		re := regexp.MustCompile(`(.*)=(.*)`)
		for _, item := range postPayload {
			itemParsed := re.FindStringSubmatch(item)

			bodyPayload[itemParsed[1]] = itemParsed[2]
		}
	}

	var formatPayLoad string

	for k, v := range bodyPayload {
		switch k {
		case `uip`:
		case `ua`:
		default:
			formatPayLoad = formatPayLoad + k + `=` + v + `&`
		}
	}

	formatPayLoad = formatPayLoad + `ua=` + url.QueryEscape(r.Header.Get(`User-Agent`)) + `&uip=` + url.QueryEscape(strings.Split(r.RemoteAddr, `:`)[0])

	switch r.Method {
	case `GET`:
		req, err = http.NewRequest(`GET`, clientURL+`?`+formatPayLoad, nil)
		if err != nil {
			fmt.Println(`Experienced problems on redirecting collect to google (GET). Aborting.`)

			return
		}
		fmt.Println(clientURL + `?` + formatPayLoad)
	case `POST`:
		req, err = http.NewRequest(`POST`, clientURL, bytes.NewBuffer([]byte(formatPayLoad)))
		if err != nil {
			fmt.Println(`Experienced problems on redirecting collect to google (POST). Aborting.`)

			return
		}
	}

	req.Header.Set(`User-Agent`, `GoGtmGaProxy 1.0.0; github.com/blaumedia/go-gtm-ga-proxy`)
	req.Header.Set(`Accept`, `*/*`)
	req.Header.Set(`Content-Type`, `text/plain;charset=UTF-8`)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(`Experienced problems on redirecting collect to google. Aborting.`)

		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(`Experienced problems on redirecting collect to google. Aborting.`)

		return
	}

	setResponseHeaders(w, resp.Header)
	w.WriteHeader(resp.StatusCode)

	w.Write([]byte(body))
}
