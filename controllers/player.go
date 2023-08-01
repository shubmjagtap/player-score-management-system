package controllers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/shubmjagtap/player-score-management-system/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type PlayerController struct {
	session *mgo.Session
}

func NewPlayerController(s *mgo.Session) *PlayerController {
	return &PlayerController{s}
}

func (pc PlayerController) GetAllPlayers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var players []models.Player
	err := pc.session.DB("mongo-golang").C("players").Find(nil).All(&players)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error fetching players: %s", err)
		return
	}

	sort.Slice(players, func(i, j int) bool {
		return players[i].Score > players[j].Score
	})

	playersJSON, err := json.Marshal(players)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error converting players to JSON: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", playersJSON)
}

func (pc PlayerController) GetRankedPlayer(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var players []models.Player
	err := pc.session.DB("mongo-golang").C("players").Find(nil).All(&players)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error fetching players: %s", err)
		return
	}

	// Sort the "players" in descending order based on the "Score" field
	sort.Slice(players, func(i, j int) bool {
		return players[i].Score > players[j].Score
	})

	rank, err := strconv.Atoi(p.ByName("val"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid rank parameter: %s", err)
		return
	}

	// Check if the rank is valid (within the range of players slice)
	if rank < 1 || rank > len(players) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Rank not found")
		return
	}

	// Decrement the rank by 1 to get the index of player in the sorted slice
	rank--
	rankedPlayer := players[rank]

	rankedPlayerJSON, err := json.Marshal(rankedPlayer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error converting player to JSON: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", rankedPlayerJSON)
}

func (pc PlayerController) GetRandomPlayer(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var players []models.Player
	err := pc.session.DB("mongo-golang").C("players").Find(nil).All(&players)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error fetching players: %s", err)
		return
	}

	// Generate a random index within the range of "players" slice
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(players))

	randomPlayer := players[randomIndex]

	randomPlayerJSON, err := json.Marshal(randomPlayer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error converting player to JSON: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", randomPlayerJSON)
}

func (pc PlayerController) CreatePlayer(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Parse the JSON data from the request body to create the player
	var ps models.Player
	err := json.NewDecoder(r.Body).Decode(&ps)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error parsing request body: %s", err)
		return
	}

	// Check if all the fields are present in the request body
	if ps.Name == "" || ps.Country == "" || ps.Score == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "All fields are mandatory")
		return
	}

	// Generate a new ObjectId for the player
	ps.Id = bson.NewObjectId()

	// Insert the new player into the database
	err = pc.session.DB("mongo-golang").C("players").Insert(ps)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating player: %s", err)
		return
	}

	// Convert the player to JSON
	pj, err := json.Marshal(ps)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error converting player to JSON: %s", err)
		return
	}

	// Set the response headers and body
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s\n", pj)
	fmt.Printf("Created player: %s\n", pj)
}

func (pc PlayerController) UpdatePlayer(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Parse the JSON data from the request body to update the player attributes
	var updatedPlayer models.Player
	err := json.NewDecoder(r.Body).Decode(&updatedPlayer)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error parsing request body: %s", err)
		return
	}

	// Check if any of the fields in the updatedPlayer are empty
	if updatedPlayer.Name == "" || updatedPlayer.Country == "" || updatedPlayer.Score == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "All fields are mandatory")
		return
	}

	// Get the player ID from the URL parameter
	id := p.ByName("id")

	// Check if the ID is a valid BSON ObjectID
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Invalid player ID")
		return
	}

	// Convert the ID to a BSON ObjectID
	oid := bson.ObjectIdHex(id)

	// Retrieve the player with the given ID from the database
	player := models.Player{}
	err = pc.session.DB("mongo-golang").C("players").FindId(oid).One(&player)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Player not found")
		return
	}

	if player.Country == updatedPlayer.Country && string(player.Id) == string(updatedPlayer.Id) {
		// Only allow updating Name and Score if the Country and ID are the same
		player.Name = updatedPlayer.Name
		player.Score = updatedPlayer.Score
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: Only Name and Score can be changed")
		return
	}

	// Save the updated player back to the database
	err = pc.session.DB("mongo-golang").C("players").UpdateId(oid, player)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error updating player: %s", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Player updated successfully")
}

func (pc PlayerController) DeletePlayer(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}

	oid := bson.ObjectIdHex(id)

	if err := pc.session.DB("mongo-golang").C("players").RemoveId(oid); err != nil {
		w.WriteHeader(404)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Deleted user", oid, "\n")
}
