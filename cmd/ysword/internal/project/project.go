package project

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var CmdNew = &cobra.Command{
	Use:   "new",
	Short: "Create a service template",
	Long:  "Create a service project using the repository template. Example: ysword new helloworld",
	Run:   run,
}

var (
	repoURL string
	branch  string
	timeout string
)

type Project struct {
	Name string
	Path string
}

func (p *Project) New(ctx context.Context, dir string, layout string, branch string) error {
	to := path.Join(dir, p.Name)
	if _, err := os.Stat(to); !os.IsNotExist(err) {
		fmt.Printf(" %s already exists\n", p.Name)
		override := false
		prompt := &survey.Confirm{
			Message: "Do you want to override the folder ?",
			Help:    "Delete the existing folder and create the project.",
		}
		e := survey.AskOne(prompt, &override)
		if e != nil {
			return e
		}
		if !override {
			return err
		}
		os.RemoveAll(to)
	}
	fmt.Printf("🚀 Creating service %s, layout repo is %s, please wait a moment.\n\n", p.Name, layout)
	repo := NewRepo(layout, branch)
	if err := repo.CopyTo(ctx, to, p.Path, []string{".git", ".github"}); err != nil {
		return err
	}

	// clear repo folder
	os.RemoveAll("repo/")

	e := os.Rename(
		path.Join(to, "cmd", "web"),
		path.Join(to, "cmd", p.Name),
	)
	if e != nil {
		return e
	}
	Tree(to, dir)

	// clean cacche dir

	fmt.Printf("\n🍺 Project creation succeeded %s\n", p.Name)
	fmt.Print("💻 Use the following command to start the project 👇:\n\n")

	fmt.Printf("$ cd %s", p.Name)
	fmt.Println("$ go generate ./...")
	fmt.Println("$ go build -o ./bin/ ./... ")
	fmt.Printf("$ ./bin/%s -conf ./configs\n", p.Name)
	fmt.Println("			🤝 Thanks for using ysword")
	return nil
}

func init() {
	repoURL = "https://github.com/zlog-fun/ysword-layout.git"

	timeout = "60s"
	CmdNew.Flags().StringVarP(&repoURL, "repo-url", "r", repoURL, "layout repo")
	CmdNew.Flags().StringVarP(&branch, "branch", "b", branch, "repo branch")
	CmdNew.Flags().StringVarP(&timeout, "timeout", "t", timeout, "time out")
}

func run(cmd *cobra.Command, args []string) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	t, err := time.ParseDuration(timeout)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()
	name := ""
	if len(args) == 0 {
		prompt := &survey.Input{
			Message: "What is project name ?",
			Help:    "Created project name.",
		}
		err = survey.AskOne(prompt, &name)
		if err != nil || name == "" {
			return
		}
	} else {
		name = args[0]
	}

	p := &Project{Name: path.Base(name), Path: name}
	done := make(chan error, 1)

	go func() {
		done <- p.New(ctx, wd, repoURL, branch)
	}()

	select {
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			fmt.Fprint(os.Stderr, "\033[31mERROR: project creation timed out\033[m\n")
		} else {
			fmt.Fprintf(os.Stderr, "\033[31mERROR: failed to create project(%s)\033[m\n", ctx.Err().Error())
		}
	case err = <-done:
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mERROR: Failed to create project(%s)\033[m\n", err.Error())
		}
	}
}
