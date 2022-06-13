package data

import (
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/upper/db/v4"
	"github.com/ownagepe/hcf/ownage/claim"
	"github.com/ownagepe/hcf/ownage/faction"
)

// factionData is the data structure that is used to store factions in the database.
type factionData struct {
	// Name is the name of the faction.
	Name string `db:"name"`
	// Description is the description of the faction.
	Description string `db:"description"`
	// Members are the members of the faction.
	Members []factionMemberData `db:"members"`
	// Money is the amount of money the faction has in balance.
	Money int `db:"money"`
	// DTR is the deaths till raidable a faction has.
	DTR float64 `db:"dtr"`
	// Regen is the DTR regeneration time the faction has.
	Regen time.Time `db:"regen`
	// Points is the amount of points a faction has (usually used to contest FTOP).
	Points int `db:"points"`
	// Home is the home of the faction.
	Home mgl64.Vec3 `db:"home"`
	// Claim is the claim of the faction.
	Claim factionClaimData `db:"claim"`
}

type factionMemberData struct {
	// Name is the name of the faction member
	Name string `db:"name"`
	// Role is the role of the faction member within the faction (i.e Member, Captain, Leader)
	Role int `db:"role"` // TODO: rfc?
	// Leader is whether the faction member owns the faction
	Leader bool `db:"role"` // Might not need this, it's just easier to work with
}

type factionClaimData struct {
	// MaxX is the max x-coordinate of the claim
	MaxX float64 `db:"maxX"`
	// MaxZ is the max z-coordinate of the claim
	MaxZ float64 `db:"maxZ"`
	// MaxX is the min x-coordinate of the claim
	MinX float64 `db:"minX"`
	// MaxX is the min z-coordinate of the claim
	MinZ float64 `db:"minZ"`
}


func LoadFactions() {
	var factions []factionData
	err := sess.Collection("factions").Find().All(&factions)
	if err != nil {
		panic(err)
	}
	for _, f := range factions {
		var members []*faction.Member
		for _, m := range f.Members {
			members = append(members, faction.NewMember(m.Name, m.Role, m.Leader))
		}
		c := cube.Box(f.Claim.MinX, 0, f.Claim.MinZ, f.Claim.MaxX, 256, f.Claim.MaxZ)
		faction.NewFaction(
			f.Name,
			f.Description,
			members,
			f.Money,
			f.DTR,
			time.Now(),
			f.Points,
			f.Home,
			c,
		)
		if c.Min().X() != 0 && c.Min().Z() != 0 && c.Max().X() != 0 && c.Max().Z() != 0 {
			claim.Write(&claim.Claim{
				Name: f.Name,
				Faction: true,
			}, int(c.Min().X()), int(c.Min().Z()), int(c.Max().X()), int(c.Max().Z()))
		}
		
	}
}

func NewFaction(name string, p *player.Player) *faction.Faction {
	member := faction.NewMember(p.Name(), faction.CAPTAIN, true)
	factions := sess.Collection("factions")
	f := faction.NewFaction(
		name,
		"",
		[]*faction.Member{member},
		0,
		1.01,
		time.Now(),
		0,
		mgl64.Vec3{},
		cube.Box(0, 0, 0, 0, 0, 0),
	)
	var members []factionMemberData
	for _, m := range f.Members() {
		members = append(members, factionMemberData{
			Name:  m.Name(),
			Role:   m.Role(),
			Leader: m.Leader(),
		})
	}
	factions.Insert(factionData{
		Name:       f.Name(),
		Description: f.Description(),
		Members:     members,
		Money:       f.Money(),
		DTR:         f.Dtr(),
		Points:      f.Points(),
		Home:        f.Home(),
		Claim:       factionClaimData{
			MaxX: f.Claim().Max().X(),
			MaxZ: f.Claim().Max().Z(),
			MinX: f.Claim().Min().X(),
			MinZ: f.Claim().Min().Z(),
		},
	})
	return f
}

// SaveFaction saves a *faction.Faction to the database. If an error occurs, it will be returned to the second return value.
func SaveFaction(f *faction.Faction) error {
	factions := sess.Collection("factions")

	var members []factionMemberData
	for _, m := range f.Members() {
		members = append(members, factionMemberData{
			Name:  m.Name(),
			Role:   m.Role(),
			Leader: m.Leader(),
		})
	}

	data := factionData{
		Name: f.Name(),
		Description: f.Description(),
		Members: members,
		Money: f.Money(),
		DTR: f.Dtr(),
		Regen: f.Regen(),
		Points: f.Points(),
		Home: f.Home(),
		Claim: factionClaimData{
			MaxX: f.Claim().Max().X(),
			MaxZ: f.Claim().Max().Z(),
			MinX: f.Claim().Min().X(),
			MinZ: f.Claim().Min().Z(),
		},
	}

	entry := factions.Find(db.Or(db.Cond{"name": strings.ToLower(f.Name())}))
	if ok, _ := entry.Exists(); ok {
		return entry.Update(data)
	}
	_, err := factions.Insert(data)
	return err
}