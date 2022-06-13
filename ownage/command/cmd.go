package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/ownagepe/hcf/ownage/faction"
	"github.com/ownagepe/hcf/ownage/user"
	"github.com/vasar-network/vails"
	con "github.com/vasar-network/vails/console"
	"github.com/vasar-network/vails/role"
	"golang.org/x/text/language"
)

// locale returns the locale of a cmd.Source.
func locale(s cmd.Source) language.Tag {
	if p, ok := s.(*player.Player); ok {
		return p.Locale()
	}
	return language.English
}

// allow is a helper function for command allowers. It allows users to easily check for the specified roles.
func allow(src cmd.Source, console bool, roles ...vails.Role) bool {
	if _, ok := src.(con.Source); ok && console {
		return true
	}
	p, ok := src.(*player.Player)
	if !ok {
		return false
	}
	u, ok := user.Lookup(p)
	return ok && u.Roles().Contains(append(roles, role.Operator{})...)
}

// faction is a helper function for command allowers. It allows users to easily check if they are in a faction
func inFaction(src cmd.Source) bool {
	p, ok := src.(*player.Player)
	if !ok {
		return false
	}
	u, ok := user.Lookup(p)
	if !ok {
		return false
	}
	return u.HasFaction()
}

func member(src cmd.Source) bool {
	if !inFaction(src) {
		return false
	}
	p, ok := src.(*player.Player)
	if !ok {
		return false
	}
	u, ok := user.Lookup(p)
	if !ok {
		return false
	}
	_, ok = faction.Lookup(u.Faction())
	if !ok {
		return false
	}
	return true
}

func captain(src cmd.Source) bool {
	if !inFaction(src) {
		return false
	}
	p, ok := src.(*player.Player)
	if !ok {
		return false
	}
	u, ok := user.Lookup(p)
	if !ok {
		return false
	}
	f, ok := faction.Lookup(u.Faction())
	if !ok {
		return false
	}
	m, ok := f.Member(u.Player().Name())
	if !ok {
		return false
	}
	return m.Captain()
}

func leader(src cmd.Source) bool {
	if !inFaction(src) {
		return false
	}
	p, ok := src.(*player.Player)
	if !ok {
		return false
	}
	u, ok := user.Lookup(p)
	if !ok {
		return false
	}
	f, ok := faction.Lookup(u.Faction())
	if !ok {
		return false
	}
	m, ok := f.Member(u.Player().Name())
	if !ok {
		return false
	}
	return m.Leader()
}