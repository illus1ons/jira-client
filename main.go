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
	//getIssues(client, "project = 'ET Taxtron'")
	getIssues(client, "project = 'ET Taxtron' and Status = 'Done'")
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
func getIssues(client *jira.Client, jql string) ([]jira.Issue, error) {

	// lastIssue is the index of the last issue returned
	lastIssue := 0
	// Make a loop through amount of issues
	var result []jira.Issue
	for {
		// Add a Search option which accepts maximum amount (1000)
		opt := &jira.SearchOptions{
			MaxResults: 1000,      // Max amount
			StartAt:    lastIssue, // Make sure we start grabbing issues from last checkpoint
		}
		issues, resp, err := client.Issue.Search(jql, opt)
		if err != nil {
			return nil, err
		}
		// Grab total amount from response
		total := resp.Total
		if issues == nil {
			// init the issues array with the correct amount of length
			result = make([]jira.Issue, 0, total)
		}

		// Append found issues to result
		result = append(result, issues...)
		// Update checkpoint index by using the response StartAt variable
		lastIssue = resp.StartAt + len(issues)
		// Check if we have reached the end of the issues
		if lastIssue >= total {
			break
		}
	}

	for _, i := range result {
		fmt.Printf("%s (%s/%s): %+v -> %s\n", i.Key, i.Fields.Type.Name, i.Fields.Priority.Name, i.Fields.Summary, i.Fields.Status.Name)
		if i.Fields.Assignee != nil {
			fmt.Printf("Assignee : %v\n", i.Fields.Assignee.DisplayName)
		}
		fmt.Printf("Reporter: %v\n", i.Fields.Reporter.DisplayName)
	}
	return result, nil
}
