package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/doneill/er-cli-go/config"
)

// ----------------------------------------------
// stucts
// ----------------------------------------------

type UserResponse struct {
	Data struct {
		Username  string `json:"username"`
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		ID        string `json:"id"`
		Pin       string `json:"pin"`
		Subject   struct {
			ID string `json:"id"`
		} `json:"subject"`
	} `json:"data"`
	ErrorDescription string `json:"error_description"`
}

// ----------------------------------------------
// exported funtions
// ----------------------------------------------

func User() (*UserResponse, error) {
	client := &http.Client{}
	clientReq := getClientRequest(config.Sitename(), API_USER_ME, config.Token())

	res, err := client.Do(&clientReq)
	if err != nil {
		fmt.Println("Error making request:", err)
	}
	var responseData UserResponse
	err = json.NewDecoder(res.Body).Decode(&responseData)
	if err != nil {
		fmt.Println("Error decoding response:", err)
	}

	if res.StatusCode == 200 {
		return &responseData, nil
	}

	fmt.Println("Error:", res.StatusCode)
	fmt.Println("Error Description:", responseData.ErrorDescription)

	return nil, err
}
