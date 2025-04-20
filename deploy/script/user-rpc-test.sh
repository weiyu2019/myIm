#! /bin/bash

repo_addr='crpi-g6muopbpy3n1uhbz.cn-hangzhou.personal.cr.aliyuncs.com/my-im-1/user-rpc-dev'
tag='latest'

container_name='myIm-user-rpc-test'
network_name='myIm-network'    #使用自定义网络myIm-network

docker stop ${container_name}

docker rm ${container_name}

docker rmi ${container_name}

docker pull ${repo_addr}:${tag}

docker run -p 10000:10000 --name=${container_name} --network=${network_name} -d ${repo_addr}:${tag}