package kit

import (
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
)

const (
	DIAMOND = 1
	BARD = 3
	ARCHER = 4
	MAGE = 5
	ROGUE = 6
	MINER = 7
)

type Add = []item.Stack
type Slots = map[int]item.Stack

type Items struct {
	Slots Slots
	Add Add
}

type Armour struct {
	Helmet, Chestplate, Leggings, Boots item.Stack
}

type Kit interface {
	Name() string
	Items() Items
	Armour() Armour
	Effects() []effect.Effect
}

func Give(p *player.Player, kit Kit) {
	inv := p.Inventory()
	armr := p.Armour()

	for _, ef := range p.Effects() {
		p.AddEffect(ef)
	}
	
	a := kit.Armour()
	armr.Set(a.Helmet, a.Chestplate, a.Leggings, a.Boots)

	for slot, i := range kit.Items().Slots{
		inv.SetItem(slot, i)
	}

	for _, i := range kit.Items().Add{
		inv.AddItem(i)
	}
}

func Determine(p *player.Player) int {
	return 0
}