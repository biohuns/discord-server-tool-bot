%define debug_package %{nil}

Name: discord-servertool
Version: %{_version}
Release: %{_release}%{?_dist}

Summary: ServerTool for Discord
License: MIT
Group: Applications/Internet
URL: https://github.com/biohuns/discord-servertool

Source0: %{name}-%{version}.tar.gz
BuildRoot: %{_tmppath}/%{name}-%{version}-%{release}-root

%description
GCE + Steam Dedicated Server を Discord で管理する Bot

%prep
%setup -q

%build

%install
install -D %{name} %{buildroot}%{_bindir}/%{name}
install -D config.json %{buildroot}%{_sysconfdir}/%{name}/config.json
install -D credential.json %{buildroot}%{_sysconfdir}/%{name}/credential.json
install -D systemd.service %{buildroot}/etc/systemd/system/%{name}.service

%files
%attr(0755,root,root) %{_bindir}/%{name}
%defattr(0644,root,root, 0755)
%config(noreplace) %{_sysconfdir}/%{name}/config.json
%config(noreplace) %{_sysconfdir}/%{name}/credential.json
%config(noreplace) %{_sysconfdir}/systemd/system/%{name}.service

%post
if [ $1 -eq 1 ]; then
    /bin/systemctl enable %{name}.service
fi
/bin/systemctl stop %{name}.service || :
/bin/systemctl daemon-reload >/dev/null 2>&1 || :
/bin/systemctl start %{name}.service

%preun
if [ $1 -eq 0 ]; then
    /bin/systemctl stop %{name}.service || :
fi

%clean
rm -rf %{buildroot}
