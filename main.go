package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/xanzy/go-gitlab"
)

var git *gitlab.Client

var True, False = true, false

type projectInfo struct {
	*gitlab.Project
	hasCode      bool
	hasIssues    bool
	hasPipelines bool
}

func main() {
	token := os.Getenv("GITLAB_TOKEN")
	if token == "" {
		log.Fatal("GITLAB_TOKEN is not present")
	}

	baseURL := os.Getenv("GITLAB_URL")
	if baseURL == "" {
		log.Fatal("GITLAB_URL is not present")
	}

	git = gitlab.NewClient(nil, token)
	git.SetBaseURL(baseURL)

	pageEmpty := false
	opts := &gitlab.ListProjectsOptions{
		Archived:   &False,
		Statistics: &True,
	}
	projects := map[string]*projectInfo{}
	projectNames := []string{}

	for page := 0; !pageEmpty; page++ {
		opts.Page = page
		projectList, _, err := git.Projects.ListProjects(opts)
		if err != nil {
			log.Fatalf("unexpected error: %+v", err)
		}

		pageEmpty = len(projectList) == 0

		for idx, project := range projectList {
			printStatus(len(projects), len(projects)+len(projectList)-idx)
			projectNames = append(projectNames, project.NameWithNamespace)
			projects[project.NameWithNamespace] = &projectInfo{
				Project:      project,
				hasCode:      checkHasCode(project),
				hasIssues:    checkHasIssues(project),
				hasPipelines: checkHasPipelines(project),
			}
		}
	}

	sort.Strings(projectNames)
	for idx, projectName := range projectNames {
		project := projects[projectName]
		fmt.Printf("#%d ID:%d Name:%q HasIssues:%v HasCode:%v HasPipelines:%v LastActivity: %q\n", idx, project.ID, project.NameWithNamespace, project.hasIssues, project.hasCode, project.hasPipelines, project.LastActivityAt)
	}
}

func checkHasIssues(p *gitlab.Project) bool {
	issues, res, err := git.Issues.ListProjectIssues(p.ID, &gitlab.ListProjectIssuesOptions{})
	if res.StatusCode == 403 {
		// we fallback to the most conservative option
		return true
	}
	if err != nil {
		log.Fatalf("cannot list issues for project %q: %+v", p.NameWithNamespace, err)
	}

	return len(issues) > 0
}

func checkHasCode(p *gitlab.Project) bool {
	// Apparently old projects with low activity have no stats. It is reasonable to report them as "does not have code".
	if p.Statistics == nil {
		return false
	}
	// We assume a project with just a few commits has no code (maybe just issue templates, readmes, etc).
	// A project could have a little bit of real code in just a few commits, but it is ok to assume that is irrelevant.
	return p.Statistics.CommitCount > 10
}

func checkHasPipelines(p *gitlab.Project) bool {
	pipelines, res, err := git.Pipelines.ListProjectPipelines(p.ID, &gitlab.ListProjectPipelinesOptions{})
	if res.StatusCode == 403 {
		// we fallback to the most conservative option
		return true
	}
	if err != nil {
		log.Fatalf("cannot list pipelines for project %q: %+v", p.NameWithNamespace, err)
	}

	return len(pipelines) > 0
}

func printStatus(partial, total int) {
	fmt.Printf("\rFetched %d/%d projects", partial, total)
}
