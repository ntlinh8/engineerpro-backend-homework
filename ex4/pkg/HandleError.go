package pkg

import (
	"fmt"
	"net/http"
)

func HandleErrorForRequest(w http.ResponseWriter, err error, message string) {
	if err != nil {
		http.Error(w, message, http.StatusBadRequest)
		return
	}
}

func HandleError(err error, message string){
	if err != nil {
		fmt.Println("Error occurs: ", message, err.Error())
		return
	}
}