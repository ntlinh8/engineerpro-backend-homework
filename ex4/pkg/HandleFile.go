package pkg

import (
	"encoding/json"
	"homework/ex4/models"
	"io"
	"os"
)

const path = "ex4/data/user.json"

func LoadFileToMemory() (userList []models.User){
	file, err := os.Open(path)
	HandleError(err, "Unable to open file: " + path)
	defer file.Close()

	bytes, err := io.ReadAll(file)
	HandleError(err, "Unable to read file: " + path)

	err = json.Unmarshal(bytes, &userList)
	HandleError(err, "Error unmarshalling JSON")
	return userList
}

func LoadDataFromMemoryToFile(userList []models.User){
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	HandleError(err, "Unable to open file: " + path)

	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(userList)
	HandleError(err, "Error writing user file")
}