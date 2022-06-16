package ownage

import (
	"math"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"github.com/vasar-network/vails/lang"
	"golang.org/x/text/language"
)

func (v *Ownage) startCrates() {
	w := v.srv.World()
	w.AddEntity(entity.NewText(text.Colourf("<bold><purple>KOTH Crate</purple></bold>"), mgl64.Vec3{10, -56, 10}))
	w.AddEntity(entity.NewText(lang.Translatef(language.English, "crate.information"), mgl64.Vec3{10, -57, 10}))
	l := world.NewLoader(6, w, world.NopViewer{})
	l.Move(w.Spawn().Vec3Middle())
	l.Load(int(math.Round(math.Pi * 36)))
	
	w.SetBlock(cube.Pos{10, -59, 10}, block.NewChest(), &world.SetOpts{})
}