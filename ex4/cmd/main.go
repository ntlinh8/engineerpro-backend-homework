package main

import (
	"encoding/json"
	"fmt"
	"homework/ex4/models"
	"homework/ex4/pkg"
	"net/http"
	"strconv"
	"golang.org/x/crypto/bcrypt"
)
var userMap = make(map[string]models.User)
var userList []models.User

func main() {
	LoadDataToMemory()

	http.HandleFunc("/signin/", handlerSignIn)
	http.HandleFunc("/login/", handlerLogin)
	http.HandleFunc("/getUsers/", handleGetUserList)
	http.HandleFunc("/getUser/", handleGetUserById)
	http.HandleFunc("/editUser/", handleEditUser)

	fmt.Println("Start server on localhost:8080 ...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server", err.Error())
	}
	

}

func LoadDataToMemory(){
	userList = pkg.LoadFileToMemory()
	for _, user := range userList {
		userMap[user.Email] = user
	}
}

// Create a new user
func handlerSignIn(w http.ResponseWriter, r *http.Request){
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		pkg.HandleErrorForRequest(w, err, "Unable to parse form")

		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")
		address := r.FormValue("address")

		if name == "" || email == "" || password == "" || address == ""{
			http.Error(w, "Please fill in all required fields", http.StatusBadRequest)
		}

		
		if _, exist := userMap[email]; exist {
			http.Error(w, "Email was exist on database", http.StatusBadRequest)
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		pkg.HandleError(err, "Cannot hash password")

		user := models.User{ID: len(userList) + 1, Name: name, Email: email, Password: hashedPassword, Address: address}
		userMap[email] = user

		userList = append(userList, user)
		pkg.LoadDataFromMemoryToFile(userList)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)

	} else{
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// Login
func handlerLogin(w http.ResponseWriter, r *http.Request){
	if r.Method == http.MethodPost{
		err := r.ParseForm()
		pkg.HandleErrorForRequest(w, err, "Unable to parse form")

		email := r.FormValue("email")
		password := r.FormValue("password")
		if email == "" || password == "" {
			http.Error(w, "Please fill in all required fields", http.StatusBadRequest)
		}

		if _, exist := userMap[email]; exist != true{
			http.Error(w, "User is not exist on this system", http.StatusBadRequest)
			return
		}

		user := userMap[email]
		err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))

		if err != nil{
			http.Error(w, "Email or password was wrong", http.StatusBadRequest)
			return
		} else{
			http.ServeFile(w, r, "ex4/static/home.html")
			return
		}

	}else{
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func handleGetUserList(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodGet{
		http.Error(w, "Invalid request method", http.StatusBadRequest)
	} else{
		jsonData, err := json.Marshal(userList)
		pkg.HandleErrorForRequest(w, err, "Unable to encode userList to JSON")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	}
}

func handleGetUserById(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodGet{
		http.Error(w, "Invalid request method", http.StatusBadRequest)
	}else{
		id, err:= strconv.Atoi(r.URL.Path[len("/getUser/"):])
		pkg.HandleErrorForRequest(w, err, "Unable to fetch data from path")

		if id > len(userList) || id <= 0{
			http.Error(w, "Index is out of range", http.StatusBadRequest)
		}

		user := userList[id - 1]
		jsonData, err := json.Marshal(user)
		pkg.HandleErrorForRequest(w, err, "Unable to envode user to JSON")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)

	}
}

func handleEditUser(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPut{
		http.Error(w, "Unable request method", http.StatusBadRequest)
		return
	} else {
		id, err:= strconv.Atoi(r.URL.Path[len("/editUser/"):])
		pkg.HandleErrorForRequest(w, err, "Unable to fetch data from path")

		if id > len(userList) || id <= 0{
			http.Error(w, "Index is out of range", http.StatusBadRequest)
		}

		user := userList[id - 1]

		err = r.ParseForm()
		pkg.HandleErrorForRequest(w, err, "Unable to parse form")

		name := r.FormValue("name")
		email := r.FormValue("email")
		address := r.FormValue("address")

		if name == "" || email == "" || address == ""{
			http.Error(w, "Please fill in all required fields", http.StatusBadRequest)
		}

		oldEmail := user.Email
		user.Name = name
		user.Email = email
		user.Address = address
		delete(userMap, oldEmail)
		userMap[email] = user
		userList[id-1] = user
		pkg.LoadDataFromMemoryToFile(userList)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}
}