package docker_registry

import (
	"github.com/tunein/terraform-provider-docker/provider/docker-registry/ecr"
	"github.com/tunein/terraform-provider-docker/provider/docker-registry/gcr"
	"github.com/tunein/terraform-provider-docker/provider/docker-registry/hub"
	"github.com/tunein/terraform-provider-docker/provider/docker-registry/quay"
	"regexp"
)

type HttpClient struct {
	ecrAuthStr  string
	hubLogin    string
	hubPassword string
}

func NewRegistryHttpClient(ecrAuthStr, hubLogin, hubPassword string) *HttpClient {
	return &HttpClient{
		ecrAuthStr:  ecrAuthStr,
		hubLogin:    hubLogin,
		hubPassword: hubPassword,
	}
}

type RegistryProvider interface {
	Login() error
	DoesImageExist(repo, tag string) (bool, error)
}

func (c HttpClient) DoesImageExist(repo, tag string) (bool, error) {
	provider := c.getProvider(repo)
	err := provider.Login()
	if err != nil {
		return false, err
	}
	return provider.DoesImageExist(repo, tag)
}

func (c HttpClient) getProvider(repo string) RegistryProvider {
	var provider RegistryProvider

	regexpHub := regexp.MustCompile(`^docker.io`)
	regexpECR := regexp.MustCompile(`^[0-9]+.dkr.ecr.[a-z-0-9]+.amazonaws.com`)
	regexpGCR := regexp.MustCompile(`^k8s.gcr.io`)
	regexpQuay := regexp.MustCompile(`^quay.io`)

	switch {
	case regexpHub.MatchString(repo):
		provider = hub.NewDockerHubClient(c.hubLogin, c.hubPassword)
	case regexpECR.MatchString(repo):
		provider = ecr.NewEcrClient(c.ecrAuthStr)
	case regexpGCR.MatchString(repo):
		provider = gcr.NewGCRClient()
	case regexpQuay.MatchString(repo):
		provider = quay.NewQuayClient()
	default:
		provider = hub.NewDockerHubClient(c.hubLogin, c.hubPassword)
	}

	return provider
}
