# frozen_string_literal: true

require 'turbine'

class MyApp
  def call(app)
    # app = Turbine::App.new('myapp') # Don't really need this if we have Turbine.register

    database = app.resource(name: 'demopg')

    # pipeline demo
    # records = database.records(collection: 'events')
    # database.write(records: records, collection: 'events_copy')

    # procedural API
    records = database.records(collection: 'events')
    processed_records = app.process(records: records, process: Passthrough.new) # Passthrough just has to match the signature
    # processed_records = records.process_with(process: Foo.new)
    database.write(records: processed_records, collection: "events_copy")

    # out_records = processed_records.join(records, key: "user_id", window: 1.day) # stream joins

    # chaining API
    database.records(collection: "events").
      process_with(process: Passthrough).
      write_to(resource: database, collection: "events_copy")
  end
end

# Turbine.run(MyApp.new) # Turbine needs to discover where the app is.
# Turbine.register(MyApp.new) # this is the App that should be run
# Turbine.register(MyApp.new, name: app_name) # optional name arg otherwise name is pulled from app.json

class Passthrough < Turbine::Process # might be useful to signal that this is a special Turbine call
  def call(records:)
    puts "got records: #{records}"
    # records.map { |r| r.value = 'hi there' }
    records
  end
end

# Need to wrap this and call the named function
f = Passthrough.new
f.serve