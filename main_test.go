package main

import (
	"LegacyRoot/matchpb"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseMatchJSON(t *testing.T) {
	match, err := parseMatch("test_match_0.json")
	assert.NoError(t, err)

	assert.Equal(t, match.GetPlayers()[0].GetType(), matchpb.FactionType_RIVERFOLK)
	assert.Equal(t, match.GetPlayers()[0].GetName(), "Riverfolk Company")

	assert.Equal(t, match.GetBots()[0].GetType(), matchpb.FactionType_CORVID)
	assert.Equal(t, match.GetBots()[1].GetType(), matchpb.FactionType_ALLIANCE)

	assert.Equal(t, match.GetBots()[0].GetName(), "Corvid Conspiracy")
	assert.Equal(t, match.GetBots()[1].GetName(), "Woodland Alliance")

	assert.Equal(t, match.GetMap().GetType(), matchpb.MapType_AUTUMN)
	assert.Equal(t, match.GetLandmarks()[0].GetType(), matchpb.LandmarkType_FORGE)
}
