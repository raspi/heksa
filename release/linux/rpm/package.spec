Name: heksa
Version: %{_version}
Release: 1%{?dist}
Summary: CLI hex dumper with colors
URL: https://github.com/raspi/heksa

Group: Applications/Utilities
License: Apache-2.0

%description
heksa is a command line hex binary dumper which uses ANSI colors

%setup -q

%clean

%files
%license /usr/share/licenses/%{NAME}/LICENSE
%doc /usr/share/doc/%{NAME}/README.md

/usr/bin/heksa

%install
install -Dm755 "usr/bin/%{NAME}" -t "/usr/bin"
