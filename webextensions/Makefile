NPM_MOD_DIR := $(CURDIR)/node_modules
NPM_BIN_DIR := $(NPM_MOD_DIR)/.bin

.PHONY: xpi install_dependency install_hook lint format host

all: xpi

install_dependency:
	[ -e "$(NPM_BIN_DIR)/eslint" -a -e "$(NPM_BIN_DIR)/jsonlint-cli" ] || npm install

install_hook:
	echo '#!/bin/sh\nmake lint' > "$(CURDIR)/../.git/hooks/pre-commit" && chmod +x "$(CURDIR)/../.git/hooks/pre-commit"

lint: install_dependency
	"$(NPM_BIN_DIR)/eslint" . --ext=.js --report-unused-disable-directives
	find . -type d -name node_modules -prune -o -type f -name '*.json' -print | xargs "$(NPM_BIN_DIR)/jsonlint-cli"

format: install_dependency
	"$(NPM_BIN_DIR)/eslint" . --ext=.js --report-unused-disable-directives --fix

xpi: lint
	rm -f ./*.xpi
	zip -r -9 list-addons-in-win-programs-we.xpi manifest.json common background dialog options resources extlib _locales -x '*/.*' >/dev/null 2>/dev/null

host:
	native-messaging-host/build.sh
	rm -f list-addons-in-win-programs.zip
	cd native-messaging-host && zip -r -9 ../list-addons-in-win-programs-host.zip 386 amd64 *.bat com.clear_code.*.json
