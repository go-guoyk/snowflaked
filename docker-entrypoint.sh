#!/bin/bash

if [ "$1" = "snowflaked" ]; then
  shift
  CMD_OPTS=""
  if [ -n "${BIND}" ]; then
    CMD_OPTS="${CMD_OPTS} -bind ${BIND}"
  fi
  if [ -n "${CLUSTER_ID}" ]; then
    CMD_OPTS="${CMD_OPTS} -cluster-id ${CLUSTER_ID}"
  fi
  STATEFULSET_SEQ=$(hostname | sed -n "s/^.*-\([0-9]\+\)$/\1/p")
  if [ -n "${STATFULSET_SEQ}" ]; then
    CMD_OPTS="${CMD_OPTS} --worker-id ${STATEFULSET_SEQ}"
  fi
  set -- snowflaked $@ ${CMD_OPTS}
  echo "$@"
fi

exec "$@"
