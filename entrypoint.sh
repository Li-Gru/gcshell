#!/usr/bin/env bash
: ${CSHELL_PKCS11_LIB:='RUTOKEN'}
[ "${CSHELL_PKCS11_LIB}" == 'RUTOKEN' ] && export CSHELL_PKCS11_LIB="/usr/lib/librtpkcs11ecp.so"
[ "${CSHELL_PKCS11_LIB}" == 'ETOKEN'  ] && export CSHELL_PKCS11_LIB="/usr/lib/libeToken.so"
: ${CSHELL_PKCS11_ID:=""}
: ${CSHELL_PKCS11_PIN:=""}
: ${CSHELL_PKCS11_SELECTOR:="example@localhost.local"}
: ${CSHELL_SNX_PREFIX:="/sslvpn"}
: ${CSHELL_SNX_GATEWAY:="vpn.localhost.local"}
: ${CSHELL_SNX_REALM:=""}
pcscd
udevadm control --reload-rules ; echo "initialized"
[ "$1" == "list" ] && exec pkcs11-tool --module ${CSHELL_PKCS11_LIB} -L | grep -v '(empty)'
exec $@