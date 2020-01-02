APPNAME?=heksa
# version from last tag
VERSION := $(shell git describe --abbrev=0 --always --tags)
BUILD := $(shell git rev-parse $(VERSION))
BUILDDATE := $(shell git log -1 --format=%aI $(VERSION))
BUILDFILES?=$$(find . -mindepth 1 -maxdepth 1 -type f \( -iname "*${APPNAME}-v*" -a ! -iname "*.shasums" \))
LDFLAGS := -ldflags "-s -w -X=main.VERSION=$(VERSION) -X=main.BUILD=$(BUILD) -X=main.BUILDDATE=$(BUILDDATE)"
SCREENSHOTCMD := ./heksa -f hex,asc,bit -l 0x200 heksa.exe
TMPDIR := $(shell mktemp -d -t heksa-rel-XXXXX)

LINUX_ARCHS := amd64 arm arm64 ppc64 ppc64le
WINDOWS_ARCHS := amd64
DARWIN_ARCHS := amd64
FREEBSD_ARCHS := amd64 arm
NETBSD_ARCHS := amd64 arm
OPENBSD_ARCHS := amd64 arm arm64

default: build
# Helper for taking screenshot when releasing new version
screenshot:
	cd $(BUILDDIR); echo "% $(SCREENSHOTCMD)" > scr.txt && $(SCREENSHOTCMD) >> scr.txt && echo "% " >> scr.txt && konsole --notransparency --noclose --hide-tabbar -e cat scr.txt

build:
	@echo "GO BUILD..."
	@CGO_ENABLED=0 go build $(LDFLAGS) -v -o ./bin/${APPNAME} .

linux-build:
	@for arch in $(LINUX_ARCHS); do \
	  echo "GNU/Linux build... $$arch"; \
	  CGO_ENABLED=0 GOOS=linux GOARCH=$$arch go build $(LDFLAGS) -v -o ./bin/linux-$$arch/${APPNAME} . 2>/dev/null; \
	done

darwin-build:
	@for arch in $(DARWIN_ARCHS); do \
	  echo "Darwin build... $$arch"; \
	  CGO_ENABLED=0 GOOS=darwin GOARCH=$$arch go build $(LDFLAGS) -v -o ./bin/darwin-$$arch/${APPNAME} . ; \
	done

freebsd-build:
	@for arch in $(FREEBSD_ARCHS); do \
	  echo "FreeBSD build... $$arch"; \
	  CGO_ENABLED=0 GOOS=freebsd GOARCH=$$arch go build $(LDFLAGS) -v -o ./bin/freebsd-$$arch/${APPNAME} . 2>/dev/null; \
	done

netbsd-build:
	@for arch in $(NETBSD_ARCHS); do \
	  echo "NetBSD build... $$arch"; \
	  CGO_ENABLED=0 GOOS=netbsd GOARCH=$$arch go build $(LDFLAGS) -v -o ./bin/netbsd-$$arch/${APPNAME} . 2>/dev/null; \
	done

openbsd-build:
	@for arch in $(OPENBSD_ARCHS); do \
	  echo "OpenBSD build... $$arch"; \
	  CGO_ENABLED=0 GOOS=openbsd GOARCH=$$arch go build $(LDFLAGS) -v -o ./bin/openbsd-$$arch/${APPNAME} . 2>/dev/null; \
	done

windows-build:
	@for arch in $(WINDOWS_ARCHS); do \
	  echo "MS Windows build... $$arch"; \
	  CGO_ENABLED=0 GOOS=windows GOARCH=$$arch go build $(LDFLAGS) -v -o ./bin/windows-$$arch/${APPNAME}.exe . 2>/dev/null; \
	done

upx-pack:
	@upx -v -9 ./bin/linux-amd64/${APPNAME}
	@upx -v -9 ./bin/linux-arm/${APPNAME}
	@upx -v -9 ./bin/windows-amd64/${APPNAME}.exe

release: linux-build darwin-build freebsd-build openbsd-build netbsd-build windows-build upx-pack tar-everything shasums
	@echo "release done..."

shasums:
	@pushd bin && shasum -a 256 $(BUILDFILES) > ${APPNAME}-${VERSION}.shasums

copycommon:
	@echo "Copying common files to $(TMPDIR)"
	@mkdir "$(TMPDIR)/bin"
	@cp LICENSE "$(TMPDIR)"
	@cp README.md "$(TMPDIR)"

tar-everything: copycommon
	@echo "tar-everything..."
	for arch in $(LINUX_ARCHS); do \
	  echo "GNU/Linux tar... $$arch"; \
	  cp -v "$(PWD)/bin/linux-$$arch/${APPNAME}" "$(TMPDIR)/bin"; \
	  cd "$(TMPDIR)"; \
	  tar -zcvf "$(PWD)/bin/${APPNAME}-${VERSION}-linux-$$arch.tar.gz" . ; \
	  rm "$(TMPDIR)/bin/${APPNAME}"; \
	done

	for arch in $(DARWIN_ARCHS); do \
	  echo "Darwin tar... $$arch"; \
	  cp -v "$(PWD)/bin/darwin-$$arch/${APPNAME}" "$(TMPDIR)/bin"; \
	  cd "$(TMPDIR)"; \
	  tar -zcvf "$(PWD)/bin/${APPNAME}-${VERSION}-darwin-$$arch.tar.gz" . ; \
	  rm "$(TMPDIR)/bin/${APPNAME}"; \
	done

	for arch in $(FREEBSD_ARCHS); do \
	  echo "FreeBSD tar... $$arch"; \
	  cp -v "$(PWD)/bin/freebsd-$$arch/${APPNAME}" "$(TMPDIR)/bin"; \
	  cd "$(TMPDIR)"; \
	  tar -zcvf "$(PWD)/bin/${APPNAME}-${VERSION}-freebsd-$$arch.tar.gz" . ; \
	  rm "$(TMPDIR)/bin/${APPNAME}"; \
	done

	for arch in $(OPENBSD_ARCHS); do \
	  echo "OpenBSD tar... $$arch"; \
	  cp -v "$(PWD)/bin/openbsd-$$arch/${APPNAME}" "$(TMPDIR)/bin"; \
	  cd "$(TMPDIR)"; \
	  tar -zcvf "$(PWD)/bin/${APPNAME}-${VERSION}-openbsd-$$arch.tar.gz" . ; \
	  rm "$(TMPDIR)/bin/${APPNAME}"; \
	done

	for arch in $(NETBSD_ARCHS); do \
	  echo "NetBSD tar... $$arch"; \
	  cp -v "$(PWD)/bin/netbsd-$$arch/${APPNAME}" "$(TMPDIR)/bin"; \
	  cd "$(TMPDIR)"; \
	  tar -zcvf "$(PWD)/bin/${APPNAME}-${VERSION}-netbsd-$$arch.tar.gz" . ; \
	  rm "$(TMPDIR)/bin/${APPNAME}"; \
	done

	for arch in $(WINDOWS_ARCHS); do \
	  echo "MS Windows zip... $$arch"; \
	  cp -v "$(PWD)/bin/windows-$$arch/${APPNAME}.exe" "$(TMPDIR)/bin"; \
	  cd "$(TMPDIR)"; \
	  zip -9 -y -r $(PWD)/bin/${APPNAME}-${VERSION}-windows-$$arch.zip . ; \
	  rm "$(TMPDIR)/bin/${APPNAME}.exe"; \
	done

.PHONY: all clean test default