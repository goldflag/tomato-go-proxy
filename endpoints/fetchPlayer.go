package endpoints

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"

	"os"

	"github.com/mitchellh/mapstructure"
)

// Endpoint: api.worldoftanks.com/wot/tanks/stats
type TankStat struct {
	all struct {
		spotted                int
		hits                   int
		losses                 int
		draws                  int
		wins                   int
		avg_damage_blocked     float32
		capture_points         int
		battles                int
		damage_dealt           int
		damage_received        int
		piercings              int
		shots                  int
		frags                  int
		tanking_factor         float32
		xp                     int
		survived_battles       int
		dropped_capture_points int
	}
	tank_id int
}

// Endpoint: api.worldoftanks.com/wot/account/info
type PlayerStat struct {
	statistics struct {
		all struct {
			battles int
		}
	}
	nickname string
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func fetchData(url string) <-chan map[string]interface{} {

	r := make(chan map[string]interface{})

	go func() {
		defer close(r)

		response, err := http.Get(url)

		if err != nil {
			fmt.Printf(err.Error())
			r <- make(map[string]interface{})
		}

		var parsed map[string]interface{}
		body, err := ioutil.ReadAll(response.Body)

		if err != nil {
			fmt.Printf(err.Error())
			r <- make(map[string]interface{})
		}

		err = json.Unmarshal([]byte(body), &parsed)

		r <- parsed
	}()

	return r

}

type Server string

const (
	NA   Server = "com"
	EU   Server = "eu"
	ASIA Server = "asia"
)

func getLinks(server string, id string) []string {
	return []string{
		fmt.Sprint(
			"https://api.worldoftanks.",
			server,
			"/wot/account/info/?application_id=",
			os.Getenv("API_KEY"),
			"&account_id=",
			id,
			"&fields=statistics.all.battles%2Cnickname",
		),
		fmt.Sprint(
			"https://api.worldoftanks.",
			server,
			"/wot/tanks/stats/?application_id=",
			os.Getenv("API_KEY"),
			"&account_id=",
			id,
			"&fields=tank_id%2call.draws%2call.wins%2call.losses%2call.xp%2call.dropped_capture_points%2call.spotted%2call.battles%2call.capture_points%2call.survived_battles%2call.damage_dealt%2call.damage_received%2call.frags%2call.tanking_factor%2call.avg_damage_blocked%2call.shots%2call.hits%2call.piercings",
		),
	}
}

type response struct {
	tankStats []TankStat
	id        string
}

func FetchPlayer(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	id, _ := vars["id"]
	server, _ := vars["server"]

	links := getLinks(server, id)

	generalStatsCh, tankStatsCh := fetchData(links[0]), fetchData(links[1])
	generalStats, tankStats := <-generalStatsCh, <-tankStatsCh
	_, _ = generalStats, tankStats

	tankStatsInner := tankStats["data"]
	tankStatsInner2 := tankStatsInner.(map[string]interface{})[id].([]any)

	// fmt.Println(prettyPrint(tankStatsInner2))

	var processedTankStats []TankStat

	for _, v := range tankStatsInner2 {
		fmt.Println(prettyPrint(v))

		jsonStr1, _ := json.Marshal(v)
		fmt.Println(jsonStr1)

		var parsedStats2 TankStat
		// json.Unmarshal(jsonStr1, &parsedStats2)

		// fmt.Println(parsedStats2)
		mapstructure.Decode(jsonStr1, &parsedStats2)
		fmt.Println(parsedStats2)
		processedTankStats = append(processedTankStats, parsedStats2)
	}

	generalStatsInner := generalStats["data"]
	generalStatsInner2 := generalStatsInner.(map[string]interface{})[id]

	jsonStr, err := json.Marshal(generalStatsInner2)
	if err != nil {
		fmt.Println(err)
	}

	var parsedStats PlayerStat
	if err := json.Unmarshal(jsonStr, &parsedStats); err != nil {
		fmt.Println(err)
	}


	respInJson, err := json.Marshal(response{processedTankStats, id})
	w.Write(respInJson)
	/*

		app.get("/fetchPlayer/:server/:id", async (req, res) => {
		    let currentTime = parseInt(Date.now()/60000);
		    const server = req.params.server;
		    const id = req.params.id;
		    let battles = 0;
		    let stats = {};
		    let data1;
		    let data2;
		    await Promise.all([
		        fetch(`https://api.worldoftanks.${server}/wot/account/info/?application_id=${APsIKey}&account_id=${id}`),
		        fetch(`https://api.worldoftanks.${server}/wot/tanks/stats/?application_id=${APIKey}&account_id=${id}&fields=tank_id%2call.draws%2call.wins%2call.losses%2call.xp%2call.dropped_capture_points%2call.spotted%2call.battles%2call.capture_points%2call.survived_battles%2call.damage_dealt%2call.damage_received%2call.frags%2call.tanking_factor%2call.avg_damage_blocked%2call.shots%2call.hits%2call.piercings`)
		    ])
		    .then(([res1, res2]) => Promise.all([res1.json(), res2.json()]))
		    .then(([d1, d2]) =>
		    {
		        data1 = d1;
		        data2 = d2;
		    });
		    //number of battles overall an account has
		    battles = data1.data[id].statistics.all.battles;
		    if (battles > 0 && data2.data[id]) {
		        //array of overall tank stats
		        stats = data2.data[id];
		        const obj = {
		            status: "success",
		            tankStats: []
		        };
		        obj.id = id;
		        obj.username = data1.data[id].nickname;
		        console.log(obj.username);
		        for (let i = 0; i < stats.length; ++i) {
		            const row = stats[i].all;
		            const tankStats = [
		                stats[i].tank_id,
		                row.battles,
		                row.damage_dealt,
		                row.damage_received,
		                row.frags,
		                row.survived_battles,
		                row.wins,
		                row.losses,
		                row.draws,
		                row.capture_points,
		                row.dropped_capture_points,
		                row.xp,
		                row.spotted,
		                row.tanking_factor,
		                row.avg_damage_blocked,
		                row.shots,
		                row.hits,
		                row.piercings
		            ];
		            obj.tankStats.push(tankStats);
		        }
		        res.json(obj);
		    }
		    else {
		        console.log("failed");
		        res.send({
		            status: "fail"
		        });
		    }
		});
	*/
}
