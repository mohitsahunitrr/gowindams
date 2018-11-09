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

const filePathTempl = "%s/%s%s"

func main() {

	user, err := user2.Current()
	if err != nil {
		log.Fatalf("Error getting current user:%s\n", err)
	}
	dfltConfigFile := fmt.Sprintf("%s/.windams/environments.yaml", user.HomeDir)

	environmentsConfigFile := flag.String("c", dfltConfigFile, "Environments config file, defaults to ~/.windams/environments.yaml")
	environmentName := flag.String("env", "Local Dev", "Environment to connect to")
	outDir := flag.String("o", "", "Directory into which catagorization directories will be created and images will be downloaded.")
	siteId := flag.String("s", "", "Unique ID of the site to download images for")
	flag.Parse()

	if *outDir == "" {
		log.Fatal("Output directory is required")
	}
	damagedDir := *outDir + "/" + "damaged"
	undamagedDir := *outDir + "/" + "undamaged"
	err = os.MkdirAll(damagedDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Unable to create directory %s", damagedDir)
	}
	err = os.MkdirAll(undamagedDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Unable to create directory %s", undamagedDir)
	}
	if *siteId == "" {
		log.Fatal("Site ID is required")
	}

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

	ierparams := gowindams.InspectionEventResourceSearchCriteria{}
	rparams := gowindams.ResourceSearchCriteria{
		SiteId: siteId,
	}
	resources, err := env.ResourceServiceClient().Search(&rparams)
	if err != nil {
		log.Fatalf("Unable to search for resources belonging to the site \"%s\": %s", *siteId, err)
	}

	var ext string
	var filepath string
	var ierlist []gowindams.InspectionEventResource
	var input *io.ReadCloser
	var target *os.File
	for _, rmeta := range resources {
		ext = extension(rmeta.ContentType)
		if ext != "" {
			ierparams.ResourceId = rmeta.ResourceId
		}
		ierlist, err = env.InspectionEventResourceServiceClient().Search(&ierparams)
		if err != nil {
			log.Fatalf("Unable to search for inspection event resources for resource \"%s\": %s", *rmeta.ResourceId, err)
		}
		if len(ierlist) == 0 {
			filepath = fmt.Sprintf(filePathTempl, undamagedDir, *rmeta.ResourceId, ext)
		} else {
			filepath = fmt.Sprintf(filePathTempl, damagedDir, *rmeta.ResourceId, ext)
		}
		if _, err = os.Stat(filepath); os.IsNotExist(err) {
			target, err = os.Create(filepath)
			if err != nil {
				log.Fatalf("Unable to create the file \"%s\" to copy to: %s", filepath, err)
			}
			input, err = env.ResourceServiceClient().Download(*rmeta.ResourceId)
			if err != nil {
				log.Fatal("Unable to download resource \"%s\": %s", *rmeta.ResourceId, err)
			}
			_, err = io.Copy(target, *input)
			if err != nil {
				log.Fatal("Unable to write contents of resource \"%s\": %s", *rmeta.ResourceId, err)
			}
			target.Close()
			(*input).Close()
			log.Printf("The resource \"%s\" has been downloaded.", *rmeta.ResourceId)
		} else {
			log.Printf("The resource \"%s\" has already been downloaded.", *rmeta.ResourceId)
		}
	}
}

func extension(contentType *string) string {
	if contentType != nil {
		if "image/gif" == *contentType {
			return ".gif"
		}
		if "image/jpeg" == *contentType {
			return ".jpg"
		}
		if "image/png" == *contentType {
			return ".png"
		}
	}
	return ""
}