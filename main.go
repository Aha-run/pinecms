//go:generate go run main.go src
package main

import "github.com/xiusin/pinecms/cmd"

func main() {

	// debug.SetCrashOutput() 设置崩溃日志输出

	cmd.Execute()
}
