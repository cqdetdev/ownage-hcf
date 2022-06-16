package ownage

import (
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/ownagepe/hcf/ownage/kit"
	"github.com/ownagepe/hcf/ownage/user"
)

type InventoryHandler struct{
	p *player.Player
}

func (i InventoryHandler) HandleTake(ctx *event.Context, slot int, it item.Stack)  {
	handleRemove(i.p)
}
func (i InventoryHandler) HandlePlace(ctx *event.Context, slot int, it item.Stack) {
	handlePlace(i.p, it)

}
func (i InventoryHandler) HandleDrop(ctx *event.Context, slot int, it item.Stack)  {
	handleRemove(i.p)
}

func handlePlace(p *player.Player, it item.Stack) {
	if _, ok := it.Item().(item.Armour); ok {
		fakeContainer := *p.Armour()
		fakeContainer.Inventory().AddItem(it)
		k, ok := kit.Determine(p)
		if !ok {
			p.Message("U don't have a kit on, clearing effects")
			return
		}
		if k.Type() == kit.DIAMOND {
			p.Message("U got diamond kit on")
		}

		if k.Type() == kit.BARD {
			p.Message("U got bard on")
		}
		
		u, _ := user.Lookup(p)
		u.SetKit(k)
		kit.ApplyEffects(k, p)
	}
}

func handleRemove(p *player.Player) {
	u, _ := user.Lookup(p)
	if u.Kit() != nil {
		kit.RemoveEffects(u.Kit(), p)
		p.Messagef("You removed a kit: %d", u.Kit().Type())
	}
}