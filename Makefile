include commons.mk

SERVICE_NAME = dashboard
SERVICE_LANG = ts

.PHONY: all build run

ui-build:
	${NPM_DR_LINUX} "${NPM_BUILD_CMD}"

ui-run:
	${NPM_RUN_LINUX} "${NPM_PKG_STAGE_CMD}"

ui-update:
	${NPM_UPDATE_ANGULAR} "${UPDATE_ANGULAR}"

image:
	docker-compose build
	doker-compose push

# Target pattern to match any from parent
%: ;