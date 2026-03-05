#=====所有需要修改的变量=====
# 服务名称，也就是包名，$1可以传值，或者在这里写死
server_name := passwordserver
#=====所有需要修改的变量=====


#=====保存路径=====
buildPath := buildSave
buildSavePath := $(buildPath)/$(server_name)
#=====保存路径=====


#=====版本信息=====
goserverPackageName := gitee.com/ouhaoqiang/goserver
goVersion := $(shell go version)
gitHash := $(shell git show -s --format=%H)
buildTime := $(shell date)
gitHashShort := $(shell git rev-parse --short HEAD)
platform := $(shell uname -i)
#=====版本信息=====


#=====构建信息=====
goVersionPackage := $(goserverPackageName)/src/version
buildParameter := -o $(buildSavePath)/$(server_name)
buildParameter += -ldflags "
buildParameter += -X '$(goVersionPackage).goVersion=$(goVersion)'
buildParameter += -X '$(goVersionPackage).platform=$(platform)'
buildParameter += -X '$(goVersionPackage).gitHash=$(gitHash)'
buildParameter += -X '$(goVersionPackage).buildTime=$(buildTime)'
buildParameter += "
ifdef $(gitHashShort)
pt:=$(date  +"%Y%m%d")
packName:=$(server_name)_$(pt)_$(platform).tgz
else
packName:=$(server_name)_$(gitHashShort)_$(platform).tgz
endif
#=====构建信息=====


all: tar
	echo $(goVersion) $(gitHash) $(buildTime) $(gitHashShort) $(platform)

# 打包
tar: build cpMaterial
	cd $(buildPath) && tar -zcvf $(packName) $(server_name)

# 复制相关依赖物料
cpMaterial: mkdirBuildSavePath
	if [ -d "deploy" ]; then cp -r deploy ${buildSavePath}/; fi
	if [ -f "start.sh" ]; then cp start.sh ${buildSavePath}/; fi
	if [ -f "stop.sh" ]; then cp stop.sh ${buildSavePath}/; fi

# 构建
build: clean mkdirBuildSavePath
	go build $(buildParameter) .
	./$(buildSavePath)/$(server_name) version

# 创建构建物料保存路径
mkdirBuildSavePath: checkVar
	mkdir -p $(buildSavePath)

# 清理构建物料
clean: checkVar
	rm -rf $(buildPath)

# 检测必要变量
checkVar:
ifndef server_name
	echo "server_name不能为空"
	exit 1
endif
ifndef buildPath
	echo "buildPath不能为空"
	exit 1
endif
ifndef buildSavePath
	echo "buildSavePath不能为空"
	exit 1
endif

# 启动编译后的服务
run: start

# 启动编译后的服务
start:
	./$(buildSavePath)/start.sh