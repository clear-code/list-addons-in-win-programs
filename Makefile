PACKAGE_NAME = list-addons-in-win-programs

.PHONY: xpi lint host

all: xpi

xpi: makexpi/makexpi.sh
	cp extlib/js-codemodule-registry/registry.jsm modules/
	cp makexpi/makexpi.sh ./
	makexpi/makexpi.sh -n $(PACKAGE_NAME) -o
	cd webextensions && make && cp ./*.xpi ../

host:
	cd webextensions && make host && cp ./*host.zip ../

makexpi/makexpi.sh:
	git submodule update --init

lint:
	cd webextensions && make lint
