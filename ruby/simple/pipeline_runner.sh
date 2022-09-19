#! /bin/sh
trap 'kill %1' SIGINT SIGTERM EXIT

mode=local

if [ "$1" = "platform" ]; then
    mode=platform
fi

turbine-core --mode=$mode &
ruby -I ../proto -I .. ../pipeline_runner.rb
