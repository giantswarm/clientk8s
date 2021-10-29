package fake

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apiextensionsclientfake "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/fake"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	dynamicfake "k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/kubernetes"
	kubernetesfake "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	clientfake "sigs.k8s.io/controller-runtime/pkg/client/fake" //nolint:staticcheck // v0.6.4 has a deprecation on pkg/client/fake that was removed in later versions

	"github.com/giantswarm/k8sclient/v6/pkg/k8sclient"
)

type Clients struct {
	logger micrologger.Logger

	ctrlClient client.Client
	dynClient  dynamic.Interface
	extClient  apiextensionsclient.Interface
	k8sClient  kubernetes.Interface
	restClient rest.Interface
	restConfig *rest.Config
}

func NewClients(config k8sclient.ClientsConfig, objects ...runtime.Object) (*Clients, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	var err error

	var restConfig *rest.Config
	{
		restConfig = config.RestConfig
	}

	var extClient apiextensionsclient.Interface
	{
		extClient = apiextensionsclientfake.NewSimpleClientset()
	}

	var ctrlClient client.Client
	{
		if config.SchemeBuilder != nil {
			// Extend the global client-go scheme which is used by all the tools under
			// the hood. The scheme is required for the controller-runtime controller to
			// be able to watch for runtime objects of a certain type.
			schemeBuilder := runtime.SchemeBuilder(config.SchemeBuilder)

			err = schemeBuilder.AddToScheme(scheme.Scheme)
			if err != nil {
				return nil, microerror.Mask(err)
			}
		}

		ctrlClient = clientfake.NewClientBuilder().
			WithScheme(scheme.Scheme).
			WithRuntimeObjects(objects...).
			Build()
	}

	var dynClient dynamic.Interface
	{
		dynClient = dynamicfake.NewSimpleDynamicClient(scheme.Scheme, objects...)
	}

	var k8sClient kubernetes.Interface
	{
		k8sClient = kubernetesfake.NewSimpleClientset(objects...)
	}

	c := &Clients{
		logger: config.Logger,

		ctrlClient: ctrlClient,
		dynClient:  dynClient,
		extClient:  extClient,
		k8sClient:  k8sClient,
		restConfig: restConfig,
	}

	return c, nil
}

func (c *Clients) CtrlClient() client.Client {
	return c.ctrlClient
}

func (c *Clients) DynClient() dynamic.Interface {
	return c.dynClient
}

func (c *Clients) ExtClient() apiextensionsclient.Interface {
	return c.extClient
}

func (c *Clients) K8sClient() kubernetes.Interface {
	return c.k8sClient
}

func (c *Clients) RESTClient() rest.Interface {
	return c.restClient
}

func (c *Clients) RESTConfig() *rest.Config {
	return rest.CopyConfig(c.restConfig)
}

func (c *Clients) Scheme() *runtime.Scheme {
	return scheme.Scheme
}
