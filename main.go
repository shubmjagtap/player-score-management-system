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
	fmt.Println("Listening and serving on port 9000...")
	http.ListenAndServe("localhost:9000", r)
}

func getSession() *mgo.Session {
	s, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	return s
}
