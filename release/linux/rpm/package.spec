%define _localbindir /usr/local/bin

Name:   heksa
Version:    %{_version}
Release:    1
Summary:    CLI hex dumper with colors
URL: https://github.com/raspi/heksa

Group:  DemoGroup
License:    Apache-2.0
BuildRoot:  %(mktemp -ud %{name}-%{version}-%{release}-XXXXXX)
BuildArch: %{buildarch
BuildRoot: %{_tmppath}/%{name}-%{version}-%{release}-root

%description
heksa is a command line hex binary dumper which uses ANSI colors

%prep

%setup -q

%build

%clean

rm -rf $RPM_BUILD_ROOT

%files

/usr/bin/heksa

%install

install -Dm644 "LICENSE" -t "$RPM_BUILD_ROOT/usr/share/licenses/%{name}"
install -Dm644 "README.md" -t "$RPM_BUILD_ROOT/usr/share/doc/%{name}"
install -Dm755 "bin/%{name}" -t "$RPM_BUILD_ROOT/usr/bin"

%changelog