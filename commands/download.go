package commands

import (
	"fmt"
	"github.com/alexgustafsson/systembolaget-api/systembolaget"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"path/filepath"
	"os"
	"strings"
)

func download(sourceName string, context *cli.Context) error {
	output := context.String("output")
	pretty := context.Bool("pretty")
	outputFormat := strings.ToLower(context.String("format"))
	outputExtension := strings.ToLower(filepath.Ext(output))
	if outputFormat != "" {
		outputExtension = "." + outputFormat
	}

	log.Debugf("Attempting to download source %s", sourceName)
	var source systembolaget.Source
	if sourceName == "assortment" {
		source = &systembolaget.Assortment{}
	} else if sourceName == "inventory" {
		source = &systembolaget.Inventory{}
	} else if sourceName == "stores" {
		source = &systembolaget.Stores{}
	}

	err := source.Download()
	if err != nil {
		return err
	}

	log.Debug("Downloaded data, converting to output format")
	var outputBytes []byte
	if outputExtension == ".xml" {
		outputBytes, err = source.ConvertToXML(pretty)
	} else if outputExtension == ".json" {
		outputBytes, err = source.ConvertToJSON(pretty)
	} else if output != "" {
		return fmt.Errorf("Unsupported output extension: %s", outputExtension)
	}
	if err != nil {
		return err
	}

	log.Debug("Converted data, writing to target")
	if output == "" {
		os.Stdout.Write(outputBytes)
	} else {
		err = ioutil.WriteFile(output, outputBytes, 0644)
	}
	if err != nil {
		return err
	}

	return nil
}

func downloadAssortmentCommand(context *cli.Context) error {
	return download("assortment", context)
}

func downloadInventoryCommand(context *cli.Context) error {
	return download("inventory", context)
}

func downloadStoresCommand(context *cli.Context) error {
	return download("stores", context)
}
