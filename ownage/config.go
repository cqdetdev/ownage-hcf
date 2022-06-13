package ownage

import (
	"time"

	"github.com/df-mc/dragonfly/server"
)

// Config is an extension of the Dragonfly server config to include fields specific to Ownage.
type Config struct {
	server.Config
	// Ownage contains fields specific to Ownage.
	Ownage struct {
		// Whitelisted is true if the server is whitelisted.
		Whitelisted bool
	}
	MapInfo struct {
		// Season indicates the map number and name
		Season string
		// Number is the season number
		Number int
		// FactionSize is the maximum faction size of the map
		FactionSize int
		// Allying is whether allying is allowed
		Allying     bool
		// Protection is the max kit protection level
		Protection  int
		// Sharpness is the max kit sharpness level
		Sharpness   int
		// SOTWStart is the start of SOTW
		SOTWSTart int64
		// SOTWEnd is the end of SOTW
		SOTWEnd     int64
	}
	// Sentry contains fields used for Sentry.
	Sentry struct {
		// Release is the release name.
		Release string
		// Dsn is the Sentry Dsn.
		Dsn string
	}
}

// DefaultConfig returns a default config for the server.
func DefaultConfig() Config {
	c := Config{}
	c.Config = server.DefaultConfig()
	c.Ownage.Whitelisted = true
	c.Ownage.Whitelisted = false
	c.MapInfo.Season = "Gaia"
	c.MapInfo.Number = 1
	c.MapInfo.FactionSize = 5
	c.MapInfo.Allying = false
	c.MapInfo.Protection = 2
	c.MapInfo.Sharpness = 2
	c.MapInfo.SOTWSTart = time.Now().Unix()
	c.MapInfo.SOTWEnd = time.Now().Add(time.Hour * 1).Unix()

	return c
}
