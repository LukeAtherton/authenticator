#!/bin/bash

# Fail hard and fast
#set -eo pipefail

export ETCD=$HOST_IP:4001

echo "[authenticator] booting container. ETCD: $ETCD"

# Loop until confd has updated the authenticator config
until confd -onetime -node $ETCD -config-file /etc/confd/conf.d/authenticator.toml; do
  echo "[authenticator] waiting for confd to refresh authenticator.conf (waiting for message queue to be available)"
  sleep 5
done

sed -i "s/HOST_IP/$HOST_IP/g" /etc/authenticator.cfg

# Start authenticator
echo "[authenticator] starting authenticator..."
/go/src/github.com/lukeatherton/authenticator/authenticator --config="/etc/authenticator.cfg"
