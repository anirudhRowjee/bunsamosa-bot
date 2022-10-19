package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/anirudhRowjee/bunsamosa-bot/globals"
)

func Leaderboard_allrecords(response http.ResponseWriter, request *http.Request) {
	records, err := globals.Myapp.Leaderboard_GetAllRecords()
	if err != nil {
		log.Println("[ERR][LEADERBOARD_HANDLER] Could not get all records ->", err)
		response.WriteHeader(http.StatusInternalServerError)
	} else {
		// Marshal into JSON
		json_string, err := json.Marshal(records)
		if err != nil {
			log.Println("[ERR][LEADERBOARD_HANDLER] Failed to Marshal records into JSON ->", err)
			response.WriteHeader(http.StatusInternalServerError)
		} else {
			response.Header().Set("Content-Type", "application/json")
			response.Write(json_string)
		}
	}
}

func Leaderboard_materialized(response http.ResponseWriter, request *http.Request) {

	records, err := globals.Myapp.Leaderboard_GetMaterialized()
	if err != nil {
		log.Println("[ERROR][LEADERBOARD_HANDLER][MATERIALIZED] Could not get all records ->", err)
		response.WriteHeader(http.StatusInternalServerError)
	} else {
		// Marshal into JSON
		json_string, err := json.Marshal(records)
		if err != nil {
			log.Println("[ERROR][LEADERBOARD_HANDLER][MATERIALIZED] Failed to Marshal records into JSON ->", err)
			response.WriteHeader(http.StatusInternalServerError)
		} else {
			response.Header().Set("Content-Type", "application/json")
			response.Write(json_string)
		}
	}
}

func Leaderboard_userspecific(response http.ResponseWriter, request *http.Request) {
	records, err := globals.Myapp.Leaderboard_GetAllRecords()
	if err != nil {
		log.Println("[ERROR][LEADERBOARD_HANDLER][USERSPECIFIC] Could not get all records ->", err)
		response.WriteHeader(http.StatusInternalServerError)
	} else {
		// Marshal into JSON
		json_string, err := json.Marshal(records)
		if err != nil {
			log.Println("[ERROR][LEADERBOARD_HANDLER][USERSPECIFIC] Failed to Marshal records into JSON ->", err)
			response.WriteHeader(http.StatusInternalServerError)
		} else {
			response.Header().Set("Content-Type", "application/json")
			response.Write(json_string)
		}
	}
}
