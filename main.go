package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/go-git/go-git/v5"
)

type Repository struct {
	Name     string    `json:"repo"`
	Packages []Package `json:"packages"`
}
type Package struct {
	Name    string   `json:"name"`
	Git     string   `json:"git"`
	Authors []string `json:"authors"`
	Needs   []string `json:"needs"`
	Build   []string `json:"build"`
	Remove  []string `json:"remove"`
}

const (
	repo_file     = false
	repo_location = "https://turtley12.github.io/pluto/packagelist.json"
)
const help = `Welcome to the pluto package manager. Try:

pluto install <name>
`

var (
	repo Repository
)

func printHelp() {
	fmt.Println(help)
}
func remove(project Package) {
	//fmt.Println("Removing " + project.Name)
	home, _ := os.UserHomeDir()
	dir := home + "/.pluto/" + strings.ToLower(project.Name)
	os.Chdir(dir)
	for _, command := range project.Remove {
		runcommand(command)
	}
	fmt.Println("Removing source directory...")
	os.RemoveAll(dir)

}
func build(project Package) {
	home, _ := os.UserHomeDir()
	dir := home + "/.pluto/" + strings.ToLower(project.Name)
	git.PlainClone(dir, false, &git.CloneOptions{
		URL:               project.Git,
		Depth:             1,
		SingleBranch:      true,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Progress:          os.Stdout,
	})
	os.Chdir(dir)
	fmt.Println("Begining build...")
	for _, command := range project.Build {
		runcommand(command)
	}
}
func runcommand(command string) {
	fmt.Println("Command: " + command)
	args := strings.Split(command, " ")
	last_index := len(args)
	cmd := exec.Command(args[0], args[1:last_index]...)

	stdout, err := cmd.StdoutPipe()
	checkError(err)

	err = cmd.Start()
	checkError(err)

	//log program output
	scan := bufio.NewScanner(stdout)
	for scan.Scan() {
		t := scan.Text()
		fmt.Println(t)
	}
	cmd.Wait()
}
func uninstall(name string) bool {
	fmt.Println("Searching package lists...")
	loadRepo()
	fmt.Println("Repo: " + repo.Name)

	for i, project := range repo.Packages {
		if strings.ToLower(project.Name) == strings.ToLower(name) {
			fmt.Println("Removing " + project.Name + " (Package " + strconv.Itoa(i) + ")")
			remove(project)
			return true

		}
	}
	return false
}
func list() {
	loadRepo()
	for _, project := range repo.Packages {
		fmt.Println(project.Name)
	}
}
func install(name string) bool {
	fmt.Println("Searching package lists...")
	loadRepo()
	fmt.Println("Repo: " + repo.Name)

	for i, project := range repo.Packages {
		if strings.ToLower(project.Name) == strings.ToLower(name) {
			fmt.Println("Downloading " + project.Name + " (Package " + strconv.Itoa(i) + ")")
			build(project)
			return true

		}
	}

	return false //if no package found
}
func main() {
	if len(os.Args) <= 1 {
		printHelp()
		return
	}
	switch os.Args[1] {
	case "install":
		if len(os.Args) != 3 {
			fmt.Println("Please specify a single package name to install")
			break
		}
		//fmt.Println("Installing")
		found := install(os.Args[2])
		if !found {
			fmt.Println("Package " + os.Args[2] + " not found")
		}
	case "remove":
		if len(os.Args) != 3 {
			fmt.Println("Please specify a single package name to uninstall")
			break
		}
		found := uninstall(os.Args[2])
		if !found {
			fmt.Println("Package " + os.Args[2] + " not found")
		}
	case "list":
		list()
	default:
		printHelp()
	}
}
func loadRepo() {
	switch repo_file {
	case true:
		loadRepoFile(repo_location)
	case false:
		loadRepoUrl(repo_location)
	}
}
func loadRepoBytes(repo_bytes []byte) {
	json.Unmarshal(repo_bytes, &repo)
}
func loadRepoFile(file string) {
	repo_bytes, _ := ioutil.ReadFile(file)
	loadRepoBytes(repo_bytes)
}
func loadRepoUrl(url string) {
	response, err := http.Get(url)
	checkError(err)
	defer response.Body.Close()
	bytes, err := io.ReadAll(response.Body)
	checkError(err)
	loadRepoBytes(bytes)
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
