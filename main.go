package main

import (
	"fmt"
	"math/rand"
)

type Faction int

const (
	Marquise Faction = iota
	Eyrie
	Alliance
	Vagabond
	Riverfolk
	Lizard
	Underground
	Corvid
	Hundreds
	Keepers
	Bandits
	Protector
	Band
)

type Map int

const (
	Autumn Map = iota
	Winter
	Lake
	Mountain
)

type Landmark int

const (
	Tower = iota
	Ferry
	City
	Treetop
	Forge
	Market
)

const (
	MarquiseDeCat     = "Marquise de Cat"
	EyrieDynasties    = "Eyrie Dynasties"
	TheVagabond       = "The Vagabond"
	RiverfolkCompany  = "Riverfolk Company"
	LizardCult        = "Lizard Cult"
	UndergroundDuchy  = "Underground Duchy"
	CorvidConspiracy  = "CorvidConspiracy"
	LordOfTheHundreds = "Lord Of The Hundreds"
	KeepersInIron     = "KeepersInIron"
)

func isBotAvailable(faction Faction) bool {
	if faction > Hundreds {
		return false
	}
	switch faction {
	case Keepers, Hundreds:
		return false
	}
	return true
}

func getLandmarkName(landmark Landmark) string {
	switch landmark {
	case Tower:
		return "The Tower"
	case Ferry:
		return "The Ferry"
	case City:
		return "Lost City"
	case Forge:
		return "Legendary Forge"
	case Treetop:
		return "Elder Treetop"
	case Market:
		return "Black Market"
	}
	return ""
}

func getFactionName(faction Faction, hireling bool, demoted bool) string {
	switch faction {
	case Marquise:
		return MarquiseDeCat
	case Eyrie:
		return EyrieDynasties
	case Vagabond:
		return TheVagabond
	case Riverfolk:
		return RiverfolkCompany
	case Lizard:
		return LizardCult
	case Underground:
		return UndergroundDuchy
	case Corvid:
		return CorvidConspiracy
	case Hundreds:
		return LordOfTheHundreds
	case Keepers:
		return KeepersInIron
	default:
		return ""
	}
}

// Struct for holding item and its weight (probability)
type Item struct {
	Name   int
	Weight float64
}

// Function to choose an item randomly based on the given probabilities
func weightedRandom(items []Item) int {
	// Calculate the total weight
	totalWeight := 0.0
	for _, item := range items {
		totalWeight += item.Weight
	}

	random := rand.Float64() * totalWeight

	// Select the item based on cumulative weight
	cumulativeWeight := 0.0
	for _, item := range items {
		cumulativeWeight += item.Weight
		if random < cumulativeWeight {
			return item.Name
		}
	}

	// In case something goes wrong, return the last item
	return items[len(items)-1].Name
}

type Match struct {
	PlayerFaction Faction
	BotFactions   [2]Faction
	Hirelings     [3]Faction
	Landmarks     [3]Landmark
}

func randomBetween(min, max int) int {
	return rand.Intn(max-min+1) + min // Generate random number in range [min, max]
}

func getNewMatch(prev Match, factions []Faction, bots []Faction) {
	//nHirelings := randomBetween(1, 3)
	//nLandmarks := randomBetween(1, 3)
	newMatch := Match{}

	playerFactions := []Item{}
	for f := range factions {
		if f == int(prev.PlayerFaction) {
			playerFactions = append(playerFactions, Item{Name: f, Weight: 0.28})
		} else {
			playerFactions = append(playerFactions, Item{Name: f, Weight: 0.08})
		}
	}

	player := weightedRandom(playerFactions)

	botFactions := []Item{}
	for f := range factions {
		if f == player {
			continue
		}
		for bot := range prev.BotFactions {
			if f == bot {
				botFactions = append(botFactions, Item{Name: bot, Weight: 0.15})
			} else {
				botFactions = append(botFactions, Item{Name: bot, Weight: 0.1})
			}
		}
	}

	newMatch.BotFactions[0] = int(weightedRandom(playerFactions))
}

func main() {
	fmt.Println("Running")
	playerFactions := []Faction{Marquise, Eyrie, Alliance, Vagabond, Riverfolk, Lizard, Underground, Corvid, Hundreds, Keepers}
	botFactions := []Faction{Marquise, Eyrie, Alliance, Vagabond, Riverfolk, Lizard, Underground, Corvid}
	/*hirelings := map[Faction][]string{
		Marquise:    {"Forest Patrol", "Feline Physicians"},
		Eyrie:       {"Last Dynasties", "Bluebird Nobles"},
		Alliance:    {"Spring Uprising", "Rabbit Scouts"},
		Vagabond:    {"The Exile", "The Bandit"},
		Riverfolk:   {"Riverfolk Flotilla", "Otter Divers"},
		Lizard:      {"Warm Sun Prophets", "Lizard Envoys"},
		Underground: {"Sunward Expedition", "Mole Artisans"},
		Corvid:      {"Corvid Spies", "Raven Sentinels"},
		Hundreds:    {"Flame Bearers", "Rat Smugglers"},
		Keepers:     {"Vault Keepers", "Badger Bodyguards"},
		Bandits:     {"Highway Bandits", "Bandit Gangs"},
		Protector:   {"Furious Protector", "Stoic Protector"},
		Band:        {"Popular Band", "Street Band"},
	}*/

	prev := Match{
		PlayerFaction: Marquise,
		BotFactions:   [2]Faction{Alliance, Eyrie},
	}

	getNewMatch(prev, playerFactions, botFactions)
}
