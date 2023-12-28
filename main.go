package main

import (
	"cli-thing/internal/models"
	"encoding/json"
	"log"
	"os"
)

func main() {
	actions := mustFile()
	actions = models.HandleActions(actions)
	saveJson(actions)
}

func mustFile() []models.Action {
	jsonName := "input.json"
	data, err := os.ReadFile(jsonName)
	if err != nil {
		log.Fatal(err)
	}

	actions := make([]models.Action, 0, 1000)

	err = json.Unmarshal(data, &actions)
	if err != nil {
		log.Fatal(err)
	}

	return actions
}

func saveJson(a []models.Action) {
	j, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("input.json", j, models.DefaultPerm)
	if err != nil {
		panic(err)
	}
}
