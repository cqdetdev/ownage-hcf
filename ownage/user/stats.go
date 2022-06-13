package user

// Stats contains all the stats of a user.
type Stats struct {
	// Kills is the amount of players the user has killed.
	Kills uint32 `bson:"kills"`
	// Deaths is the amount of times the user has died.
	Deaths uint32 `bson:"deaths"`

	// KillStreak is the current streak of kills the user has without dying.
	KillStreak uint32 `bson:"kill_streak"`
	// BestKillStreak is the highest kill-streak the user has ever gotten.
	BestKillStreak uint32 `bson:"best_kill_streak"`
}


// DefaultStats returns the default stats of a user.
func DefaultStats() Stats {
	s := Stats{}
	return s
}
