package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/shubmjagtap/player-score-management-system/controllers"
	"gopkg.in/mgo.v2"
)

func main() {
	r := httprouter.New()
	pc := controllers.NewPlayerController(getSession())
	r.GET("/players", pc.GetAllPlayers)
	r.GET("/players/rank/:val", pc.GetRankedPlayer)
	r.GET("/players/random", pc.GetRandomPlayer)
	r.POST("/players", pc.CreatePlayer)
	r.PUT("/players/:id", pc.UpdatePlayer)
	r.DELETE("/players/:id", pc.DeletePlayer)
	http.ListenAndServe("localhost:9000", r)
	fmt.Println("Listening and serving on port 9000")
}

func getSession() *mgo.Session {
	s, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	return s
}

// 1. POST /players – Creates a new entry for a player
// 2. PUT /players/:id – Updates the player attributes. Only name and
// score can be updated
// 3. DELETE /players/:id – Deletes the player entry
// 4. GET /players – Displays the list of all players in descending order
// 5. GET /players/rank/:val – Fetches the player ranked “val”
// 6. GET /players/random – Fetches a random player
