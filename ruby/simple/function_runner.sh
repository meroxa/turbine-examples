#! /bin/sh
app=$1
function=$2
# bundle exec ruby  -e "load('$app'); Turbine::Support::FunctionServer.run '$function'"
bundle exec ruby ./app.rb 
