package main

import (
	"context"
	"fmt"
	jira "github.com/andygrunwald/go-jira"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

func init() {
	if envErr := godotenv.Load(".env"); envErr != nil {
		fmt.Println(".env file missing")
	}
}

func jiraClient() *jira.Client {
	jt := jira.BasicAuthTransport{
		Username: os.Getenv("JIRA_USER"),
		Password: os.Getenv("JIRA_TOKEN"),
	}

	client, err := jira.NewClient(jt.Client(), os.Getenv("JIRA_URL"))
	if err != nil {
		fmt.Println(err)
	}

	me, _, err := client.User.GetSelf()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(me.AccountID)

	return client
}

var projectCode string
var email string
var sprint string

func main() {

	fmt.Println("Enter project code: ")
	fmt.Scanln(&projectCode)

	fmt.Println("Enter email address: ")
	fmt.Scanln(&email)

	fmt.Println("Enter Sprint ID: ")
	fmt.Scanln(&sprint)

	if projectCode != "" && email != "" && sprint != "" {
		getTickets(projectCode)
	}
}

func getTickets(projectCode string) {
	client := jiraClient()

	opt := &jira.SearchOptions{
		MaxResults: 1000, // Max results can go up to 1000
		StartAt:    0,
	}

	var accountID string
	var assigneeName string

	users, _, _ := client.User.FindWithContext(context.Background(), email)
	if len(users) < 1 {
		fmt.Println("Sorry, user not found!")
		return
	}

	for _, user := range users {
		accountID = user.AccountID
		assigneeName = user.DisplayName
	}

	issues, _, _ := client.Issue.Search("project = "+projectCode+" AND assignee = "+accountID+" AND sprint = "+sprint+" ORDER BY created DESC", opt)
	totalPoint := 0

	fmt.Println("Assignee :", assigneeName)
	for _, issue := range issues {
		storyPoint, _ := issue.Fields.Unknowns.Value("customfield_10028")

		stringPoint, _ := strconv.Atoi(fmt.Sprintf("%v", storyPoint))
		totalPoint += stringPoint

		issueURL := fmt.Sprintf(os.Getenv("JIRA_URL")+"browse/%v (%v)", issue.Key, storyPoint)

		issueDetail := fmt.Sprintf("%v", issueURL)

		fmt.Println(issueDetail)
	}

	fmt.Printf("Total Tickets: %v, Total Story Points: %v ", len(issues), totalPoint)
}
