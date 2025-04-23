package apiclient

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
	"gopkg.in/yaml.v3"

	"github.com/zydee3/stockdb/internal/api/types/crd"
)

var applyYamlCommand = cli.Command{
	Name:        "apply",
	Usage:       "Apply a YAML file to the StockDB server",
	ArgsUsage:   "<yaml-file>",
	Description: `Apply a YAML file to the StockDB server.`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "file, f",
			Usage: "Path to the YAML file",
		},
	},
	Before: func(c *cli.Context) error {
		filename := c.String("file")

		if filename == "" {
			return cli.NewExitError("no yaml file provided", 1)
		}

		// Check if the filename is provided
		if filename == "" && c.NArg() == 0 {
			return cli.NewExitError("no yaml file provided", 1)
		}

		// check if the file exists
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			return cli.NewExitError(fmt.Sprintf("file %s does not exist", filename), 1)
		}

		return nil
	},
	Action: func(c *cli.Context) error {
		filename := c.String("file")

		crd, err := loadYaml(filename)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		fmt.Println("Parsed DataCollection:", crd)
		fmt.Println("Name:", crd.Metadata.Name)
		fmt.Println("Source Type:", crd.Spec.Source.Type)
		fmt.Println("Securities:", crd.Spec.Targets.Securities)

		return nil
	},
}

func loadYaml(filename string) (*crd.DataCollection, error) {
	fmt.Println("Applying YAML file:", filename)

	// Read the YAML file
	yamlData, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %s", err.Error())
	}

	dataCollection := &crd.DataCollection{}

	// Unmarshal the YAML data into our DataCollection struct
	if err := yaml.Unmarshal(yamlData, dataCollection); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %s", err.Error())
	}

	return dataCollection, nil
}
