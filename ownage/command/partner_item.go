package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	pi "github.com/ownagepe/hcf/ownage/item/partner"
	"github.com/ownagepe/hcf/ownage/user"
)

type PartnerItem struct {
	Item items `name:"item"`
}

// Run ...
func (r PartnerItem) Run(s cmd.Source, o *cmd.Output) {
	u, _ := user.Lookup(s.(*player.Player))
	l := locale(s)
	switch r.Item {
	case "strength_powder":
		u.Player().SetHeldItems(pi.NewPartnerItem(pi.StrengthPowder{Locale: l}, 1, l), item.Stack{})
	}
}

type items string

// Type ...
func (items) Type() string {
	return "items"
}

// Options ...
func (items) Options(cmd.Source) []string {
	return []string{
		"strength_powder",
	}
}
