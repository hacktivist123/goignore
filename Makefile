VERSION=1.0.4
OSES=linux darwin windows
ARCHS=arm64 amd64
SRCS=$(shell find . -iname '*.go')
ALLSRCS=$(SRCS) version/version.go

default: $(ALLSRCS)
	go build ./cmd/goignore

_build:
	mkdir $@

version/version.go: version/version.go.m4 Makefile
	m4 -D__VERSION__=$(VERSION) < $< > $@

define make-target
ifneq ($1, windows)
_build/$(VERSION)/$1_$2/goignore: $(ALLSRCS) | _build
	GOOS=$1 GOARCH=$2 go build -o _build/$(VERSION)/$1_$2/goignore ./cmd/goignore

all:: _build/$(VERSION)/$1_$2/goignore

_build/$(VERSION)/goignore_$(VERSION)_$1_$2.tar.gz: _build/$(VERSION)/$1_$2/goignore | _build
	tar -zcf _build/$(VERSION)/goignore_$(VERSION)_$1_$2.tar.gz -C _build/$(VERSION)/$1_$2 goignore

archive:: _build/$(VERSION)/goignore_$(VERSION)_$1_$2.tar.gz
else
_build/$(VERSION)/$1_$2/goignore.exe: $(ALLSRCS)  | _build
	GOOS=$1 GOARCH=$2 go build -o _build/$(VERSION)/$1_$2/goignore.exe ./cmd/goignore

all:: _build/$(VERSION)/$1_$2/goignore.exe

_build/$(VERSION)/goignore_$(VERSION)_$1_$2.zip: _build/$(VERSION)/$1_$2/goignore.exe | _build
	zip _build/$(VERSION)/goignore_$(VERSION)_$1_$2.zip -j _build/$(VERSION)/$1_$2/goignore.exe

archive:: _build/$(VERSION)/goignore_$(VERSION)_$1_$2.zip
endif
endef

$(foreach os,$(OSES), $(foreach arch, $(ARCHS), $(eval $(call make-target,$(os),$(arch)))))

_build/$(VERSION)/darwin_universal/goignore: _build/$(VERSION)/darwin_amd64/goignore _build/$(VERSION)/darwin_arm64/goignore | _build
	-mkdir _build/$(VERSION)/darwin_universal
	lipo -create -output $@ $+

all:: _build/$(VERSION)/darwin_universal/goignore

_build/$(VERSION)/goignore_$(VERSION)_darwin_universal.tar.gz: _build/$(VERSION)/darwin_universal/goignore | _build
	tar -zcf _build/$(VERSION)/goignore_$(VERSION)_darwin_universal.tar.gz -C _build/$(VERSION)/darwin_universal goignore

archive:: _build/$(VERSION)/goignore_$(VERSION)_darwin_universal.tar.gz

clean:
	-rm -rf _build version/version.go

.PHONY: all clean default archive