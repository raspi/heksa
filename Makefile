APPNAME?=heksa

# ????
# makepkg --printsrcinfo > .SRCINFO; git add PKGBUILD .SRCINFO; git commit; git push
#

# version from last tag
VERSION := $(shell git describe --abbrev=0 --always --tags)
BUILD := $(shell git rev-parse $(VERSION))
BUILDDATE := $(shell git log -1 --format=%aI $(VERSION))
BUILDFILES?=$$(find . -mindepth 1 -maxdepth 1 -type f \( -iname "*${APPNAME}-v*" -a ! -iname "*.shasums" \))
LDFLAGS := -ldflags "-s -w -X=main.VERSION=$(VERSION) -X=main.BUILD=$(BUILD) -X=main.BUILDDATE=$(BUILDDATE)"
RELEASETMPDIR := $(shell mktemp -d -t ${APPNAME}-rel-XXXXXX)
APPANDVER := ${APPNAME}-$(VERSION)
RELEASETMPAPPDIR := $(RELEASETMPDIR)/$(APPANDVER)

UPXFLAGS := -v -9
XZCOMPRESSFLAGS := --verbose --keep --compress --threads 0 --extreme -9

# https://golang.org/doc/install/source#environment
LINUX_ARCHS := amd64 arm arm64 ppc64 ppc64le
WINDOWS_ARCHS := amd64
DARWIN_ARCHS := amd64
FREEBSD_ARCHS := amd64 arm
NETBSD_ARCHS := amd64 arm
OPENBSD_ARCHS := amd64 arm arm64

default: build

build:
	@echo "GO BUILD..."
	@CGO_ENABLED=0 go build $(LDFLAGS) -v -o ./bin/${APPNAME} .

# Update go module(s)
modup:
	@go get -u github.com/raspi/go-PKGBUILD@v0.0.5
	@go mod vendor
	@go mod tidy

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

# Compress executables
upx-pack:
	@upx $(UPXFLAGS) ./bin/linux-amd64/${APPNAME}
	@upx $(UPXFLAGS) ./bin/linux-arm/${APPNAME}
	@upx $(UPXFLAGS) ./bin/windows-amd64/${APPNAME}.exe

release: linux-build darwin-build freebsd-build openbsd-build netbsd-build windows-build upx-pack compress-everything shasums release-ldistros
	@echo "release done..."

# Linux distributions
release-ldistros: ldistro-arch
	@echo "Linux distros release done..."

shasums:
	@echo "Checksumming..."
	@pushd "release/${VERSION}" && shasum -a 256 $(BUILDFILES) > $(APPANDVER).shasums

# Copy common files to release directory
copycommon:
	@echo "Copying common files to temporary release directory '$(RELEASETMPAPPDIR)'.."
	@mkdir -p "$(RELEASETMPAPPDIR)/bin"
	@cp -v "./LICENSE" "$(RELEASETMPAPPDIR)"
	@cp -v "./README.md" "$(RELEASETMPAPPDIR)"
	@mkdir --parents "$(PWD)/release/${VERSION}"

# Compress files: FreeBSD
compress-freebsd:
	@for arch in $(FREEBSD_ARCHS); do \
	  echo "FreeBSD xz... $$arch"; \
	  cp -v "$(PWD)/bin/freebsd-$$arch/${APPNAME}" "$(RELEASETMPAPPDIR)/bin"; \
	  cd "$(RELEASETMPDIR)"; \
	  tar --numeric-owner --owner=0 --group=0 -cf - . | xz $(XZCOMPRESSFLAGS) - > "$(PWD)/release/${VERSION}/$(APPANDVER)-freebsd-$$arch.tar.xz" ; \
	  rm "$(RELEASETMPAPPDIR)/bin/${APPNAME}"; \
	done

# Compress files: OpenBSD
compress-openbsd:
	@for arch in $(OPENBSD_ARCHS); do \
	  echo "OpenBSD xz... $$arch"; \
	  cp -v "$(PWD)/bin/openbsd-$$arch/${APPNAME}" "$(RELEASETMPAPPDIR)/bin"; \
	  cd "$(RELEASETMPDIR)"; \
	  tar --numeric-owner --owner=0 --group=0 -cf - . | xz $(XZCOMPRESSFLAGS) - > "$(PWD)/release/${VERSION}/$(APPANDVER)-openbsd-$$arch.tar.xz" ; \
	  rm "$(RELEASETMPAPPDIR)/bin/${APPNAME}"; \
	done

# Compress files: NetBSD
compress-netbsd:
	@for arch in $(NETBSD_ARCHS); do \
	  echo "NetBSD xz... $$arch"; \
	  cp -v "$(PWD)/bin/netbsd-$$arch/${APPNAME}" "$(RELEASETMPAPPDIR)/bin"; \
	  cd "$(RELEASETMPDIR)"; \
	  tar --numeric-owner --owner=0 --group=0 -cf - . | xz $(XZCOMPRESSFLAGS) - > "$(PWD)/release/${VERSION}/$(APPANDVER)-netbsd-$$arch.tar.xz" ; \
	  rm "$(RELEASETMPAPPDIR)/bin/${APPNAME}"; \
	done

# Compress files: GNU/Linux
compress-linux:
	@for arch in $(LINUX_ARCHS); do \
	  echo "GNU/Linux tar... $$arch"; \
	  cp -v "$(PWD)/bin/linux-$$arch/${APPNAME}" "$(RELEASETMPAPPDIR)/bin"; \
	  cd "$(RELEASETMPDIR)"; \
	  tar --numeric-owner --owner=0 --group=0 -zcvf "$(PWD)/release/${VERSION}/$(APPANDVER)-linux-$$arch.tar.gz" . ; \
	  rm "$(RELEASETMPAPPDIR)/bin/${APPNAME}"; \
	done

# Compress files: Darwin
compress-darwin:
	@for arch in $(DARWIN_ARCHS); do \
	  echo "Darwin tar... $$arch"; \
	  cp -v "$(PWD)/bin/darwin-$$arch/${APPNAME}" "$(RELEASETMPAPPDIR)/bin"; \
	  cd "$(RELEASETMPDIR)"; \
	  tar --owner=0 --group=0 -zcvf "$(PWD)/release/${VERSION}/$(APPANDVER)-darwin-$$arch.tar.gz" . ; \
	  rm "$(RELEASETMPAPPDIR)/bin/${APPNAME}"; \
	done

# Compress files: Microsoft Windows
compress-windows:
	@for arch in $(WINDOWS_ARCHS); do \
	  echo "MS Windows zip... $$arch"; \
	  cp -v "$(PWD)/bin/windows-$$arch/${APPNAME}.exe" "$(RELEASETMPAPPDIR)/bin"; \
	  cd "$(RELEASETMPAPPDIR)"; \
	  mv "LICENSE" "LICENSE.txt" && \
	  pandoc --standalone --to rtf --output LICENSE.rtf LICENSE.txt && \
	  rm "LICENSE.txt" ; \
	  cd "$(RELEASETMPDIR)" ; \
	  zip -v -9 -r -o -9 "$(PWD)/release/${VERSION}/$(APPANDVER)-windows-$$arch.zip" . ; \
	  rm "$(RELEASETMPAPPDIR)/LICENSE.rtf"; \
	  cp -v "$(PWD)/LICENSE" "$(RELEASETMPAPPDIR)" ; \
	  rm "$(RELEASETMPAPPDIR)/bin/${APPNAME}.exe"; \
	done

# Move all to temporary directory and compress with common files
compress-everything: copycommon compress-linux compress-windows compress-freebsd compress-netbsd compress-openbsd
	@echo "$@ ..."
	rm -rf "$(RELEASETMPDIR)/*"

# Distro: Arch linux - https://www.archlinux.org/
# Generates multi-arch PKGBUILD
ldistro-arch:
	pushd release/linux/arch && go run . -version ${VERSION} > "$(PWD)/release/${VERSION}/$(APPANDVER)-linux-Arch.PKGBUILD"

.PHONY: all clean test default