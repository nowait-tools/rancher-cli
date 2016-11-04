package cmd

import (
	"os"
	"time"
)

var (
	cattleUrl       string
	cattleAccessKey string
	cattleSecret    string

	defaultUpgradeInterval time.Duration
)

func init() {
	cattleUrl = os.Getenv("CATTLE_URL")
	cattleAccessKey = os.Getenv("CATTLE_ACCESS_KEY")
	cattleSecret = os.Getenv("CATTLE_SECRET_KEY")

	defaultUpgradeInterval = 10 * time.Second
}
