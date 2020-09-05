package main

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"log"
	"os"
)

func CheckoutFromScm()  {

	// Clone the given repository to the given directory
	Info("git clone %s %s", cloneUrl, repoClonePath)

	var authUsername = ""
	switch scmProvider {
	case "GitHub":
		authUsername = "x-oauth-basic"
	case "GitLab":
		authUsername = "oauth2"
	case "Bitbucket":
		authUsername = "x-token-auth"
	}

	r, err := git.PlainClone(repoClonePath, false, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: authUsername,
			Password: token,
		},
		URL:      cloneUrl,
		Progress: os.Stdout,
	})
	CheckIfError(err)
	// ... retrieving the commit being pointed by HEAD
	Info("git show-ref --head HEAD")
	ref, err := r.Head()
	CheckIfError(err)
	log.Println(ref.Hash())
	w, err := r.Worktree()
	CheckIfError(err)
	// ... checking out to commit
	Info("git checkout %s", commit)
	err = w.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(commit),
	})
	CheckIfError(err)

	// ... retrieving the commit being pointed by HEAD, it shows that the
	// repository is pointing to the giving commit in detached mode
	Info("git show-ref --head HEAD")
	ref, err = r.Head()
	CheckIfError(err)
	log.Println(ref.Hash())
}
