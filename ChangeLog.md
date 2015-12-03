# ChangeLog

# 0.5.2 (2015-12-03)
1. 支持swarm集群的创建，swarm master的默认端口3376
2. 支持主机配置的定制，增加 —ucloud-cpu-cores, —ucloud-disk-space, —ucloud-memory-size。默认CPU: 1核，Memory: 1G, Disk: 20G 
3. 支持付费类型的定制，增加 —ucloud-charge-type 默认Month, 支持Year, Month, Dynamic, Trial
4. 调用的API请求设置User-Agent 为 docker-machine/v0.5.2

# 0.5.0-rc3 (2015-11-03)
1. 支持UHost的创建
