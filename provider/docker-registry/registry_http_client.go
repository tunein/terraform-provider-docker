package docker_registry

import (
	"github.com/tunein/terraform-provider-docker/provider/docker-registry/ecr"
	"github.com/tunein/terraform-provider-docker/provider/docker-registry/gcr"
	"github.com/tunein/terraform-provider-docker/provider/docker-registry/hub"
	"github.com/tunein/terraform-provider-docker/provider/docker-registry/quay"
	"regexp"
)

// 602401143452.dkr.ecr.us-west-2.amazonaws.com/eks/kube-proxy
// quay.io/coreos/kube-state-metrics
// docker.io/newrelic/k8s-metadata-injection
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

func (c RegistryHttpClient) IfImageExist(repo, tag string) error {
	var provider ContainerRegistryProvider
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
	}
	provider.Login()
	return provider.IfImageExist(repo, tag)
}
