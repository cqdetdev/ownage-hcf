package enchant

import (
	"github.com/df-mc/dragonfly/server/item"
)

type GlintEnchant struct{}

func (e GlintEnchant) Level() int {
	return 1
}

func (e GlintEnchant) Name() string {
	return ""
}

func (e GlintEnchant) MaxLevel() int {
	return 1
}

func (e GlintEnchant) CompatibleWith(s item.Stack) bool {
	return true
}

func NewGlintEnchant() item.Enchantment {
	return item.NewEnchantment(GlintEnchant{}, 1)
}
