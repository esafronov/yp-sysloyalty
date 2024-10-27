package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

type AppParams struct {
	RunAddress           *string `env:"RUN_ADDRESS"`            //address and port to run service
	DatabaseURI          *string `env:"DATABASE_URI"`           //database uri
	AccrualSystemAddress *string `env:"ACCRUAL_SYSTEM_ADDRESS"` //accrual system address
	AccessTokenSecret    *string `env:"ACCESS_TOKEN_SECRET"`    //access token secret for authorization
	ExpireAccessToken    *int    `env:"EXPIRE_ACCESS_TOKEN"`    //expire access token in hours
	GrabInterval         *int    `env:"GRAB_INTERVAL"`          //orders grab interval for update in seconds
	ProcessRate          *int    `env:"PROCESS_RATE"`           //desired speed of processing, orders per second
}

func parseEnv(p *AppParams) error {
	if err := env.Parse(p); err != nil {
		return err
	}
	return nil
}

func parseFlags(p *AppParams) {

	runAddressFlag := flag.String("a", "localhost:8080", "address and port to run service")
	if p.RunAddress == nil {
		p.RunAddress = runAddressFlag
	}

	databaseURIFlag := flag.String("d", "", "database uri (required)")
	if p.DatabaseURI == nil {
		p.DatabaseURI = databaseURIFlag
	}

	accrualSystemAddressFlag := flag.String("r", "http://localhost:8081", "accrual system address")
	if p.AccrualSystemAddress == nil {
		p.AccrualSystemAddress = accrualSystemAddressFlag
	}

	AccessTokenSecretFlag := flag.String("s", "1234", "access token secret for authorization")
	if p.AccessTokenSecret == nil {
		p.AccessTokenSecret = AccessTokenSecretFlag
	}

	ExpireAccessTokenFlag := flag.Int("e", 1, "expire access token in hours")
	if p.ExpireAccessToken == nil {
		p.ExpireAccessToken = ExpireAccessTokenFlag
	}

	GrabIntervalFlag := flag.Int("i", 10, "orders grab interval for update in seconds")
	if p.GrabInterval == nil {
		p.GrabInterval = GrabIntervalFlag
	}

	ProcessRateFlag := flag.Int("p", 100, "desired speed of processing, orders per second")
	if p.ProcessRate == nil {
		p.ProcessRate = ProcessRateFlag
	}

	flag.Parse()
}

func GetAppParams() (params *AppParams, err error) {
	params = &AppParams{}
	if err := parseEnv(params); err != nil {
		return nil, err
	}
	parseFlags(params)
	if *params.DatabaseURI == "" {
		flag.PrintDefaults()
		return nil, fmt.Errorf("database uri is empty")
	}
	return
}
