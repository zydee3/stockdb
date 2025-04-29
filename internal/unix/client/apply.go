package client

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/urfave/cli"
	"gopkg.in/yaml.v3"

	"github.com/zydee3/stockdb/internal/common/crd"
	"github.com/zydee3/stockdb/internal/common/logger"
	"github.com/zydee3/stockdb/internal/unix/messages"
	"github.com/zydee3/stockdb/internal/unix/socket"
)

//nolint:gochecknoglobals // gochecknoglobals
var applyYamlCommand = cli.Command{
	Name:        "apply",
	Usage:       "Apply a YAML file to the StockDB server",
	ArgsUsage:   "<yaml-file>",
	Description: `Apply a YAML file to the StockDB server.`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "file, f",
			Usage: "(-f <yaml-file>)",
		},
	},
	Before: onBefore,
	Action: onAction,
}

func loadDataCollectionYaml(filename string) (*crd.DataCollection, error) {
	// Read the YAML file
	yamlData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	dataCollection := &crd.DataCollection{}

	// Unmarshal the YAML data into our DataCollection struct
	if unmarshalError := yaml.Unmarshal(yamlData, dataCollection); unmarshalError != nil {
		return nil, unmarshalError
	}

	return dataCollection, nil
}

func onBefore(c *cli.Context) error {
	filename := c.String("file")

	if filename == "" {
		return cli.NewExitError("no yaml file provided", 1)
	}

	// check if the file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return cli.NewExitError(fmt.Sprintf("file %s does not exist", filename), 1)
	}

	return nil
}

func onAction(c *cli.Context) error {
	filename := c.String("file")

	crd, err := loadDataCollectionYaml(filename)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	conn, err := net.Dial("unix", socket.SocketPath)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	defer conn.Close()

	cmd := messages.Command{
		Type:       messages.CommandTypeApply,
		Parameters: make(map[string]string),
		Data:       crd,
	}

	encoder := json.NewEncoder(conn)
	if encodeError := encoder.Encode(cmd); encodeError != nil {
		logger.Error("%w", encodeError)
		return cli.NewExitError(encodeError, 1)
	}

	// Receive and parse response
	response := messages.Response{}

	decoder := json.NewDecoder(conn)
	if decodeError := decoder.Decode(&response); decodeError != nil {
		return cli.NewExitError(decodeError, 1)
	}

	logger.Info("Response received from server:", response)

	return nil
}
