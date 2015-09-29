# docker-machine-ucloud

docker-machine-ucloud is a plugin of docker-machine for UCloud.

This plugin is a plugin for Docker Machine, which works with new plugin model for Docker Machine v0.5.0.
Hard developing is still in progress, so don't use it in product environment. Please feel free to send feedback and issues.


# install

The new plugin mechanism of docker-machine which in still in development, so please try the branch of 
[nathanleclaire/machine/libmachine_rpc_plugins](https://github.com/nathanleclaire/machine/tree/libmachine_rpc_plugins)

```
# @nathanleclaire developpnig libmachine-rpc branch
go get github.com/nathanleclaire/machine
cd $GOPATH/src/github.com/nathanleclaire/machine
git checkout nathanleclaire/libmachine_rpc_plugins
# Make libmachine rpc include docker-machine_darwin-amd64 binary
script/build
```
Then, you could install `docker-machine-ucloud` in the $GOPATH and add $GOPATH/bin to the $PATH env. 

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

# TODO

- [ ] `ls` command panic occasionally.

- [ ]  SecurityGroup of UNet can't works well.

- [ ]  Testing
    - [ ] testing for more situations.
    - [ ] integration testing

- [ ]  Swarm Support
