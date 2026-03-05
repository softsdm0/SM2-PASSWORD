# passwordserver

## 项目命令

### 编译服务

```shell
make

# 测试运行
make start
```

> 编译成功后会生成`buildSave`目录和`passwordserver_xxxxxxx_x86_64.tgz`部署包

### 部署服务

> 部署编译后的包

```shell
mkdir -p /data/soft/
tar -xvf passwordserver_08afd0c_x86_64.tgz -C /data/soft/
```

### 启动服务

```shell
cd /data/soft/passwordserver/
./start.sh
```

### 停止服务

```shell
cd /data/soft/passwordserver/
./stop.sh
```

### 查看版本

```shell
cd /data/soft/passwordserver/
./passwordserver version
```

### 更新构建工具

- 由于gs-cli会更新，带来了`makefile` `start.sh` `stop.sh`脚本的更新，所以需要用到`gs-cli update`来对这些文件更新。

- `gs-cli update`只能在`goserver`创建的项目根目录中使用。依赖于`makefile`里面的`server_name`变量。

```shell
cd /data/soft/mygoserver/
gs-cli update
```