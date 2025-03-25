# nginx-operator

## Technical requirements
* go version 1.16+
* An operator-sdk binary installed locally
* Make sure your user is authorized with cluster-admin permissions.

## Step 1 : Setting up your project

First, create an empty project directory and cd into it.
```bash
mkdir nginx-operator
cd nginx-operator
```

Now, initialize a boilerplate project structure with the following:

```bash
# we'll use a domain of example.com
# so all API groups will be <group>.example.com
operator-sdk init --domain example.com --repo github.com/saeed-mcu/nginx-operator
```

`--domain` will be used as the prefix of the API group your custom resources will be created in. API groups are a mechanism to group portions of the Kubernetes API. API groups are used internally to version your Kubernetes resources and are thus used for many things. Importantly, you should name your domain to group your resource types in meaningful group(s) for ease of understanding and because these groups determine how access can be controlled to your resource types using RBAC.


## Step 2 : Defining an API
The Operator's API will be the definition of how it is represented within a Kubernetes
cluster. The API is directly translated to a generated CRD, which describes the blueprint
for the custom resource object that users will consume to interact with the Operator.
Therefore, creating this API is a necessary first step before writing other logic for the Operator
Building an Operator API is done by writing a Go struct to represent the object.

```bash
operator-sdk create api --group operator --version v1alpha1 --kind NginxOperator --resource --controller
```
This command does the following:
1. Creates the API types in a new directory called `api/`
2. Defines these types as belonging to the API group `operator.example.com`
3. Creates the initial version of the API named `v1alpha1`
4. Names these types after our Operator, `NginxOperator`
5. Instantiates boilerplate controller code under a new directory called `controllers/`
6. Updates main.go to add boilerplate code for starting the new controller

## Step 3 : Modify Operator types
Looking at nginxoperator_types.go, there are already some empty structs with instructions to fill in additional fields.
The three most important types in this file are `NginxOperator`, `NginxOperatorSpec`, and `NginxOperatorStatus`:

* `NginxOperatorSpec` defines the **desired state**
* `NginxOperatorStatus` defines the **observed state**

Modify  `NginxOperatorSpec` with your desired fields.
```go
type NginxOperatorSpec struct {
	// Port is the port number to expose on the Nginx Pod
	Port *int32 `json:"port,omitempty"`

	// Replicas is the number of deployment replicas to scale
	Replicas *int32 `json:"replicas,omitempty"`

	// ForceRedploy is any string, modifying this field
	// instructs the Operator to redeploy the Operand
	ForceRedploy string `json:"forceRedploy,omitempty"`
}
```
Once the Operator types have been modified, it is sometimes necessary to run `make generate` from the project root.
This updates generated files, such as the zz_generated.deepcopy.go. It is good practice to develop the habit of
regularly running this command whenever making changes to the API, even if it does not always produce any changes.

## Step 4 : Adding resource manifests
With the API types defined, it is now possible to generate an equivalent CRD manifest,
The first step is to generate a CRD from the Go types.
```bash
make manifests
```
The underlying command for `make manifests` is actually calling an additional tool, `controller-gen`
Running `controller-gen` manually is an acceptable way to generate files and code in non-default ways.

The `make manifests` command also creates a corresponding **Role-Based Access Control (RBAC)** role
with can be bound to the Operator's ServiceAccount to give the Operator access to their own custom object.

### Additional manifests
The `ClusterRole` can be conveniently generated with `Kubebuilder tags` in the code
To do that, we need to update the RBAC role for the Operator.
This can be done automatically using Kubebuilder markers on the Reconcile() function
```go
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
```

### Using go:embed to access resources
The `go:embed` marker was included as a compiler directive in Go 1.16 to provide native
resource embedding without the need for external tools such as `go-bindata`.
To start with  this approach, create the resource manifest files under a new directory called assets/.

To do this, we will keep the existing assets/ directory (to use as an importable
Go module path that holds helper functions for loading and processing the files) and place
a new manifests/ directory underneath it (which will hold the actual manifest files).


### what events will trigger this loop to run?
This is set up in SetupWithManager():
```go
func (r *NginxOperatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
            For(&operatorv1alpha1.NginxOperator{}).
            Owns(&appsv1.Deployment{}).
            Complete(r)
}
```
Finally run:
```bash
make manifest
```

## STEP 5: Operator CRD conditions
As part of the Kubernetes API conventions, objects (including custom resources) should
include both a `spec` and `status` field. In the case of Operators, we are using `spec`
as an **input for configuring the Operator's parameters** already.

```go
// NginxOperatorStatus defines the observed state of
NginxOperator
type NginxOperatorStatus struct {
 // Conditions is the list of status condition updates
 Conditions []metav1.Condition `json:"conditions"`
}
```
Then, we will run `make generate` (to update the generated client code) and
`make manifests` (to update the Operator's CRD with the new field),


Now that the Operator's CRD has a field to report the latest status conditions, the code
can be updated to implement them. For this, we can use the `SetStatusCondition()` helper function.
