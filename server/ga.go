package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/tdewolff/minify/v2"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type gaSourceCodeCache struct {
	lastUpdate int64
	src        []byte
	headers    http.Header
}

var (
	srcGaCache      = gaSourceCodeCache{lastUpdate: 0}
	srcGaDebugCache = gaSourceCodeCache{lastUpdate: 0}
)

func generateGACookie() string {
	rand.Seed(time.Now().UnixNano())
	return `GA` + GaCookieVersion + `.2.` + strconv.FormatInt(int64(rand.Intn(888888888)+111111111), 10) + `.` + strconv.FormatInt(time.Now().Unix(), 10)
}

func googleAnalyticsJsHandle(w http.ResponseWriter, r *http.Request, debug bool) {
	var sourceCodeToReturn []byte
	var statusCodeToReturn int = 200
	var headersToReturn http.Header
	var usedCache bool
	var srcCachePointer *gaSourceCodeCache

	if debug {
		srcCachePointer = &srcGaDebugCache
	} else {
		srcCachePointer = &srcGaCache
	}

	if srcCachePointer.lastUpdate > (time.Now().Unix() - GaCacheTime) {
		if DebugOutput {
			fmt.Println(`Fetching GA-JS Request from Cache...`)
		}

		sourceCodeToReturn = srcCachePointer.src
		headersToReturn = srcCachePointer.headers
		usedCache = true
	} else {
		if DebugOutput {
			fmt.Println(`Refreshing Cache for GA-JS Request...`)
		}

		client := &http.Client{}

		var req *http.Request
		var err error
		if debug == false {
			req, err = http.NewRequest(`GET`, `https://www.google-analytics.com/analytics.js`, nil)

			if DebugOutput {
				fmt.Println(`REQUESTING: https://www.google-analytics.com/analytics.js`)
			}
		} else {
			req, err = http.NewRequest(`GET`, `https://www.google-analytics.com/analytics_debug.js`, nil)

			if DebugOutput {
				fmt.Println(`REQUESTING: https://www.google-analytics.com/analytics_debug.js`)
			}
		}

		if err != nil {
			fmt.Println(`Experienced problems on requesting analytics.js from google. Aborting.`)

			return
		}

		req.Header.Set(`User-Agent`, `GoGtmGaProxy `+os.Getenv(`APP_VERSION`)+`; github.com/blaumedia/go-gtm-ga-proxy`)

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
		body = re.ReplaceAll([]byte(body), []byte(EndpointURI))

		re = regexp.MustCompile(`\/gtm.js`)
		body = re.ReplaceAll([]byte(body), []byte(`/`+string(JsSubdirectoryWithoutLeadingSlash)+GtmFilename))

		re = regexp.MustCompile(`www.google-analytics.com`)
		body = re.ReplaceAll([]byte(body), []byte(EndpointURI))

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

		if JsEnableMinify {
			m := minify.New()
			m.AddCmd(`application/javascript`, exec.Command("uglifyjs"))

			var previousLengthOfJs int
			if DebugOutput {
				previousLengthOfJs = len(body)
			}

			body, err = m.Bytes(`application/javascript`, body)
			if err != nil {
				panic(err)
			}

			if DebugOutput {
				afterLengthOfJs := len(body)
				compressChange := fmt.Sprintf(`%f`, (float64(previousLengthOfJs-afterLengthOfJs)/float64(previousLengthOfJs))*float64(100))
				fmt.Println(`Compressed the Google Analytics JS File and reduced it by ` + compressChange + `%.`)
			}
		}

		if resp.StatusCode == 200 {
			srcCachePointer.headers = resp.Header
			srcCachePointer.src = body
			srcCachePointer.lastUpdate = time.Now().Unix()
		}

		headersToReturn = resp.Header
		statusCodeToReturn = resp.StatusCode
		sourceCodeToReturn = body
		usedCache = false
	}

	if EnableServerSideGaCookies {
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

		if CookieSecure {
			if DebugOutput {
				fmt.Println(`Set-Cookie: ` + ServerSideGaCookieName + `=` + newCookieContent + `; Domain=` + CookieDomain + `; Secure; HttpOnly; SameSite=Lax; Path=/; Max-Age=63072000`)
			}
			w.Header().Add(`Set-Cookie`, ServerSideGaCookieName+`=`+newCookieContent+`; Domain=`+CookieDomain+`; Secure; HttpOnly; SameSite=Lax; Path=/; Max-Age=63072000`)
			w.Header().Add(`Set-Cookie`, ClientSideGaCookieName+`=`+newCookieDecodedContent+`; Domain=`+CookieDomain+`; Secure; SameSite=Lax; Path=/; Max-Age=63072000`)
		} else {
			if DebugOutput {
				fmt.Println(`Set-Cookie: ` + ServerSideGaCookieName + `=` + newCookieContent + `; Domain=` + CookieDomain + `; HttpOnly; SameSite=Lax; Path=/; Max-Age=63072000`)
			}
			w.Header().Add(`Set-Cookie`, ServerSideGaCookieName+`=`+newCookieContent+`; Domain=`+CookieDomain+`; HttpOnly; SameSite=Lax; Path=/; Max-Age=63072000`)
			w.Header().Add(`Set-Cookie`, ClientSideGaCookieName+`=`+newCookieDecodedContent+`; Domain=`+CookieDomain+`; SameSite=Lax; Path=/; Max-Age=63072000`)
		}
	}

	setResponseHeaders(w, headersToReturn)

	if usedCache {
		w.Header().Add(`X-Cache-Hit`, `true`)
	} else {
		w.Header().Add(`X-Cache-Hit`, `false`)
	}

	for _, f := range pluginEngine.dispatcher[`after_ga_js`] {
		f(&w, r, &statusCodeToReturn, &sourceCodeToReturn)
	}

	w.WriteHeader(statusCodeToReturn)

	w.Write([]byte(sourceCodeToReturn))
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

	if DebugOutput {
		fmt.Println(`Collect-Redirect:`)
		fmt.Println(`Payload:`)
		fmt.Println(bodyPayload)
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
		if DebugOutput {
			fmt.Println(clientURL + `?` + formatPayLoad)
		}
	case `POST`:
		req, err = http.NewRequest(`POST`, clientURL, bytes.NewBuffer([]byte(formatPayLoad)))
		if err != nil {
			fmt.Println(`Experienced problems on redirecting collect to google (POST). Aborting.`)

			return
		}
		if DebugOutput {
			fmt.Println(`Format Payload:`)
			fmt.Println(formatPayLoad)
		}
	}

	req.Header.Set(`User-Agent`, `GoGtmGaProxy `+os.Getenv(`APP_VERSION`)+`; github.com/blaumedia/go-gtm-ga-proxy`)
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
