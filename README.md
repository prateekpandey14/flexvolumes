# flexvolume-digitalocean

This is a Kubernetes FlexVolume plugin for DigitalOcean volumes.

### Usage

Build with `go build .` and copy `digitalocean` and `flexvolume-digitalocean` to
`/usr/libexec/kubernetes/kubelet-plugins/volume/exec/digitalocean/digitalocean`.

Note that CoreOS mounts `/usr` as read-only so you'll want to add something like
`--volume-plugin-dir=/etc/kubernetes/volumeplugins` to `KUBELET_ARGS` in
`/etc/kubernetes/kubelet.env` and then `systemctl restart kubelet.service`.

