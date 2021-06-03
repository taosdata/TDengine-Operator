# Setup K8s Cluster with Rancher

NOTE: I build this at May 26 2021 in Beijing(UTC+8), China. The problems I've met may not fit with you, keep your eyes open and clear.

## Install RancherD to deploy Rancher

For most of the cases, just run the rancherd installer.

```sh
curl -sfL https://get.rancher.io | sh -
```

If failed (like me), download the latest rancherd package from github releases assets.

```sh
# fill the proxy url if you use one
export https_proxy=
curl -s https://api.github.com/repos/rancher/rancher/releases/latest \
  |jq '.assets[] |
    select(.browser_download_url|contains("rancherd-amd64.tar.gz")) |
    .browser_download_url' -r \
  |wget -ci -
```

And install it.

```sh
tar xzf rancherd-amd64.tar.gz -C /usr/local
```

Then start the rancherd service.

```sh
systemctl enable rancherd-server
systemctl start rancherd-server
```

Keep tracking with the service.

```sh
journalctl -fu rancherd-server
```

End with log **successfully**:

```log
"Event occurred" object="cn120" kind="Node" apiVersion="v1" \ 
type="Normal" reason="Synced" message="Node synced successfully"
```

## Setup kubeconfig and kubectl

Once the Kubernetes cluster is up, set up RancherD’s kubeconfig file and kubectl:

```sh
export KUBECONFIG=/etc/rancher/rke2/rke2.yaml
export PATH=$PATH:/var/lib/rancher/rke2/bin
```

Check rancher status with kubectl:

```sh
kubectl get daemonset rancher -n cattle-system
kubectl get pod -n cattle-system
```

Result:

```
NAME      DESIRED   CURRENT   READY   UP-TO-DATE   AVAILABLE   NODE SELECTOR                         AGE
rancher   1         1         1       1            1           node-role.kubernetes.io/master=true   36m
NAME                               READY   STATUS      RESTARTS   AGE
helm-operation-5c2wd               0/2     Completed   0          34m
helm-operation-bdxlx               0/2     Completed   0          33m
helm-operation-cgcvr               0/2     Completed   0          34m
helm-operation-cj4g4               0/2     Completed   0          33m
helm-operation-hq282               0/2     Completed   0          34m
helm-operation-lp5nn               0/2     Completed   0          33m
rancher-kf592                      1/1     Running     0          36m
rancher-webhook-65f558c486-vrjz9   1/1     Running     0          33m
```

## Set Rancher Password

```sh
rancherd reset-admin
```

You would see like this:

```text
INFO[0000] Server URL: https://*.*.*.*:8443      
INFO[0000] Default admin and password created. Username: admin, Password: ****
```

## HA Settings

Check the token in `/var/lib/rancher/rke2/server/node-token`.

Install rancherd-server in other nodes like first node:

```sh
tar xzf rancherd-amd64.tar.gz -C /usr/local
systemctl enable rancherd-server
```

Prepare config dir:

```sh
mkdir -p /etc/rancher/rke2
```

Change the config file in `/etc/rancher/rke2/config.yaml`.

```yaml
server: https://192.168.60.120:9345
token: <the token in /var/lib/rancher/rke2/server/node-token>
```

Start rancherd

```sh
systemctl start rancherd-server
journalctl -fu rancherd-server
```

Other nodes just copy the config.yaml and start rancherd, and those will be joined to cluster automatically.

Type `kubectl get daemonset rancher -n cattle-system`：

```text
NAME      DESIRED   CURRENT   READY   UP-TO-DATE   AVAILABLE   NODE SELECTOR                         AGE
rancher   3         3         3       3            3           node-role.kubernetes.io/master=true   129m
```

Three nodes rancher+k8s cluster are avalibale now.
