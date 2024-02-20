package main

import "github.com/Oluwatunmise-olat/WaveDeploy/internal"

func main() {
	internal.BootstrapApp()
}

// TODO's

// Security
// Validate Incoming Webhook Delivery Security (sha 256)

// Cli Stuff
// Add cli Auto Suggestion
// Add More Commands
// Auto-Scale Based on metrics from prometheus command
// wave-deploy get all projects
// update envs (should retrigger app redeployment)
// Domains stuff (CNAME)
// Disconnect Github
