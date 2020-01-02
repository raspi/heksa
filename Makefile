APPNAME?="heksa"
# version from last tag
VERSION := $(shell git describe --abbrev=0 --always --tags)
BUILD := $(shell git rev-parse $(VERSION))
BUILDDATE := $(shell git log -1 --format=%aI $(VERSION))
BUILDFILES?=$$(find . -mindepth 1 -maxdepth 1 -type f \( -iname "*heksa-v*" -a ! -iname "*.shasums" \))
LDFLAGS := -ldflags "-s -w -X=main.VERSION=$(VERSION) -X=main.BUILD=$(BUILD) -X=main.BUILDDATE=$(BUILDDATE)"
SCREENSHOTCMD := ./heksa -f hex,asc,bit -l 0x200 heksa.exe
TMPDIR := $(shell mktemp -d -t heksa-rel-XXXXX)

default: build
# Helper for taking screenshot when releasing new version
screenshot:
	cd $(BUILDDIR); echo "% $(SCREENSHOTCMD)" > scr.txt && $(SCREENSHOTCMD) >> scr.txt && echo "% " >> scr.txt && konsole --notransparency --noclose --hide-tabbar -e cat scr.txt

build:
	@echo "GO BUILD..."
	@CGO_ENABLED=0 go build $(LDFLAGS) -v -o ./bin/${APPNAME} .

linux-build:
	@echo "linux build... amd64"
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -v -o ./bin/linux-amd64/${APPNAME} . 2>/dev/null
	@echo "linux build... arm"
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm go build $(LDFLAGS) -v -o ./bin/linux-arm/${APPNAME} . 2>/dev/null

darwin-build:
	@echo "darwin build... amd64"
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -v -o ./bin/darwin-amd64/${APPNAME} . 2>/dev/null

freebsd-build:
	@echo "freebsd build... amd64"
	@CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build $(LDFLAGS) -v -o ./bin/freebsd-amd64/${APPNAME} . 2>/dev/null

windows-build:
	@echo "windows build... amd64"
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -v -o ./bin/windows-amd64/${APPNAME}.exe . 2>/dev/null

upx-pack:
	@upx -v -9 ./bin/linux-amd64/${APPNAME}
	@upx -v -9 ./bin/linux-arm/${APPNAME}
	@upx -v -9 ./bin/windows-amd64/${APPNAME}.exe

release: linux-build darwin-build freebsd-build windows-build upx-pack copycommon tar-everything shasums
	@echo "release done..."

shasums:
	@pushd bin && shasum -a 256 $(BUILDFILES) > ${APPNAME}-${VERSION}.shasums

copycommon:
	@mkdir "$(TMPDIR)/bin"
	@cp LICENSE "$(TMPDIR)"
	@cp README.md "$(TMPDIR)"

tar-everything:
	@echo "tar-everything..."
	@cp "$(PWD)/bin/linux-amd64/${APPNAME}" "$(TMPDIR)/bin" && cd "$(TMPDIR)" && tar -zcvf $(PWD)/bin/${APPNAME}-${VERSION}-linux-amd64.tar.gz . && rm "$(TMPDIR)/bin/${APPNAME}"
	@cp "$(PWD)/bin/linux-arm/${APPNAME}" "$(TMPDIR)/bin" && cd "$(TMPDIR)" && tar -zcvf $(PWD)/bin/${APPNAME}-${VERSION}-linux-arm.tar.gz . && rm "$(TMPDIR)/bin/${APPNAME}"
	@cp "$(PWD)/bin/freebsd-amd64/${APPNAME}" "$(TMPDIR)/bin" && cd "$(TMPDIR)" && tar -zcvf $(PWD)/bin/${APPNAME}-${VERSION}-freebsd-amd64.tar.gz . && rm "$(TMPDIR)/bin/${APPNAME}"
	@cp "$(PWD)/bin/darwin-amd64/${APPNAME}" "$(TMPDIR)/bin" && cd "$(TMPDIR)" && tar -zcvf $(PWD)/bin/${APPNAME}-${VERSION}-darwin-amd64.tar.gz . && rm "$(TMPDIR)/bin/${APPNAME}"
	@cp "$(PWD)/bin/windows-amd64/${APPNAME}.exe" "$(TMPDIR)/bin" && cd "$(TMPDIR)" && zip -9 -y -r $(PWD)/bin/${APPNAME}-${VERSION}-windows-amd64.zip . && rm "$(TMPDIR)/bin/${APPNAME}.exe"

.PHONY: all clean test default