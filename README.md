<h2 align="center">
  <picture>
    <img alt="DataInfra Logo" src="https://raw.githubusercontent.com/datainfrahq/.github/main/images/logo.svg">
  </picture>
  <br>
  Operator Runtime
</h2>


<div align="center">

![Build Status](https://github.com/datainfrahq/operator-runtime/actions/workflows/go.yml/badge.svg) [![Slack](https://img.shields.io/badge/slack-brightgreen.svg?logo=slack&label=Community&style=flat&color=%2373DC8C&)](https://launchpass.com/datainfra-workspace)
[![Go Reference](https://pkg.go.dev/badge/github.com/datainfrahq/operator-runtime.svg)](https://pkg.go.dev/github.com/datainfrahq/operator-runtime)
![GitHub issues](https://img.shields.io/github/issues/datainfrahq/operator-runtime) [![Go Report Card](https://goreportcard.com/badge/github.com/datainfrahq/operator-runtime)](https://goreportcard.com/report/github.com/datainfrahq/operator-runtime)

</div>

Operator runtime is a library to build kubernetes operators which adhere to [Dsoi-Spec](https://github.com/datainfrahq/dsoi-spec). This library provides  high level abstractions and standardisation to focus on building robust reconcile loops for data application with multiple node types. 
Operator's built using operator runtime
- [Control Plane For Apache Pinot On Kubernetes](https://github.com/datainfrahq/pinot-control-plane-k8s)
- [Control Plane For Parseable On Kubernetes](https://github.com/parseablehq/operator)

## :dart: Motivation

- At DataInfra, we are building data centric control planes to power self served data platforms, we build a lot of kubernetes operator's for distributed systems. Building reconcile loops for different nodetypes/components is extremely time consuming, error prone and becomes repetitive for mulitple applications. There arn't any useful abstractions which can be consumed by the reconcilation loops, while building operator's we wanted to focus on application building blocks, using the operator runtime we abstracted out the underlying k8s object reconcilation internals.

Using this library we introduce standardisation across 

    1. Builder Abstractions to build kubernetes objects and reconcile.
    2. Triggering Reconcilation on state changes using hashes.
    3. Building an internal store for reducing k8s API calls.
    4. Event Emitters 

## Abstractions

- This library abstracts out controller runtime ```client.Client``` by wrapping it with CRUD methods and inbuilt event recorders.
- To build objects initalise the ```Builder``` and leverage the ```ReconcileInterface``` to reconcile objects.

### Build Objects 
- Operator builder exposes k8s objects, these objects can be built and passed to the NewBuilder. Example:
```
  // construct builder
	builder := builder.NewBuilder(
		builder.ToNewBuilderConfigMap(pinotConfigMap),
		builder.ToNewBuilderDeploymentStatefulSet(pinotDeploymentOrStatefulset),
		builder.ToNewBuilderStorageConfig(pinotStorage),
		builder.ToNewBuilderRecorder(builder.BuilderRecorder{Recorder: r.Recorder, ControllerName: "pinotOperator"}),
		builder.ToNewBuilderContext(builder.BuilderContext{Context: ctx}),
		builder.ToNewBuilderService(pinotService),
		builder.ToNewBuilderStore(*builder.NewStore(ib.client, ib.commonLabels, pt.Namespace, pt)),
	)
```
- The library makes configmap hashes and passes them as env to deployment and sts objects, which force trigger rollout of pods when the intenral k8s client updates objects.

### Reconcile Objects
- once builder is constructed, the ```ReconcileInterface``` can be called in the build and reoncile objects
```
  // reconcile configmap
	result, err := builder.ReconcileConfigMap()
	if err != nil {
		return err
	}

```
### Internal Store
- Reconcile store
```
if err := builder.ReconcileStore(); err != nil {
	return err
}
```
- Interal Store is an in memory which maintains key values for an object name and its kind.
  This store acts as an abstracted desired state, and is used as a matchmaker against the current state of CR.
```
type InternalStore struct {
	ObjectNameKind map[string]string
	CommonBuilder
}
```

