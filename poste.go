package main

import (
	"github.com/urfave/cli"
	"os"
	"poste/mailman"
	"os/signal"
	"syscall"
	"poste/consul"
	"poste/dispather"
	"poste/data"
	"poste/api"
)

func main() {

	var (
		host string
		port int
		serverType string
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
	}

	app := cli.NewApp()
	app.Name = "poste"
	app.Version = "0.0.1"
	app.Description = "a lightweight, distributed, realtime message server"

	app.Commands = []cli.Command{
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
					consul.Deregister(consul.Dispatcher, dispather.D.Host, dispather.D.Port)
					os.Exit(1)
				}()
				dispather.Serve(host, port)
				return nil
			},
		},
		{
			Name:    "mailman",
			Aliases: []string{"m"},
			Usage:   "start a mailman server",
			Flags: append(flags, cli.StringFlag{
				Name: "type",
				Value: string(mailman.WsType),
				Usage: "mailman server type : ws or tcp",
				Destination: &serverType,
			}),
			Action:  func(c *cli.Context) error {
				ch := make(chan os.Signal, 2)
				signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
				go func() {
					<-ch
					consul.Deregister(consul.Mailman, mailman.M.Host, mailman.M.Port)
					os.Exit(1)
				}()
				mailman.Serve(host, port, data.ServerType(serverType))
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
					consul.Deregister(consul.Api, api.A.Host, api.A.Port)
					os.Exit(1)
				}()
				api.Serve(host, port)
				return nil
			},
		},
	}

	app.Run(os.Args)
}