#!/usr/bin/env make

### Submodules
.PHONY: init
## fetch all submodules
init:
	git submodule init
	git submodule update
	git submodule foreach git submodule update --init --remote --recursive

.PHONY: master
## checkout master on all submodules
master: 
	git submodule foreach --recursive git checkout master

# .PHONY: pull
## pull submodules
pull: 
	git submodule foreach --recursive git pull

## update submodules
update: init master pull
update:
	echo "All submodules HEAD are now branch main latest"
