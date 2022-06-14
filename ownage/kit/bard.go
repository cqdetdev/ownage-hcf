package kit

import (
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/ownagepe/hcf/ownage/enchant"
	it "github.com/ownagepe/hcf/ownage/item"
)

type Bard struct {}

// Items ...
func (n Bard) Items(*player.Player) [36]item.Stack {
	items := [36]item.Stack{
		item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithEnchantments(item.NewEnchantment(enchantment.Unbreaking{}, 10)),
		item.NewStack(it.VasarPearl{}, 16),
	}
	for i := 2; i < 36; i++ {
		items[i] = item.NewStack(it.VasarPotion{Type: potion.StrongHealing()}, 1)
	}

	items[9] = item.NewStack(item.Sugar{}, 64).WithValue("bard", true).WithEnchantments(enchant.NewGlintEnchant())
	items[10] = item.NewStack(item.IronIngot{}, 64).WithValue("bard", true).WithEnchantments(enchant.NewGlintEnchant())
	items[11] = item.NewStack(item.BlazePowder{}, 64).WithValue("bard", true).WithEnchantments(enchant.NewGlintEnchant())
	items[12] = item.NewStack(item.Feather{}, 64).WithValue("bard", true).WithEnchantments(enchant.NewGlintEnchant())
	return items
}

// Armour ...
func (Bard) Armour(*player.Player) [4]item.Stack {
	durability := item.NewEnchantment(enchantment.Unbreaking{}, 10)
	return [4]item.Stack{
		item.NewStack(item.Helmet{Tier: item.ArmourTierGold}, 1).WithEnchantments(durability),
		item.NewStack(item.Chestplate{Tier: item.ArmourTierGold}, 1).WithEnchantments(durability),
		item.NewStack(item.Leggings{Tier: item.ArmourTierGold}, 1).WithEnchantments(durability),
		item.NewStack(item.Boots{Tier: item.ArmourTierGold}, 1).WithEnchantments(durability),
	}
}

// Effects ...
func (n Bard) Effects(*player.Player) []effect.Effect {
	return []effect.Effect{
	}
}