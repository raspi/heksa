APPNAME?=heksa

# version from last tag
VERSION := $(shell git describe --abbrev=0 --always --tags)
BUILD := $(shell git rev-parse $(VERSION))
BUILDDATE := $(shell git log -1 --format=%aI $(VERSION))
BUILDFILES?=$$(find . -mindepth 1 -maxdepth 1 -type f \( -iname "*${APPNAME}-v*" -a ! -iname "*.shasums" \))
LDFLAGS := -ldflags "-s -w -X=main.VERSION=$(VERSION) -X=main.BUILD=$(BUILD) -X=main.BUILDDATE=$(BUILDDATE)"
TMPDIR := $(shell mktemp -d -t ${APPNAME}-rel-XXXXXX)

UPXFLAGS := -v -9

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

upx-pack:
	@upx $(UPXFLAGS) ./bin/linux-amd64/${APPNAME}
	@upx $(UPXFLAGS) ./bin/linux-arm/${APPNAME}
	@upx $(UPXFLAGS) ./bin/windows-amd64/${APPNAME}.exe

release: linux-build darwin-build freebsd-build openbsd-build netbsd-build windows-build upx-pack tar-everything shasums release-ldistros
	@echo "release done..."

# Linux distributions
release-ldistros: ldistro-arch
	@echo "Linux distros release done..."

shasums:
	@echo "Checksumming..."
	@pushd "release/${VERSION}" && shasum -a 256 $(BUILDFILES) > ${APPNAME}-${VERSION}.shasums

# Copy common files to release directory
copycommon:
	@echo "Copying common files to $(TMPDIR)"
	@mkdir "$(TMPDIR)/bin"
	@cp -v LICENSE "$(TMPDIR)"
	@cp -v README.md "$(TMPDIR)"

# Move all to temporary directory and compress with common files
tar-everything: copycommon
	@mkdir --parents "$(PWD)/release/${VERSION}"

	@echo "tar-everything..."
	# GNU/Linux
	@for arch in $(LINUX_ARCHS); do \
	  echo "GNU/Linux tar... $$arch"; \
	  cp -v "$(PWD)/bin/linux-$$arch/${APPNAME}" "$(TMPDIR)/bin"; \
	  cd "$(TMPDIR)"; \
	  tar -zcvf "$(PWD)/release/${VERSION}/${APPNAME}-${VERSION}-linux-$$arch.tar.gz" . ; \
	  rm "$(TMPDIR)/bin/${APPNAME}"; \
	done

	# Darwin
	@for arch in $(DARWIN_ARCHS); do \
	  echo "Darwin tar... $$arch"; \
	  cp -v "$(PWD)/bin/darwin-$$arch/${APPNAME}" "$(TMPDIR)/bin"; \
	  cd "$(TMPDIR)"; \
	  tar -zcvf "$(PWD)/release/${VERSION}/${APPNAME}-${VERSION}-darwin-$$arch.tar.gz" . ; \
	  rm "$(TMPDIR)/bin/${APPNAME}"; \
	done

	# FreeBSD
	@for arch in $(FREEBSD_ARCHS); do \
	  echo "FreeBSD tar... $$arch"; \
	  cp -v "$(PWD)/bin/freebsd-$$arch/${APPNAME}" "$(TMPDIR)/bin"; \
	  cd "$(TMPDIR)"; \
	  tar -zcvf "$(PWD)/release/${VERSION}/${APPNAME}-${VERSION}-freebsd-$$arch.tar.gz" . ; \
	  rm "$(TMPDIR)/bin/${APPNAME}"; \
	done

	# OpenBSD
	@for arch in $(OPENBSD_ARCHS); do \
	  echo "OpenBSD tar... $$arch"; \
	  cp -v "$(PWD)/bin/openbsd-$$arch/${APPNAME}" "$(TMPDIR)/bin"; \
	  cd "$(TMPDIR)"; \
	  tar -zcvf "$(PWD)/release/${VERSION}/${APPNAME}-${VERSION}-openbsd-$$arch.tar.gz" . ; \
	  rm "$(TMPDIR)/bin/${APPNAME}"; \
	done

	# NetBSD
	@for arch in $(NETBSD_ARCHS); do \
	  echo "NetBSD tar... $$arch"; \
	  cp -v "$(PWD)/bin/netbsd-$$arch/${APPNAME}" "$(TMPDIR)/bin"; \
	  cd "$(TMPDIR)"; \
	  tar -zcvf "$(PWD)/release/${VERSION}/${APPNAME}-${VERSION}-netbsd-$$arch.tar.gz" . ; \
	  rm "$(TMPDIR)/bin/${APPNAME}"; \
	done

	#Windows
	@for arch in $(WINDOWS_ARCHS); do \
	  echo "MS Windows zip... $$arch"; \
	  cp -v "$(PWD)/bin/windows-$$arch/${APPNAME}.exe" "$(TMPDIR)/bin"; \
	  cd "$(TMPDIR)"; \
	  zip -9 -y -r "$(PWD)/release/${VERSION}/${APPNAME}-${VERSION}-windows-$$arch.zip" . ; \
	  rm "$(TMPDIR)/bin/${APPNAME}.exe"; \
	done

	rm -rf "$(TMPDIR)/*"

# Distro: Arch linux - https://www.archlinux.org/
# Generates multi-arch PKGBUILD
ldistro-arch:
	pushd release/linux/arch && go run . -version ${VERSION} > "$(PWD)/release/${VERSION}/${APPNAME}-${VERSION}-linux-Arch.PKGBUILD"

# RPM
ldistro-rpm:
	@for arch in $(LINUX_ARCHS); do \
	  echo "Generating RPM... $$arch"; \
	  cd "$(TMPDIR)"; \
	  echo "  Extracting source package.." ; \
	  tar -xzf "$(PWD)/release/${VERSION}/${APPNAME}-${VERSION}-linux-$$arch.tar.gz" . ; \
	  rpmbuild --define "_builddir $(TMPDIR)" --define "_version ${VERSION}" --define "buildarch $$arch" -bb "$(PWD)/release/linux/rpm/package.spec" ; \
	  rm -rf "$(TMPDIR)/*"; \
	  echo ""; \
	  echo "------------------------------------------------------------"; \
	  echo ""; \
	done

# Create FreeBSD binary release package
# uses FreeBSD's pkg https://github.com/freebsd/pkg
# pkg help create
bsd-freebsd:
	@for arch in $(FREEBSD_ARCHS); do \
	  echo "Generate FreeBSD package... $$arch"; \
	  cd "$(TMPDIR)"; \
	  echo "  Extracting source package.." ; \
	  tar -xzf "$(PWD)/release/${VERSION}/${APPNAME}-${VERSION}-freebsd-$$arch.tar.gz" . ; \
	  echo "  Creating directory structure for package.." ; \
	  mkdir -p ./usr/local/bin ; \
	  mv ./bin/${APPNAME} ./usr/local/bin ; \
	  rm -rf ./bin ; \
	  cp "$(PWD)/release/freebsd/manifest.sh" /tmp/__manifest.sh; \
	  sed -i 's/<VERSION>/${VERSION}/' /tmp/__manifest.sh ; \
	  sed -i "s/<ARCH>/$$arch/" /tmp/__manifest.sh ; \
	  cat /tmp/__manifest.sh ; \
	  echo "  Creating pkg binary release package.." ; \
	  pkg create --verbose --format txz --root-dir $(TMPDIR) --manifest /tmp/__manifest.sh && \
	  cp "${APPNAME}-${VERSION}.txz" "$(PWD)/release/${VERSION}/${APPNAME}-${VERSION}-freebsd-pkg-$$arch.txz" ; \
	  rm -rf "$(TMPDIR)/*"; \
	  echo ""; \
	  echo "------------------------------------------------------------"; \
	  echo ""; \
	done


.PHONY: all clean test default