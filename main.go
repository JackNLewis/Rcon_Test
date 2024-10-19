package main

import (
	"log"
	"net/http"
	"rcon_test/client"
	"regexp"
)

func main() {
	//Connect to minecraft server
	err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer client.ShutDown()

	go client.Start()

	http.HandleFunc("/playerlist", getPlayerListHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))

}

var playerNumbersRegex = regexp.MustCompile(`\d+`)

type playerListResponse struct {
	players int
	isEmpty bool
}

func getPlayerListHandler(w http.ResponseWriter, r *http.Request) {
	playerList, _ := client.GetPlayerCommand()

	players := playerNumbersRegex.FindAllString(playerList, -1)

	response := playerListResponse{
		players: players[0],
		isEmpty: players[0] > 0,
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(numbers[0]))
}
