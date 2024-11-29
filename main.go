package main

import (
	"fmt"
	"math/rand"
)

const (
	Marquise int = iota
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

type Map struct {
	val  int
	name string
}

const (
	Autumn int = iota
	Winter
	Lake
	Mountain
)

const (
	Tower int = iota
	Ferry
	City
	Treetop
	Forge
	Market
)

const (
	MarquiseDeCat     = "Marquise de Cat"
	EyrieDynasties    = "Eyrie Dynasties"
	WoodlandAlliance  = "Woodland Alliance"
	TheVagabond       = "The Vagabond"
	RiverfolkCompany  = "Riverfolk Company"
	LizardCult        = "Lizard Cult"
	UndergroundDuchy  = "Underground Duchy"
	CorvidConspiracy  = "Corvid Conspiracy"
	LordOfTheHundreds = "Lord Of The Hundreds"
	KeepersInIron     = "KeepersInIron"
)

type Faction struct {
	val  int
	name string
}

func NewFaction(enum int) Faction {
	f := Faction{val: enum, name: getFactionName(enum)}
	return f
}

func isBotAvailable(int int) bool {
	if int > Hundreds {
		return false
	}
	switch int {
	case Keepers, Hundreds:
		return false
	}
	return true
}

func getLandmarkName(landmark int) string {
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

func getFactionName(int int) string {
	switch int {
	case Marquise:
		return MarquiseDeCat
	case Eyrie:
		return EyrieDynasties
	case Alliance:
		return WoodlandAlliance
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
func pickRandom(items []Item) int {
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

type Landmark struct {
	val  int
	name string
}

type Match struct {
	playerFactions Faction
	BotFactions    [2]Faction
	Hirelings      [3]Faction
	Landmarks      [3]Landmark
	Map            Map
}

type MatchCfg struct {
	UseHirelings bool
	UseLandmarks bool
	BotEnemies   int
	Players      int
}

func randomBetween(min, max int) int {
	return rand.Intn(max-min+1) + min // Generate random number in range [min, max]
}

func removeFromPool(e int, pool []Item) []Item {
	for i, bot := range pool {
		if bot.Name == e {
			pool = append(pool[:i], pool[i+1:]...)
			break
		}
	}
	return pool
}

func generateNewMatch(prev Match, factions map[int]string, bots map[int]string, hirelings map[int][]string, cfg *MatchCfg) Match {
	//nLandmarks := randomBetween(1, 3)
	newMatch := Match{}

	// Create player faction pool probabilities with less weight for previous played faction
	playerFactions := []Item{}
	for f := range factions {
		if f == int(prev.playerFactions.val) {
			playerFactions = append(playerFactions, Item{Name: f, Weight: 0.28})
		} else {
			playerFactions = append(playerFactions, Item{Name: f, Weight: 0.08})
		}
	}
	factionId := pickRandom(playerFactions)
	playerFaction := NewFaction(factionId)
	newMatch.playerFactions = playerFaction

	delete(hirelings, factionId)

	// Pick first bot
	delete(bots, factionId)
	BotFactions := []Item{}
	for f := range bots {
		for prevBot := range prev.BotFactions {
			if f == prevBot {
				BotFactions = append(BotFactions, Item{Name: f, Weight: 0.15})
				break
			} else {
				BotFactions = append(BotFactions, Item{Name: f, Weight: 0.1})
				break
			}
		}
	}
	for i := range cfg.BotEnemies {
		botId := int(pickRandom(BotFactions))
		newMatch.BotFactions[i] = NewFaction(botId)
		removeFromPool(botId, BotFactions)
		delete(hirelings, botId)
	}

	nHirelings := randomBetween(0, 3)
	if nHirelings > 0 {
		hirelingFactions := []Item{}
		prevCount := 0
		for k := range hirelings {
			for _, prevHireling := range prev.Hirelings {
				if prevHireling.val == k {
					prevCount += 1
					hirelingFactions = append(hirelingFactions, Item{Name: k, Weight: 0.15})
				}
			}
		}

		weightAll := 1.0
		for k := range hirelings {
			if k == newMatch.BotFactions[0].val || k == newMatch.playerFactions.val || k == newMatch.BotFactions[1].val {
				continue
			}
			hirelingFactions = append(hirelingFactions, Item{Name: k, Weight: float64((weightAll - (0.15 * float64(prevCount))) / 10)})
		}

		for i := range nHirelings {
			rank := randomBetween(0, 1)
			h := pickRandom(hirelingFactions)
			newMatch.Hirelings[i] = Faction{val: h, name: hirelings[h][rank]}
			for j, h := range hirelingFactions {
				if h.Name == newMatch.Hirelings[i].val {
					hirelingFactions = append(hirelingFactions[:j], hirelingFactions[j+1:]...)
				}
			}
		}
	}

	maps := map[int]string{Autumn: "Autumn", Winter: "Winter", Lake: "Lake", Mountain: "Mountain"}
	//maps := [4]int{Autumn, Winter, Lake, Mountain}
	mapSelection := []Item{}
	for k := range maps {
		if k == prev.Map.val {
			mapSelection = append(mapSelection, Item{Name: k, Weight: 0.34})
		} else {
			mapSelection = append(mapSelection, Item{Name: k, Weight: 0.22})
		}
	}
	m := pickRandom(mapSelection)
	newMatch.Map = Map{val: m, name: maps[m]}

	nLandmarks := randomBetween(0, 3)
	if nLandmarks > 0 {
		landmarks := [6]int{Tower, Ferry, Treetop, City, Market, Forge}
		landmarkSelection := []Item{}
		for _, v := range landmarks {
			landmarkSelection = append(landmarkSelection, Item{Name: v, Weight: float64(1.0 / 6.0)})
		}

		for i := range nLandmarks {
			l := pickRandom(landmarkSelection)
			newMatch.Landmarks[i] = Landmark{val: l, name: getLandmarkName(l)}

			for j, v := range landmarks {
				if v == newMatch.Landmarks[i].val {
					landmarkSelection = append(landmarkSelection[:j], landmarkSelection[j+1:]...)
				}
			}
		}
	}

	return newMatch
}

func main() {
	fmt.Println("Running")
	//var hFlag = flag.Bool("h", true, "Use hirelings")
	//playerFactions := []int{Marquise, Eyrie, Alliance, Vagabond, Riverfolk, Lizard, Underground, Corvid, Hundreds, Keepers}
	playerFactions := map[int]string{Marquise: MarquiseDeCat,
		Eyrie:       EyrieDynasties,
		Alliance:    WoodlandAlliance,
		Vagabond:    TheVagabond,
		Riverfolk:   RiverfolkCompany,
		Lizard:      LizardCult,
		Underground: UndergroundDuchy,
		Corvid:      CorvidConspiracy,
		Hundreds:    LordOfTheHundreds,
		Keepers:     KeepersInIron,
	}

	BotFactions := map[int]string{Marquise: MarquiseDeCat,
		Eyrie:       EyrieDynasties,
		Alliance:    WoodlandAlliance,
		Vagabond:    TheVagabond,
		Riverfolk:   RiverfolkCompany,
		Lizard:      LizardCult,
		Underground: UndergroundDuchy,
		Corvid:      CorvidConspiracy,
	}

	hirelings := map[int][]string{
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
	}

	prev := Match{
		playerFactions: NewFaction(Riverfolk),
		BotFactions:    [2]Faction{NewFaction(Corvid), NewFaction(Alliance)},
		Hirelings:      [3]Faction{},
		Map:            Map{val: Autumn, name: "Autumn"},
		Landmarks:      [3]Landmark{},
	}

	cfg := MatchCfg{UseHirelings: false, UseLandmarks: true, Players: 1, BotEnemies: 1}
	newMatch := generateNewMatch(prev, playerFactions, BotFactions, hirelings, &cfg)
	fmt.Printf("Player Faction: %v\n", newMatch.playerFactions.name)
	fmt.Printf("Enemies: %v %v\n", newMatch.BotFactions[0].name, newMatch.BotFactions[1].name)
	fmt.Printf("Hirelings: ")
	for i := range newMatch.Hirelings {
		fmt.Printf("%v ", newMatch.Hirelings[i].name)
	}
	fmt.Println("")
	fmt.Printf("Map: %v\n", newMatch.Map.name)
	fmt.Printf("Landmarks: ")
	for i := range newMatch.Landmarks {
		fmt.Printf("%v ", newMatch.Landmarks[i].name)
	}
	fmt.Println("")
}
