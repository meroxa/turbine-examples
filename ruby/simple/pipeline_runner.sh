#! /bin/sh
trap 'kill %1' SIGINT SIGTERM EXIT

mode=local

if [ "$1" = "platform" ]; then
    mode=platform
fi

app=$2

turbine-core --mode=$mode &
bundle exec ruby -e "require_relative('./app'); Turbine.run;"