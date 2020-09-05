package main

import (
	"fmt"
	"os"
	"strings"
)

var scmProvider = os.Getenv("SCM_PROVIDER")
var cloneUrl = os.Getenv("CLONEURL")
var token = os.Getenv("OAUTH_TOKEN")
var commit = os.Getenv("COMMITID")
var dataDir = "/data"
var repoClonePath = dataDir + "/repo"
var dockerFilePath = ""

// Info should be used to describe the example commands that are about to run.
func Info(format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

func main() {
	CheckoutFromScm()

	if len(os.Getenv("DOCKERFILE_PATH")) != 0 {
		dockerFilePath = " -f " + os.Getenv("DOCKERFILE_PATH") + " "
	}
	var registryFqdn string
	if len(os.Getenv("DOCKER_CLOUDOPS")) != 0 {
		dockerCloudOps = os.Getenv("DOCKER_CLOUDOPS")
		if len(os.Getenv("CLOUD_WRAPPER_HOST_PORT")) != 0 {
			cloudWrapperHost = os.Getenv("CLOUD_WRAPPER_HOST_PORT")
		}
		if strings.Contains(dockerCloudOps, "getDockerLoginFile") {
			registryFqdn = GenerateDockerLoginFile()
			if len(os.Getenv("REPO_NAME")) != 0 && strings.Contains(dockerCloudOps,"getDockerBuildPushFile") {
				GenerateDockerBuildPushFile(registryFqdn, os.Getenv("REPO_NAME"))
			}
		}
	}
}