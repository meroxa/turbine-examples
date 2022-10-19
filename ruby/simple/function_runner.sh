#! /bin/sh
function=$1
bundle exec ruby -e "require_relative('./app'); $function.new.serve"
