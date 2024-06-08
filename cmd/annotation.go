package cmd

import (
	annotation "github.com/YReshetko/go-annotation/pkg"
	"github.com/spf13/cobra"

	_ "github.com/xiusin/pinecms/cmd/util/annotations/rest"
)

var annotationsCmd = &cobra.Command{
	Use:   "src",
	Short: "注解路由生成",
	Run: func(cmd *cobra.Command, args []string) {
		annotation.Process()
	},
}
