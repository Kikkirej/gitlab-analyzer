package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/xanzy/go-gitlab"
	"gitlabAnalyzer/settings"
	"log"
	"os"
	"strconv"
	"strings"
)

func Clone(project *gitlab.Project) (string, *git.Repository) {
	cloneUrl := cloneUrlOf(project)
	destinationPath := destinationPathOf(project)
	repo, err := git.PlainClone(destinationPath, false, &git.CloneOptions{
		URL:             cloneUrl,
		Progress:        os.Stdout,
		InsecureSkipTLS: true,
	})
	if err != nil {
		log.Println("error while pulling", project.HTTPURLToRepo, ":", err)
		return "error", nil
	}
	return destinationPath, repo
}

func destinationPathOf(project *gitlab.Project) string {
	builder := strings.Builder{}
	builder.WriteString(settings.Struct.WorkingDir)
	builder.WriteString(strconv.Itoa(project.ID))
	return builder.String()
}

func cloneUrlOf(project *gitlab.Project) string {
	cloneUrl := strings.ReplaceAll(
		project.HTTPURLToRepo,
		"://",
		"://oauth2:"+settings.Struct.GitlabPersonalToken+"@",
	)
	cloneUrl = strings.ReplaceAll(
		cloneUrl,
		"https://",
		"http://",
	)
	return cloneUrl
}

func CheckoutBranch(repo *git.Repository, branch *gitlab.Branch) error {
	worktree, err := repo.Worktree()
	if err != nil {
		log.Println("could not checkout branch", branch.Name, ":", err)
		return err
	}
	referenceName := plumbing.NewRemoteReferenceName("origin", branch.Name)
	errCheckout := worktree.Checkout(&git.CheckoutOptions{
		Branch: referenceName,
	})
	if errCheckout != nil {
		log.Println("could not checkout branch", branch.Name, ":", errCheckout)
		return errCheckout
	}
	return nil
}
