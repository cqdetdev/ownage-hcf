package command

import (
	"fmt"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/sandertv/gophertunnel/minecraft/text"
	cl "github.com/ownagepe/hcf/ownage/claim"
	"github.com/ownagepe/hcf/ownage/data"
	"github.com/ownagepe/hcf/ownage/faction"
	"github.com/ownagepe/hcf/ownage/user"
	"github.com/vasar-network/vails/lang"
	"github.com/vasar-network/vails/role"
)

type FactionCreate struct {
	Sub  create
	Name string `name:"name"`
}

func (c FactionCreate) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	u, _ := user.Lookup(s.(*player.Player))
	if _, ok := faction.Lookup(u.Faction()); ok {
		o.Errorf(text.Colourf(lang.Translatef(l, "command.faction.create.exists")))
	} else {
		if len(c.Name) < 3 {
			o.Errorf(text.Colourf(lang.Translatef(l, "command.faction.create.short")))
			return
		}
		if len(c.Name) > 15 {
			o.Errorf(text.Colourf(lang.Translatef(l, "command.faction.create.long")))
			return
		}
		f := data.NewFaction(c.Name, u.Player())
		u.SetFaction(f.Name())
		data.SaveUser(u)
		o.Printf(text.Colourf(lang.Translatef(l, "command.faction.create.success", f.Name())))
	}
}

func (FactionCreate) Allow(s cmd.Source) bool {
	return allow(s, false, role.Default{})
}

var claiming map[string]mgl64.Vec3 = make(map[string]mgl64.Vec3)

type FactionClaim struct {
	Sub claim
	Option claimOptions `optional:"start" name:"option"`
}

func (c FactionClaim) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	u, _ := user.Lookup(s.(*player.Player))

	f, _ := faction.Lookup(u.Faction())
	if f.HasClaim() {
		o.Print(lang.Translatef(l, "command.faction.claim.hasclaim"))
		return
	}

	if len(c.Option) == 0 || string(c.Option) == "start" {
		if _, ok := claiming[u.Player().Name()]; ok {
			o.Print(lang.Translatef(l, "command.faction.claim.started"))
		} else {
			pos := u.Player().Position()
			if cl.CanClaimXZ(int(pos.X()), int(pos.Z())) {
				claiming[u.Player().Name()] = pos
				o.Print(lang.Translatef(l, "command.faction.claim.finish"))
			} else {
				o.Print(lang.Translatef(l, "command.faction.claim.cannotclaim"))
			}
		}
	}

	if string(c.Option) == "finish" {
		if i, ok := claiming[u.Player().Name()]; !ok {
			o.Print(lang.Translatef(l, "command.faction.claim.nofinalize"))
		} else {
			p := u.Player().Position()
			f, _ := faction.Lookup(u.Faction())
			box := cube.Box(i.X(), 0, i.Z(), p.X(), 256, p.Z())
			area := int(box.Width() * box.Height())
			if area > 2500 {
				o.Print(lang.Translatef(l, "command.faction.claim.toolarge"))
			}
			if cl.CanClaim(
				int(box.Min().X()),
				int(box.Min().Z()),
				int(box.Max().X()),
				int(box.Max().Z()),
			) {
				f.SetClaim(box)
				go data.SaveFaction(f)
				o.Print(lang.Translatef(l, "command.faction.claim.success"))
				delete(claiming, u.Player().Name())
			} else {
				o.Print(lang.Translatef(l, "command.faction.claim.cannotclaim"))
			}
		}
	}

	if string(c.Option) == "reset" {
		if _, ok := claiming[u.Player().Name()]; !ok {
			o.Print(lang.Translatef(l, "command.faction.claim.noreset"))
		} else {
			delete(claiming, u.Player().Name())
			o.Print(lang.Translatef(l, "command.faction.claim.reset.success"))
		}
	}

	if string(c.Option) == "info" {
		if i, ok := claiming[u.Player().Name()]; !ok {
			o.Print(lang.Translatef(l, "command.faction.claim.noinfo"))
		} else {
			o.Print(lang.Translatef(l, "command.faction.claim.info", int(i.X()), int(i.Y()), int(i.Z())))
		}
	}	
}

func (FactionClaim) Allow(s cmd.Source) bool {
	return allow(s, false, role.Default{}) && (captain(s) || leader(s))
}

type FactionWho struct {
	Sub who
	Name factionName `optional:"" name:"name"`
}

func (c FactionWho) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	u, _ := user.Lookup(s.(*player.Player))
	var f *faction.Faction
	if c.Name == "" {
		fac, ok := faction.Lookup(u.Faction())
		if !ok {
			u.Player().Message("Please specifiy an arguemnt")
			return
		}
		f = fac
	} else {
		fac, ok := faction.Lookup(string(c.Name))
		if !ok {
			u.Player().Message("Doesn't exist")
			return
		}
		f = fac
	}
	var home string
	if f.Home().X() == 0 && f.Home().Y() == 0 && f.Home().Z() == 0 {
		home = "No Home"
	} else {
		home = fmt.Sprintf("%d, %d, %d", int(f.Home().X()), int(f.Home().Y()), int(f.Home().Z()))
	}
	u.Player().Message(lang.Translatef(l, "command.faction.who.display", f.Name(), home, f.Leader().Name(), "None", f.Dtr()))
}

func (FactionWho) Allow(s cmd.Source) bool {
	return allow(s, false, role.Default{})
}


type (
	create string
	claim  string
	unclaim string
	who    string
	top    string
)

func (create) SubName() string {
	return "create"
}

func (claim) SubName() string {
	return "claim"
}

func (who) SubName() string {
	return "who"
}

func (top) SubName() string {
	return "top"
}

type claimOptions string

// Type ...
func (claimOptions) Type() string {
	return "claimOptions"
}

// Options ...
func (claimOptions) Options(cmd.Source) []string {
	return []string{
		"start",
		"finish",
		"reset",
		"info",
	}
}

type factionName string

// Type ...
func (factionName) Type() string {
	return "factionName"
}

// Options ...
func (factionName) Options(cmd.Source) []string {
	var names []string
	for _, f := range faction.All() {
		names = append(names, f.Name())
	}
	return names
}