package main

import (
	"cli-thing/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {
	actions := mustFile()
	corruptedSet := make(map[int]struct{}, 0)
	for i := range actions {
		_, ok := corruptedSet[i]
		if ok {
			continue
		}
		switch actions[i].Type {
		case models.CreateFile:
			actions[i].CreateFile()
		case models.ChangeFileName:
			actions[i].ChangeFileName()
		case models.DeleteFile:
			actions[i].DeleteFile()
		case models.GetFileDate:
			actions[i].GetFileDate()
		case models.WriteToFile:
			actions[i].WriteToFile()
		case models.DateMoreThan:
			isMoreThan, err := actions[i].DateMoreThan()
			if err != nil {
				actions[i].Result = fmt.Sprintf("could not compare: %v", err.Error())
				corruptedSet[i+1] = struct{}{}
				corruptedSet[i+2] = struct{}{}
			}
			if isMoreThan {
				corruptedSet[i+2] = struct{}{}
				continue
			}
			corruptedSet[i+1] = struct{}{}

		}
	}
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

// TODO save result array into input json
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
