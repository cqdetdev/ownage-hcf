package data

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ownagepe/hcf/ownage/user"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/mongo"
)

// userData is the data structure that is used to store the user data in the database.
type userData struct {
	// XUID is the XUID of the user.
	XUID string `bson:"xuid"`
	// Name is the last username of the user.
	Name string `bson:"name"`
	// DisplayName is the name displayed to other users.
	DisplayName string `bson:"display_name"`
	// DeviceID is the device ID of the last device the user logged in from.
	DeviceID string `bson:"did"`
	// SelfSignedID is the self-signed ID of the last client session of the user.
	SelfSignedID string `bson:"ssid"`
	// Address is the hashed IP address of the user.
	Address string `bson:"address"`
	// Whitelisted is true if the user is whitelisted.
	Whitelisted bool `bson:"whitelisted"`
	// FirstLogin is the time the user first logged in.
	FirstLogin time.Time `bson:"first_login"`
	// PlayTime is the duration the user has played VasarHCF for.
	PlayTime time.Duration `bson:"playtime"`

	// Money is the amount of money of the User has.
	Money int
	
	// Roles is a list of roles that the user has.
	Roles []roleData `bson:"roles"`

	// Cooldowns is a list of cooldowns that user has.
	Cooldowns []cooldownData `bson:"cooldown"`

	// Settings is a list of settings that the user has.
	Settings user.Settings `bson:"settings"`
	// Stats is a list of user statistics specific to VasarHCF.
	Stats user.Stats `bson:"stats"`

	// Timer is the timer data of a user.
	Timer timerData `bson:"timer"`

	// Faction is the name of user's faction.
	Faction string `bson:"faction"`
}

// roleData is the data structure that is used to store roles in the database.
type roleData struct {
	// Name is the name of the role.
	Name string `bson:"name"`
	// Expires is true if the role expires.
	Expires bool `bson:"expires"`
	// Expiration is the expiration time of the role.
	Expiration time.Time `bson:"expiration"`
}

// cooldownData is a data structure that is used to store cooldowns in the database.
type cooldownData struct {
	// Name is the name of the cooldown.
	Name string `bson:"name"`
	// Expires is when the cooldown expires.
	Expires time.Time `bson:"expires"`
}

// timerData is a data structure that is used to store pvp timer data in the database
type timerData struct {
	// Has is true if the user has PVP timer
	Has bool `bson:"has"`
	// Expires is the time when the PVP timer expires
	Expires time.Time `bson:"expires"`
}

// salt contains the salt that starts with used for hashing.
const salt = "ERTYUIOFGHJKNBVERFGHJK"

// sess is the Upper database session.
var sess db.Session

// init creates the Upper database connection.
func init() {
	path := os.Getenv("VASAR_DB")
	if len(path) == 0 {
		panic("vasar: mongo environment variable is not set")
	}

	b, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("vasar: %s", err))
	}

	var settings mongo.ConnectionURL
	err = json.Unmarshal(b, &settings)
	if err != nil {
		panic(fmt.Sprintf("vasr: %s", err))
	}

	sess, err = mongo.Open(settings)
	if err != nil {
		panic(fmt.Sprintf("failed to start mongo connection: %v", err))
	}
}
