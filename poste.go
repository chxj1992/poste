package main

import (
	"github.com/urfave/cli"
	"os"
	"poste/mailman"
	"os/signal"
	"syscall"
	"poste/dispather"
	"poste/api"
	"poste/util"
	"poste/register"
)

func main() {

	var (
		host string
		port int
	)

	var flags = []cli.Flag{
		cli.StringFlag{
			Name: "host",
			Value: "127.0.0.1",
			Usage: "server host",
			Destination: &host,
		},
		cli.IntFlag{
			Name: "port",
			Value: 0,
			Usage: "server port",
			Destination: &port,
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "display debug info",
			Destination: &util.Debugging,
		},
	}

	app := cli.NewApp()
	app.Name = "poste"
	app.Version = "0.0.1"
	app.Description = "a lightweight, distributed, realtime message server"

	app.Commands = []cli.Command{
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "init configuration to consul service",
			Action:  func(c *cli.Context) error {
				if util.Debugging {
					util.LogDebug("debug mode")
				}
				register.Init()
				return nil
			},
		},
		{
			Name:    "dispatcher",
			Aliases: []string{"d"},
			Usage:   "start a dispatcher server",
			Flags: flags,
			Action:  func(c *cli.Context) error {
				ch := make(chan os.Signal, 2)
				signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
				go func() {
					<-ch
					dispather.OnShutDown()
					os.Exit(1)
				}()

				if util.Debugging {
					util.LogDebug("debug mode")
				}
				dispather.Serve(host, port)
				return nil
			},
		},
		{
			Name:    "mailman",
			Aliases: []string{"m"},
			Usage:   "start a mailman server",
			Flags: flags,
			Action:  func(c *cli.Context) error {
				ch := make(chan os.Signal, 2)
				signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
				go func() {
					<-ch
					mailman.OnShutDown()
					os.Exit(1)
				}()

				if util.Debugging {
					util.LogDebug("debug mode")
				}
				mailman.Serve(host, port)
				return nil
			},
		},
		{
			Name:    "api",
			Aliases: []string{"a"},
			Usage:   "start an api server",
			Flags: flags,
			Action:  func(c *cli.Context) error {
				ch := make(chan os.Signal, 2)
				signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
				go func() {
					<-ch
					api.OnShutDown()
					os.Exit(1)
				}()

				if util.Debugging {
					util.LogDebug("debug mode")
				}
				api.Serve(host, port)
				return nil
			},
		},
	}

	app.Run(os.Args)
}