package agent

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/google/go-github/v32/github"
	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netns"
	"golang.org/x/oauth2"
)

// GitClient function
func (a *Agent) GitClient() {
	log.Printf("NewGithub ...")

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: a.Config.YangConfig.Token.Value},
	)
	tc := oauth2.NewClient(ctx, ts)
	//cm := "commit message"
	bb := "master"
	//prs := "PR SRL Git"
	//prd := "PR Description"
	f := "/etc/opt/srlinux/config.json"

	a.Github.ctx = ctx
	a.Github.baseBranch = &bb
	//a.Github.commitMessage = &cm
	//a.Github.prSubject = &prs
	//a.Github.prDescription = &prd
	a.Github.file = &f
	a.Github.client = github.NewClient(tc)
	a.Github.state = new(gitClientState)
}

// GetRef function
func (a *Agent) GetRef(commitBranch *string) (err error) {
	log.Info("GetRef ...")

	// Get NS
	var nsName string
	if a.Config.YangConfig.NetworkInstance.Value == "" {
		nsName = "mgmt"
	} else {
		nsName = a.Config.YangConfig.NetworkInstance.Value
	}
	ns, err := netns.GetFromName("srbase-" + nsName)
	if err != nil {
		log.Fatal(err)
	}
	// Set NS
	err = netns.Set(ns)
	if err != nil {
		log.Fatal(err)
	}

	o := ""
	if a.Config.YangConfig.Organization.Value != "" {
		o = a.Config.YangConfig.Organization.Value
	}
	if o == "" {
		if a.Config.YangConfig.Owner.Value != "" {
			o = a.Config.YangConfig.Owner.Value
		} else {
			log.Printf("Error: organization or owner should be set")
			return nil
		}
	}

	if a.Github.Ref, _, err = a.Github.client.Git.GetRef(a.Github.ctx, o, a.Config.YangConfig.Repo.Value, "refs/heads/"+*commitBranch); err == nil {
		return nil
	}

	// We consider that an error means the branch has not been found and needs to
	// be created.
	if *commitBranch == *a.Github.baseBranch {
		return errors.New("The commit branch does not exist but `-base-branch` is the same as `-commit-branch`")
	}

	if *a.Github.baseBranch == "" {
		return errors.New("The `-base-branch` should not be set to an empty string when the branch specified by `-commit-branch` does not exists")
	}

	var baseRef *github.Reference
	if baseRef, _, err = a.Github.client.Git.GetRef(a.Github.ctx, o, a.Config.YangConfig.Repo.Value, "refs/heads/"+*a.Github.baseBranch); err != nil {
		return err
	}
	newRef := &github.Reference{Ref: github.String("refs/heads/" + *commitBranch), Object: &github.GitObject{SHA: baseRef.Object.SHA}}
	a.Github.Ref, _, err = a.Github.client.Git.CreateRef(a.Github.ctx, o, a.Config.YangConfig.Repo.Value, newRef)
	return err

}

// GetTree function
func (a *Agent) GetTree() (err error) {
	// Create a tree with what to commit.
	entries := []*github.TreeEntry{}

	file := a.Config.YangConfig.FileName.Value
	content, err := ioutil.ReadFile(*a.Github.file)

	entry := &github.TreeEntry{
		Path:    github.String(file),
		Type:    github.String("blob"),
		Content: github.String(string(content)),
		Mode:    github.String("100644"),
	}

	// log.Printf("File: %v", entry)

	entries = append(entries, entry)

	var nsName string
	if a.Config.YangConfig.NetworkInstance.Value == "" {
		nsName = "mgmt"
	} else {
		nsName = a.Config.YangConfig.NetworkInstance.Value
	}
	ns, err := netns.GetFromName("srbase-" + nsName)
	if err != nil {
		log.Fatal(err)
	}
	// Set NS
	err = netns.Set(ns)
	if err != nil {
		log.Fatal(err)
	}

	o := ""
	if a.Config.YangConfig.Organization.Value != "" {
		o = a.Config.YangConfig.Organization.Value
	}
	if o == "" {
		if a.Config.YangConfig.Owner.Value != "" {
			o = a.Config.YangConfig.Owner.Value
		} else {
			log.Printf("Error: organization or owner should be set")
			return nil
		}
	}

	// log.Printf("Organization: %s, Repo: %s, Ref: %s", o, a.Config.YangConfig.Repo.Value, *a.Github.Ref.Object.SHA)

	a.Github.Tree, _, err = a.Github.client.Git.CreateTree(a.Github.ctx, o, a.Config.YangConfig.Repo.Value, *a.Github.Ref.Object.SHA, entries)
	return err
}

// PushCommit creates the commit in the given reference using the given tree.
func (a *Agent) PushCommit(args *Args, ref *github.Reference, tree *github.Tree) (err error) {
	var nsName string
	if a.Config.YangConfig.NetworkInstance.Value == "" {
		nsName = "mgmt"
	} else {
		nsName = a.Config.YangConfig.NetworkInstance.Value
	}
	ns, err := netns.GetFromName("srbase-" + nsName)
	if err != nil {
		log.Fatal(err)
	}
	// Set NS
	err = netns.Set(ns)
	if err != nil {
		log.Fatal(err)
	}
	o := ""
	if a.Config.YangConfig.Organization.Value != "" {
		o = a.Config.YangConfig.Organization.Value
	}
	if o == "" {
		if a.Config.YangConfig.Owner.Value != "" {
			o = a.Config.YangConfig.Owner.Value
		} else {
			log.Printf("Error: organization or owner should be set")
			return nil
		}
	}
	// Get the parent commit to attach the commit to.
	parent, _, err := a.Github.client.Repositories.GetCommit(a.Github.ctx, o, a.Config.YangConfig.Repo.Value, *a.Github.Ref.Object.SHA)
	if err != nil {
		return err
	}
	// This is not always populated, but is needed.
	parent.Commit.SHA = parent.SHA

	// Create the commit using the tree.
	date := time.Now()
	author := &github.CommitAuthor{Date: &date, Name: &a.Config.YangConfig.Author.Value, Email: &a.Config.YangConfig.AuthorEmail.Value}
	commit := &github.Commit{Author: author, Message: &args.Comment, Tree: tree, Parents: []*github.Commit{parent.Commit}}
	newCommit, _, err := a.Github.client.Git.CreateCommit(a.Github.ctx, o, a.Config.YangConfig.Repo.Value, commit)
	if err != nil {
		return err
	}

	// Attach the commit to the master branch.
	ref.Object.SHA = newCommit.SHA
	_, _, err = a.Github.client.Git.UpdateRef(a.Github.ctx, o, a.Config.YangConfig.Repo.Value, ref, false)
	return err
}

// CreatePR creates a pull request. Based on: https://godoc.org/github.com/google/go-github/github#example-PullRequestsService-Create
func (a *Agent) CreatePR(args *Args, commitBranch *string) (err error) {
	var nsName string
	if a.Config.YangConfig.NetworkInstance.Value == "" {
		nsName = "mgmt"
	} else {
		nsName = a.Config.YangConfig.NetworkInstance.Value
	}
	ns, err := netns.GetFromName("srbase-" + nsName)
	if err != nil {
		log.Fatal(err)
	}
	// Set NS
	err = netns.Set(ns)
	if err != nil {
		log.Fatal(err)
	}

	if args.Subject == "" {
		return errors.New("missing `-pr-title` flag; skipping PR creation")
	}

	newPR := &github.NewPullRequest{
		Title:               &args.Subject,
		Head:                commitBranch,
		Base:                a.Github.baseBranch,
		Body:                &args.Comment,
		MaintainerCanModify: github.Bool(true),
	}

	o := ""
	if a.Config.YangConfig.Organization.Value != "" {
		o = a.Config.YangConfig.Organization.Value
	}
	if o == "" {
		if a.Config.YangConfig.Owner.Value != "" {
			o = a.Config.YangConfig.Owner.Value
		} else {
			log.Printf("Error: organization or owner should be set")
			return nil
		}
	}

	pr, _, err := a.Github.client.PullRequests.Create(a.Github.ctx, o, a.Config.YangConfig.Repo.Value, newPR)
	if err != nil {
		return err
	}

	fmt.Printf("PR created: %s\n", pr.GetHTMLURL())
	return nil
}
