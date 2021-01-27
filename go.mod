module github.com/imjasonh/tektoncd-concurrency

go 1.15

require (
	github.com/tektoncd/pipeline v0.20.1-0.20210127155341-d4a3aa2ec86e
	k8s.io/api v0.20.2
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	knative.dev/pkg v0.0.0-20210125222030-6040b3af4803
)

replace k8s.io/client-go => k8s.io/client-go v0.20.2
