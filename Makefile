PACKAGE_NAME = list-addons-in-win-programs

all: xpi

xpi: makexpi/makexpi.sh
	cp extlib/js-codemodule-registry/registry.jsm modules/
	cp makexpi/makexpi.sh ./
	./makexpi.sh -n $(PACKAGE_NAME) -o
	rm ./makexpi.sh

makexpi/makexpi.sh:
	git submodule update --init
