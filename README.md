# Operator-Builder

![Build Status](https://github.com/datainfrahq/operator-builder/actions/workflows/go.yml/badge.svg) [![Slack](https://img.shields.io/badge/slack-brightgreen.svg?logo=slack&label=Community&style=flat&color=%2373DC8C&)](https://launchpass.com/datainfra-workspace)
[![Go Reference](https://pkg.go.dev/badge/github.com/datainfrahq/operator-builder.svg)](https://pkg.go.dev/github.com/datainfrahq/operator-builder)
![GitHub issues](https://img.shields.io/github/issues/datainfrahq/operator-runtime) [![Go Report Card](https://goreportcard.com/badge/github.com/datainfrahq/operator-runtime)](https://goreportcard.com/report/github.com/datainfrahq/operator-runtime)



## Introduction
- Operator builder is a library to build kubernetes operators which adhere to [dsoi-spec](https://github.com/datainfrahq/dsoi-spec).
- The library provides users with high level abstractions and standardisation to focus on building robust reconcile loops for data application with multiple node types.

## Motivation

- At datainfra, we are building cloud native data infrastructure toolings to power self served data platforms, we build a lot of kubernetes operator's for various large distributed systems. Building reconcile loops for different nodetypes/components in a distributed systems is extremely time consuming, error prone and becomes repetitive for mulitple applications. 
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
		builder.ToNewBuilderConfigMap(configMap),
		builder.ToNewBuilderConfigMapHash(configMapHash),
		builder.ToNewBuilderDeploymentStatefulSet(deploymentOrStatefulset),
		builder.ToNewBuilderStorageConfig(storage),
		builder.ToNewBuilderRecorder(builder.BuilderRecorder{Recorder: r.Recorder, ControllerName: "ParseableOperator"}),
		builder.ToNewBuilderContext(builder.BuilderContext{Context: ctx}),
		builder.ToNewBuilderService(service),
		builder.ToNewBuilderStore(*builder.NewStore(r.client, r.commonLabels, cr.Namespace, cr)),
	)
```
- The library makes configmap hashes and passes them as env to deployment and sts objects, which force trigger rollout of pods when the intenral k8s client updates objects.

### Reconcile Objects
- once builder is constructed, the ```ReconcileInterface``` can be called in the build and reoncile objects
```
  // reconcile configmap
	_, err := builder.ReconcileConfigMap()
	if err != nil {
		return err
	}

```
### Interal Store
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
