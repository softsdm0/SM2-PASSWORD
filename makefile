# 配置
server_name=passwordserver
web_name=passwordweb

server_port=30005
web_port=30006

server_config_path=./server/deploy/config/config.docker.yaml
web_config_path=./web/public/config.json
web_nginx_config_path=./web/nginx.conf

buildSave_dir=./buildSave

help:
	# 构建镜像:
	#   make all
	#   make build-all-image
	#   make build-server-image
	#   make build-web-image
	
	# 保存镜像:
	#   make build-all-image-save
	#   make build-server-image-save
	#   make build-web-image-save
	
	# 运行docker容器:
	#   make docker-run
	#   make docker-run-server
	#   make docker-run-web
	#   make docker-stop
	#   make docker-stop-server
	#   make docker-sotp-web
	
	# 清理docker:
	#   make docker-clean
	
	# 构建docker compose必要的文件:
	#   make docker-compose-build

# 构建镜像
all build-all-image: build-server-image build-web-image
build-server-image:
	docker build -t ${server_name} -f dockerfile.server .
build-web-image:
	docker build -t ${web_name} -f dockerfile.web .
# 构建镜像并保持
build-all-image-save: build-server-image-save build-web-image-save
build-server-image-save: build-server-image
	mkdir -p ${buildSave_dir}
	docker save -o ${buildSave_dir}/${server_name}.tar ${server_name}
build-web-image-save: build-web-image
	mkdir -p ${buildSave_dir}
	docker save -o ${buildSave_dir}/${web_name}.tar ${web_name}

# 运行所有docker容器
docker-run: docker-run-server docker-run-web
docker-stop: docker-stop-server docker-sotp-web
docker-run-server: docker-stop-server
	docker run -itd \
	  --restart=always \
	  -p ${server_port}:80 \
	  -v ${server_config_path}:/data/passwordserver/deploy/config/config.yaml \
	  --name ${server_name} \
	  ${server_name}
docker-stop-server:
	docker rm -f ${server_name}
docker-run-web: docker-sotp-web
	docker run -itd \
	  --restart=always \
	  -p ${web_port}:80 \
	  -v ${web_config_path}:/data/web/config.json \
	  -v ${web_nginx_config_path}:/etc/nginx/nginx.conf \
	  --name ${web_name} \
	  ${web_name}
docker-sotp-web:
	docker rm -f ${web_name}

# 清理docker
docker-clean:
	docker system prune -a -f

# 构建docker compose必要的文件
compose_dir_name=${buildSave_dir}/${server_name}-docker-compose
docker-compose-build:
	rm -rf ${compose_dir_name}
	mkdir -p ${compose_dir_name}
	mkdir -p ${compose_dir_name}/server
	mkdir -p ${compose_dir_name}/web
	cp -r ./middleware/casdoor ${compose_dir_name}
	cp -r ./middleware/postgresql ${compose_dir_name}
	cp -r ./middleware/redis ${compose_dir_name}
	cp docker-compose.yaml ${compose_dir_name}/docker-compose.yaml
	cp ./server/deploy/config/config.yaml.example ${compose_dir_name}/server/
	cp ./web/public/config.json.example ${compose_dir_name}/web/
	cd ${buildSave_dir} && tar -zcvf ${server_name}-docker-compose.tar.gz ${server_name}-docker-compose

	
