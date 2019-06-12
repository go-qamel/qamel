package cmd

import (
	"go/build"
	"log"
	"os"
	fp "path/filepath"

	ap "github.com/muesli/go-app-paths"
	"github.com/spf13/cobra"
)

var (
	qamelDir   = fp.Join("github.com", "RadhiFadlillah", "qamel")
	configPath = "config.json"
)

func init() {
	// Get qamel directory
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	qamelDir = fp.Join(gopath, "src", qamelDir)

	// Get config path in ${XDG_CONFIG_HOME}/qamel/config.json
	var err error
	userScope := ap.NewScope(ap.User, "qamel", "qamel")
	configPath, err = userScope.ConfigPath("config.json")
	if err != nil {
		log.Fatalln(err)
	}
}

// QamelCmd returns the root command for qamel
func QamelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "qamel",
		Short: "qamel is tools and binding for creating GUI app in Go using Qt + QML",
		Args:  cobra.NoArgs,
	}

	cmd.AddCommand(
		buildCmd(),
		dockerCmd(),
		profileCmd(),
	)

	return cmd
}
