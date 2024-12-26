package main

import (
	"fmt"
	"math/rand"
	"os"

	"LegacyRoot/matchpb"

	"github.com/a-h/templ"
	"github.com/labstack/echo"
	"google.golang.org/protobuf/encoding/protojson"
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
	return FactionNames[id]
}

var FactionNames = map[int32]string{
	Marquise:    MarquiseDeCat,
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

var Hirelings = map[int32][]string{
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

func NewFaction(enum int32) *matchpb.Faction {
	f := &matchpb.Faction{Type: matchpb.FactionType(enum), Name: getFactionName(enum)}
	return f
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

func generateNewMatch(
	prev *matchpb.Match,
	factions map[int32]string,
	bots map[int32]string,
	hirelings map[int32][]string,
	cfg *MatchCfg,
) *matchpb.Match {
	newMatch := &matchpb.Match{}

	// Pick player factions.
	newMatch.Players = []*matchpb.Faction{pickPlayerFactions(prev, factions)}

	// Remove player factions from bot and hirelings pools.
	delete(hirelings, int32(newMatch.GetPlayers()[0].GetType()))
	delete(bots, int32(newMatch.GetPlayers()[0].GetType()))

	// Pick Bots
	newMatch.Bots = pickBotFactions(prev, cfg.BotEnemies, bots)

	// Remove non compatible hirelings based on bot factions.
	for _, bot := range newMatch.GetBots() {
		delete(hirelings, int32(bot.GetType()))
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

func pickLandmarks(n int32, landmarks []int32) []*matchpb.Landmark {
	pickedLandmarks := []*matchpb.Landmark{{}, {}, {}}
	if n > 0 {
		landmarkSelection := []Item{}
		for _, v := range landmarks {
			landmarkSelection = append(landmarkSelection, Item{Name: v, Weight: 1.0 / float64(len(landmarks))})
		}

		for i := range n {
			landmarkId := pickRandom(landmarkSelection)
			pickedLandmarks[i] = &matchpb.Landmark{
				Type: matchpb.LandmarkType(landmarkId),
				Name: getLandmarkName(landmarkId),
			}
			landmarkSelection = removeFromPool(landmarkId, landmarkSelection)
		}
	}
	return pickedLandmarks
}

func pickMap(prev *matchpb.Match, maps map[int32]string) *matchpb.MapVal {
	mapSelection := []Item{}
	for k := range maps {
		if k == int32(prev.Map.GetType()) {
			mapSelection = append(mapSelection, Item{Name: k, Weight: 0.34})
		} else {
			mapSelection = append(mapSelection, Item{Name: k, Weight: 0.22})
		}
	}
	m := pickRandom(mapSelection)

	return &matchpb.MapVal{Type: matchpb.MapType(m), Name: maps[m]}
}

func pickHirelings(prev *matchpb.Match, hirelings map[int32][]string) []*matchpb.Faction {
	nHirelings := randomBetween(0, 3)
	pickedHirelings := []*matchpb.Faction{{}, {}, {}}
	if nHirelings > 0 {
		hirelingFactions := []Item{}
		prevCount := 0
		for k := range hirelings {
			for _, prevHireling := range prev.Hirelings {
				if int32(prevHireling.GetType()) == k {
					prevCount += 1
					hirelingFactions = append(hirelingFactions, Item{Name: k, Weight: 0.15})
				}
			}
		}

		weightAll := 1.0
		for k := range hirelings {
			hirelingFactions = append(
				hirelingFactions,
				Item{Name: k, Weight: float64((weightAll - (0.15 * float64(prevCount))) / 10)},
			)
		}

		for i := range nHirelings {
			rank := randomBetween(0, 1)
			h := pickRandom(hirelingFactions)
			pickedHirelings[i] = &matchpb.Faction{Type: matchpb.FactionType(h), Name: hirelings[h][rank]}
			for j, h := range hirelingFactions {
				if h.Name == int32(pickedHirelings[i].Type) {
					hirelingFactions = append(hirelingFactions[:j], hirelingFactions[j+1:]...)
				}
			}
		}
	}
	return pickedHirelings
}

func pickPlayerFactions(prev *matchpb.Match, factions map[int32]string) *matchpb.Faction {
	playerFactions := []Item{}
	for f := range factions {
		if f == int32(prev.Players[0].GetType()) {
			playerFactions = append(playerFactions, Item{Name: f, Weight: 0.28})
		} else {
			playerFactions = append(playerFactions, Item{Name: f, Weight: 0.08})
		}
	}
	factionId := pickRandom(playerFactions)
	playerFaction := NewFaction(factionId)
	return playerFaction
}

func pickBotFactions(prev *matchpb.Match, n int32, factions map[int32]string) []*matchpb.Faction {
	BotFactions := []Item{}
	for f := range factions {
		for _, prevBot := range prev.Bots {
			if f == int32(prevBot.GetType()) {
				BotFactions = append(BotFactions, Item{Name: f, Weight: 0.15})
				break
			} else {
				BotFactions = append(BotFactions, Item{Name: f, Weight: 0.1})
				break
			}
		}
	}
	bots := []*matchpb.Faction{{}, {}}
	for i := range n {
		botId := int32(pickRandom(BotFactions))
		bots[i] = NewFaction(botId)
		removeFromPool(botId, BotFactions)
	}
	return bots
}

/*
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

	previous, err := parseMatch("match.json")
	if err != nil {
		fmt.Printf("failed to parse previous match: %v", err)
		return
	}

	cfg := MatchCfg{UseHirelings: false, UseLandmarks: true, Players: 1, BotEnemies: 1}
	newMatch := generateNewMatch(previous, playerFactions, BotFactions, Hirelings, &cfg)
	fmt.Printf("Player Faction: %v\n", newMatch.GetPlayers()[0].Name)
	fmt.Printf("Enemies: %v %v\n", newMatch.GetBots()[0].Name, newMatch.GetBots()[1].Name)
	fmt.Printf("Hirelings: ")
	for i := range newMatch.Hirelings {
		fmt.Printf("%v ", newMatch.Hirelings[i].Name)
	}
	fmt.Println("")
	fmt.Printf("Map: %v\n", newMatch.Map.Name)
	fmt.Printf("Landmarks: ")
	for i := range newMatch.Landmarks {
		fmt.Printf("%v ", newMatch.GetLandmarks()[i].Name)
	}
	fmt.Println("")
}*/

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return render(c, hello("John"))
	})
	e.Logger.Fatal(e.Start(":1323"))
}

func render(ctx echo.Context, cmp templ.Component) error {
	return cmp.Render(ctx.Request().Context(), ctx.Response())
}

func parseMatch(filename string) (*matchpb.Match, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read match file: %w", err)
	}

	match := &matchpb.Match{}
	err = protojson.Unmarshal(data, match)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize match: %w", err)
	}

	return match, nil
}
