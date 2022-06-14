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

func Apply(kit vails.Kit, p *player.Player) {
	for _, eff := range p.Effects() {
		p.RemoveEffect(eff.Type())
	}

	inv := p.Inventory()
	armourInv := p.Armour()
	items := kit.Items(p)
	effects := kit.Effects(p)
	armour := kit.Armour(p)
	for slot, it := range items {
		if i, ok := inv.Item(slot); i.Empty() && ok == nil {
			_ = inv.SetItem(slot, it)
		} else {
			p.Drop(it)
		}
	}
	for _, eff := range effects {
		p.AddEffect(eff)
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

// Determine is a function to determine what kit a player has on
func Determine(p *player.Player) (vails.Kit, bool) {
	slots := p.Armour().Slots()
	helmet := slots[0]
	chest := slots[1]
	legs := slots[2]
	boots := slots[3]

	if 
		helmet.Equal(item.NewStack(item.Helmet{Tier: item.ArmourTierDiamond}, 1)) &&
		chest.Equal(item.NewStack(item.Chestplate{Tier: item.ArmourTierDiamond}, 1)) &&
		legs.Equal(item.NewStack(item.Leggings{Tier: item.ArmourTierDiamond}, 1)) &&
		boots.Equal(item.NewStack(item.Boots{Tier: item.ArmourTierDiamond}, 1)) {
		return Diamond{}, true
	}

	if 
		helmet.Equal(item.NewStack(item.Helmet{Tier: item.ArmourTierGold}, 1)) &&
		chest.Equal(item.NewStack(item.Chestplate{Tier: item.ArmourTierGold}, 1)) &&
		legs.Equal(item.NewStack(item.Leggings{Tier: item.ArmourTierGold}, 1)) &&
		boots.Equal(item.NewStack(item.Boots{Tier: item.ArmourTierGold}, 1)) {
		return Bard{}, true
	}

	return nil, false
}