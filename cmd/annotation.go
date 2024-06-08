package cmd

import (
	"github.com/spf13/cobra"

	annotation "github.com/xiusin/go-annotation/pkg"
	_ "github.com/xiusin/pinecms/cmd/util/annotations/rest"
)

var annotationsCmd = &cobra.Command{
	Use:   "src",
	Short: "注解路由生成",
	Run: func(cmd *cobra.Command, args []string) {
		annotation.Process()
	},
}
