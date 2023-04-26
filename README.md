<h2 align="center">
  <picture>
    <img alt="DataInfra Logo" src="https://raw.githubusercontent.com/datainfrahq/.github/main/images/logo.svg" width="500" height="100">
  </picture>
  <br>
  Operator Runtime
  </br>
</h2>


<div align="center">

![Build Status](https://github.com/datainfrahq/operator-runtime/actions/workflows/go.yml/badge.svg) [![Slack](https://img.shields.io/badge/slack-brightgreen.svg?logo=slack&label=Community&style=flat&color=%2373DC8C&)](https://launchpass.com/datainfra-workspace)
[![Go Reference](https://pkg.go.dev/badge/github.com/datainfrahq/operator-runtime.svg)](https://pkg.go.dev/github.com/datainfrahq/operator-runtime)
![GitHub issues](https://img.shields.io/github/issues/datainfrahq/operator-runtime) [![Go Report Card](https://goreportcard.com/badge/github.com/datainfrahq/operator-runtime)](https://goreportcard.com/report/github.com/datainfrahq/operator-runtime)

</div>

Operator runtime is a Go-based library that facilitates the development of Kubernetes operators that conform to the [Dsoi-Spec](https://github.com/datainfrahq/dsoi-spec). This library offers high-level abstractions and standardization to simplify the process of building resilient reconcile loops for data applications with multiple node types. Operators created with operator runtime benefit from consistent, well-defined behavior and functionality.
Operator's built using operator runtime
- [Control Plane For Apache Pinot On Kubernetes](https://github.com/datainfrahq/pinot-control-plane-k8s)
- [Control Plane For Parseable On Kubernetes](https://github.com/parseablehq/operator)

## :dart: Motivation

- At DataInfra, we specialize in developing data-centric control planes to facilitate self-service data platforms. Our team has extensive experience building Kubernetes operators for distributed systems. However, creating reconcile loops for various node types and components can be a tedious and error-prone task, particularly when dealing with multiple applications. Unfortunately, there are few abstractions available that can be readily leveraged by the reconciliation loops. To address this challenge, we have developed a unique approach that focuses on application building blocks. By leveraging the operator runtime, we have abstracted out the underlying Kubernetes object reconciliation internals. This approach has enabled us to streamline our development process, reduce errors, and enhance the overall efficiency of our Kubernetes operator implementations.

Using this library we introduce standardisation across
 
- **Builder Abstractions** - Our library includes robust abstractions that enable streamlined building of Kubernetes objects and reconciliation processes.

- **State Change Triggering** - We utilize hashes to trigger reconciliation on state changes, improving the accuracy and efficiency of the process.

- **Internal Store** - Our library includes an internal store that reduces the number of Kubernetes API calls required, resulting in improved performance and resource utilization.

- **Event Emitters** - We have incorporated event emitters that provide real-time feedback on system activity, enhancing visibility and facilitating effective troubleshooting.

## :bricks: Abstractions

- This library abstracts out controller runtime ```client.Client``` by wrapping it with CRUD methods and inbuilt event recorders.
- To build objects initalise the ```Builder``` and leverage the ```ReconcileInterface``` to reconcile objects.

### Example - Build and Reconcile a ConfigMap Object

- Operator builder exposes k8s objects, these objects can be built and passed to the NewBuilder. Example:
```
	getOwnerRef := makeOwnerRef(
		env.APIVersion,
		env.Kind,
		env.Name,
		env.UID,
	)

	cm := makeEnvConfigMap(env, client, getOwnerRef, env.Spec)

	build := builder.NewBuilder(
		builder.ToNewBuilderConfigMap([]builder.BuilderConfigMap{*cm}),
		builder.ToNewBuilderRecorder(builder.BuilderRecorder{Recorder: record, ControllerName: "envoperator"}),
		builder.ToNewBuilderContext(builder.BuilderContext{Context: ctx}),
		builder.ToNewBuilderStore(
			*builder.NewStore(client, map[string]string{"app": env.Name}, env.Namespace, env),
		),
	)

	resp, err := build.ReconcileConfigMap()
	if err != nil {
		return err
	}
```

- Construct a configmap and owner ref function

```
func makeEnvConfigMap(
	env *v1.Environment,
	client client.Client,
	ownerRef *metav1.OwnerReference,
	data interface{},
) *builder.BuilderConfigMap {

	dataByte, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}

	configMap := &builder.BuilderConfigMap{
		CommonBuilder: builder.CommonBuilder{
			ObjectMeta: metav1.ObjectMeta{
				Name:      env.GetName(),
				Namespace: env.GetNamespace(),
			},
			Client:   client,
			CrObject: env,
			OwnerRef: *ownerRef,
		},
		Data: map[string]string{
			"data": string(dataByte),
		},
	}

	return configMap
}

// create owner ref ie parseable tenant controller
func makeOwnerRef(apiVersion, kind, name string, uid types.UID) *metav1.OwnerReference {
	controller := true

	return &metav1.OwnerReference{
		APIVersion: apiVersion,
		Kind:       kind,
		Name:       name,
		UID:        uid,
		Controller: &controller,
	}
}

```

## :stethoscope: Support

- For questions and feedback please feel free to reach out to us on [Slack ↗︎](https://launchpass.com/datainfra-workspace).
- For bugs, please create issue on [GitHub ↗︎](https://github.com/datainfrahq/operator-runtime/issues).
- For commercial support and consultation, please reach out to us at [`hi@datainfra.io` ↗︎](mailto:hi@datainfra.io).
