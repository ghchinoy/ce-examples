package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/ghchinoy/ce-go/ce"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/viper"
)

var profileFilter string

func init() {
	flag.StringVar(&profileFilter, "filter", "", "profile filter")
	flag.Parse()
}

func main() {
	// load a configuration file
	loadConfigFile()
	// get all profiles
	allprofiles := getAllProfiles()
	// filter profiles
	profiles := filterList(allprofiles, profileFilter)
	log.Printf("Querying %v/%v profiles: %s\n", len(profiles), len(allprofiles), profileFilter)

	data := [][]string{}
	for _, v := range profiles {
		users, err := getAllUsers(v)
		if err != nil {
			log.Printf("%s %s", v, err.Error())
		}
		for _, u := range users {
			data = append(data, []string{
				v,
				u.EMail,
			})
		}
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Profile", "EMail"})
	table.SetBorder(false)
	table.AppendBulk(data)
	table.Render()
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
