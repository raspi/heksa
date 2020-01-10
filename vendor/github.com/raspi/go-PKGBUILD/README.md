# go-PKGBUILD

Generate Arch Linux PKGBUILD string from template struct

## Projects using this library

* [Json2ArchPkgBuild](https://github.com/raspi/Json2ArchPkgBuild) - This library as CLI application
* [heksa](https://github.com/raspi/heksa) - Cross-platform command line hex dumper
* Yours?
  * Send pull request or open new issue
 
 ## Example JSON
 
 ```json
 {
  "_meta": {
    "ver": "v1.0.0"
  },
  "maintainer": "John Doe",
  "maintainer_email": "jd@example.org",
  "name": [
    "exampleapp"
  ],
  "version": "v1.0.0",
  "release": 1,
  "release_time": "1970-01-01T02:00:00+02:00",
  "short_description": "my example application",
  "licenses": [
    "Apache 2.0"
  ],
  "url": "https://github.com/examplerepo/exampleapp",
  "changelog_file": "",
  "groups": null,
  "dependencies": {
    "": {
      "packages": [
        "example-core"
      ],
      "build_packages": [
        "example-dev"
      ],
      "test_packages": [
        "example-test"
      ]
    },
    "x86_64": {
      "packages": [
        "example-core-x86"
      ]
    }
  },
  "optional_packages": {
    "": [
      {
        "package": "php",
        "reason": "because PHP is EPIC!"
      }
    ]
  },
  "provides": null,
  "options": [
    "!strip",
    "docs",
    "libtool",
    "staticlibs",
    "emptydirs",
    "!zipman",
    "!ccache",
    "!distcc",
    "!buildflags",
    "makeflags",
    "!debug"
  ],
  "install": "$pkgname.install",
  "files": {
    "aarch64": [
      {
        "url": "https://github.com/examplerepo/exampleapp/releases/download/$pkgver/$pkgname-$pkgver-linux-arm64.tar.gz",
        "checksums": {
          "sha256": "11d2b36d6b320dfee489d475635b53206b59288537554ea8bc24f97d06139d64"
        }
      }
    ],
    "arm": [
      {
        "url": "https://github.com/examplerepo/exampleapp/releases/download/$pkgver/$pkgname-$pkgver-linux-arm.tar.gz",
        "checksums": {
          "sha256": "5e79210655a9a71a7b77a3168194e9ead024a120182fa8560348a24dc87da159"
        }
      }
    ],
    "ppc64": [
      {
        "url": "https://github.com/examplerepo/exampleapp/releases/download/$pkgver/$pkgname-$pkgver-linux-ppc64.tar.gz",
        "checksums": {
          "sha256": "f744e32caf67a609aa435df9f8c519460b1856f7968c057e6ba61397cf79ec15"
        }
      }
    ],
    "ppc64le": [
      {
        "url": "https://github.com/examplerepo/exampleapp/releases/download/$pkgver/$pkgname-$pkgver-linux-ppc64le.tar.gz",
        "checksums": {
          "sha256": "6baef7ee046ceb4450e703a87f05fa5662708d4c3562c26abb427d34b4c82819"
        }
      }
    ],
    "x86_64": [
      {
        "url": "https://github.com/examplerepo/exampleapp/releases/download/$pkgver/$pkgname-$pkgver-linux-amd64.tar.gz",
        "checksums": {
          "sha256": "de3edfb94d5d0ae3d027c6c743e27290fa0500da4777da57154f2acab52775bf"
        }
      }
    ]
  },
  "commands": {
    "prepare": [
      "echo foo \u003e\u003e main.c"
    ],
    "build": [
      "make"
    ],
    "test": [
      "make test"
    ],
    "install": [
      "cd \"$srcdir\"",
      "install -Dm644 \"LICENSE\" -t \"$pkgdir/usr/share/licenses/$pkgname\"",
      "install -Dm644 \"README.md\" -t \"$pkgdir/usr/share/doc/$pkgname\"",
      "install -Dm755 \"bin/$pkgname\" -t \"$pkgdir/usr/bin\""
    ]
  }
}
 ```
 
## Example PKGBUILD output:
 
 ```bash
# Maintainer: John Doe <jd@example.org>
# Generated at: 2020-01-10 00:42:46.792588521 +0200 EET m=+0.000536267 

pkgname=exampleapp
pkgver=v1.0.0
pkgrel=1
pkgdesc="my example application"
url="https://github.com/examplerepo/exampleapp"
license=('Apache 2.0')
arch=('aarch64' 'arm' 'ppc64' 'ppc64le' 'x86_64')
install=$pkgname.install
depends_x86_64=('example-core-x86')

depends=('example-core')

makedepends=('example-dev')

checkdepends=('example-test')
optdepends=('php: because PHP is EPIC!')
sha256sums_aarch64=('11d2b36d6b320dfee489d475635b53206b59288537554ea8bc24f97d06139d64')
sha256sums_arm=('5e79210655a9a71a7b77a3168194e9ead024a120182fa8560348a24dc87da159')
sha256sums_ppc64=('f744e32caf67a609aa435df9f8c519460b1856f7968c057e6ba61397cf79ec15')
sha256sums_ppc64le=('6baef7ee046ceb4450e703a87f05fa5662708d4c3562c26abb427d34b4c82819')
sha256sums_x86_64=('de3edfb94d5d0ae3d027c6c743e27290fa0500da4777da57154f2acab52775bf')
source_aarch64=("https://github.com/examplerepo/exampleapp/releases/download/$pkgver/$pkgname-$pkgver-linux-arm64.tar.gz")
source_arm=("https://github.com/examplerepo/exampleapp/releases/download/$pkgver/$pkgname-$pkgver-linux-arm.tar.gz")
source_ppc64=("https://github.com/examplerepo/exampleapp/releases/download/$pkgver/$pkgname-$pkgver-linux-ppc64.tar.gz")
source_ppc64le=("https://github.com/examplerepo/exampleapp/releases/download/$pkgver/$pkgname-$pkgver-linux-ppc64le.tar.gz")
source_x86_64=("https://github.com/examplerepo/exampleapp/releases/download/$pkgver/$pkgname-$pkgver-linux-amd64.tar.gz")

prepare() {
  echo foo >> main.c
}

build() {
  make
}

check() {
  make test
}

package() {
  cd "$srcdir"
  install -Dm644 "LICENSE" -t "$pkgdir/usr/share/licenses/$pkgname"
  install -Dm644 "README.md" -t "$pkgdir/usr/share/doc/$pkgname"
  install -Dm755 "bin/$pkgname" -t "$pkgdir/usr/bin"
}
 ```
