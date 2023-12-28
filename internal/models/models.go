package models

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/djherbis/times"
)

const DefaultPerm = 0774

type Type int

const (
	CreateFile = iota
	ChangeFileName
	DeleteFile
	GetFileDate
	WriteToFile
	DateMoreThan
)

type Action struct {
	Type   Type   `json:"type"`
	Name   string `json:"name"`
	Result string `json:"result"`
	Params param  `json:"params"`
}

type param struct {
	FileName    string `json:"file_name,omitempty"`
	NewFileName string `json:"new_file_name"`
	Message     string `json:"message"`
	CompareDate string `json:"compare_date"`
}
type CivilTime time.Time

func HandleActions(actions []Action) []Action {
	corruptedSet := make(map[int]struct{}, 0)
	for i := range actions {
		_, ok := corruptedSet[i]
		if ok {
			continue
		}
		switch actions[i].Type {
		case CreateFile:
			actions[i].CreateFile()
		case ChangeFileName:
			actions[i].ChangeFileName()
		case DeleteFile:
			actions[i].DeleteFile()
		case GetFileDate:
			actions[i].GetFileDate()
		case WriteToFile:
			actions[i].WriteToFile()
		case DateMoreThan:
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

	return actions
}
func (a *Action) CreateFile() {
	fName := a.Params.FileName
	f, err := os.Create(fName)
	if err != nil {
		a.Result = fmt.Sprintf("could not create file: %v", err.Error())
		return
	}
	defer f.Close()
	a.Result = fmt.Sprintf("Created file %v", fName)
}

func (a *Action) ChangeFileName() {
	oldName, newName := a.Params.FileName, a.Params.NewFileName
	err := os.Rename(oldName, newName)
	if err != nil {
		a.Result = fmt.Sprintf("could not rename file: %v", err.Error())
		return
	}
	a.Result = fmt.Sprintf("Renamed %v to %v", oldName, newName)
}

func (a *Action) DeleteFile() {
	fName := a.Params.FileName
	err := os.Remove(fName)

	if err != nil {
		a.Result = fmt.Sprintf("could not delete file:%v", err.Error())
		return
	}
	a.Result = fmt.Sprintf("File %v successfully deleted", fName)
}

func (a *Action) GetFileDate() {
	fName := a.Params.FileName

	//get birth time of a file
	creationTime, err := getCreationTime(fName)
	if err != nil {
		log.Printf("error getting file date: %v\n", err.Error())
		a.Result = fmt.Sprintf("created at: %v", time.Time{})
		return
	}
	a.Result = fmt.Sprintf("created at: %v", creationTime)
}

func (a *Action) WriteToFile() {
	fName, Msg := a.Params.FileName, a.Params.Message

	err := os.WriteFile(fName, []byte(Msg), DefaultPerm)
	if err != nil {
		a.Result = fmt.Sprintf("could not write to file: %v", err.Error())
		return
	}
	a.Result = fmt.Sprintf("Write message to %v", fName)
}

func (a *Action) DateMoreThan() (bool, error) {
	const layout = "2006-01-02 15:04:05 -0700"
	t, err := time.Parse(layout, a.Params.CompareDate)
	if err != nil {
		a.Result = fmt.Sprintf("could not parse json params compare_date: %v", err.Error())
		return false, err
	}
	//in case i need to compare it to som file
	// tFile, err := getCreationTime(a.Params.FileName)
	// if err != nil {
	// 	a.Result = fmt.Sprintf("could not get file time: %v", err.Error())
	// 	return false, err
	// }
	return t.After(time.Now()), nil
}

func getCreationTime(fileName string) (time.Time, error) {
	ts, err := times.Stat(fileName)
	if err != nil {
		log.Println(err.Error())
		return time.Time{}, err
	}
	if !ts.HasBirthTime() {
		return time.Time{}, errors.New("cannot get creation time of a file")
	}
	return ts.BirthTime(), nil
}
