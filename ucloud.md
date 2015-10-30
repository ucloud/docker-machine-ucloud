<!--[metadata]>
+++
title = "UCloud"
description = "UCloud driver for machine"
keywords = ["machine, ucloud, driver"]
[menu.main]
parent="smn_machine_drivers"
+++
<![end-metadata]-->

# UCloud
Create machines on [UCloud](http://ucloud.cn). To create machines on [UCloud](http://ucloud.cn), you must supply two required parameters:

 - API Public Key
 - API Private Key
 

Obtain your Keys from UCloud:

  1. Login to the UCloud console
  2. Go to **Management -> API Key -> Display**.
  3. Get the public key and private key

Then, you could pass the keys to `docker-machine create` options with `--ucloud-public-key` and `--ucloud-private-key` to create an
uhost machine at UCloud.

```
$ docker-machine create --driver ucloud --ucloud-public-key <public-key> --ucloud-private-key <private key>  uhost-01
```


### Options
 -  `--ucloud-imageid 							UHost image id`
 -  `--ucloud-private-address-only				Only use a private IP address`
 -  `--ucloud-private-key 						UCloud Private Key [$UCLOUD_PRIVATE_KEY]`
 -  `--ucloud-public-key 						UCloud Public Key [$UCLOUD_PUBLIC_KEY]`
 -  `--ucloud-region 				            Region of ucloud idc [$UCLOUD_REGION]`
 -  `--ucloud-security-group                    UCloud security group`
 -  `--ucloud-ssh-port  						SSH port`
 -  `--ucloud-ssh-user      					SSH user`
 -  `--ucloud-user-password 					Password of ucloud user`

By default, the UCloud machine driver will use image of CentOS 7.0.


Environment variables and default values:

| CLI option                          | Environment variable    | Default          |
|-------------------------------------|-------------------------|------------------|
| `--ucloud-imageid`                  | -                       | -                |
| `--ucloud-private-address-only`     | -                       |`false`           |
| **`--ucloud-private-key`**          | `UCLOUD_PRIVATE_KEY`    | -                |
| **`--ucloud-public-key`**           | `UCLOUD_PUBLIC_KEY`     | -                |
| `--ucloud-region`                   | `UCLOUD_REGION`         |`cn-north-03`     |
| `--ucloud-security-group`           | -                       |`docker-machine`  |
| `--ucloud-ssh-port`                 | -                       | `22`             |
| `--ucloud-ssh-user`                 | -                       | `root`           |
| `--ucloud-user-password`            | -                       | -                |
