package main

import (
	"flag"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"redirect/core"
	"redirect/router"
	"redirect/utils"
	"time"
)

const (
	WebReadTimeout  = 15 * time.Second
	WebWriteTimeout = 15 * time.Second
)

var (
	group     errgroup.Group
	cmdStart  string
	cmdConfig string
)

func main() {
	flag.StringVar(&cmdConfig, "c", "config.ini", "config file path")
	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, `redirect version:%s
		Usage: redirect [-s admin|portal|<omit to start both>] [-c config_file_path]`, utils.Version)
		flag.PrintDefaults()
	}
	flag.Parse()
	err := core.InitSettings(cmdConfig)
	if err != nil {
		fmt.Println("init setting error:", err)
		return
	}

	err = core.InitializeTrans()
	utils.ExitOnError("validator translator initialize failed.", err)

	routerAdmin, err := router.InitRouteAdmin()
	utils.ExitOnError("Router admin initialize failed.", err)

	admin := &http.Server{
		Addr:         fmt.Sprintf(":%d", utils.AppConfig.AdminPort),
		Handler:      routerAdmin,
		ReadTimeout:  time.Duration(utils.AppConfig.WebReadTimeout) * time.Second,
		WriteTimeout: time.Duration(utils.AppConfig.WebWriteTimeout) * time.Second,
	}

	startAdmin(group, *admin)

	err = group.Wait()
	utils.ExitOnError("Group failed,", err)
}

func startAdmin(g errgroup.Group, server http.Server) {
	group.Go(func() error {
		log.Println("[redirect] ticker store access log")
		core.RunJob()
		return nil
	})

	group.Go(func() error {
		log.Printf("[redirect] admin starts at http://localhost:%d", utils.AppConfig.AdminPort)
		return server.ListenAndServe()
	})
}
