package main

import (
	"log"
	"github.com/urfave/cli"
	"os"
	"fmt"
	"github.com/akaumov/cube-executor"
	"github.com/akaumov/cube-http-gateway"
)

func main() {
	app := cli.NewApp()
	app.Version = "0.0.5"
	app.Action = runServer
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "bus-host",
			EnvVar: "GATEWAY_BUS_HOST",
			Usage:  "bus host",
		},
		cli.IntFlag{
			Name:   "bus-port",
			EnvVar: "GATEWAY_BUS_PORT",
			Usage:  "bus port",
		},
		cli.StringFlag{
			Name:   "jwt-secret",
			EnvVar: "GATEWAY_JWT_SECRET",
			Usage:  "jwt secret",
		},
		cli.IntFlag{
			Name:   "timeout",
			EnvVar: "GATEWAY_TIMEOUT",
			Usage:  "requests timeout ms",
		},
		cli.StringFlag{
			Name:   "endpoints-map",
			EnvVar: "GATEWAY_ENDPOINTS_MAP",
			Usage:  "map url to endpoint",
		},
		cli.BoolTFlag{
			Name:   "only-authorized-requests",
			EnvVar: "GATEWAY_ONLY_AUTHORIZED_REQUESTS",
			Usage:  "handle only authorized requests",
		},
		cli.BoolFlag{
			Name:   "dev",
			EnvVar: "GATEWAY_DEV",
			Usage:  "log all requests",
		},
		cli.StringFlag{
			Name:   "port",
			EnvVar: "GATEWAY_PORT",
			Usage:  "port to listen",
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func runServer(c *cli.Context) error {

	busHost := c.String("bus-host")
	if busHost == "" {
		return fmt.Errorf("bus host is required")
	}

	busPort := c.Int("bus-port")
	if busPort == 0 {
		return fmt.Errorf("db port is required")
	}

	jwtSecret := c.String("jwt-secret")
	if jwtSecret == "" {
		return fmt.Errorf("jwt secret is required")
	}

	timeoutMs := c.String("timeout")
	port := c.String("port")
	endpointsMap := c.String("endpoints-map")

	onlyAuthorizedRequests := "true"
	if c.Bool("only-authorized-requests") {
		onlyAuthorizedRequests = "true"
	} else {
		onlyAuthorizedRequests = "false"
	}

	dev := "false"
	if c.Bool("dev") {
		dev = "true"
	} else {
		dev = "false"
	}

	cube, err := cube_executor.NewCube(cube_executor.CubeConfig{
		BusPort: busPort,
		BusHost: busHost,
		Params: map[string]string{
			"jwtSecret":              jwtSecret,
			"timeoutMs":              timeoutMs,
			"endpointsMap":           endpointsMap,
			"onlyAuthorizedRequests": onlyAuthorizedRequests,
			"dev":                    dev,
			"port":                   port,
		},
	}, &cube_http_gateway.Handler{})

	if err != nil {
		return fmt.Errorf("can't start: %v", err)
	}

	return cube.Start()
}
