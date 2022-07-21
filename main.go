package main

import (
	"github.com/andygrunwald/go-jira"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	godotenv.Load(".env")

	tp := jira.BasicAuthTransport{
		Username: os.Getenv("JIRA_USER"),
		Password: os.Getenv("JIRA_TOKEN"),
	}

	client, err := jira.NewClient(tp.Client(), os.Getenv("JIRA_URL"))
	if err != nil {
		log.Fatalln(err)
	}

	me, _, err := client.User.GetSelf()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(me)
}
