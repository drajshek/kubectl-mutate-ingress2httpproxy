# kubectl mutate dc2deployment

This is a kubectl plugin to mutate OpenShift deployment config yamls to Kubernetes deployment.apps yamls.

Useful tool if you are migrating from OpenShift 3.11 to any K8s.

## Usage

```bash
kubectl mutate dc2deployment -h
Usage of kubectl-mutate-dc2deployment:
  -debug
        debug flag, it turns verbose on
  -f string
        the name of the input file with deployment config. The default is stdin.
  -json
        json parameter, it accepts JSON format as input. The output will be JSON as well.
  -o string
        the name of the output file of deployment.apps. The default is stdout.
```

The plugin warns if there is any deploymentconfig field that cannot be mutated to deployment.apps.

*ATTENTION:* Debug, errors and warning messages go to stderr. The deployment custom resource is stdout.

## Installation

How to install (requires kubectl):

```bash
mkdir -p ~/go/src/github.com/brito-rafa; cd ~/go/src/github.com/brito-rafa
git clone https://github.com/brito-rafa/kubectl-mutate-dc2deployment.git
cd kubectl-mutate-dc2deployment
make install
```

or

```bash
# download https://github.com/brito-rafa/kubectl-mutate-dc2deployment/archive/v0.1.tar.gz
# wget https://github.com/brito-rafa/kubectl-mutate-dc2deployment/archive/v0.1.tar.gz # need authentication
tar -xvzf kubectl-mutate-dc2deployment-0.1.tar.gz 
mv kubectl-mutate-dc2deployment/bin/kubectl-mutate-dc2deployment /usr/local/bin
chmod +x /usr/local/bin/kubectl-mutate-dc2deployment
```

## Example

Example of execution with YAML:

```bash
 $ kubectl mutate dc2deployment -f example-dc.yaml -o output-example-deployment.yaml
WARN[0000] [kubectl-mutate-dc2deployment] DeploymentConfig.Spec.Triggers unsupported 

$ head output-example-deployment.yaml 
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    kubectl-mutate-dc2deployment/DeploymentConfig.Spec.Test: unsupported
    kubectl-mutate-dc2deployment/DeploymentConfig.Spec.Triggers: unsupported
  creationTimestamp: null
  labels:
    docker-registry: default
  name: docker-registry
(...)
```

Example of execution with JSON (note the `-json` command line parameter):

```bash
$ bin/kubectl-mutate-dc2deployment -f example-dc.json  -o output-example-deployment.json -json
WARN[0000] [kubectl-mutate-dc2deployment] DeploymentConfig.Spec.Triggers unsupported 

$ head output-example-deployment.json 
{
 "kind": "Deployment",
 "apiVersion": "apps/v1",
 "metadata": {
  "name": "docker-registry",
(...)
```