package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Package struct {
	Success bool `json:"Success"`
}

type User struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
	Token    string `json:"Token"`
}

func handleRequests() {
	http.HandleFunc("/login", loginEndpoint)
	log.Fatal(http.ListenAndServe(":12345", nil))
}

func main() {
	handleRequests()
}

func loginEndpoint(w http.ResponseWriter, r *http.Request) {

	var loginData User

	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Accept")

	returnData := Package{
		Success: false,
	}

	if checkLoginInput(loginData) {
		returnData.Success = true
	}

	json.NewEncoder(w).Encode(returnData)

	fmt.Println("Endpoint Hit: loginEndpoint")
}

func checkLoginInput(data User) bool {
	if data.Token != strconv.Itoa(time.Now().Hour())+strconv.Itoa(time.Now().Minute()) {
		return false
	}
	jsonFile, err := os.Open("users.json")
	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var users []User

	json.Unmarshal(byteValue, &users)

	for i := 0; i < len(users); i++ {
		if users[i].Username == data.Username {
			//compare stored hash to newly hashed input
			fmt.Println(data.Password)
			fmt.Println(checkPasswordHash(data.Password, users[i].Password))
			return checkPasswordHash(data.Password, users[i].Password)
		}
	}

	return false
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// func hashPassword(password string) (string, error) {
// 	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
// 	return string(bytes), err
// }
