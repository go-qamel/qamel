package cmd

import (
	"log"

	ap "github.com/muesli/go-app-paths"
	"github.com/spf13/cobra"
)

var (
	configPath = "config.json"
)

func init() {
	// Get config path in ${XDG_CONFIG_HOME}/qamel/config.json
	var err error
	userScope := ap.NewScope(ap.User, "qamel")
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
