package flags

import (
	"github.com/spf13/cobra"
)

var projectName string

func InitializeProjectNameFlag(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&projectName, "name", "n", "", "Project name")
	cmd.MarkFlagRequired("name")
}

func GetProjectName() string {
	return projectName
}
