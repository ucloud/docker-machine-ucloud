# docker-machine-ucloud

docker-machine-ucloud是由[UCloud](https://www.ucloud.cn)提供的基于`docker-machine`的插件, 该插件提供在UCloud的云平台上提供可靠的docker服务。

# 安装

## 最小要求
   要运行UCloud的插件，需要依赖于
   * golang 1.4.3+
   * docker-machine v0.5.0rc2

## 安装docker-machine
   docker-machine的安装可以参照docker/machine的[release文档](https://github.com/docker/machine/releases/tag/v0.5.0-rc2)，
   或者直接通过源码安装的方式:
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


## 安装docker-machine-ucloud插件
   完成`docker-machine`安装后，可以通过以下步骤安装通过源码的方式docker-machine-ucloud的插件，安装之前，需要准备Golang的开发环境，并且
   设置好GOPATH环境变量，并将$GOPATH/bin添加到$PATH环境变量中（docker-machine运行时候需要能够在PATH中找到docker-machine-driver-ucloud插件）
   ```
   $ export GOPATH=<path_to_gopath> && export PATH=$GOPATH/bin:$PATH
   $ go get github.com/ucloud/docker-machine-ucloud
   $ cd $GOPATH/src/github.com/ucloud/docker-machine-ucloud
   $ make
   $ make install
   ```

# 运行
  安装完成后，可以通过运行`docker-machine create --help -d ucloud`查看如何使用`docker-machine`创建UCloud云主机，正如帮助信息所示，创建云主机的时候
  需要提供UCloud API的公钥和私钥。更多内容可以参照[UCloud API文档](https://docs.ucloud.cn/api/index.html)中关于如何密钥的介绍。

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
    $ ./docker-machine create -d ucloud ucloud-machine
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
    $ ./docker-machine_darwin-amd64 ls
    NAME             ACTIVE   DRIVER   STATE     URL                        SWARM
    ucloud-machine   -        ucloud   Running   tcp://123.59.66.163:2376
    ```
