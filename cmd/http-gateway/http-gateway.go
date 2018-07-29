package main

import (
	"log"
	"github.com/urfave/cli"
	"os"
	"fmt"
	"github.com/akaumov/cube_executor"
	"github.com/akaumov/cube-http-gateway"
)

func main() {
	app := cli.NewApp()
	app.Version = "0.0.1"
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
		cli.IntFlag{
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
			Name: "endpoints-map",
			EnvVar: "GATEWAY_ENDPOINTS_MAP",
			Usage: "map url to endpoint",
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
	endpointsMap := c.String("endpoints-map")

	cube, err := cube_executor.NewCube(cube_executor.CubeConfig{
		BusPort: busPort,
		BusHost: busHost,
		Params: map[string]string {
			"jwtSecret":    jwtSecret,
			"timeoutMs":    timeoutMs,
			"endpointsMap": endpointsMap,
		},
	}, cube_http_gateway.Handler{})


	if err != nil {
		return fmt.Errorf("can't start: %v", err)
	}

	return cube.Start()
}
