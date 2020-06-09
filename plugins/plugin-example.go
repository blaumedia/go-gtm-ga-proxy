package main

import (
	"fmt"
	"net/http"
)

// Dispatcher saves all dispatch-functions in it.
var Dispatcher map[string][]func(*http.ResponseWriter, *http.Request, *int, *[]byte)

// Main function integrates the trigger and functions into the Dispatcher map.
func Main() {
	Dispatcher = make(map[string][]func(*http.ResponseWriter, *http.Request, *int, *[]byte))

	Dispatcher[`after_ga_js`] = append(Dispatcher[`after_ga_js`], func(w *http.ResponseWriter, r *http.Request, statusCodeToReturn *int, sourceCodeToReturn *[]byte) {
		fmt.Println("PLUGIN: Executed on event 'after_ga_js'!")
		(*w).Header().Add(`X-Plugin`, `Injected`)
	})

	fmt.Println("PLUGIN: Injected!")
}

func main() {}
