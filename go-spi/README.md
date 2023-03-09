# 安装依赖

```shell
go mod tidy
```

# 热启动

```shell
go install github.com/kataras/rizla
rizla main.go
```

# clickhouse使用
> 源码地址：https://github.com/ClickHouse/clickhouse-go
## 安装
```shell
# 使用脚本下载yum源
curl -s https://packagecloud.io/install/repositories/altinity/clickhouse/script.rpm.sh | sudo bash
 
# yum 安装 server 以及 client
sudo yum install -y clickhouse-server clickhouse-client
 
# 查看是否安装完成
sudo yum list installed 'clickhouse*'
 
# 启动server端clickhouse-server
service clickhouse-server start
 
# 查看server端服务开启/关闭状态 
service clickhouse-server status
 
# 进入client端clickhouse-client
clickhouse-client
```

## 卸载

```shell
# 卸载及删除安装文件（需root权限）
yum list installed | grep clickhouse
 
yum remove -y clickhouse-common-static
 
yum remove -y clickhouse-server-common
 
rm -rf /var/lib/clickhouse
 
rm -rf /etc/clickhouse-*
 
rm -rf /var/log/clickhouse-server
```

## 修改配置文件
`etc/cliclhouse-server/config.xml`
将下面的注释打开，允许远程连接
```xml
<listen_host>::</listen_host>
```


# `clickhouse-go`连接clickhouse
```shell
go get github.com/ClickHouse/clickhouse/v2
go get github.com/jmoiron/sqlx
```
<font color="red">PS:使用v2版本，golang版本为1.18大概率会报`内部编译失败的错误`，可以降/升级版本即可解决</font>

