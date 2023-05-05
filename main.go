package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"

	model "d-yuji/openai-quickstart-go/model"

	"github.com/joho/godotenv"
)

const endpoint = "https://api.openai.com/v1/completions"

func main() {
	log.Println("Start Server")
	dir, _ := os.Getwd()
	http.HandleFunc("/", mainHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(dir+"/static/"))))
	http.ListenAndServe(":5000", nil)
	defer log.Println("Finish Server")
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	var result string
	if r.Method == http.MethodPost {
		animal := r.FormValue("animal")
		p := generatePrompt(animal)
		var err error
		result, err = postOpenAIAPI(p)
		if err != nil {
			panic(err.Error())
		}
	}
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		panic(err.Error())
	}
	if err := t.Execute(w, result); err != nil {
		panic(err.Error())
	}
}

func postOpenAIAPI(prompt string) (string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return "", fmt.Errorf("can't load .env file: %v", err)
	}
	log.Println("loaded .env")
	OPENAI_API_KEY := os.Getenv("OPENAI_API_KEY")

	j := model.OpenAIRequest{
		Model:       "text-davinci-003",
		Temperature: 0.6,
		Prompt:      prompt,
	}

	jsonString, err := json.Marshal(j)
	if err != nil {
		return "", err
	}

	req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonString))
	req.Header.Set("Authorization", "Bearer "+OPENAI_API_KEY)
	req.Header.Set("Content-Type", "application/json")
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		if resp.Body != nil {
			defer resp.Body.Close()
		}
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Response Body=%v\n", string(body))
		return "", fmt.Errorf("HTTP Error, StatusCode=%v, Body=%v", resp.StatusCode, string(body))
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}
	body, _ := io.ReadAll(resp.Body)
	log.Println(string(body))

	res := model.OpenAIResponse{}

	if err := json.Unmarshal(body, &res); err != nil {
		return "", err
	}
	return res.Choices[0].Text, nil
}

func generatePrompt(animal string) string {
	return fmt.Sprintf(`Suggest three names for an animal that is a superhero.

	Animal: Cat
	Names: Captain Sharpclaw, Agent Fluffball, The Incredible Feline
	Animal: Dog
	Names: Ruff the Protector, Wonder Canine, Sir Barks-a-Lot
	Animal: %s
	Names: 
	`, animal)
}
