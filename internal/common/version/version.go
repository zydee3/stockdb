package version

import "fmt"

//nolint:gochecknoglobals // gochecknoglobals
var (
	version   = ""
	gitCommit = ""
)

func GetVersion() string {
	return fmt.Sprintf("%s-%s", version, gitCommit)
}
