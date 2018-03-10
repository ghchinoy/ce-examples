package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/ghchinoy/ce-go/ce"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/viper"
)

var (
	profileFilter string
	elementFilter string
	showUsers     bool
)

func init() {
	flag.StringVar(&profileFilter, "filter", "", "profile filter")
	flag.StringVar(&elementFilter, "element", "", "element filter")
	flag.BoolVar(&showUsers, "users", false, "show account e-mails")
	flag.Parse()
}

func main() {

	// load a configuration file
	loadConfigFile()
	// get all profiles
	allprofiles := getAllProfiles()
	// filter profiles
	profiles := filterList(allprofiles, profileFilter)
	if profileFilter != "" {
		log.Printf("Querying %v/%v profiles: %s\n", len(profiles), len(allprofiles), profileFilter)
	} else {
		log.Printf("Querying %v profiles", len(profiles))
	}
	if elementFilter != "" {
		log.Printf("Filtering by Element key: %s", elementFilter)
	}

	// print out the instances in the profile which have events enabled
	data := [][]string{}
	for _, profile := range profiles { // todo: goroutine
		instances, err := getInstancesWithEvents(profile)
		if err != nil {
			log.Println(profile, "instances", err.Error())
			//os.Exit(1)
		}

		if len(instances) > 0 { // this could be a goroutine returning []string
			var userlist []string
			// conditional for user account e-mail output
			if showUsers {
				users, err := getAllUsers(profile)
				if err != nil {
					log.Printf("%s users %s", profile, err.Error())
				}

				for _, user := range users {
					userlist = append(userlist, user.EMail)
				}
			}
			// filter Elements
			if elementFilter != "" {
				instances = filterElements(instances)
			}
			for _, i := range instances {
				// conditional for user account e-mail output
				if showUsers {
					data = append(data, []string{
						profile,
						fmt.Sprintf("%v", userlist),
						i.Element.Key,
						strconv.Itoa(i.ID),

						i.Name,
						strconv.FormatBool(i.Disabled),
						i.Configuration.EventVendorType,
						i.Configuration.EventPollerRefreshInterval,
					})
				} else {
					data = append(data, []string{
						profile,
						i.Element.Key,
						strconv.Itoa(i.ID),

						i.Name,
						strconv.FormatBool(i.Disabled),
						i.Configuration.EventVendorType,
						i.Configuration.EventPollerRefreshInterval,
					})
				}
			}
		}
	}
	table := tablewriter.NewWriter(os.Stdout)
	if showUsers {
		table.SetHeader([]string{"Profile", "EMail", "Element", "ID", "Name", "Disabled", "EventType", "Interval"})
	} else {
		table.SetHeader([]string{"Profile", "Element", "ID", "Name", "Disabled", "EventType", "Interval"})
	}
	table.SetBorder(true)
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.AppendBulk(data)
	table.Render()
}

func filterElements(instances []ce.ElementInstance) []ce.ElementInstance {
	var filteredInstances []ce.ElementInstance
	for _, i := range instances {
		if i.Element.Key == elementFilter {
			filteredInstances = append(filteredInstances, i)
		}
	}
	return filteredInstances
}

func getAllUsers(profile string) ([]ce.User, error) {
	var users []ce.User
	if profile == "profile" {
		return users, fmt.Errorf("'profile' is not a real profile")
	}
	base := viper.Get(profile + ".base").(string)
	auth := fmt.Sprintf("User %s, Organization %s",
		viper.Get(profile+".user").(string),
		viper.Get(profile+".org").(string),
	)
	bodybytes, status, _, err := ce.GetAllUsers(base, auth)
	if err != nil {
		return users, err
	}

	if status != 200 {
		return users, fmt.Errorf("Non-200 Status: %v", status)
	}
	err = json.Unmarshal(bodybytes, &users)
	if err != nil {
		return users, err
	}

	return users, nil
}

// getInstancesWithEvents returns Instances with events on for a profile
func getInstancesWithEvents(profile string) ([]ce.ElementInstance, error) {

	var instances []ce.ElementInstance

	if profile == "profile" {
		return instances, fmt.Errorf("'profile' is not a real profile")
	}

	/*
		fmt.Println("profile", profile)
		fmt.Println("org", viper.Get(profile+".org").(string))
		fmt.Println("user", viper.Get(profile+".user").(string))
		fmt.Println("base", viper.Get(profile+".base").(string))
	*/

	base := viper.Get(profile + ".base").(string)
	auth := fmt.Sprintf("User %s, Organization %s",
		viper.Get(profile+".user").(string),
		viper.Get(profile+".org").(string),
	)
	bodybytes, status, _, err := ce.GetAllInstances(base, auth)
	if err != nil {
		return instances, err
	}

	if status != 200 {
		return instances, fmt.Errorf("Non-200 Status: %v", status)
	}
	err = json.Unmarshal(bodybytes, &instances)
	if err != nil {
		return instances, err
	}

	// here's where the filtering of only event-enabled Instances occur
	var eventInstances []ce.ElementInstance
	for _, i := range instances {
		if i.EventsEnabled == true {
			eventInstances = append(eventInstances, i)
		}
	}

	return eventInstances, nil
}

// filterList returns a list containing only the filter term
func filterList(originallist []string, filterterm string) []string {
	var filteredlist []string
	if filterterm == "" {
		return originallist
	}
	for _, v := range originallist {
		if strings.Contains(v, filterterm) {
			filteredlist = append(filteredlist, v)
		}
	}
	return filteredlist
}

// getAllProfiles lists all the profile names within the settings
func getAllProfiles() []string {
	settings := viper.AllSettings()
	var profiles []string
	for k := range settings {
		profiles = append(profiles, k)
	}
	sort.Strings(profiles)
	return profiles
}

// loadConfigFile loads a toml config file
func loadConfigFile() {
	var cfgfile string
	viper.SetConfigName("cectl")
	viper.AddConfigPath(os.Getenv("HOME") + "/.config/ce")
	if err := viper.ReadInConfig(); err == nil {
		//fmt.Println("Using config file:", viper.ConfigFileUsed())
		cfgfile = viper.ConfigFileUsed()
	}
	//fmt.Printf("%s\n", cfgfile)
	viper.SetConfigFile(cfgfile)
}
