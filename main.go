package main

import (
	"log"
	"net/http"

	"github.com/RestartFU/gophig"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/ownagepe/hcf/ownage"
	"github.com/ownagepe/hcf/ownage/command"
	"github.com/sirupsen/logrus"

	_ "net/http/pprof"
	"os"
)

func main() {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{ForceColors: true}
	logger.Level = logrus.InfoLevel

	config, err := readConfig()
	if err != nil {
		logger.Fatalln(err)
	}

	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:1111", nil))
	}()

	chat.Global.Subscribe(chat.StdoutSubscriber{})

	v := ownage.New(logger, config)
	registerCommands(v)
	if err = v.Start(); err != nil {
		logger.Fatalln(err)
	}
}

// readConfig reads the configuration from the config.toml file, or creates the file if it does not yet exist.
func readConfig() (ownage.Config, error) {
	c := ownage.DefaultConfig()
	g := gophig.NewGophig("./config", "toml", 0777)

	err := g.GetConf(&c)
	if os.IsNotExist(err) {
		err = g.SetConf(c)
	}

	return c, err
}

func registerCommands(v *ownage.Ownage) {
	for _, c := range []cmd.Command{
		cmd.New("faction", "The main faction command", []string{"f"}, command.FactionCreate{}, command.FactionClaim{}, command.FactionWho{}),
		cmd.New("partneritem", "The main partner item faction", []string{"pi", "pp"}, command.PartnerItem{}),
		cmd.New("pvp", "The main PVP command", nil, command.PvpEnable{}),
		cmd.New("i", "The main utility commands for developer testing", nil, command.Kill{}),
		cmd.New("kit", "The main kit command", []string{"kits"}, command.Kit{}),
	} {
		cmd.Register(c)
	}
}
