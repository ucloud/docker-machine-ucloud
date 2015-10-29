# docker-machine-ucloud

docker-machine-ucloud是由[UCloud](https://www.ucloud.cn)提供的`docker-machine`插件。`docker-machine`可以利用该插件方便的创建UCloud云主机并部署docker和swarm服务。

## 安装

`docker-machine-ucloud`的安装主要有两种方式：
### 1. 直接下载binary

在https://github.com/ucloud/docker-machine-ucloud/releases/ 可以找到release版本的信息，当前为`docker-machine-ucloud-v0.5.0-rc3`, 目前仅仅提供linux系统的安装包。

安装包下载后需要copy到$PATH目录下，可以通过以下命令安装
```
$ curl -L https://github.com/ucloud/docker-machine-ucloud/releases/download/v0.5.0-rc3/docker-machine_linux-amd64_v0.5.0-ucloud.tar.gz > machine.tar.gz &&
$ tar -xzf machine.tar.gz && rm machine.tar.gz
$ mv docker-machine* /usr/local/bin/
```
成功下载后，可以在安装目录看到有 `docker-machine-driver-ulcoud`文件， 如：
```
$ ll /usr/local/bin 
-rwxr-xr-x 1 root root  7685920 Oct 26 17:02 docker-machine
-rwxr-xr-x 1 root root  7568896 Oct 26 17:02 docker-machine-driver-amazonec2
-rwxr-xr-x 1 root root  8141120 Oct 26 17:02 docker-machine-driver-azure
-rwxr-xr-x 1 root root  7369856 Oct 26 17:02 docker-machine-driver-digitalocean
-rwxr-xr-x 1 root root  7234688 Oct 26 17:02 docker-machine-driver-exoscale
-rwxr-xr-x 1 root root  7091328 Oct 26 17:02 docker-machine-driver-generic
-rwxr-xr-x 1 root root  8512672 Oct 26 17:02 docker-machine-driver-google
-rwxr-xr-x 1 root root  7169152 Oct 26 17:02 docker-machine-driver-hyperv
-rwxr-xr-x 1 root root  7083136 Oct 26 17:02 docker-machine-driver-none
-rwxr-xr-x 1 root root  7788032 Oct 26 17:02 docker-machine-driver-openstack
-rwxr-xr-x 1 root root  7820800 Oct 26 17:02 docker-machine-driver-rackspace
-rwxr-xr-x 1 root root  7185536 Oct 26 17:02 docker-machine-driver-softlayer
-rwxr-xr-x 1 root root 10683142 Oct 26 17:23 docker-machine-driver-ucloud
-rwxr-xr-x 1 root root  7242880 Oct 26 17:02 docker-machine-driver-virtualbox
-rwxr-xr-x 1 root root  7189632 Oct 26 17:02 docker-machine-driver-vmwarefusion
-rwxr-xr-x 1 root root  7704064 Oct 26 17:02 docker-machine-driver-vmwarevcloudair
-rwxr-xr-x 1 root root  7246976 Oct 26 17:02 docker-machine-driver-vmwarevsphere
```


### 2. 源码安装

#### 编译要求

   要编译UCloud的插件，依赖于
   * golang 1.4.3+
   * docker-machine v0.5.0rc2

#### 安装docker-machine
   docker-machine的安装主要有两种方法，直接下载编译好的二进制或者通过源码安装:
   
   1. 直接下载
   
    可以参照docker-machine的官方[安装文档](https://docs.docker.com/machine/install-machine/)。
    docker-machine 0.5版本目前还未正式release，可以参照docker/machine的[release文档](https://github.com/docker/machine/releases/tag/v0.5.0-rc3)安装。
    
   2. 源码安装

   * 本地编译
      
    ```
    $ go get github.com/docker/machine
    $ cd $GOPATH/src/github.com/docker/machine
    $ make build
    ```
   * 使用容器编译，需要运行docker实例，并设置 `export USE_CONTAINER=true`
      
   ```
    $ go get github.com/docker/machine
    $ cd $GOPATH/src/github.com/docker/machine
    $ export USE_CONTAINER=true
    $ make build
   ```
   更多通过源码编译的方法，请参考docker/machine的[COMTRIBUTING文档](https://github.com/docker/machine/blob/master/CONTRIBUTING.md#building)


#### 安装docker-machine-driver-ucloud插件
   编译前将$GOPATH/bin添加到$PATH环境变量中（docker-machine运行时候需要能够在PATH中找到docker-machine-driver-ucloud插件).
   `docker-machine-driver-ucloud`可以参照如下命令安装:
   ```
   $ export GOPATH=<path_to_gopath> && export PATH=$GOPATH/bin:$PATH
   $ go get github.com/ucloud/docker-machine-ucloud
   $ cd $GOPATH/src/github.com/ucloud/docker-machine-ucloud
   $ make
   $ make install
   ```

## 运行
  安装完成后，可以通过运行`docker-machine create --help -d ucloud`查看如何使用`docker-machine`创建UCloud云主机，正如帮助信息所示，创建云主机的时候
  需要提供UCloud API的公钥和私钥,可以通过控制台的[API密钥](https://consolev3.ucloud.cn/apikeyv3)创建。更多内容可以参照[UCloud API文档](https://docs.ucloud.cn/api-docs/index.html)中关于如何密钥的介绍。

  要创建UCloud云主机，需要将Public Key和Private Key通过环境变量，或者命令行参数传递给docker-machine，可以参照下面的例子创建运行docker的云主机：

   ```
    $ export UCLOUD_PUBLIC_KEY=<public-key-of-ucloud-api>
    $ export UCLOUD_PUBLIC_KEY=<private-key-of-ucloud-api>
    $ docker-machine create -d ucloud <machine-name>
   ```

  或者运行

   ```
    $ docker-machine create -d ucloud --ucloud-public-key=<public-key> --ucloud-private-key=<private-key> <machine-name>
   ```

  运行命令后，可以看到创建云主机的过程，如：

   ```
    $ docker-machine create -d ucloud ucloud-machine
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

   成功创建云主机后，可以通过`docker-machine ls`命令查看主机的运行情况

   ```
    $ docker-machine ls
    NAME             ACTIVE   DRIVER   STATE     URL                        SWARM
    ucloud-machine   -        ucloud   Running   tcp://123.59.66.163:2376
   ```
