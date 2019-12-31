LAST_TAG := $(shell git describe --abbrev=0 --always --tags)
BUILD := $(shell git rev-parse $(LAST_TAG))
BUILDDATE := $(shell git log -1 --format=%aI $(LAST_TAG))

BINARY := heksa
UNIXBINARY := $(BINARY)
WINBINARY := $(UNIXBINARY).exe
BUILDDIR := build

LINUXRELEASE := $(BINARY)-$(LAST_TAG)-linux-x64.tar.gz
WINRELEASE := $(BINARY)-$(LAST_TAG)-windows-x64.zip

LDFLAGS := -ldflags "-s -w -X=main.VERSION=$(LAST_TAG) -X=main.BUILD=$(BUILD) -X=main.BUILDDATE=$(BUILDDATE)"

SCREENSHOTCMD := ./heksa -f hex,asc,bit -l 0x100 heksa.exe

bin:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -v -o $(BUILDDIR)/$(UNIXBINARY)
	upx -v -9 $(BUILDDIR)/$(UNIXBINARY)

bin-windows:
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -v -o $(BUILDDIR)/$(WINBINARY)
	upx -v -9 $(BUILDDIR)/$(WINBINARY)

release:
	cd $(BUILDDIR); tar cvzf $(LINUXRELEASE) $(UNIXBINARY)

release-windows:
	cd $(BUILDDIR); zip -v -9 $(WINRELEASE) $(WINBINARY)

screenshot:
	cd build; echo "% $(SCREENSHOTCMD)" > scr.txt && $(SCREENSHOTCMD) >> scr.txt && konsole --notransparency --noclose -e cat scr.txt

.PHONY: all clean test