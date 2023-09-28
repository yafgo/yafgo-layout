#!/usr/bin/env bash

#########################################
# 本机 docker 环境
# 在终端执行该文件, 参数传要执行的函数名即可
# 如: ./start.sh usage
#########################################

# docker compose run [OPTIONS] SERVICE [COMMAND] [ARGS...]

cwd=$(pwd)
NC='\033[0m'
C_RED='\033[1;31m'
C_GREEN='\033[1;32m'
C_YELLOW='\033[1;33m'
C_BLUE='\033[1;34m'
C_PURPLE='\033[1;35m'
C_SKY='\033[36m'

# 使用说明，用来提示输入参数
usage() {
    echo -e "${C_PURPLE}本地 docker 开发环境${NC}"
    echo -e "${C_BLUE}Usage:${NC}\n  ${C_GREEN}./start.sh${NC} [command]\n"
    echo -e "${C_BLUE}Available Commands:${NC}"
    echo -e " ${C_GREEN} start ${NC}              启动 docker 容器"
    echo -e " ${C_GREEN} exec ${NC}               进入容器命令行"
    echo -e " ${C_GREEN} restart <服务名> ${NC}    重启指定的/所有容器"
    echo -e " ${C_GREEN} stop <服务名> ${NC}       停止指定的/所有容器"
    echo -e " ${C_GREEN} rm <服务名> ${NC}         移除指定的/所有容器"
    echo
}

base() {
    if ! [ -x "$(command -v docker-compose)" ]; then
        # shopt -s expand_aliases
        alias docker-compose='docker compose'
    fi
}

start() {
    docker-compose up -d --build yafgo
}

exec() {
    docker-compose exec -it yafgo /bin/sh
}

restart() {
    docker-compose restart "$@"
}

stop() {
    docker-compose stop "$@"
}

rm() {
    docker-compose down "$@"
}

prod() {
    docker-compose -f docker-compose.prod.yml up yafgo "$@"
}

###############################################################################

if [ $# -eq 0 ]; then
    usage
    echo -e "${C_YELLOW}请输入要执行的函数名${NC}"
    exit 1
fi
# 判断函数是否存在
func_name=$1
type $func_name >/dev/null 2>&1
if [ $? -ne 0 ]; then
    if [ $func_name != '-h' ]; then
        echo -e "${C_RED}参数 ${func_name} 不存在${NC}"
        echo
    fi
    usage
    exit 1
fi
# 执行前置操作
base
# 执行指定函数
shift
$func_name "$@"
