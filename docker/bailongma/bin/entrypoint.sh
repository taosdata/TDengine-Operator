#!/bin/bash
set -e
if [ "$TZ" != "" ]; then
    ln -sf /usr/share/zoneinfo/$TZ /etc/localtime
    echo $TZ > /etc/timezone
fi
# write config to file
env-to-cfg > /etc/taos/taos.cfg

$@
