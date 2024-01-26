

ROOT := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
BASH_CMD = /bin/bash -c
CUR_DIR = $(shell pwd)
NODE_IMAGE = node:16
NPM_BUILD_CMD = npm i
NPM_BUILD_PROD_CMD = npm run build:prod
NPM_PKG_STAGE_CMD = npm run start
USER_ID := $(shell id -u)
GROUP_ID := $(shell id -g)
ENV_VAR := -e API_KEY=127.0.0.1


# NPM_CMN_ARG_DR_LINUX =	--net=host 
NPM_CMN_ARG_DR_LINUX =	--name web \
	-v ${ROOT}:${ROOT} \
	-w ${ROOT}/${SERVICE_NAME} \

NPM_PXY_ARG_DR_LINUX = -e HOME=. \
	-e NPM_CONFIG_CACHE="${ROOT}/dashboard/.npm" \
	-u ${USER_ID}:${GROUP_ID}

NPM_ARG_DR_LINUX = ${ENV_VAR} ${NPM_CMN_ARG_DR_LINUX} ${NPM_PXY_ARG_DR_LINUX}

NPM_DR_LINUX = docker run --rm \
	${NPM_ARG_DR_LINUX} \
	${NODE_IMAGE} \
	${BASH_CMD}

NPM_PORT = -p 4200:4200 -p 49153:49153

NPM_RUN_LINUX = docker run ${NPM_PORT} --rm \
	${NPM_ARG_DR_LINUX} \
	${NODE_IMAGE} \
	${BASH_CMD}

UPDATE_ANGULAR = npm i && ng update @angular/cli @angular/core

NPM_UPDATE_ANGULAR = docker run --rm \
	${NPM_ARG_DR_LINUX} \
	${NODE_IMAGE} \
	${BASH_CMD}
