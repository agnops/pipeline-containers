package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	reqHttp "net/http"
	"os"
	"strings"
)

var cloudWrapperHost = "localhost:5000"
var dockerCloudOps = ""

type CloudRegistryDetails struct {
	User    	string    `json:"user"`
	Password    string    `json:"password"`
	Url    		string    `json:"url"`
}

func GenerateDockerLoginFile() string {

	response, err := reqHttp.Get("http://"+cloudWrapperHost+"/get_docker_login")
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var crd CloudRegistryDetails
	_ = json.Unmarshal(responseData, &crd)

	dockerLoginCmd := `export DOCKER_PASS="`+crd.Password+`"
export DOCKER_USER="`+crd.User+`"
export DOCKER_URL="`+crd.Url+`"
export REGISTRY_DOMAIN="`+strings.Replace(crd.Url, "https://", "", -1)+`"
echo "$DOCKER_PASS" | docker login -u "$DOCKER_USER" "$DOCKER_URL" --password-stdin`

	f, err := os.Create(dataDir + "/dockerLogin.sh")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	_, err = f.WriteString(dockerLoginCmd)
	if err != nil {
		log.Println(err)
		_ = f.Close()
		os.Exit(1)
	}
	return strings.Replace(crd.Url, "https://", "", -1)
}

func GenerateDockerBuildPushFile(registryFqdn string, repoName string) {

	if strings.Contains(dockerCloudOps,"ensureRepoExistence") {

		requestBody, err := json.Marshal(map[string]string{
			"repo_name": repoName,
		})
		if err != nil {
			log.Println(err)
			return
		}

		response, err := reqHttp.Post("http://"+cloudWrapperHost+"/create_docker_repository", "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		defer response.Body.Close()
	}
	dockerBuildPushCmd := `docker build `+dockerFilePath+`-t `+registryFqdn+`/`+repoName+`:`+commit+` .`+`
docker push `+registryFqdn+`/`+repoName+`:`+commit

	f, err := os.Create(dataDir + "/dockerBuildPush.sh")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	_, err = f.WriteString(dockerBuildPushCmd)
	if err != nil {
		log.Println(err)
		_ = f.Close()
		os.Exit(1)
	}

	deploymentEnvs := "export REGISTRY_DOMAIN="+registryFqdn

	f, err = os.Create(dataDir + "/deploymentEnvs")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	_, err = f.WriteString(deploymentEnvs)
	if err != nil {
		log.Println(err)
		_ = f.Close()
		os.Exit(1)
	}
}