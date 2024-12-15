package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"

	"github.com/gecko05/legacyRoot/matchpb"
)

const (
	Marquise int32 = iota
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
	None
)

type Map struct {
	val  int32
	name string
}

const (
	Autumn int32 = iota
	Winter
	Lake
	Mountain
)

const (
	Tower int32 = iota
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

func getLandmarkName(landmark int32) string {
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

func FactionId(name string) int32 {
	switch name {
	case MarquiseDeCat:
		return Marquise
	case EyrieDynasties:
		return Eyrie
	case WoodlandAlliance:
		return Alliance
	case TheVagabond:
		return Vagabond
	case RiverfolkCompany:
		return Riverfolk
	case LizardCult:
		return Lizard
	case UndergroundDuchy:
		return Underground
	case CorvidConspiracy:
		return Corvid
	case LordOfTheHundreds:
		return Hundreds
	case KeepersInIron:
		return Keepers
	default:
		return None
	}
}

func getFactionName(id int32) string {
	switch id {
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
	Name   int32
	Weight float64
}

// Function to choose an item randomly based on the given probabilities
func pickRandom(items []Item) int32 {
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
	val  int32
	name string
}

type Faction struct {
	val  int32
	name string
}

func NewFaction(enum int32) Faction {
	f := Faction{val: enum, name: getFactionName(enum)}
	return f
}

type Match struct {
	PlayerFactions Faction     `json:"players"`
	BotFactions    [2]Faction  `json:"bots"`
	Hirelings      [3]Faction  `json:"hirelings"`
	Landmarks      [3]Landmark `json:"landmarks"`
	Map            Map         `json:"map"`
}

type MatchCfg struct {
	UseHirelings bool
	UseLandmarks bool
	BotEnemies   int32
	Players      int32
}

func randomBetween(min, max int32) int32 {
	return int32(rand.Intn(int(max - min + 1 + min))) // Generate random number in range [min, max]
}

func removeFromPool(e int32, pool []Item) []Item {
	for i, bot := range pool {
		if bot.Name == e {
			pool = append(pool[:i], pool[i+1:]...)
			break
		}
	}
	return pool
}

func generateNewMatch(prev Match, factions map[int32]string, bots map[int32]string, hirelings map[int32][]string, cfg *MatchCfg) Match {
	newMatch := Match{}

	// Pick player factions.
	newMatch.PlayerFactions = pickPlayerFactions(prev, factions)

	// Remove player factions from bot and hirelings pools.
	delete(hirelings, newMatch.PlayerFactions.val)
	delete(bots, newMatch.PlayerFactions.val)

	// Pick Bots
	newMatch.BotFactions = pickBotFactions(prev, cfg.BotEnemies, bots)

	// Remove non compatible hirelings based on bot factions.
	for _, bot := range newMatch.BotFactions {
		delete(hirelings, bot.val)
	}

	// Pick hireings.
	newMatch.Hirelings = pickHirelings(prev, hirelings)

	// Pick Map
	maps := map[int32]string{Autumn: "Autumn", Winter: "Winter", Lake: "Lake", Mountain: "Mountain"}
	newMatch.Map = pickMap(prev, maps)

	// Pick Landmarks
	nLandmarks := randomBetween(0, 3)
	landmarks := []int32{Tower, Ferry, Treetop, City, Market, Forge}
	newMatch.Landmarks = pickLandmarks(nLandmarks, landmarks)

	return newMatch
}

func pickLandmarks(n int32, landmarks []int32) [3]Landmark {
	pickedLandmarks := [3]Landmark{}
	if n > 0 {
		landmarkSelection := []Item{}
		for _, v := range landmarks {
			landmarkSelection = append(landmarkSelection, Item{Name: v, Weight: float64(1.0 / len(landmarks))})
		}

		for i := range n {
			landmarkId := pickRandom(landmarkSelection)
			pickedLandmarks[i] = Landmark{val: landmarkId, name: getLandmarkName(landmarkId)}
			removeFromPool(landmarkId, landmarkSelection)
		}
	}
	return pickedLandmarks
}

func pickMap(prev Match, maps map[int32]string) Map {
	mapSelection := []Item{}
	for k := range maps {
		if k == prev.Map.val {
			mapSelection = append(mapSelection, Item{Name: k, Weight: 0.34})
		} else {
			mapSelection = append(mapSelection, Item{Name: k, Weight: 0.22})
		}
	}
	m := pickRandom(mapSelection)

	return Map{val: m, name: maps[m]}
}

func pickHirelings(prev Match, hirelings map[int32][]string) [3]Faction {
	nHirelings := randomBetween(0, 3)
	pickedHirelings := [3]Faction{}
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
			hirelingFactions = append(hirelingFactions, Item{Name: k, Weight: float64((weightAll - (0.15 * float64(prevCount))) / 10)})
		}

		for i := range nHirelings {
			rank := randomBetween(0, 1)
			h := pickRandom(hirelingFactions)
			pickedHirelings[i] = Faction{val: h, name: hirelings[h][rank]}
			for j, h := range hirelingFactions {
				if h.Name == pickedHirelings[i].val {
					hirelingFactions = append(hirelingFactions[:j], hirelingFactions[j+1:]...)
				}
			}
		}
	}
	return pickedHirelings
}

func pickPlayerFactions(prev Match, factions map[int32]string) Faction {
	playerFactions := []Item{}
	for f := range factions {
		if f == int32(prev.PlayerFactions.val) {
			playerFactions = append(playerFactions, Item{Name: f, Weight: 0.28})
		} else {
			playerFactions = append(playerFactions, Item{Name: f, Weight: 0.08})
		}
	}
	factionId := pickRandom(playerFactions)
	playerFaction := NewFaction(factionId)
	return playerFaction
}

func pickBotFactions(prev Match, n int32, factions map[int32]string) [2]Faction {
	BotFactions := []Item{}
	for f := range factions {
		for prevBot := range prev.BotFactions {
			if f == int32(prevBot) {
				BotFactions = append(BotFactions, Item{Name: f, Weight: 0.15})
				break
			} else {
				BotFactions = append(BotFactions, Item{Name: f, Weight: 0.1})
				break
			}
		}
	}
	bots := [2]Faction{}
	for i := range n {
		botId := int32(pickRandom(BotFactions))
		bots[i] = NewFaction(botId)
		removeFromPool(botId, BotFactions)
	}
	return bots
}

func main() {
	fmt.Println("Running")
	//var hFlag = flag.Bool("h", true, "Use hirelings")
	playerFactions := map[int32]string{Marquise: MarquiseDeCat,
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

	BotFactions := map[int32]string{Marquise: MarquiseDeCat,
		Eyrie:       EyrieDynasties,
		Alliance:    WoodlandAlliance,
		Vagabond:    TheVagabond,
		Riverfolk:   RiverfolkCompany,
		Lizard:      LizardCult,
		Underground: UndergroundDuchy,
		Corvid:      CorvidConspiracy,
	}

	hirelings := map[int32][]string{
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

	data, err := os.ReadFile("match.json")
	if err != nil {
		fmt.Printf("failed to read match file: %v", err)
	}

	previous := matchpb.Match{}
	err = json.Unmarshal(data, &previous)
	if err != nil {
		fmt.Printf("failed to deserialize match: %v", err)
	}

	cfg := MatchCfg{UseHirelings: false, UseLandmarks: true, Players: 1, BotEnemies: 1}
	newMatch := generateNewMatch(previous, playerFactions, BotFactions, hirelings, &cfg)
	fmt.Printf("Player Faction: %v\n", newMatch.PlayerFactions.name)
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
