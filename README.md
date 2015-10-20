# docker-machine-ucloud

docker-machine-ucloud is a plugin of docker-machine for UCloud. It requires Docker Machine's version is greater than v0.5.0-rc1.

（[中文文档](README-zh.md)） 

# Install and Run

You could install run docker-machine-driver-ucloud as following steps.

First, docker-machine v0.5.0 rc2 is required, documentation for how to install `docker-machine`
[is available here](https://github.com/docker/machine/releases/tag/v0.5.0-rc2#Installation).

or you can install `docker-machine` from source code by running these commands
```
$ go get github.com/docker/machine
$ cd $GOPATH/src/github.com/docker/machine
$ make build
```

Then, you could install `docker-machine-ucloud` driver in the $GOPATH and add $GOPATH/bin to the $PATH env. 

```
go get github.com/ucloud/docker-machine-ucloud
cd $GOPATH/src/github.com/ucloud/docker-machine-ucloud
make
make install
```

Now, you can run `docker-machine create --help -d ucloud` to see how to create a UCloud machine. Public and private keys of UCloud API
are needed to create machine. Both options and environment variable are available to set that:

```
$ export UCLOUD_PUBLIC_KEY=<public-key-of-ucloud-api>
$ export UCLOUD_PUBLIC_KEY=<private-key-of-ucloud-api>
$ docker-machine create -d ucloud <machine-name>
```
or  run 

```
$ docker-machine create -d ucloud --ucloud-public-key=<public-key> --ucloud-private-key=<private-key> <machine-name>
```

for example,

```
$ ./docker-machine_darwin-amd64 create -d ucloud ucloud-machine
Running pre-create checks...
Creating machine...
Waiting for machine to be running, this may take a few minutes...
Machine is running, waiting for SSH to be available...
Detecting operating system of created instance...
Provisioning created instance...
Copying certs to the local machine directory...
Copying certs to the remote machine...
Setting Docker configuration on the remote daemon...
To see how to connect Docker to this machine, run: docker-machine_darwin-amd64 env ucloud-machine
```

After run this command, a machine with name `ucloud-machine` is created.

```
$ ./docker-machine_darwin-amd64 ls
NAME             ACTIVE   DRIVER   STATE     URL                        SWARM
ucloud-machine   -        ucloud   Running   tcp://123.59.66.163:2376
```

