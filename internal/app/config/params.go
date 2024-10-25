package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

type AppParams struct {
	RunAddress           *string `env:"RUN_ADDRESS"`
	DatabaseURI          *string `env:"DATABASE_URI"`
	AccrualSystemAddress *string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	AccessTokenSecret    *string `env:"ACCESS_TOKEN_SECRET"`
	ExpireAccessToken    *int    `env:"EXPIRE_ACCESS_TOKEN"`
	GrabInterval         *int    `env:"GRAB_INTERVAL"`
	ProcessRate          *int    `env:"PROCESS_RATE"`
	//PollWorkerCount      *int    `env:"POLL_WORKER_COUNT"`
	//	RefreshTokenSsecret   *string `env:"REFRESH_TOKEN_SECRET"`
	//	ExpireRefreshToken   *int    `env:"EXPIRE_REFRESH_TOKEN"`
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

	AccessTokenSecretFlag := flag.String("s", "1234", "access token secret")
	if p.AccessTokenSecret == nil {
		p.AccessTokenSecret = AccessTokenSecretFlag
	}

	/*	RefreshTokenSecretFlag := flag.String("refresh_secret", "1234", "refresh token secret")
		if p.RefreshTokenSecret == nil {
			p.RefreshTokenSecret = RefreshTokenSecretFlag
		}
	*/

	ExpireAccessTokenFlag := flag.Int("e", 1, "expire access token in hours")
	if p.ExpireAccessToken == nil {
		p.ExpireAccessToken = ExpireAccessTokenFlag
	}

	/*
		ExpireRefreshTokenFlag := flag.Int("expire_refresh", 3, "expire refresh token (hours)")
		if p.ExpireRefreshToken == nil {
			p.ExpireRefreshToken = ExpireRefreshTokenFlag
		}
	*/

	GrabIntervalFlag := flag.Int("i", 10, "grab orders for update interval in seconds")
	if p.GrabInterval == nil {
		p.GrabInterval = GrabIntervalFlag
	}

	ProcessRateFlag := flag.Int("p", 100, "speed of processing, orders per second")
	if p.ProcessRate == nil {
		p.ProcessRate = ProcessRateFlag
	}

	/*	PollWorkerCountFlag := flag.Int("worker_count", 2, "poll status worker count")
		if p.PollWorkerCount == nil {
			p.PollWorkerCount = PollWorkerCountFlag
		}
	*/
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
