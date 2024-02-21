package main

import "github.com/Oluwatunmise-olat/WaveDeploy/internal"

func main() {
	internal.BootstrapApp()
	//	add extra check for spa or web api
	// use varying tmpl files
}

// TODO's

// Security
// Validate Incoming Webhook Delivery Security (sha 256)

// Cli Stuff
// Add cli Auto Suggestion - improvement
// Add More Cli Commands - improvement
// Auto-Scale Based on metrics from prometheus command -
// update envs (should retrigger app redeployment)
// Domains stuff (CNAME)
// Disconnect Github
