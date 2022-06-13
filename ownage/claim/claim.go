package claim

import "sync"

var (
	claimMu sync.Mutex
	claims  = map[[2]int]string{}
)

type Claim struct {
	// name is the name of the claim (it can be a faction or landmark)
	Name string
	// faction is true if the claim is a faction
	Faction bool
	// wilderness is true if the claim is in the wilderness
	Wilderness bool
	// warzone is true if the claim is in the warzone
	Warzone bool
	// koth is true if the claim is in the koth
	Koth bool
	// Spawn is true if the claim is the spawn
	Spawn bool
	// Road is true if the claim is a road
	Road bool
}

var (
	spawnMinX = -100
	spawnMinZ = -100
	spawnMaxX = 100
	spawnMaxZ = 100

	wzMinX = -400
	wzMinZ = -400
	wzMaxX = 400
	wzMaxZ = 400

	kothMinX = -500
	kothMinZ = -500
	kothMaxX = -450
	kothMaxZ = -450

	southRoadMinX = -10
	southRoadMinZ = 100
	southRoadMaxX = 10
	southRoadMaxZ = 2000
	
	eastRoadMinX = 100
	eastRoadMinZ = -10
	eastRoadMaxX = 2000
	eastRoadMaxZ = 10

	northRoadMinX = -10
	northRoadMinZ = -2000
	northRoadMaxX = 10
	northRoadMaxZ = -100

	westRoadMinX = -2000
	westRoadMinZ = -10
	westRoadMaxX = -100
	westRoadMaxZ = 10
)

func init() {
	// Use this to initialize the constant claims like KOTH, warzone, and spawn
	Write(&Claim{
		Name: "Warzone",
		Warzone: true,
	}, wzMinX, wzMinZ, wzMaxX, wzMaxZ)
	// We purposely overwrize warzone because that's how it works, spawn is another inner claim inside warzone
	Write(&Claim{
		Name:  "Spawn",
		Spawn: true,
	}, spawnMinX, spawnMinZ, spawnMaxX, spawnMaxZ)
	Write(&Claim{
		Name: "South Road",
		Road: true,
	}, southRoadMinX, southRoadMinZ, southRoadMaxX, southRoadMaxZ)
	Write(&Claim{
		Name: "East Road",
		Road: true,
	}, eastRoadMinX, eastRoadMinZ, eastRoadMaxX, eastRoadMaxZ)
	Write(&Claim{
		Name: "North Road",
		Road: true,
	}, northRoadMinX, northRoadMinZ, northRoadMaxX, northRoadMaxZ)
	Write(&Claim{
		Name: "West Road",
		Road: true,
	}, westRoadMinX, westRoadMinZ, westRoadMaxX, westRoadMaxZ)
	Write(&Claim{
		Name:    "Kingdom",
		Koth: 	true,
	}, kothMinX, kothMinZ, kothMaxX, kothMaxZ)

}

// Lookup will look up the claim by its positional hash key
func Lookup(x int, z int) (*Claim, bool) {
	claimMu.Lock()
	defer claimMu.Unlock()
	c, ok := claims[[2]int{x, z}]
	if !ok {
		return byName("")
	}
	return byName(c)
}

// Write will write the claim to the claims map
func Write(claim *Claim, minX int, minZ int, maxX int, maxZ int) {
	claimMu.Lock()
	defer claimMu.Unlock()
	for x := minX; x < maxX; x++ {
		for z := minZ; z < maxZ; z++ {
			claims[[2]int{x, z}] = claim.Name
		}
	}
}

// CanClaim will return whether a certain area can be claimed
func CanClaim(minX int, minZ int, maxX int, maxZ int) bool {
	claimMu.Lock()
	defer claimMu.Unlock()
	for x := minX; x < maxX; x++ {
		for z := minZ; z < maxZ; z++ {
			if _, ok := claims[[2]int{x, z}]; ok {
				return false
			}
		}
	}
	return true
}

// CanClaimXZ will return whether a certain point can be claimed
func CanClaimXZ(x int, z int) bool {
	claimMu.Lock()
	defer claimMu.Unlock()
	if _, ok := claims[[2]int{x, z}]; ok {
		return false
	} else {
		return true
	}
}

// byName will convert the name of the claim to a Claim struct
func byName(c string) (*Claim, bool) {
	switch c {
	case "Spawn":
		return &Claim{
			Name: c,
			Spawn: true,
		}, true
	case "Warzone":
		return &Claim{
			Name:    "Warzone",
			Warzone: true,
		}, true
	case "Kingdom", "Cove", "Citadel", "Temple", "Fortress", "Sanctuary":
		return &Claim{
			Name: c,
			Koth: true,
		}, true
		case "South Road", "West Road", "North Road", "East Road":
		return &Claim{
			Name: c,
			Road: true,
		}, true
	case "":
		return &Claim{
			Name:       "Wilderness",
			Wilderness: true,
		}, true
	default:
		return &Claim{
			Name:    c,
			Faction: true,
		}, true
	}
}
