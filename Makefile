GOMOD=hredirect
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
HASH=$(shell git log -n1 --pretty=format:%h)
REVS=$(shell git log --oneline|wc -l)

debug: setver compdbg
release: setver comprel
setver:
	cp verinfo.tpl version.go
	sed -i 's/{_G_HASH}/$(HASH)/' version.go
	sed -i 's/{_G_REVS}/$(REVS)/' version.go
comprel:
	mkdir -p bin && go build -ldflags="-s -w" .  && mv $(GOMOD) bin
	upx --best --lzma bin/$(GOMOD)
compdbg:
	mkdir -p bin && go build -race -gcflags=all=-d=checkptr=0 .  && mv $(GOMOD) bin
clean:
	rm -fr bin version.go
