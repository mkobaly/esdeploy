package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mkobaly/esdeploy/elastic"

	"log"

	"github.com/alecthomas/kingpin"
	"github.com/fatih/color"
)

var (
	app = kingpin.New("esdeploy", "A command-line deployment tool to version Elastic Search.")

	drCmd  = app.Command("dryrun", "Performs a dry run listing out changes that would be made")
	drURL  = drCmd.Arg("url", "Elastic Search URL to run against").Required().String()
	drPath = drCmd.Flag("folder", "Folder containing schema js files").Short('f').Default(".").String()

	deployCmd = app.Command("deploy", "Deploy elastic search changes")
	dURL      = deployCmd.Arg("url", "Elastic Search URL to run against").Required().String()
	dPath     = deployCmd.Flag("folder", "Folder containing schema js files").Short('f').Default(".").String()
	dSilent   = deployCmd.Flag("silent", "Don't prompt for confirmation, run silently").Short('s').Bool()
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case drCmd.FullCommand():
		if *drPath == "" {
			*drPath, _ = os.Getwd()
		}

		color.Yellow("Running dry run against %v", *drURL)
		color.Yellow("Folder containing schema files is %v", *drPath)

		schemaChanger := elastic.NewEsSchemaChanger(*drURL)
		esRunner := elastic.NewRunner(*drPath, schemaChanger)
		results, err := esRunner.DryRun()
		if err != nil {
			log.Fatal(err)
		}
		for _, r := range results {
			color.Green("%v", r)
		}
		color.Green("Dry Run completed")

	case deployCmd.FullCommand():
		if *dPath == "" {
			*dPath, _ = os.Getwd()
		}
		color.Yellow("About to perform dry run against %v", *dURL)
		color.Yellow("Folder containing schema files is %v", *dPath)

		if *dSilent == false {
			color.Cyan("Do you want to proceed? Yes(Y) or No(N)")
			var input string
			fmt.Scanln(&input)
			if strings.ToUpper(input) == "N" {
				os.Exit(0)
			}
		}

		schemaChanger := elastic.NewEsSchemaChanger(*dURL)
		esRunner := elastic.NewRunner(*dPath, schemaChanger)
		results, err := esRunner.Deploy()
		if err != nil {
			color.Red(err.Error())
			os.Exit(1)
		}
		for _, r := range results {
			color.Green("%v", r)
		}
		color.Green("Deploy completed")
	}
}
