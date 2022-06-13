package partner

import (
	"time"

	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/ownagepe/hcf/ownage/enchant"
	"github.com/ownagepe/hcf/ownage/user"
	"golang.org/x/text/language"
)

type PartnerItem interface {
	world.Item

	// Run is called when the item is used.
	// user will always be defined as the person who used the partner item
	// and on may be nil, depending on if the partner item debuffs another player
	Run(user *user.User, on *user.User)

	// Name will return the formatted name of the partner item.
	Name() string

	// Meta will return the internal identifier of the partneritem
	Meta() string

	// Description will return the formatted description of the partner item.
	Description() string

	// Attack will return whether the partner item is used when attacked (this is false if it's a "usable" partner item).
	Attack() bool

	// Cooldown will return the length of the cooldown of the specific partneritem
	Cooldown() time.Duration
}

// countTarget is used as a data structure in partner items that require multiple hits to use
type countTarget struct {
	target string
	times int
}

func NewPartnerItem(pi PartnerItem, count int, l language.Tag) item.Stack {
	return item.NewStack(
		pi,
		count,
	).WithCustomName(pi.Name()).WithLore(pi.Description()).WithEnchantments(enchant.NewGlintEnchant())
}

var partnerItems []PartnerItem
var metaToItem map[string]PartnerItem

func init() {
	partnerItems = []PartnerItem{
		StrengthPowder{},
	}
	metaToItem = map[string]PartnerItem{}

	for _, pi := range partnerItems {
		metaToItem[pi.Meta()] = pi
	}
}

func ItemByMeta(meta string) (PartnerItem, bool) {
	pi, ok := metaToItem[meta]
	return pi, ok
}