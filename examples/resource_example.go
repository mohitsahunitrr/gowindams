package main

import (
	"flag"
	"fmt"
	"github.com/Inspectools/gowindams"
	user2 "os/user"
)

func main() {

	user, err := user2.Current()
	if err != nil {
		fmt.Errorf("Error getting current user:%s\n", err)
	}
	dfltConfigFile := fmt.Sprintf("%s/.windams/environments.yaml", user.HomeDir)

	environmentsConfigFile := flag.String("c", dfltConfigFile, "Environments config file, defaults to ~/.windams/environments.yaml")
	environmentName := flag.String("e", "Local Dev", "Environment to connect to")
	resourceId := flag.String("r", "", "Unique ID of the resource")
	flag.Parse()

	var environments *gowindams.Environments
	environments, err = gowindams.LoadEnvironments(*environmentsConfigFile)
	if err != nil {
		fmt.Errorf("Unable to load environments config file from %s:\t%s\n", *environmentsConfigFile, err)
	}

	var env *gowindams.Environment = nil
	for _, e := range *environments {
		if e.Name == *environmentName {
			env = &e
			break
		}
	}

	if env == nil {
		fmt.Errorf("Unable to locate environment with name \"%s\"\n", environmentName)
	}

	fmt.Printf("Loading resource \"%s\"\n", *resourceId)
	rmeta, err := env.ResourceServiceClient().Get(*resourceId)
	if err != nil {
		fmt.Errorf("Error loading image resource \"%s\":\t%s\n", resourceId, err)
	}
	fmt.Printf("Resource ID:\t%s\n", *(rmeta.ResourceId))
	fmt.Printf("Download URL:\t%s\n", *(rmeta.DownloadURL))
	fmt.Printf("Zoomify ID:\t%s\n", *(rmeta.ZoomifyId))
	fmt.Printf("Zoomify URL:\t%s\n", *(rmeta.ZoomifyURL))
}
