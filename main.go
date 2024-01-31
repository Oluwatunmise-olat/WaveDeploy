package main

import (
	"github.com/Oluwatunmise-olat/WaveDeploy/cmd"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal"
)

func main() {
	internal.BootstrapApp()
	cmd.Execute()
}
