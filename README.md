# nginx-operator

## Technical requirements
* go version 1.16+
* An operator-sdk binary installed locally
* Make sure your user is authorized with cluster-admin permissions.

## STEP 1 : Setting up your project

First, create an empty project directory and cd into it.
```bash
mkdir nginx-operator
cd nginx-operator
```

Now, initialize a boilerplate project structure with the following:

```bash
operator-sdk init --domain example.com --repo github.com/saeed-mcu/nginx-operator
```

`--domain` will be used as the prefix of the API group your custom resources will be created in. API groups are a mechanism to group portions of the Kubernetes API. API groups are used internally to version your Kubernetes resources and are thus used for many things. Importantly, you should name your domain to group your resource types in meaningful group(s) for ease of understanding and because these groups determine how access can be controlled to your resource types using RBAC.


## STEP2: Defining an API
