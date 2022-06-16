package ownage

import (
	"net/netip"
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/entity/damage"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/player/skin"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/ownagepe/hcf/ownage/data"
	"github.com/ownagepe/hcf/ownage/module"
	"github.com/ownagepe/hcf/ownage/user"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"github.com/vasar-network/vails"
	h "github.com/vasar-network/vails/handler"
	"github.com/vasar-network/vails/lang"
	"github.com/vasar-network/vails/role"
)

// handler is a base handler that forwards all events to their respective modules.
// TODO: Merge this with user.User or do something similar.
type handler struct {
	player.NopHandler

	u   *user.User
	srv *Ownage

	c  *module.Claim
	cl *h.Click
	co *module.Combat
	cr *module.Crate
	f *module.Faction
	i  *module.Inventory
	r  *module.Rods
	p  *module.Protection
	s  *module.Settings
	m  *module.Moderation
	pi *module.PartnerItem
	k *module.Koth
}

var (
	// tlds is a list of top level domains used for checking for advertisements.
	tlds = []string{".me", ".club", "www.", ".com", ".net", ".gg", ".cc", ".net", ".co", ".co.uk", ".ddns", ".ddns.net", ".cf", ".live", ".ml", ".gov", "http://", "https://", ",club", "www,", ",com", ",cc", ",net", ",gg", ",co", ",couk", ",ddns", ",ddns.net", ",cf", ",live", ",ml", ",gov", ",http://", "https://", "gg/"}
	// emojis is a map between emojis and their unicode representation.
	emojis = map[string]string{
		":l:":     "\uE107",
		":skull:": "\uE105",
		":fire:":  "\uE108",
		":eyes:":  "\uE109",
	}
)

// newHandler ...
func newHandler(u *user.User, v *Ownage) *handler {
	p := u.Player()
	cl := h.NewClick(func(cps int) {
		for _, w := range u.ClickWatchers() {
			w.Player().SendTip(text.Colourf("<white>%v CPS</white>", cps))
		}
		if u.WatchingClicks() != nil {
			return
		}
		if u.Settings().Display.CPS {
			p.SendTip(text.Colourf("<white>%v CPS</white>", cps))
		}
	})
	ha := &handler{
		srv: v,
		u:   u,

		c:  module.NewClaim(u),
		cl: cl,
		co: module.NewCombat(u),
		cr: module.NewCrate(u),
		f: module.NewFaction(u),
		i:  module.NewInventory(u),
		m:  module.NewModeration(u),
		p:  module.NewProtection(p),
		r:  module.NewRods(u),
		s:  module.NewSettings(u),
		pi: module.NewPartnerItem(u),
		k: module.NewKoth(u),
	}
	vails.SetHandlers(p, cl)
	ha.m.HandleJoin()
	ha.s.HandleJoin()

	s := p.Skin()
	if s.Persona {
		p.SetSkin(s)
	} else if percent, err := searchTransparency(s); err != nil || percent >= 0.05 || s.Persona {
		p.SetSkin(s)
	}
	return ha
}

// HandleChat ...
func (h *handler) HandleChat(ctx *event.Context, message *string) {
	p := h.u.Player()
	if h.srv.GlobalMuted() {
		p.Message(lang.Translatef(p.Locale(), "user.globalmuted"))
		ctx.Cancel()
		return
	}

	if msg := strings.TrimSpace(*message); len(msg) > 0 {
		for _, word := range strings.Split(msg, " ") {
			if emoji, ok := emojis[strings.ToLower(word)]; ok {
				msg = strings.ReplaceAll(msg, word, emoji)
			}
		}
		formatted := h.u.Roles().Highest().Chat(h.u.DisplayName(), msg)
		if _, ok := h.u.Roles().Highest().(role.Plus); ok {
			formatted = strings.ReplaceAll(formatted, "§0", h.u.Settings().Advanced.VasarPlusColour)
		}
		if !h.u.CanSendMessage() {
			p.Message(formatted)
			ctx.Cancel()
			return
		}
		if !h.u.Roles().Contains(role.Operator{}) {
			for _, tld := range tlds {
				if strings.Contains(strings.ToLower(msg), tld) {
					p.Message(formatted)
					ctx.Cancel()
					return
				}
			}
		}
		_, _ = chat.Global.WriteString(formatted)
		h.u.RenewLastMessage()
	}
	ctx.Cancel()
}

// HandleSkinChange ...
func (h *handler) HandleSkinChange(ctx *event.Context, s *skin.Skin) {
	if s.Persona {
		*s = steve
	} else if percent, err := searchTransparency(*s); err != nil || percent >= 0.05 {
		*s = steve
	}
	h.s.HandleSkinChange(ctx, s)
}

// HandleCommandExecution ...
func (h *handler) HandleCommandExecution(ctx *event.Context, command cmd.Command, args []string) {
	h.m.HandleCommandExecution(ctx, command, args)
}

// HandleFoodLoss ...
func (h *handler) HandleFoodLoss(ctx *event.Context, from, to int) {
	h.p.HandleFoodLoss(ctx, from, to)
}

// HandleHurt ...
func (h *handler) HandleHurt(ctx *event.Context, dmg *float64, imm *time.Duration, src damage.Source) {
	h.m.HandleHurt(ctx, dmg, src)
	h.p.HandleHurt(ctx, dmg, src)
	h.co.HandleHurt(ctx, dmg, src)

	if (h.u.Player().Health()-h.u.Player().FinalDamageFrom(*dmg, src) <= 0 || (src == damage.SourceVoid{})) && !ctx.Cancelled() {
		ctx.Cancel()
		h.HandleDeath(src)
	}
}

// HandleDeath ...
func (h *handler) HandleDeath(source damage.Source) {
	h.s.HandleDeath(source)
	h.co.HandleDeath(source)
	h.f.HandleDeath(source)
}

// HandleBlockPlace ...
func (h *handler) HandleBlockPlace(ctx *event.Context, pos cube.Pos, b world.Block) {
	h.p.HandleBlockPlace(ctx, pos, b)
}

// HandleBlockBreak ...
func (h *handler) HandleBlockBreak(ctx *event.Context, pos cube.Pos, drops *[]item.Stack) {
	h.p.HandleBlockBreak(ctx, pos, drops)
}

// HandleItemUseOnBlock ...
func (h *handler) HandleItemUseOnBlock(ctx *event.Context, pos cube.Pos, face cube.Face, clickPos mgl64.Vec3) {
	h.i.HandleItemUseOnBlock(ctx, pos, face, clickPos)
	h.cr.HandleItemUseOnBlock(ctx, pos, face, clickPos)
}

// HandleItemUse ...
func (h *handler) HandleItemUse(ctx *event.Context) {
	h.m.HandleItemUse(ctx)
	h.i.HandleItemUse(ctx)
	h.r.HandleItemUse(ctx)
	h.pi.HandleItemUse(ctx)
}

// HandleAttackEntity ...
func (h *handler) HandleAttackEntity(ctx *event.Context, e world.Entity, force, height *float64, critical *bool) {
	if !h.srv.PvP() {
		h.u.Player().Message(lang.Translatef(h.u.Player().Locale(), "pvp.disabled"))
		ctx.Cancel()
	}
	h.m.HandleAttackEntity(ctx, e, force, height, critical)
	h.cl.HandleAttackEntity(ctx, e, force, height, critical)
	h.p.HandleAttackEntity(ctx, e, force, height, critical)
}

// HandlePunchAir ...
func (h *handler) HandlePunchAir(ctx *event.Context) {
	h.cl.HandlePunchAir(ctx)
	if pl := h.u.Player(); pl.World() == h.srv.srv.World() && !pl.OnGround() {
		h.u.Launch()
	}
}

// HandleMove ...
func (h *handler) HandleMove(ctx *event.Context, newPos mgl64.Vec3, newYaw, newPitch float64) {
	h.s.HandleMove(ctx, newPos, newYaw, newPitch)
	h.c.HandleMove(ctx, newPos, newYaw, newPitch)
	h.k.HandleMove(ctx, newPos, newYaw, newPitch)
}

// HandleItemDrop ...
func (h *handler) HandleItemDrop(ctx *event.Context, e *entity.Item) {
}

// HandleQuit ...
func (h *handler) HandleQuit() {
	h.u.StopWatchingClicks()
	for _, w := range h.u.ClickWatchers() {
		w.StopWatchingClicks()
	}

	addr, _ := netip.ParseAddrPort(h.u.Address().String())
	ip := addr.Addr()

	connectionsMu.Lock()
	if connections[ip] <= 1 {
		delete(connections, ip)
	} else {
		connections[ip]--
	}
	connectionsMu.Unlock()

	h.u.Close()
	h.co.HandleQuit()

	_ = data.SaveUser(h.u)
	vails.CloseHandlers(h.u.Player())
	h.srv.srv.World().SetPlayerSpawn(h.u.Player().UUID(), cube.PosFromVec3(h.u.Player().Position()))
}

// HandleItemDamage ...
func (h *handler) HandleItemDamage(ctx *event.Context, i item.Stack, d int) {
}
