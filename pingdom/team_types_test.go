package pingdom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTeamValid(t *testing.T) {
	team := Team{
		Name:      "fake team",
		MemberIDs: nil,
	}
	params := team.Valid()

	assert.Equal(t, nil, params, "Team.Valid() should return nil if valid")
}

func TestTeamNotValid(t *testing.T) {
	team := Team{
		Name:      "",
		MemberIDs: []int{1, 3},
	}
	params := team.Valid()

	assert.NotEqual(t, nil, params, "Team.Valid() should return not nil if not valid")
}
