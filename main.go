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
	version string
)

var (
	app         = kingpin.New("esdeploy", "A command-line deployment tool to version ElasticSearch.")
	appUser     = app.Flag("username", "Username to authenticate with").Short('u').String()
	appPassword = app.Flag("password", "Password to authenticat with").Short('p').String()
	appInsecure = app.Flag("insecure", "Ignore SSL certificate warnings").Short('k').Bool()

	drCmd  = app.Command("dryrun", "Only lists out changes that would be made to ElasticSearch.")
	drURL  = drCmd.Arg("url", "Elastic Search URL to run against").Required().String()
	drPath = drCmd.Flag("folder", "Folder containing schema js files").Short('f').Default(".").String()

	validateCmd  = app.Command("validate", "Performs a validation of all files to ensure they are properly formatted")
	validatePath = validateCmd.Flag("folder", "Folder containing schema js files").Short('f').Default(".").String()

	deployCmd = app.Command("deploy", "Deploy elastic search changes")
	dURL      = deployCmd.Arg("url", "Elastic Search URL to run against").Required().String()
	dPath     = deployCmd.Flag("folder", "Folder containing schema js files").Short('f').Default(".").String()
	dSilent   = deployCmd.Flag("silent", "Don't prompt for confirmation, run silently").Short('s').Bool()

	seedCmd  = app.Command("seed", "Seed elastic search with data stored in json files")
	seedURL  = seedCmd.Arg("url", "Elastic Search URL to run against").Required().String()
	seedPath = seedCmd.Flag("folder", "Folder containing json data files").Short('f').Default(".").String()

	versionCmd = app.Command("version", "Display version of esdeploy")
)

func main() {
	var cred elastic.Creds

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {

	case versionCmd.FullCommand():
		color.Cyan("version %v", version)
		os.Exit(0)

	//Validation
	case validateCmd.FullCommand():
		if *validatePath == "" {
			*validatePath, _ = os.Getwd()
		}
		exit := 0
		color.Cyan("Running validation against folder %v", *validatePath)
		esRunner := elastic.NewRunner(*validatePath, nil)
		results := esRunner.Validate()
		for _, r := range results {
			if !r.IsValid {
				color.Red("FILE INVALID: %s", r.File)
				exit = 1
				continue
			}
			color.Green("File Valid: %s", r.File)
		}

		color.Cyan("Validation completed")
		os.Exit(exit)

	//Dry run
	case drCmd.FullCommand():
		if *drPath == "" {
			*drPath, _ = os.Getwd()
		}

		color.Cyan("Running dry run against %v", *drURL)
		color.Cyan("Folder containing schema files is %v", *drPath)

		schemaChanger := elastic.NewEsSchemaChanger(*drURL, cred, *appInsecure)
		esRunner := elastic.NewRunner(*drPath, schemaChanger)
		results, err := esRunner.DryRun()
		if err != nil {
			log.Fatal(err)
		}
		for _, r := range results {
			color.Green("%v", r)
		}
		color.Cyan("Dry Run completed")

	//Full deployment
	case deployCmd.FullCommand():
		if *dPath == "" {
			*dPath, _ = os.Getwd()
		}

		if *appUser != "" && *appPassword != "" {
			cred = elastic.Creds{Username: *appUser, Password: *appPassword}
		}

		color.Cyan("About to perform deployment against %v", *dURL)
		color.Cyan("Folder containing schema files is %v", *dPath)

		if *dSilent == false {
			color.Yellow("Do you want to proceed? Yes(Y) or No(N)")
			var input string
			fmt.Scanln(&input)
			if strings.ToUpper(input) == "N" {
				os.Exit(0)
			}
		}

		schemaChanger := elastic.NewEsSchemaChanger(*dURL, cred, *appInsecure)
		esRunner := elastic.NewRunner(*dPath, schemaChanger)
		results, err := esRunner.Deploy()
		if err != nil {
			for _, r := range results {
				color.Red("%v", r)
			}
			color.Red(err.Error())
			os.Exit(1)
		}
		for _, r := range results {
			color.Green("%v", r)
		}
		color.Cyan("Deploy completed")
	//Seed data
	case seedCmd.FullCommand():
		if *seedPath == "" {
			*seedPath, _ = os.Getwd()
		}

		color.Cyan("Seeding data against %v", *seedURL)
		color.Cyan("Folder containing data files is %v", *seedPath)

		seeder := elastic.NewSeeder(*seedPath, *seedURL, cred)
		results, err := seeder.Seed()
		if err != nil {
			log.Fatal(err)
		}
		for _, r := range results {
			color.Green("%v", r)
		}
		color.Cyan("Seeding completed")
	}
}
