package client

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/urfave/cli"
	"github.com/zydee3/stockdb/internal/unix/messages"
	"gopkg.in/yaml.v3"

	"github.com/zydee3/stockdb/internal/unix/socket"
	"github.com/zydee3/stockdb/internal/unix/types/crd"

	"github.com/zydee3/stockdb/internal/common/logger"
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

		conn, err := net.Dial("unix", socket.SOCKET_PATH)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("failed to connect to socket: %s", err.Error()), 1)
		}

		defer conn.Close()

		cmd := messages.Command{
			Type:       messages.CommandTypeApply,
			Parameters: make(map[string]string),
			Data:       crd,
		}

		encoder := json.NewEncoder(conn)
		if err := encoder.Encode(cmd); err != nil {
			logger.Errorf("Error encoding command: %s", err.Error())
			return cli.NewExitError(err.Error(), 1)
		}

		// Receive and parse response
		response := messages.Response{}

		decoder := json.NewDecoder(conn)
		if err := decoder.Decode(&response); err != nil {
			return cli.NewExitError(fmt.Sprintf("failed to decode response: %s", err.Error()), 1)
		}

		logger.Info("Response received from server:", response)

		return nil
	},
}

func loadYaml(filename string) (*crd.DataCollection, error) {
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
