module github.com/giantswarm/k8sclient/v5

go 1.14

require (
	github.com/giantswarm/apiextensions/v3 v3.19.1-0.20210302205925-820fb372e72e
	github.com/giantswarm/backoff v0.2.0
	github.com/giantswarm/microerror v0.3.0
	github.com/giantswarm/micrologger v0.5.0
	github.com/google/go-cmp v0.5.4
	k8s.io/api v0.20.2
	k8s.io/apiextensions-apiserver v0.20.2
	k8s.io/apimachinery v0.20.2
	k8s.io/client-go v0.20.2
	sigs.k8s.io/controller-runtime v0.8.2
)

replace sigs.k8s.io/cluster-api => sigs.k8s.io/cluster-api v0.3.11-0.20210302171319-f7351b165992
