`Makefile` in the root directory generates version releases in this 
directory when running `make release`. 
It generates directory in `vX.Y.Z` format (for example `v1.11.0`).
The `vX.Y.Z` directory then contains compressed files which then contains common files such as `README.md` and `LICENSE` and of course the binary for 
different operating systems and architectures (for example `heksa-v1.11.0-linux-amd64.tar.gz`).

This directory also contains possible helper scripts, 
manifests and such for extra packaging such as Arch Linux PKGBUILD format.