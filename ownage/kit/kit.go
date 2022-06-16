package kit

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/vasar-network/vails"
)

const (
	DIAMOND = 1
	BARD = 3
	ARCHER = 4
	MAGE = 5
	ROGUE = 6
	MINER = 7
)

type Kit interface {
	vails.Kit

	Type() int
}

func Apply(kit Kit, p *player.Player) {
	for _, eff := range p.Effects() {
		p.RemoveEffect(eff.Type())
	}

	inv := p.Inventory()
	armourInv := p.Armour()
	items := kit.Items(p)
	armour := kit.Armour(p)
	for slot, it := range items {
		if i, ok := inv.Item(slot); i.Empty() && ok == nil {
			_ = inv.SetItem(slot, it)
		} else {
			p.Drop(it)
		}
	}

	for i, arm := range armour {
		switch i {
		case 0:
			if armourInv.Helmet().Empty() {
				armourInv.SetHelmet(arm)
			} else {
				p.Drop(arm)
			}
		case 1:
			if armourInv.Chestplate().Empty() {
				armourInv.SetChestplate(arm)
			} else {
				p.Drop(arm)
			}
		case 2:
			if armourInv.Leggings().Empty() {
				armourInv.SetLeggings(arm)
			} else {
				p.Drop(arm)
			}
		case 3:
			if armourInv.Boots().Empty() {
				armourInv.SetBoots(arm)
			} else {
				p.Drop(arm)
			}
		}
	}
}

// ApplyEffects is a function that applies ONLY the kit effect to a player (this is useful when a player crafts their own set of the armor and wants the effects, like archer or miner)
func ApplyEffects(kit vails.Kit, p *player.Player) {
	for _, eff := range p.Effects() {
		p.RemoveEffect(eff.Type())
	}
	effects := kit.Effects(p)
	for _, eff := range effects {
		p.AddEffect(eff)
	}
}

func RemoveEffects(kit vails.Kit, p *player.Player) {
	effects := kit.Effects(p)
	for _, eff := range effects {
		p.RemoveEffect(eff.Type())
	}
}

// Determine is a function to determine what kit a player has on
func Determine(p *player.Player) (Kit, bool) {
	slots := p.Armour().Slots()
	helmet := slots[0]
	chest := slots[1]
	legs := slots[2]
	boots := slots[3]

	if 
		(helmet.Item() == item.Helmet{Tier: item.ArmourTierDiamond}) &&
		(chest.Item() == item.Chestplate{Tier: item.ArmourTierDiamond}) &&
		(legs.Item() == item.Leggings{Tier: item.ArmourTierDiamond}) &&
		(boots.Item() == item.Boots{Tier: item.ArmourTierDiamond}) {
		return Diamond{}, true
	}

	if 
		(helmet.Item() == item.Helmet{Tier: item.ArmourTierGold}) &&
		(chest.Item() == item.Chestplate{Tier: item.ArmourTierGold}) &&
		(legs.Item() == item.Leggings{Tier: item.ArmourTierGold}) &&
		(boots.Item() == item.Boots{Tier: item.ArmourTierGold}) {
		return Bard{}, true
	}

	return nil, false
}