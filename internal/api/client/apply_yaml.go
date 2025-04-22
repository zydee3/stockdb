package apiclient

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
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
		fmt.Println("Applying YAML file:", filename)

		// Read the YAML file
		yamlData, err := os.ReadFile(filename)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("failed to read file: %v", err), 1)
		}

		// verify the yaml data
		// send the yaml data to the server

		// // Create a map to store the YAML content
		// var data map[string]any

		// // Unmarshal the YAML data using yaml.v3
		// err = yaml.Unmarshal(yamlData, &data)
		// if err != nil {
		// 	return cli.NewExitError(fmt.Sprintf("failed to parse YAML: %v", err), 1)
		// }

		// // Print out the parsed data for debugging
		// fmt.Printf("yaml contents: %+v\n", data)

		return nil
	},
}
