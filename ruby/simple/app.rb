require 'rubygems'
require 'bundler/setup'
require 'turbine_rb'

class MyApp
  def call(app)
    database = app.resource(name: 'demopg')

    # ELT pipeline example
    # records = database.records(collection: 'events')
    # database.write(records: records, collection: 'events_copy')

    records = database.records(collection: 'events',configs:{"incrementing.column.name" => "id"})

    # This register the secret to be available in the turbine application
    app.register_secrets("MY_ENV_TEST") 

    # you can also register several secrets at once
    # app.register_secrets(["MY_ENV_TEST", "MY_OTHER_ENV_TEST"])

    processed_records = app.process(records: records, process: Passthrough.new) # Passthrough just has to match the signature
    database.write(records: processed_records, collection: "events_copy")
  end
end

class Passthrough < TurbineRb::Process 
  def call(records:)
    puts "got records: #{records}"
    # records.map { |r| r.value = 'hi there' }
    records
  end
end

TurbineRb.register(MyApp.new)
