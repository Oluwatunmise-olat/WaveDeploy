package main

import (
	"github.com/Oluwatunmise-olat/WaveDeploy/internal"
)

func main() {
	internal.BootstrapApp()
}

// TODO's

// Cli Stuff -*
// Add cli Auto Suggestion - Improvement 18
// Logs streaming - Improvement 19
// PR Links - Improvement 19
// Auto-Scale Based on metrics from prometheus command - Improvement 20
// better error handling

// stop running project
// delete project +
// get all projects +
// all db actions should revert if deployment fails

// later once the above is done
// 1. check for unique env keys before update and on save
// 2. check app is actually live redeploy
