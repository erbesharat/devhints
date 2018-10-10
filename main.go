package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"

	"github.com/fatih/color"
	"github.com/russross/blackfriday"
)

func main() {
	// Get user instance (Usage: getting user's home directory for cloning cheatsheets)
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	CACHE := user.HomeDir + "/.config/devhints/"

	sync := flag.Bool("sync", false, "Sync the cheatsheets")
	flag.Parse()
	if *sync {
		syncCheatSheets(CACHE)
		os.Exit(1)
	}
	if len(flag.Args()) < 1 {
		log.Fatalf(color.RedString("Please provide a chatsheet name as a subcommand. Check -all for the list"))
	}
	content, err := ioutil.ReadFile(CACHE + "cheatsheets/" + flag.Args()[0] + ".md")
	if err != nil {
		log.Fatalf(color.RedString("Couldn't open the cheatsheet file. Error: %e", err))
	}

	output := blackfriday.MarkdownCommon(content)
	os.Stdout.Write(output)
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func syncCheatSheets(CACHE string) {
	if exists(CACHE) {
		cmd := exec.Command("git", []string{"pull", "origin", "master"}...)
		cmd.Dir = CACHE + "cheatsheets"
		if err := cmd.Run(); err != nil {
			log.Fatalf(color.RedString("Couldn't update the cheatsheets. Error: %s", err.Error()))
		}
		color.Green("Synced the cheatsheets successfully. Check %s to be sure.", CACHE)
		os.Exit(1)
	}
	err := os.MkdirAll(CACHE, os.ModePerm)
	if err != nil {
		log.Fatalf(color.RedString("Couldn't create the directory to keep the cheatsheets. Error: %s", err.Error()))
	}
	cmd := exec.Command("git", []string{"clone", "https://github.com/rstacruz/cheatsheets.git"}...)
	cmd.Dir = CACHE
	if err := cmd.Run(); err != nil {
		log.Fatal(color.RedString("Couldn't get the cheatsheets. Error: %s", err.Error()))
	}
	color.Green("Synced the cheatsheets successfully. Check %s to be sure.", CACHE)
}
