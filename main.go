package main

import (
	"fmt"

	jira "github.com/andygrunwald/go-jira"
)

var (
	jiraURL  = "https://issues.redhat.com/"
	username = "XXX"
	password = "XXX"
	filterID = 12348916
	boardID  = "5296"
)

func main() {
	tp := jira.BasicAuthTransport{
		Username: username,
		Password: password,
	}

	jiraClient, err := jira.NewClient(tp.Client(), jiraURL)

	filter, _, err := jiraClient.Filter.Get(filterID)
	if err != nil {
		fmt.Printf("err getting filter: %v", err)
	}

	sprint, err := getCurrentSprint(jiraClient)
	if err != nil {
		fmt.Printf("err getting current sprint: %v", err)
	}

	issues, _, err := jiraClient.Issue.Search(filter.Jql, nil)
	if err != nil {
		fmt.Printf("err getting list of issues: %v", err)
	}
	ids := getIssueIDs(issues)

	jiraClient.Sprint.MoveIssuesToSprint(sprint, ids)
	if err != nil {
		fmt.Printf("err moving issues to current sprint: %v\n", err)
	}
	fmt.Printf("Moved %d issues!\n", len(ids))
}

func getCurrentSprint(client *jira.Client) (int, error) {
	sprints, _, err := client.Board.GetAllSprints(boardID)
	if err != nil {
		return 0, err
	}

	for _, sprint := range sprints {
		if sprint.CompleteDate == nil {
			return sprint.ID, nil
		}
	}
	return 0, fmt.Errorf("unable to find current sprint")
}

func getIssueIDs(issues []jira.Issue) []string {
	ids := make([]string, 0)
	for k, issue := range issues {
		ids = append(ids, issue.ID)
	}
	return ids
}
