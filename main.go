package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("args", "Testing command line arguments.")
	env = app.Flag("env", "environment").Short('e').Default("").String()

	list = app.Command("list", "Lists all timeslices for an environment.").Default()

	info          = app.Command("info", "Gets information about a timeslice.")
	infoTimeslice = info.Arg("timeslice", "The timeslice to get information about.").Required().String()

	set          = app.Command("set", "Sets flags in shard specifications for a timeslice.")
	setTimeslice = set.Arg("timeslice", "The timeslice to change settings.").Required().String()
	setEnabled   = set.Flag("enabled", "Sets the enabled flag.").Default("").String()
	setRunning   = set.Flag("running", "Sets the running flag.").Default("").String()
	setWedged    = set.Flag("wedged", "Sets the wedged flag.").Default("").String()

	delete          = app.Command("delete", "Deletes a timeslice from shard specifications.")
	deleteTimeslice = delete.Arg("timeslice", "The timeslice to delete.").Required().String()
	deleteDryRun    = delete.Flag("dry-run", "Dry run?").Default("true").Bool()

	reassign              = app.Command("reassign", "Reassigns all accounts' primary and backup shards to another timeslice.")
	reassignFromTimeslice = reassign.Arg("from-timeslice", "The timeslice to reassign accounts.").Required().String()
	reassignToTimeslice   = reassign.Arg("to-timeslice", "The timeslice to assign accounts to.").Required().String()
)

func main() {
	defaultEnv := "staging"
	hostEnv := ""

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case list.FullCommand():
		fmt.Printf("Listing timeslices\n")

	case info.FullCommand():
		fmt.Printf("Information about [%s]\n", *infoTimeslice)
		hostEnv = parseEnvFromHost(*infoTimeslice)

	case set.FullCommand():
		fmt.Printf("Setting flags on [%s]\n", *setTimeslice)
		hostEnv = parseEnvFromHost(*setTimeslice)
		if *setEnabled != "" {
			if *setEnabled != "0" && *setEnabled != "1" {
				log.Fatalf("Invalid value for enabled flag [%s]\n", *setEnabled)
			}
			fmt.Printf("\tSetting enabled to [%s]\n", *setEnabled)
		}
		if *setRunning != "" {
			if *setRunning != "0" && *setRunning != "1" {
				log.Fatalf("Invalid value for running flag [%s]\n", *setRunning)
			}
			fmt.Printf("\tSetting running to [%s]\n", *setRunning)
		}
		if *setWedged != "" {
			if *setWedged != "0" && *setWedged != "1" {
				log.Fatalf("Wedged must be 0 or 1: [%s]\n", *setWedged)
			}
			fmt.Printf("\tWedge flag = [%s]\n", *setWedged)
		}

	case delete.FullCommand():
		fmt.Printf("Deleting [%s]\n", *deleteTimeslice)
		hostEnv = parseEnvFromHost(*deleteTimeslice)
		fmt.Printf("\tDry run mode = [%t]\n", *deleteDryRun)

	case reassign.FullCommand():
		fmt.Printf("Reassigning metric shards from [%s] to [%s]\n", *reassignFromTimeslice, *reassignToTimeslice)
		hostEnv = parseEnvFromHost(*reassignFromTimeslice)
	}

	if *env == "" && hostEnv != "" {
		*env = hostEnv
	} else if *env == "" {
		*env = defaultEnv
	}
	fmt.Printf("env = [%s]\n", *env)
}

func parseEnvFromHost(host string) (env string) {
	if strings.Contains(host, "-prod-") {
		env = "prod"
	} else if strings.Contains(host, "-staging-") {
		env = "staging"
	}
	return
}
