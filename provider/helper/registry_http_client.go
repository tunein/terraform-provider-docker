package helper

import (
	"fmt"
	"net/http"
	"strings"
)

// 602401143452.dkr.ecr.us-west-2.amazonaws.com/eks/kube-proxy
// quay.io/coreos/kube-state-metrics
// docker.io/docker.io/newrelic/k8s-metadata-injection
// k8s.gcr.io/metrics-server/metrics-server - https://k8s.gcr.io/v2/external-dns/external-dns/manifests/v0.9.0

type RegistryHttpClient struct {
	ecrAuthStr  string
	hubLogin    string
	hubPassword string
}

func NewRegistryHttpClient(ecrAuthStr, hubLogin, hubPassword string) *RegistryHttpClient {
	return &RegistryHttpClient{
		ecrAuthStr:  ecrAuthStr,
		hubLogin:    hubLogin,
		hubPassword: hubPassword,
	}
}

type ContainerRegistryProvider interface {
	Login() error
	IfImageExist(repo, tag string) error
}

type ECR struct {
	authStr string
}

type GCR struct {
}

type Quay struct {
}

func (ecr ECR) IfImageExist(repo, tag string) bool {
	params := strings.Split(repo, "/")
	url := fmt.Sprintf("https://%s/v2/%s/%s/manifests/%s", params[0], params[1], params[2], tag)
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("Authorization", "Basic "+ecr.authStr)
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	fmt.Println(response.Status)
	return false
}

func (gcr GCR) IfImageExist(repo, tag string) bool {
	params := strings.Split(repo, "/")
	url := fmt.Sprintf("https://k8s.gcr.io/v2/%s/%s/manifests/%s", params[1], params[2], tag)
	return HttpCheck(url)
}

func (quay Quay) IfImageExist(repo, tag string) bool {
	params := strings.Split(repo, "/")
	url := fmt.Sprintf("https://quay.io/v2/%s/%s/manifests/%s", params[1], params[2], tag)
	return HttpCheck(url)
}

func (c RegistryHttpClient) IfImageExist(repo, tag string) error {
	var provider ContainerRegistryProvider
	switch {
	case strings.Contains(repo, "docker.io"):
		provider = NewDockerHubClient(c.hubLogin, c.hubPassword)
	}
	provider.Login()
	return provider.IfImageExist(repo, tag)
}

func HttpCheck(url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return false
	}
	return true
}
