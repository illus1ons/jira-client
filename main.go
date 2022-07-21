package main

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	// Load 환경변수
	godotenv.Load(".env")

	tp := jira.BasicAuthTransport{
		Username: os.Getenv("JIRA_USER"),
		Password: os.Getenv("JIRA_TOKEN"),
	}

	// jira client 생성
	client, err := jira.NewClient(tp.Client(), os.Getenv("JIRA_URL"))
	if err != nil {
		log.Fatalln(err)
	}

	me, _, err := client.User.GetSelf()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(me)

	getProjects(client)
	getIssues(client, "project = 'ET Taxtron'")
}

// 모든 프로젝트 목록을 가져오는 함수
func getProjects(client *jira.Client) {
	projectList, _, err := client.Project.GetList()
	if err != nil {
		log.Fatalln(err)
	}

	for _, project := range *projectList {
		fmt.Printf("%s: %s\n", project.Key, project.Name)
	}
}

// Jira JQL 문자열을 이용하여 이슈 조회
func getIssues(client *jira.Client, jql string) {
	issues, _, err := client.Issue.Search(jql, nil)
	if err != nil {
		panic(err)
	}

	for _, i := range issues {
		fmt.Printf("%s (%s/%s): %+v -> %s\n", i.Key, i.Fields.Type.Name, i.Fields.Priority.Name, i.Fields.Summary, i.Fields.Status.Name)
		if i.Fields.Assignee != nil {
			fmt.Printf("Assignee : %v\n", i.Fields.Assignee.DisplayName)
		}
		fmt.Printf("Reporter: %v\n", i.Fields.Reporter.DisplayName)
	}
}
