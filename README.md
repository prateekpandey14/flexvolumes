# flexvolume-digitalocean

This is a Kubernetes FlexVolume plugin for DigitalOcean volumes. Since
FlexVolumes are unstable, so too is this. Use at your own risk. Contributions
are welcome.

### Build

Build with `go build .`

### Install

Copy `digitalocean` and `flexvolume-digitalocean` to the Kubernetes volume
plugin directory, which is by default
`/usr/libexec/kubernetes/kubelet-plugins/volume/exec/digitalocean/digitalocean`

Note that CoreOS mounts `/usr` as read-only so instead you'll want to add
`--volume-plugin-dir=/etc/kubernetes/volumeplugins` to `KUBELET_ARGS` in
`/etc/kubernetes/kubelet.env`.

Restart kubelet with `systemctl restart kubelet.service`.

### Usage

See `example`. Fill in your DigitalOcean API key in `secret.yaml` and upload:

```
kubectl create -f secret.yaml
```

Next, create a volume on DigitalOcean if you haven't already done so, and find
its id. From the website it looks like the best way to do it is to inspect
element on the volumes page and look for a div with data-id="...", or use their
API, or if you're using terraform inspect the state. Fill in the id in pod.yaml.
Then:

```
kubectl create -f pod.yaml
```
