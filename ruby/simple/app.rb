# frozen_string_literal: true

require 'turbine'

class MyApp
  def call(app)
    database = app.resource(name: 'demopg')

    # ELT pipeline example
    # records = database.records(collection: 'events')
    # database.write(records: records, collection: 'events_copy')

    # procedural API
    records = database.records(collection: 'events')
    processed_records = app.process(records: records, process: Passthrough.new) # Passthrough just has to match the signature
    database.write(records: processed_records, collection: "events_copy")

    # out_records = processed_records.join(records, key: "user_id", window: 1.day) # stream joins

    # chaining API
    # database.records(collection: "events").
    #   process_with(process: Passthrough.new).
    #   write_to(resource: database, collection: "events_copy")
  end
end

class Passthrough < Turbine::Process # might be useful to signal that this is a special Turbine call
  def call(records:)
    puts "got records: #{records}"
    # records.map { |r| r.value = 'hi there' }
    records
  end
end

Turbine.register(MyApp.new)
