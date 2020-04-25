package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func generateGACookie() string {
	rand.Seed(time.Now().UnixNano())
	return `GA` + GaCookieVersion + `.2.` + strconv.FormatInt(int64(rand.Intn(888888888)+111111111), 10) + `.` + strconv.FormatInt(time.Now().Unix(), 10)
}

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

	if strings.ToLower(EnableServerSideGaCookies) == "true" || EnableServerSideGaCookies == "1" {
		cookieContent, errCookie := r.Cookie(ServerSideGaCookieName)

		var newCookieContent string
		var newCookieDecodedContent string

		if errCookie == nil {
			cookieDecodedContent, errCookieDecode := base64.StdEncoding.DecodeString(cookieContent.Value)

			if errCookieDecode == nil {
				newCookieContent = cookieContent.Value
				newCookieDecodedContent = string(cookieDecodedContent)
			} else {
				newCookieDecodedContent = generateGACookie()
				newCookieContent = base64.StdEncoding.EncodeToString([]byte(newCookieDecodedContent))
			}
		} else {
			if gaCookie, gaErr := r.Cookie(`_ga`); gaErr == nil {
				newCookieDecodedContent = gaCookie.Value
				newCookieContent = base64.StdEncoding.EncodeToString([]byte(newCookieDecodedContent))
			} else {
				newCookieDecodedContent = generateGACookie()
				newCookieContent = base64.StdEncoding.EncodeToString([]byte(newCookieDecodedContent))
			}
		}

		if strings.ToLower(CookieSecure) == "true" || CookieSecure == "1" {
			w.Header().Add(`Set-Cookie`, ServerSideGaCookieName+`=`+newCookieContent+`; Domain=`+CookieDomain+`; Secure; HttpOnly; SameSite=Lax; Path=/; Max-Age=63072000`)
			w.Header().Add(`Set-Cookie`, ClientSideGaCookieName+`=`+newCookieDecodedContent+`; Domain=`+CookieDomain+`; Secure; SameSite=Lax; Path=/; Max-Age=63072000`)
		} else if strings.ToLower(CookieSecure) == "false" || CookieSecure == "0" {
			w.Header().Add(`Set-Cookie`, ServerSideGaCookieName+`=`+newCookieContent+`; Domain=`+CookieDomain+`; HttpOnly; SameSite=Lax; Path=/; Max-Age=63072000`)
			w.Header().Add(`Set-Cookie`, ClientSideGaCookieName+`=`+newCookieDecodedContent+`; Domain=`+CookieDomain+`; SameSite=Lax; Path=/; Max-Age=63072000`)
		} else {
			fmt.Println(`ERROR: Environment variable 'GA_SERVER_SIDE_COOKIE_SECURE' is not true/false or 1/0. Falling back to true.`)
			w.Header().Add(`Set-Cookie`, ServerSideGaCookieName+`=`+newCookieContent+`; Domain=`+CookieDomain+`; Secure; HttpOnly; SameSite=Lax; Path=/; Max-Age=63072000`)
			w.Header().Add(`Set-Cookie`, ClientSideGaCookieName+`=`+newCookieDecodedContent+`; Domain=`+CookieDomain+`; Secure; SameSite=Lax; Path=/; Max-Age=63072000`)
		}
	}

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

	if os.Getenv(`PROXY_IP_HEADER`) != `` {
		proxyHeaderIps := strings.Split(r.Header.Get(os.Getenv(`PROXY_IP_HEADER`)), `,`)
		var proxyHeaderIPIndex = 0

		if os.Getenv(`PROXY_IP_HEADER_INDEX`) != `` {
			n, err := strconv.Atoi(os.Getenv(`PROXY_IP_HEADER_INDEX`))
			if err == nil {
				proxyHeaderIPIndex = n
			} else {
				fmt.Println(`ERROR: Couldn't convert PROXY_IP_HEADER_INDEX environment variable to type int. Falling back to 0.`)
			}
		}

		IPToRedirect := ``
		for ipIndex, ipValue := range proxyHeaderIps {
			if ipIndex == proxyHeaderIPIndex {
				IPToRedirect = strings.TrimSpace(ipValue)
			}
		}

		if IPToRedirect == `` {
			fmt.Println(`ERROR: Given PROXY_IP_HEADER_INDEX environment variable wasn't found in header range. Fallback to index 0.`)

			IPToRedirect = strings.TrimSpace(proxyHeaderIps[0])
		}

		formatPayLoad = formatPayLoad + `ua=` + url.QueryEscape(r.Header.Get(`User-Agent`)) + `&uip=` + url.QueryEscape(IPToRedirect)
	} else {
		formatPayLoad = formatPayLoad + `ua=` + url.QueryEscape(r.Header.Get(`User-Agent`)) + `&uip=` + url.QueryEscape(strings.Split(r.RemoteAddr, `:`)[0])
	}

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
