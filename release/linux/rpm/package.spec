Name:   heksa
Version:    %{_version}
Release:    1
Summary:    CLI hex dumper with colors
URL: https://github.com/raspi/heksa

Group:  DemoGroup
License:    Apache-2.0

%description
heksa is a command line hex binary dumper which uses ANSI colors

%setup -q

%clean

rm -rf %{buildroot}

%files

/usr/bin/heksa

%license /usr/share/licenses/heksa/LICENSE
%doc /usr/share/doc/heksa/README.md

%install
cd %{buildroot}
install -Dm755 "bin/%{NAME}" -t "%{buildroot}/usr/bin"
