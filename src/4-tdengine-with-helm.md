# Setup TDengine Cluster with helm

Is it simple enough? Let's do something more.

## Install Helm

```sh
curl -fsSL -o get_helm.sh \
  https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3
chmod +x get_helm.sh
./get_helm.sh
```

Helm will use kubectl and the kubeconfig setted in chapter 1.

## Install TDengine Chart

Download TDengine chart.

```sh
wget https://github.com/taosdata/kubenetes/raw/main/helm/tdengine-0.1.0.tgz
```

First, check your sotrage class name:

```sh
helm get storageclass
```

And then deploy TDengine in one line:

```sh
helm update tdengine tdengine-0.1.0.tgz \
  --set storage.className=<your storage class name>
```

## Values

TDengine support `values.yaml` append.

To see a full list of values, use `helm show values`:

```sh
helm show values tdengine-0.1.0.tgz
```

You cound save it to `values.yaml`, and do some changs on it, like replica count, storage class name, and so on. Then type:

```sh
helm install tdengine tdengine-0.1.0.tgz -f values.yaml
```

The full list of values:

```yaml
{{#include ../helm/tdengine/values.yaml }}
```

## Uninstall

```sh
helm uninstall tdengine
```

## Scale Up and Down

See the details in chapter 4.
