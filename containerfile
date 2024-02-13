from docker.io/library/golang:1.20.6-bookworm
workdir /src
copy . .
run go build -o gcshell cmd/gcshell/main.go

from docker.io/library/debian:12.1
arg SNX_URL
arg SNX_PREFIX
run set -x \
 && dpkg --add-architecture i386 \
 && apt-get update -y\
 && apt-get install -y \
    bzip2 kmod expect iptables net-tools iputils-ping iproute2 curl wget bash unzip udev procps\
    libstdc++5:i386 libpam0g:i386 libx11-6:i386
run curl -sfLko /tmp/snx_install.sh "${SNX_HOST_URL:-https://ext.vpn.vtb.ru}${SNX_PREFIX:-/sslvpn}/SNX/INSTALL/snx_install.sh" \
 && bash /tmp/snx_install.sh && rm -rf /tmp/snx_install.sh \
 && apt-get install --yes opensc pcscd openssl libengine-pkcs11-openssl gnutls-bin libssl3 \
 && curl -sfLko /tmp/librtpkcs11ecp.deb http://download.rutoken.ru/Rutoken/PKCS11Lib/2.8.1.0/Linux/x64/librtpkcs11ecp_2.8.1.0-1_amd64.deb \
 && curl -sfLko /tmp/libssl1.1.deb http://security.debian.org/debian-security/pool/updates/main/o/openssl/libssl1.1_1.1.1n-0+deb11u5_amd64.deb \
 && curl -sfLko /tmp/safenet.zip https://www.globalsign.com/en/safenet-drivers/USB/10.7/Safenet_Linux_Installer_DEB_x64.zip \
 && cd /tmp && unzip /tmp/safenet.zip && apt-get install --yes ./*.deb && rm -rf /tmp/* && cd /
copy --from=0 /src/gcshell /usr/local/bin/
copy entrypoint.sh /usr/local/bin/
entrypoint ["/usr/local/bin/entrypoint.sh"]
cmd ["/usr/local/bin/gcshell"]
