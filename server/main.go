package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"plugin"
	"strconv"
	"strings"
)

// GaCookieVersion is the cookie version in the _ga cookie
const GaCookieVersion = "1"

type pluginSystem struct {
	plugins    []*plugin.Plugin
	dispatcher map[string][]func(*http.ResponseWriter, *http.Request, *int, *[]byte)
}

type settingsStruct struct {
	EnableDebugOutput         bool
	EndpointURI               string
	JsSubdirectory            string
	GaCacheTime               int64
	GtmCacheTime              int64
	JsEnableMinify            bool
	GtmFilename               string
	GtmAFilename              string
	GaFilename                string
	GaDebugFilename           string
	GaPluginsDirectoryname    string
	GtagFilename              string
	GaCollectEndpoint         string
	GaCollectEndpointRedirect string
	GaCollectEndpointJ        string
	RestrictGtmIds            bool
	AllowedGtmIds             []string
	EnableServerSideGaCookies bool
	ServerSideGaCookieName    string
	CookieDomain              string
	CookieSecure              bool
	ClientSideGaCookieName    string
	PluginsEnabled            bool
	pluginEngine              pluginSystem
}

var settingsGGGP settingsStruct

func isInSlice(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func javascriptFilesHandle(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case `/` + settingsGGGP.JsSubdirectory + `/` + settingsGGGP.GtmFilename:
		googleTagManagerHandle(w, r, `default`)
	case `/` + settingsGGGP.JsSubdirectory + `/` + settingsGGGP.GtmAFilename:
		googleTagManagerHandle(w, r, `default_a`)
	case `/` + settingsGGGP.JsSubdirectory + `/` + settingsGGGP.GtagFilename:
		googleTagManagerHandle(w, r, `gtag`)
	case `/` + settingsGGGP.JsSubdirectory + `/` + settingsGGGP.GaFilename:
		googleAnalyticsJsHandle(w, r, `default`)
	case `/` + settingsGGGP.JsSubdirectory + `/` + settingsGGGP.GaDebugFilename:
		googleAnalyticsJsHandle(w, r, `debug`)
	default:
		if r.URL.Path[:len(settingsGGGP.JsSubdirectory+settingsGGGP.GaPluginsDirectoryname)+3] == `/`+settingsGGGP.JsSubdirectory+`/`+settingsGGGP.GaPluginsDirectoryname+`/` {
			googleAnalyticsJsHandle(w, r, r.URL.Path)
		} else {
			fmt.Println(`###################################`)
			fmt.Println(`404 Page accessed: ` + r.URL.Path)
			fmt.Println(`###################################`)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`Not found`))
		}
	}
	return
}

func collectHandle(w http.ResponseWriter, r *http.Request) {
	googleAnalyticsCollectHandle(w, r)
}

func setResponseHeaders(w http.ResponseWriter, headers http.Header) {
	// Looping through headers from request
	for headerName, headerValue := range headers {
		// Picking and sending relevant headers to client
		if headerName == `Age` || headerName == `Cache-Control` || headerName == `Content-Type` || headerName == `Date` || headerName == `Expires` || headerName == `Last-Modified` {
			w.Header().Set(headerName, headerValue[0])
		}
	}

	w.Header().Set(`X-Powered-By`, `GoGtmGaProxy `+os.Getenv(`APP_VERSION`))
}

func main() {
	// Check if required environment variables are set
	for _, envVar := range [...]string{
		"JS_SUBDIRECTORY",
		"GA_PLUGINS_DIRECTORYNAME",
		"GTM_FILENAME",
		"GTM_A_FILENAME",
		"GTAG_FILENAME",
		"GA_FILENAME",
		"GADEBUG_FILENAME",
		"GA_COLLECT_ENDPOINT",
		"GA_COLLECT_REDIRECT_ENDPOINT",
		"GA_COLLECT_J_ENDPOINT",
	} {
		if os.Getenv(envVar) == `` {
			fmt.Println(`ERROR: Seems the required environment variable '` + envVar + `' is missing. Exiting.`)
			os.Exit(1)
		}
	}

	settingsGGGP.EndpointURI = os.Getenv(`ENDPOINT_URI`)
	settingsGGGP.JsSubdirectory = os.Getenv(`JS_SUBDIRECTORY`)

	settingsGGGP.GaCacheTime, _ = strconv.ParseInt(os.Getenv(`GA_CACHE_TIME`), 10, 64)
	settingsGGGP.GtmCacheTime, _ = strconv.ParseInt(os.Getenv(`GTM_CACHE_TIME`), 10, 64)

	settingsGGGP.GtmFilename = os.Getenv(`GTM_FILENAME`)
	settingsGGGP.GtmAFilename = os.Getenv(`GTM_A_FILENAME`)
	settingsGGGP.GaFilename = os.Getenv(`GA_FILENAME`)
	settingsGGGP.GtagFilename = os.Getenv(`GTAG_FILENAME`)
	settingsGGGP.GaDebugFilename = os.Getenv(`GADEBUG_FILENAME`)

	settingsGGGP.GaPluginsDirectoryname = os.Getenv(`GA_PLUGINS_DIRECTORYNAME`)

	settingsGGGP.GaCollectEndpoint = os.Getenv(`GA_COLLECT_ENDPOINT`)
	settingsGGGP.GaCollectEndpointRedirect = os.Getenv(`GA_COLLECT_REDIRECT_ENDPOINT`)
	settingsGGGP.GaCollectEndpointJ = os.Getenv(`GA_COLLECT_J_ENDPOINT`)

	settingsGGGP.AllowedGtmIds = strings.Split(os.Getenv(`GTM_IDS`), `,`)

	settingsGGGP.ServerSideGaCookieName = os.Getenv(`GA_SERVER_SIDE_COOKIE`)
	settingsGGGP.CookieDomain = os.Getenv(`COOKIE_DOMAIN`)
	settingsGGGP.ClientSideGaCookieName = `_ga`

	settingsGGGP.pluginEngine = pluginSystem{
		dispatcher: make(map[string][]func(*http.ResponseWriter, *http.Request, *int, *[]byte)),
	}

	// Replace ClientSideGaCookieName if environment variable is set
	if os.Getenv(`GA_CLIENT_SIDE_COOKIE`) != `` {
		settingsGGGP.ClientSideGaCookieName = os.Getenv(`GA_CLIENT_SIDE_COOKIE`)
	}

	if strings.ToLower(os.Getenv(`ENABLE_DEBUG_OUTPUT`)) == `true` || strings.ToLower(os.Getenv(`ENABLE_DEBUG_OUTPUT`)) == `1` {
		settingsGGGP.EnableDebugOutput = true
	}

	if strings.ToLower(os.Getenv(`JS_MINIFY`)) == `true` || strings.ToLower(os.Getenv(`JS_MINIFY`)) == `1` {
		settingsGGGP.JsEnableMinify = true
	}

	if strings.ToLower(os.Getenv(`ENABLE_SERVER_SIDE_GA_COOKIES`)) == `true` || strings.ToLower(os.Getenv(`ENABLE_SERVER_SIDE_GA_COOKIES`)) == `1` {
		settingsGGGP.EnableServerSideGaCookies = true
	}

	if strings.ToLower(os.Getenv(`RESTRICT_GTM_IDS`)) == `true` || strings.ToLower(os.Getenv(`RESTRICT_GTM_IDS`)) == `1` {
		settingsGGGP.RestrictGtmIds = true
	}

	if strings.ToLower(os.Getenv(`COOKIE_SECURE`)) == `true` || strings.ToLower(os.Getenv(`COOKIE_SECURE`)) == `1` {
		settingsGGGP.CookieSecure = true
	}

	if strings.ToLower(os.Getenv(`ENABLE_PLUGINS`)) == `true` || strings.ToLower(os.Getenv(`ENABLE_PLUGINS`)) == `1` {
		settingsGGGP.PluginsEnabled = true
	}

	if settingsGGGP.PluginsEnabled {
		_, err := os.Stat(`/app/plugins`)

		if os.IsNotExist(err) == false {
			err := filepath.Walk(`./plugins/`, func(path string, info os.FileInfo, _ error) error {
				if info.IsDir() == false && info.Name()[len(info.Name())-2:] == "so" {
					fmt.Println(`/app/plugins/` + info.Name())
					p, err := plugin.Open(`/app/plugins/` + info.Name())

					if err != nil {
						fmt.Println(`ERROR: Failure on opening plugin!`)
						panic(err)
					}

					mainFunc, err := p.Lookup(`Main`)

					if err != nil {
						fmt.Println(`ERROR: Failure on starting main routine of plugin!`)
						panic(err)
					}

					mainFunc.(func())()

					pluginDispatcher, err := p.Lookup(`Dispatcher`)

					if err != nil {
						fmt.Println(`ERROR: Failure on reading Dispatcher of plugin!`)
						panic(err)
					}

					for k, v := range *pluginDispatcher.(*map[string][]func(*http.ResponseWriter, *http.Request, *int, *[]byte)) {
						for _, f := range v {
							settingsGGGP.pluginEngine.dispatcher[k] = append(settingsGGGP.pluginEngine.dispatcher[k], f)
						}
					}
				}
				return nil
			})

			if err != nil {
				fmt.Println(`ERROR: Failure on opening plugin!`)
				panic(err)
			}
		} else {
			fmt.Println(`ERROR: plugins-directory doesn't exist! Did you pick the right docker image that is made for plugin usage or did you forget to disable ENABLE_PLUGINS environment variable while switching images?`)
			os.Exit(1)
		}
	}

	http.HandleFunc(`/`+settingsGGGP.JsSubdirectory+`/`, javascriptFilesHandle)

	http.HandleFunc(settingsGGGP.GaCollectEndpoint, collectHandle)
	http.HandleFunc(settingsGGGP.GaCollectEndpointRedirect, collectHandle)
	http.HandleFunc(settingsGGGP.GaCollectEndpointJ, collectHandle)

	if err := http.ListenAndServe(`:8080`, nil); err != nil {
		panic(err)
	}
}
