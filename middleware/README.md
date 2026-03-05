# casdoor

## casdoor

### 配置

```shell
vim ./casdoor/conf/app.conf

# 修改dataSourceName为postgres相关的信息
# 修改redisEndpoint为redis相关的信息
```

### 创建数据库

```shell
# 需要在数据库实例创建casdoor名称数据库
```

### files

> `casdoor/files/`目录是`casdoor`提供商`Storage`本地存储类型的挂载目录
>
> 如果`Storage`的类型为`Local File System`，那么在页面上资源里面上传的文件，就是保存到这里的。

```shell
# 创建挂载目录
mkdir casdoor/files/
```

## postgresql

> 如果已经有`postgresql`实例可以不需要创建。

### 配置

- 复制环境变量文件，并且修改数据库密码。

```shell
cp example.env .env
vim .env
```

### postgres 配置文件

- 复制 postgres 配置文件，并且配置

```shell
cp ./postgresql/postgresql.conf.example ./postgresql/postgresql.conf
vim ./postgresql/postgresql.conf
```

## redis

> 如果已经有`redis`实例可以不需要创建。

### redis 需要配置内核参数

- 调整宿主机内核参数

```shell
# 永久生效
cp ./redis/redis-sysctl.conf /etc/sysctl.d/
# 临时生效
sysctl -p /etc/sysctl.d/redis-sysctl.conf
```

- 调整容器内核参数

`docker-compose.yaml`

```yaml
services:
  redis:
    ... ..
    sysctls:
      # 用于设置监听套接字（listening socket）的最大等待连接队列长度
      net.core.somaxconn: 4096
    ... ...
```

### 配置文件

```shell
cp ./redis/redis.conf.example./redis/redis.conf
vim ./redis/redis.conf
```

### 数据目录

`./redis/data/`