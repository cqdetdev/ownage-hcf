package faction

import (
	"math"
	"sync"
	"time"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/ownagepe/hcf/ownage/claim"
	"golang.org/x/exp/maps"
)

var (
	factionMu sync.Mutex
	factions  = map[string]*Faction{}
)

// Lookup will lookup the faction by its name
func Lookup(n string) (*Faction, bool) {
	factionMu.Lock()
	defer factionMu.Unlock()
	f, ok := factions[n]
	return f, ok
}

// All will return all the factions
func All() []*Faction {
	factionMu.Lock()
	defer factionMu.Unlock()
	return maps.Values(factions)
}

type Faction struct {
	name        string
	description string
	members     []*Member
	money       int
	dtr         float64
	regen time.Time
	points      int
	home        mgl64.Vec3
	claim       cube.BBox
}

func NewFaction(
	name string,
	description string,
	members []*Member,
	money int,
	dtr float64,
	regen time.Time,
	points int,
	home mgl64.Vec3,
	c cube.BBox,
) *Faction {
	f := &Faction{
		name:        name,
		description: description,
		members:     members,
		money:       money,
		dtr:         dtr,
		regen: regen,
		points:      points,
		home:        home,
		claim:       c,
	}

	factionMu.Lock()
	factions[f.name] = f
	factionMu.Unlock()

	maxX := int(f.claim.Max().X())
	minX := int(f.claim.Min().X())
	maxZ := int(f.claim.Max().Z())
	minZ := int(f.claim.Min().Z())

	claim.Write(&claim.Claim{
		Name:    f.name,
		Faction: true,
	}, minX, minZ, maxX, maxZ)

	return f
}

func (f *Faction) Name() string {
	return f.name
}

func (f *Faction) SetName(name string) {
	f.name = name
}

func (f *Faction) Description() string {
	return f.description
}

func (f *Faction) SetDescription(description string) {
	f.description = description
}

func (f *Faction) Members() []*Member {
	return f.members
}

func (f *Faction) Member(n string) (*Member, bool) {
	for _, m := range f.members {
		if m.name == n {
			return m, true
		}
	}
	return nil, false
}

func (f *Faction) AddMember(member *Member) {
	f.members = append(f.members, member)
}

func (f *Faction) Leader() *Member {
	for _, m := range f.members {
		if m.Leader() {
			return m
		}
	}
	panic("unreachable")
}

func (f *Faction) Broadcast(msg string) {
	for _, m := range f.members {
		if u, ok := m.User(); ok {
			u.Player().Message(msg)
		}
	}
}

func (f *Faction) Money() int {
	return f.money
}

func (f *Faction) SetMoney(money int) {
	f.money = money
}

func (f *Faction) Dtr() float64 {
	return f.dtr
}

func (f *Faction) SetDtr(dtr float64) {
	f.dtr = dtr
}

func (f *Faction) MaxDtr() float64 {
	return math.Round((float64(len(f.members)) * 1.01) * 100) / 100
}

func (f *Faction) DecDtr() {
	f.dtr -= 1.0
}

func (f *Faction) Regen() time.Time {
	return f.regen
}

func (f *Faction) SetRegen() time.Time {
	return f.regen
}

func (f *Faction) Raidable() bool {
	return f.dtr < 0
}

func (f *Faction) Points() int {
	return f.points
}

func (f *Faction) SetPoints(points int) {
	f.points = points
}

func (f *Faction) Home() mgl64.Vec3 {
	return f.home
}

func (f *Faction) SetHome(home mgl64.Vec3) {
	f.home = home
}

func (f *Faction) Claim() cube.BBox {
	return f.claim
}

func (f *Faction) SetClaim(c cube.BBox) {
	f.claim = c
	claim.Write(&claim.Claim{
		Name: f.name,
		Faction: true,
	}, int(c.Min().X()), int(c.Min().Z()), int(c.Max().X()), int(c.Max().Z()))
}

func (f *Faction) HasClaim() bool {
	c := f.claim
	return c.Min().X() != 0 && c.Min().Z() != 0 && c.Max().X() != 0 && c.Max().Z() != 0
}