package main

import (
	"flag"
	"fmt"
	"github.com/Inspectools/gowindams"
	"log"
	user2 "os/user"
)

func main() {

	user, err := user2.Current()
	if err != nil {
		log.Fatalf("Error getting current user:%s\n", err)
	}
	dfltConfigFile := fmt.Sprintf("%s/.windams/environments.yaml", user.HomeDir)

	environmentsConfigFile := flag.String("c", dfltConfigFile, "Environments config file, defaults to ~/.windams/environments.yaml")
	environmentName := flag.String("env", "Local Dev", "Environment to connect to")
	flag.Parse()

	environments, err := gowindams.LoadEnvironments(*environmentsConfigFile)
	if err != nil {
		log.Fatalf("Unable to load environments config file from %s:\t%s\n", *environmentsConfigFile, err)
	}

	env := environments.Find(*environmentName)
	if env == nil {
		log.Fatal("Unable to find environment \"%s\" in file %s", *environmentName, *environmentsConfigFile)
	}
	keys, err := env.ObtainSigningKeys()
	if err != nil {
		log.Fatalf("Unable to load signing keys for environment %s: %s", env.Name, err)
	}
	for kid, _ := range keys {
		log.Printf("Retrieved a key for kid %s", kid)
	}
}