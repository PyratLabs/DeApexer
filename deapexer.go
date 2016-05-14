/*
 * ############################################################
 * #                       DEAPEXER                           #
 * ############################################################
 * # Written by Xan Manning, PyratLabs 2016. http://pyrat.io/ #
 * ############################################################
 *
 */

package main

import (
    "log"
    "net/http"
    "regexp"
)

type Configuration struct {
    ListenPort          int         `json:"listen_port"`
    DefaultSubdomain    string      `json:"default_subdomain"`
}

func deapex(w http.ResponseWriter, r *http.Request) {
    log.Print("URL: ", r.Host)
    log.Print("RequestURI: ", r.RequestURI)

    requestHost := r.Host
    requestURI := r.RequestURI

    rpPort := regexp.MustCompile(":[0-9]+")
    hostName := rpPort.ReplaceAllString(requestHost, "")

    forwardAddr := "http://www." + hostName + requestURI

    http.Redirect(w, r, forwardAddr, 301)
}

func main() {
    http.HandleFunc("/", deapex)
    initErr := http.ListenAndServe(":8080", nil)

    if initErr != nil {
        log.Fatal("ListenAndServe: ", initErr)
    }
}
