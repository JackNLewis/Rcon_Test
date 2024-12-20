package main

import (
	"encoding/json"
	"log"
	"net/http"
	"rcon_test/client"
	"regexp"
	"strconv"
)

func main() {
	// Set up connection to minecraft sever
	err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer client.ShutDown()

	// Start sending commands to server
	go client.Start()

	http.HandleFunc("/playerlist", getPlayerListHandler)
	http.HandleFunc("/kick", kickPlayerHandler)
	http.HandleFunc("/tp", teleportCoord)

	log.Fatal(http.ListenAndServe(":8080", nil))

}

var playerNumbersRegex = regexp.MustCompile(`\d+`)

type playerListResponse struct {
	Players int  `json:"players"`
	IsEmpty bool `json:"is_empty"`
}

func getPlayerListHandler(w http.ResponseWriter, r *http.Request) {
	playerList, _ := client.GetPlayerCommand()

	players := playerNumbersRegex.FindString(playerList)

	playerCount, err := strconv.Atoi(players)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	response := playerListResponse{
		Players: playerCount,
		IsEmpty: playerCount == 0,
	}

	data, _ := json.Marshal(&response)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func kickPlayerHandler(w http.ResponseWriter, r *http.Request) {
	playerName := r.URL.Query().Get("player")

	err := client.KickPlayer(playerName)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)

}

func teleportCoord(w http.ResponseWriter, r *http.Request) {
	playerName := r.URL.Query().Get("player")
	x := r.URL.Query().Get("x")
	y := r.URL.Query().Get("y")
	z := r.URL.Query().Get("z")

	err := client.TeleportToCoord(playerName, x, y, z)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)

}
