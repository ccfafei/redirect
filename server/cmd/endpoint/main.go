package main

import (
	"flag"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"redirect/core"
	"redirect/router"
	"redirect/utils"
	"time"
)

var (
	cmdConfig string
	group     errgroup.Group
)

func main() {
	flag.StringVar(&cmdConfig, "c", "config.ini", "config file path")
	flag.Parse()

	err := core.InitSettings(cmdConfig)
	utils.ExitOnError("setting initialize failed.", err)

	routerEndpoint, err := router.InitRouteEndPoint()
	utils.ExitOnError("Router endpoint initialize failed.", err)

	endpoint := &http.Server{
		Addr:         fmt.Sprintf(":%d", utils.AppConfig.EndpointPort),
		Handler:      routerEndpoint,
		ReadTimeout:  time.Duration(utils.AppConfig.WebReadTimeout) * time.Second,
		WriteTimeout: time.Duration(utils.AppConfig.WebWriteTimeout) * time.Second,
	}

	endpoint.SetKeepAlivesEnabled(false)

	startEndpoint(group, *endpoint)

	err = group.Wait()
	utils.ExitOnError("Group failed,", err)
}

func startEndpoint(g errgroup.Group, server http.Server) {
	group.Go(func() error {
		log.Printf("[redirect] endpoint starts at http://localhost:%d", utils.AppConfig.EndpointPort)
		return server.ListenAndServe()
	})
}
