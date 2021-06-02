#!/bin/bash
set -ex
if [ "$TZ" != "" ]; then
    ln -sf /usr/share/zoneinfo/$TZ /etc/localtime
    echo $TZ > /etc/timezone
fi
if [[ ! "$1" = *"taosd" ]]; then
    $@
    exit
fi
DATA_DIR=${TAOS_DATA_DIR:-/var/lib/taos}

if [ "$TAOS_FQDN" = "" ]; then
    echo TAOS_FQDN must be setted!
    exit 255
fi

ORIG_FQDN=$(cat /etc/taos/taos.cfg|grep -v "^#"| grep fqdn |sed -E 's/.*fqdn\s+//')

if [ "$ORIG_FQDN" != "" ] && [ "$ORIG_FQDN" != "$TAOS_FQDN" ]; then
    echo FQDN should not be changable after initialized!
    exit 254
fi

grep "$TAOS_FQDN" /etc/hosts || echo "127.0.0.1 $TAOS_FQDN" >> /etc/hosts

# write config to file
env-to-cfg > /etc/taos/taos.cfg

CLUSTER=${CLUSTER:=}
FIRST_EP_HOST=${TAOS_FIRST_EP%:*}
SERVER_PORT=${TAOS_SERVER_PORT:-6030}
# if has mnode ep set or the host is first ep or not for cluster, just start.
if [ -f "$DATA_DIR/dnode/mnodeEpSet.json" ] || \
  [ "$TAOS_FQDN" = "$FIRST_EP_HOST" ] || [ "$CLUSTER" = "" ]; then
    $@
# others will first wait the first ep ready.
else
    if [ "$CLUSTER" != "" ] && [ "$TAOS_FIRST_EP" == "" ]; then
        echo "TAOS_FIRST_EP must be setted in cluster"
        exit
    fi
    while true
    do
        es=0
        taos -h $FIRST_EP_HOST -s "show mnodes" > /dev/null || es=$?
        if [ "$es" -eq 0 ]; then
            taos -h $FIRST_EP_HOST -s "create dnode \"$TAOS_FQDN:$SERVER_PORT\";"
            break
        fi
        sleep 1s
    done
    $@
fi
