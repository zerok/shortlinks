#!/bin/sh
if [ -z "$TOKEN" ]; then
   echo "Please set the TOKEN environment variable."
   exit 1
fi
exec /usr/local/bin/shortlinks --db /data/shortlinks.sqlite --addr 0.0.0.0:8000 --migrations /var/lib/shortlinks/migrations --token "$TOKEN"
