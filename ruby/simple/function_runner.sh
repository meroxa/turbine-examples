#! /bin/sh
app=$1
function=$2
# bundle exec ruby  -e "load('$app'); '$function.new.serve'"
bundle exec ruby -e "require_relative('$app'); $function.new.serve"
# bundle exec ruby ./app.rb 
