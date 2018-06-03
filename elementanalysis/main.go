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

	"github.com/spf13/viper"
)

// ElementOAS holds both the Element and its OAS
type ElementOAS struct {
	Element    ce.Element
	OAS        []byte
	Error      error
	Methods    map[string]int
	Paths      []string
	JSONSchema []string
}

var (
	// routefilters are the CE added routes to remove from pathlist
	routefilters = []string{"/bulk", "/{objectName}", "/objects", "/ping"}
	debug        bool
)

func init() {
	flag.BoolVar(&debug, "debug", false, "output debug information")
	flag.Parse()
}

func main() {

	loadConfigFile()
	profile := "default"
	base := viper.Get(profile + ".base").(string)
	auth := fmt.Sprintf("User %s, Organization %s",
		viper.Get(profile+".user").(string),
		viper.Get(profile+".org").(string),
	)

	if debug {
		log.Println(base)
		log.Println(debug, flag.Args())
	}

	filterelements := os.Args[1:]

	bodybytes, status, _, err := ce.GetAllElements(base, auth)
	if err != nil {
		fmt.Printf("Can't connect to %s, %s\n", base, err)
		os.Exit(1)
	}
	if status != 200 {
		fmt.Printf("Non-200 status returned, %v\n", status)
		os.Exit(1)
	}

	var allelements []ce.Element
	var elements []ce.Element
	json.Unmarshal(bodybytes, &allelements)

	if len(filterelements) != 0 {
		// select only these elements
		var foundelements []string
		for _, e := range allelements {
			for _, n := range filterelements {
				if e.Key == n {
					elements = append(elements, e)
					foundelements = append(foundelements, n)
				}
			}
		}
		if debug {
			log.Printf("%v Elements %v\n", len(elements), foundelements)
		}
	} else { // all elements, but filter out private and beta Elements
		for _, e := range allelements {
			if e.Private != true {
				if e.Beta != true {
					elements = append(elements, e)
				}
			}
		}
		if debug {
			log.Println(len(elements), "Elements (non-private, non-beta)")
		}
	}
	if len(elements) == 0 {
		os.Exit(0)
	}

	// source: list of Elements
	// sink/consumer: list of Elements + OAS

	// get all OAS
	oaschannel := make(chan ElementOAS, len(elements))
	// ... by creating a goroutine to make the HTTP call to /docs
	for _, e := range elements {
		go func(e ce.Element) {
			var oas ElementOAS
			oas.Element = e
			bodybytes, status, _, err := ce.GetElementOAI(base, auth, strconv.Itoa(e.ID))
			if err != nil {
				oas.Error = err
			}
			if status != 200 {
				oas.Error = fmt.Errorf("Non-200 response %v", status)
			}
			oas.OAS = bodybytes
			oaschannel <- oas
		}(e)
	}

	// as the channel receives Element OAS, process them
	var items []ElementOAS
	var num int
	for i := range oaschannel {
		i.Methods = make(map[string]int)
		if i.Error != nil {
			fmt.Println(num, i.Element.Name, err.Error())
		} else {
			var oas interface{}
			err := json.Unmarshal(i.OAS, &oas)
			if err != nil {
				fmt.Println(num, i.Element.Name, "couldn't create interface from bytes")
				break
			}
			// obtain a list of paths, not on excluded route list (routefilters)
			p := oas.(map[string]interface{})["paths"]
			routes := p.(map[string]interface{})
			var routekeys []string
			for k := range routes {
				if !isInList(k, routefilters) {
					routekeys = append(routekeys, k)
				}
			}
			// count methods per path
			for _, k := range routekeys {
				methods := (routes[k]).(map[string]interface{})
				for m := range methods {
					if c, ok := i.Methods[m]; ok {
						i.Methods[m] = c + 1
					} else {
						i.Methods[m] = 1
					}
				}
			}
			/*
				// remove paths that have {id} in them from the final path count
				var finalroutekeys []string
				for _, p := range routekeys {
					if !isInList(p, []string{"{id}"}) {
						finalroutekeys = append(finalroutekeys, p)
					}
				}
				i.Paths = finalroutekeys
			*/
			i.Paths = routekeys

			// list of schemas
			var schema []string
			d := oas.(map[string]interface{})["definitions"]
			if definitions, ok := d.(map[string]interface{}); ok {
				for d := range definitions {
					if !strings.Contains(d, "swagger") {
						schema = append(schema, d)
					}
				}
			}
			i.JSONSchema = schema

			// collect Element OAS
			items = append(items, i)

			if debug {
				fmt.Println(num, i.Element.Name)
				fmt.Printf("\t%+v\n", i.Methods)
				fmt.Printf("\t%+v\n", routekeys)
				fmt.Printf("\t%+v\n", schema)
			}
		}
		num++
		if num == len(elements) {
			close(oaschannel)
		}
	}

	// Collate totals
	var resources int
	var allschema int
	allmethods := make(map[string]int)
	for _, v := range items {
		for m := range v.Methods {
			if c, ok := v.Methods[m]; ok {
				allmethods[m] += c
			} else {
				allmethods[m] = 0
			}
		}
		resources += len(v.Paths)
		allschema += len(v.JSONSchema)
	}

	// Output
	format := "%6v %5v\n"
	var alphamethod []string
	for m := range allmethods {
		alphamethod = append(alphamethod, m)
	}
	sort.Strings(alphamethod)
	for m, c := range allmethods {
		fmt.Printf(format, strings.ToUpper(m), c)
	}
	fmt.Printf(format, "Paths", resources)
	fmt.Printf(format, "Schema", allschema)
}

// isInList checks to see if the string (item) is in the list of strings
// via strings.Contains()
func isInList(item string, list []string) bool {
	for _, route := range list {
		if strings.Contains(item, route) {
			return true
		}
	}
	return false
}

// loadConfigFile loads a toml config file
// looking for ~/.config/ce/cectl.toml
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
