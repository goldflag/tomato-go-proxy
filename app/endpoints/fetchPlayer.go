package endpoints

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"encoding/json"
	"os"
	"strconv"

	"github.com/mitchellh/mapstructure"
)

// Endpoint: api.worldoftanks.com/wot/tanks/stats
type TankStat struct {
	All struct {
		Spotted                int     `json:"spotted"`
		Hits                   int     `json:"hits"`
		Losses                 int     `json:"losses"`
		Draws                  int     `json:"draws"`
		Wins                   int     `json:"wins"`
		Avg_damage_blocked     float32 `json:"avg_damage_blocked"`
		Capture_points         int     `json:"capture_points"`
		Battles                int     `json:"battles"`
		Damage_dealt           int     `json:"damage_dealt"`
		Damage_received        int     `json:"damage_received"`
		Piercings              int     `json:"piercings"`
		Shots                  int     `json:"shots"`
		Frags                  int     `json:"frags"`
		Tanking_factor         float32 `json:"tanking_factor"`
		Xp                     int     `json:"xp"`
		Survived_battles       int     `json:"survived_battles"`
		Dropped_capture_points int     `json:"dropped_capture_points"`
	} `json:"all"`
	Tank_id int `json:"tank_id"`
}

// Endpoint: api.worldoftanks.com/wot/account/info
type PlayerStat struct {
	Statistics struct {
		All struct {
			Battles int `json:"battles"`
		} `json:"all"`
	} `json:"statistics"`
	Nickname string `json:"nickname"`
}

func getGeneralStats(server, id string) (PlayerStat, error) {
	link := fmt.Sprint(
		"https://api.worldoftanks.",
		server,
		"/wot/account/info/?application_id=",
		os.Getenv("API_KEY"),
		"&account_id=",
		id,
		"&fields=statistics.all.battles%2Cnickname",
	)
	resChannel, errChannel := fetchData(link)

	select {
	case generalStats := <-resChannel:

		// handles Wargaming API error
		if generalStats["status"] == "error" {
			return PlayerStat{}, errors.New(fmt.Sprint(generalStats["error"]))
		}

		generalStatsInner := generalStats["data"].(map[string]interface{})[id]

		var stats PlayerStat
		mapstructure.Decode(generalStatsInner, &stats)
		return stats, nil
	case err := <-errChannel:
		return PlayerStat{}, err
	}
}

func getTankStats(server, id string) ([]TankStat, error) {
	link := fmt.Sprint(
		"https://api.worldoftanks.",
		server,
		"/wot/tanks/stats/?application_id=",
		os.Getenv("API_KEY"),
		"&account_id=",
		id,
		"&fields=tank_id%2call.draws%2call.wins%2call.losses%2call.xp%2call.dropped_capture_points%2call.spotted%2call.battles%2call.capture_points%2call.survived_battles%2call.damage_dealt%2call.damage_received%2call.frags%2call.tanking_factor%2call.avg_damage_blocked%2call.shots%2call.hits%2call.piercings",
	)

	resChannel, errChannel := fetchData(link)

	select {
	case tankStats := <-resChannel:

		// handles Wargaming API error
		if tankStats["status"] == "error" {
			return make([]TankStat, 0), errors.New(fmt.Sprint(tankStats["error"]))
		}

		tankStatsInner := tankStats["data"].(map[string]interface{})

		// If no user exists but the user ID is in a valid format, the WG API will just return null data instead of an actual error
		if tankStatsInner[id] == nil {
			return make([]TankStat, 0), errors.New("Player tank stats null")
		}

		tankStatsInner_ := tankStatsInner[id].([]any)

		var processedTankStats []TankStat
		for _, v := range tankStatsInner_ {
			var tank TankStat
			mapstructure.Decode(v, &tank)
			processedTankStats = append(processedTankStats, tank)
		}

		return processedTankStats, nil
	case err := <-errChannel:
		return make([]TankStat, 0), err
	}
}

type Response struct {
	Status    string      `json:"status"`
	Username  string      `json:"username"`
	Id        string      `json:"id"`
	TankStats [][]float32 `json:"tankStats"`
	Error     string      `json:"error"`
}

func FetchPlayer(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, _ := vars["id"]
	server, _ := vars["server"]

	fmt.Println(id, server)
	if _, err := strconv.Atoi(id); err != nil {
		invalidId, _ := json.Marshal(Response{Status: "error", Id: id, Error: "ID must be an integer"})
		w.Write(invalidId)
		return
	}

	if server != "com" && server != "eu" && server != "asia" && server != "ru" {
		invalidServer, _ := json.Marshal(Response{Status: "error", Id: id, Error: `Invalid server. Only valid servers are "com", "asia", "eu", "ru"`})
		w.Write(invalidServer)
		return
	}

	generalStats, generalStatsErr := getGeneralStats(server, id)
	processedTankStats, processedTankStatsErr := getTankStats(server, id)

	if generalStatsErr != nil {
		fail, _ := json.Marshal(Response{Status: "error", Id: id, Error: generalStatsErr.Error()})
		w.Write(fail)
		return
	}

	if processedTankStatsErr != nil {
		fail, _ := json.Marshal(Response{Status: "error", Id: id, Error: processedTankStatsErr.Error()})
		w.Write(fail)
		return
	}

	if generalStats.Statistics.All.Battles > 0 {
		var tankStats [][]float32
		for _, v := range processedTankStats {
			arrayedTankStats := []float32{
				float32(v.Tank_id),
				float32(v.All.Battles),
				float32(v.All.Damage_dealt),
				float32(v.All.Damage_received),
				float32(v.All.Frags),
				float32(v.All.Survived_battles),
				float32(v.All.Wins),
				float32(v.All.Losses),
				float32(v.All.Draws),
				float32(v.All.Capture_points),
				float32(v.All.Dropped_capture_points),
				float32(v.All.Xp),
				float32(v.All.Spotted),
				v.All.Tanking_factor,
				v.All.Avg_damage_blocked,
				float32(v.All.Shots),
				float32(v.All.Hits),
				float32(v.All.Piercings),
			}

			tankStats = append(tankStats, arrayedTankStats)

		}
		successResp, _ := json.Marshal(Response{Status: "success", TankStats: tankStats, Id: id, Username: generalStats.Nickname})
		w.Write(successResp)
		return

	}

	failedResp, _ := json.Marshal(Response{Status: "error", Id: id, Error: "player has no battles"})
	w.Write(failedResp)
}
