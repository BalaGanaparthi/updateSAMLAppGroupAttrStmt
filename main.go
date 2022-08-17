package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("PORT")
	router := mux.NewRouter()
	secure := router.PathPrefix("/secure").Subrouter()
	secure.Use(JwtVerify)
	secure.HandleFunc("/updatesamlapp", updateSAMLApp).Methods("POST")
	http.ListenAndServe(":"+port, router)
}

func updateSAMLApp(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	log.Println("<<<=====>>>\n updateSAMLApp.")
	w.Header().Set("Content-Type", "application/json")

	var payload ApplicationPayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Println("Error parsing httpBody", err)
	}

	log.Printf("Payload value is %+v \n", payload)

	var apiToken = os.Getenv("x-api-token")
	var orgUrl = os.Getenv("x-org-url")

	appName := payload.Name
	attrValue := payload.AttrValue

	status, err := _updateSAMLApp(orgUrl, apiToken, appName, attrValue)

	if err != nil {
		log.Println("Error while executing getUserGroupScopes")
	}
	log.Printf("Update status is %s \n", status)

	w.Write([]byte("{\"response\":\"Success\"}"))

	end := time.Now()

	nanoTimeDelta := end.UnixNano() - start.UnixNano()
	millisDelta := nanoTimeDelta / 1000000
	log.Printf("Total [%d] nano(s), [%d] milli(s) taken to complete\n", nanoTimeDelta, millisDelta)
}

func _updateSAMLApp(orgUrl string, apiKey string, appName string, grpAttrStmtValue string) (string, error) {

	client := &http.Client{}

	params := url.Values{}
	params.Add("filter", fmt.Sprintf("name eq \"%s\"", appName))

	url := fmt.Sprintf("%s/api/v1/apps?%s", orgUrl, params.Encode())
	r, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatal(err)
	}
	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", fmt.Sprintf("SSWS %s", apiKey))

	res, err := client.Do(r)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("Body : %s", string(body)))

	var bodyMap []map[string]interface{}
	json.Unmarshal(body, &bodyMap)

	appBody := bodyMap[0]

	settings := appBody["settings"].(map[string]interface{})
	signon := settings["signOn"]

	attributeStatements := signon.(map[string]interface{})["attributeStatements"]

	_attributeStatements := attributeStatements.([]interface{})
	for _, attributeStatement := range _attributeStatements {
		statement := attributeStatement.(map[string]interface{})
		if statement["type"] == "GROUP" {
			statement["filterValue"] = grpAttrStmtValue
		}
	}

	updatedPayload, _ := json.Marshal(appBody)
	appID := appBody["id"]
	url = fmt.Sprintf("%s/api/v1/apps/%s", orgUrl, appID)
	req, _ := http.NewRequest("PUT", url, strings.NewReader(string(updatedPayload)))

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("SSWS %s", apiKey))

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return "Done", nil
}

func JwtVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var header = r.Header.Get("x-access-token")

		header = strings.TrimSpace(header)

		log.Println("Authenticate with x-access-token header")

		if header == "" {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode("Missing auth token")
			return
		} else {
			if header != os.Getenv("x-access-token") {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode("Unmatched auth token")
				return
			}
			log.Println("Authentication is successful...")
		}
		log.Println("Allow secure method...")
		next.ServeHTTP(w, r)
		log.Println("Done executing secure method.")
	})
}

type ApplicationPayload struct {
	Name      string `json:"name"`
	AttrValue string `json:"attributeValue"`
}
