package main

import (
	"net/http"
	"os"

	"tomato_proxy/endpoints"

	"fmt"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	r := mux.NewRouter()

	/* Returns player tank data in compressed form to save bandwidth
	 *
	 * returns each tank data in an array of this order:
	 * [
	 *	  Tank_id,
	 *	  Battles,
	 *	  Damage_dealt,
	 *	  Damage_received,
	 *	  Frags,
	 *	  Survived_battles,
	 *	  Wins,
	 *	  Losses,
	 *	  Draws,
	 *	  Capture_points,
	 *	  Dropped_capture_points,
	 *	  Xp,
	 *	  Spotted,
	 *	  Tanking_facto,
	 *	  Avg_damage_blocked,
	 *	  Shots,
	 *	  Hits,
	 *	  Piercings
	 * ]
	 */
	r.HandleFunc("/fetchPlayer/{server}/{id}", endpoints.FetchPlayer)

	http.ListenAndServe(fmt.Sprint(":", os.Getenv("PORT")), r)
}
