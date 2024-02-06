package main

import (
	"github.com/Oluwatunmise-olat/WaveDeploy/cmd"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal"
)

func main() {
	internal.BootstrapApp()
	cmd.Execute()
}

// http handler
// webhook handler

// TODO's
// Connect Github
// domains stuff (CNAME)
// Get project envs wave-deploy envs <project>
// Set secret
// wave-deploy create <project>
// wave-deploy get all projects
// scale command
// Get running processes
// Service Discovery
