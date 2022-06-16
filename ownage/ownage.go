package ownage

import (
	"fmt"
	"math"
	"net/netip"
	"sync"
	"time"
	_ "unsafe"

	"github.com/df-mc/atomic"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/dustin/go-humanize"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/ownagepe/hcf/ownage/data"
	"github.com/ownagepe/hcf/ownage/kit"
	"github.com/ownagepe/hcf/ownage/module"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"github.com/sirupsen/logrus"
	_ "github.com/vasar-network/vails/command"
	_ "github.com/vasar-network/vails/console"
	"github.com/vasar-network/vails/lang"
	"github.com/vasar-network/vails/worlds"
	"golang.org/x/text/language"
)

var (
	connectionsMu sync.Mutex
	connections   = make(map[netip.Addr]int)
)

// Ownage creates a new instance of OwnageHCF.
type Ownage struct {
	log    *logrus.Logger
	config Config

	worlds *worlds.Manager
	srv    *server.Server

	mute atomic.Bool
	pvp  atomic.Bool

	c chan struct{}
}

// New creates a new instance of Ownage.
func New(log *logrus.Logger, config Config) *Ownage {
	config.Config.WorldConfig = func(def world.Config) world.Config {
		def.PortalDestination = nil // For right now
		def.ReadOnly = false
		return def
	}
	v := &Ownage{
		srv: server.New(&config.Config, log),
		c:   make(chan struct{}),

		pvp: *atomic.NewBool(true),

		log:    log,
		config: config,
	}
	v.srv.Allow(&allower{v: v})
	v.srv.CloseOnProgramEnd()
	v.srv.PlayerProvider(&Provider{srv: v.srv})
	v.srv.SetName(text.Colourf("<bold><dark-aqua>VASAR</dark-aqua></bold>") + "§8")

	v.loadLocales()

	v.loadTexts()
	v.startCrates()
	go v.startBoards()
	go v.startBroadcasts()
	go v.startPlayerBroadcasts()
	go v.startKOTH()
	return v
}

// Start starts the server.
func (v *Ownage) Start() error {
	if err := v.srv.Start(); err != nil {
		return err
	}

	data.LoadFactions()

	w := v.srv.World()
	w.Handle(&worlds.Handler{})
	w.StopWeatherCycle()
	w.SetDefaultGameMode(world.GameModeSurvival)

	for v.srv.Accept(v.accept) {
	
	}

	close(v.c)
	return nil
}

// ToggleGlobalMute ...
func (v *Ownage) ToggleGlobalMute() (old bool) {
	return v.mute.Toggle()
}

// GlobalMuted ...
func (v *Ownage) GlobalMuted() bool {
	return v.mute.Load()
}

// TogglePvP ...
func (v *Ownage) TogglePvP() (old bool) {
	return v.pvp.Toggle()
}

// PvP ...
func (v *Ownage) PvP() bool {
	return v.pvp.Load()
}

func (v *Ownage) SOTW() bool {
	end := time.Unix(v.config.MapInfo.SOTWEnd, 0)
	return time.Now().Before(end)
}

// accept accepts the incoming player
func (v *Ownage) accept(p *player.Player) {
	addr, _ := netip.ParseAddrPort(p.Addr().String())
	ip := addr.Addr()

	connectionsMu.Lock()
	if connections[ip] >= 5 {
		p.Disconnect(lang.Translatef(p.Locale(), "user.connections.limit"))
		connectionsMu.Unlock()
		return
	}
	connections[ip]++
	connectionsMu.Unlock()

	u, err := data.LoadUser(p)
	if err != nil {
		p.Disconnect(lang.Translatef(p.Locale(), "user.account.error"))
		return
	}
	_ = data.SaveUser(u) // Ensure the user is saved on join, in case this is their first join.

	k, ok := kit.Determine(u.Player())
	if ok {
		kit.ApplyEffects(k, p)
		u.SetKit(k)
	}

	p.Handle(newHandler(u, v))
	p.Armour().Inventory().Handle(InventoryHandler{p: p})
	v.welcome(p)
}

// welcome welcomes the player provided.
func (v *Ownage) welcome(p *player.Player) {
	start := time.Unix(v.config.MapInfo.SOTWSTart, 0)
	end := time.Unix(v.config.MapInfo.SOTWEnd, 0)

	var info string
	if now := time.Now(); now.Before(end) {
		info = fmt.Sprintf(
			"SOTW commenced %s and will end %s!",
			humanize.Time(start),
			humanize.Time(end),
		)
	} else {
		info = fmt.Sprintf(
			"SOTW ended %s, the map has commenced!",
			humanize.Time(end),
		)
	}

	p.ShowCoordinates()
	// Format the welcome message.
	p.Message(lang.Translatef(p.Locale(), "vasar.welcome", v.config.MapInfo.Season, v.config.MapInfo.Number, info))
	p.Messagef(lang.Translatef(p.Locale(), "koth.current", module.Current().Name()))
}

// loadLocales loads all supported locales to Vails.
func (v *Ownage) loadLocales() {
	lang.Register(language.English)
	// TODO: More languages in the future?
}

// loadTexts loads all relevant floating texts to the lobby world.
func (v *Ownage) loadTexts() {
	w := v.srv.World()
	w.AddEntity(entity.NewText(text.Colourf("<bold><aqua>OWNAGE</aqua></bold>"), mgl64.Vec3{0.5, -56, 10.5}))
	w.AddEntity(entity.NewText(text.Colourf("<grey>https://ownage.tebex.io</grey>"), mgl64.Vec3{0.5, -57, 10.5}))
	w.AddEntity(entity.NewText(text.Colourf("<grey>discord.gg/ownage</grey>"), mgl64.Vec3{0.5, -58, 10.5}))
	l := world.NewLoader(6, w, world.NopViewer{})
	l.Move(w.Spawn().Vec3Middle())
	l.Load(int(math.Round(math.Pi * 36)))
	go v.startLeaderboards()
}
