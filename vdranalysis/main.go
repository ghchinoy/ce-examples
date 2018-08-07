package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ghchinoy/ce-go/ce"
	"github.com/spf13/viper"
)

// Creds is an object for CE credentials
type Creds struct {
	Base    string
	Org     string
	User    string
	Auth    string
	Profile string
}

var profile string

func init() {
	flag.StringVar(&profile, "profile", "default", "profile to use")
	flag.Parse()
}

func main() {

	args := flag.Args()

	// vdranalysis <resource_name> --profile <profile>
	if len(args) < 1 {
		fmt.Println("Missing resource")
		fmt.Println("Usage: vdranalysis <resource>")
		os.Exit(1)
	}
	log.Println(profile, args[0])
	vdr := args[0]

	// profile
	loadConfigFile()
	creds := setProfile(profile)

	err := getAllVDRInfo(vdr, creds)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func getAllVDRInfo(vdr string, creds Creds) error {

	// get the Resource
	resourceBytes, status, _, err := ce.GetResourceDefinition(creds.Base, creds.Auth, vdr, false)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("Non-200 status: %v", status)
	}
	var obj ce.CommonResource
	err = json.Unmarshal(resourceBytes, &obj)
	if err != nil {
		return err
	}
	// get Elements associated with the Resource
	bodybytes, status, _, err := ce.GetTransformationAssocation(creds.Base, creds.Auth, vdr)
	if err != nil {
		return err
	}
	if status != 200 {
		return err
	}
	var associations []ce.AccountElement
	err = json.Unmarshal(bodybytes, &associations)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", obj.Fields)
	fmt.Printf("%+v\n", associations)

	return nil
}

// setProfile sets up the profile info
func setProfile(profile string) Creds {
	base := viper.Get(profile + ".base").(string)
	user := viper.Get(profile + ".user").(string)
	org := viper.Get(profile + ".org").(string)
	auth := fmt.Sprintf("User %s, Organization %s", user, org)

	creds := Creds{base, user, org, auth, profile}

	return creds
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
