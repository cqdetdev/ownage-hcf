package faction

import "github.com/ownagepe/hcf/ownage/user"

const (
	DEFAULT = iota
	CAPTAIN
)

type Member struct {
	name   string
	role   int // TODO: Again rfc?
	leader bool
}

// Name ...
func (m *Member) Name() string {
	return m.name
}

// Role ...
func (m *Member) Role() int {
	return m.role
}

// Leader ...
func (m *Member) Leader() bool {
	return m.leader
}

func (m *Member) User() (*user.User, bool) {
	for _, u := range user.All() {
		if u.Player().Name() == m.name {
			return u, true
		}
	}
	return nil, false
}

// Captain returns whether the member is at least a captain in hierarchy.
func (m *Member) Captain() bool {
	return m.role == CAPTAIN || m.Leader()
}

// NewMember creates a new faction member.
func NewMember(name string, role int, leader bool) *Member {
	return &Member{
		name:   name,
		role:   role,
		leader: leader,
	}
}