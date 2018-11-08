package main

import (
	"flag"
	"fmt"
	"github.com/Inspectools/gowindams"
	"io"
	"log"
	"os"
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
	resourceId := flag.String("r", "", "Unique ID of the resource")
	outputFile := flag.String("o", "", "Path to write the resource to")
	flag.Parse()

	var environments *gowindams.Environments
	environments, err = gowindams.LoadEnvironments(*environmentsConfigFile)
	if err != nil {
		log.Fatalf("Unable to load environments config file from %s:\t%s\n", *environmentsConfigFile, err)
	}

	var env *gowindams.Environment = nil
	for _, e := range *environments {
		if e.Name == *environmentName {
			env = &e
			break
		}
	}

	if env == nil {
		log.Fatalf("Unable to locate environment with name \"%s\"\n", environmentName)
	}

	log.Printf("Loading resource \"%s\"\n", *resourceId)
	rmeta, err := env.ResourceServiceClient().Get(*resourceId)
	if err != nil {
		log.Fatalf("Error loading image resource \"%s\":\t%s\n", resourceId, err)
	}
	log.Printf("Resource ID:\t%s\n", *(rmeta.ResourceId))
	log.Printf("Download URL:\t%s\n", *(rmeta.DownloadURL))
	log.Printf("Zoomify ID:\t%s\n", *(rmeta.ZoomifyId))
	log.Printf("Zoomify URL:\t%s\n", *(rmeta.ZoomifyURL))

	if *outputFile != "" {
		// Download the file
		log.Printf("Downloading resource %s\n", *resourceId)

		indata, err := env.ResourceServiceClient().Download(*resourceId)
		if err != nil {
			log.Fatalf("Unable to download the resource \"%s\":\t%s\n", *resourceId, err)
		}
		imgFile, err := os.Create(*outputFile)
		if err != nil {
			log.Fatalf("Error downloading image %s: %s", *resourceId, err)
		}

		_, err = io.Copy(imgFile, *indata)
		if err != nil {
			log.Fatalf("Error downloading image %s: %s\n", *resourceId, err)
		} else {
			(*indata).Close()
		}

		err = imgFile.Close()
		if err != nil {
			log.Fatalf("Error downloading image %s: %s\n", *resourceId, err)
		}

		log.Printf("Resource %s written to %s\n", *resourceId, *outputFile)
	}
}
