Name:      heksa
Version:   %{_version}
Release:   1%{?dist}
Summary:   CLI hex dumper with colors
URL:       https://github.com/raspi/heksa
Source0:   src.tar.gz

License: Apache-2.0

%description
heksa is a command line hex binary dumper which uses ANSI colors

%prep
%autosetup

%install
install -Dm755 bin/%{name} -t %{buildroot}/%{_bindir}

%files
%license LICENSE
%doc README.md

%{_bindir}/%{name}

%changelog
