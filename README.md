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
