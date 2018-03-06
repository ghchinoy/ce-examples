package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/ghchinoy/ce-go/ce"
	"github.com/spf13/viper"
)

func main() {

	// load a configuration file
	loadConfigFile()

	// get all profiles
	profiles := listAllProfiles()

	// print out the instances in the profile which have events enabled
	for _, profile := range profiles {
		instances, err := getInstancesWithEvents(profile)
		if err != nil {
			log.Println(profile, err.Error())
			//os.Exit(1)
		}
		if len(instances) > 0 {
			for _, i := range instances {

				fmt.Printf("%s: %v %s/%s: %s @ %s\n",
					profile,
					i.ID,
					i.Element.Key,
					i.Name,
					i.Configuration.EventVendorType,
					i.Configuration.EventPollerRefreshInterval,
				)
			}
		}
	}
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

	var eventInstances []ce.ElementInstance
	for _, i := range instances {
		if i.EventsEnabled == true {
			eventInstances = append(eventInstances, i)
		}
	}

	return eventInstances, nil
}

// listAllProfiles lists all the profile names within the settings
func listAllProfiles() []string {
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
