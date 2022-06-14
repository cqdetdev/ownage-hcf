package kit

import (
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-mc/dragonfly/server/player"
	it "github.com/ownagepe/hcf/ownage/item"
)

type Diamond struct {}

// Items ...
func (n Diamond) Items(*player.Player) [36]item.Stack {
	items := [36]item.Stack{
		item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithEnchantments(item.NewEnchantment(enchantment.Unbreaking{}, 10)),
		item.NewStack(it.VasarPearl{}, 16),
	}
	for i := 2; i < 36; i++ {
		items[i] = item.NewStack(it.VasarPotion{Type: potion.StrongHealing()}, 1)
	}

	items[2] = item.NewStack(it.Potion{Type: potion.Swiftness()}, 1)
	items[26] = item.NewStack(it.Potion{Type: potion.Swiftness()}, 1)
	items[35] = item.NewStack(it.Potion{Type: potion.Swiftness()}, 1)
	return items
}

// Armour ...
func (Diamond) Armour(*player.Player) [4]item.Stack {
	durability := item.NewEnchantment(enchantment.Unbreaking{}, 10)
	return [4]item.Stack{
		item.NewStack(item.Helmet{Tier: item.ArmourTierDiamond}, 1).WithEnchantments(durability),
		item.NewStack(item.Chestplate{Tier: item.ArmourTierDiamond}, 1).WithEnchantments(durability),
		item.NewStack(item.Leggings{Tier: item.ArmourTierDiamond}, 1).WithEnchantments(durability),
		item.NewStack(item.Boots{Tier: item.ArmourTierDiamond}, 1).WithEnchantments(durability),
	}
}

// Effects ...
func (n Diamond) Effects(*player.Player) []effect.Effect {
	return []effect.Effect{}
}