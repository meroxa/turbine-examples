#!/usr/bin/env ruby -I ..
# frozen_string_literal: true

require 'turbine'

def app
  app = Turbine::App.new('myapp')

  database = app.resource(name: 'demopg')

  # pipeline demo
  records = database.records(collection: 'events')
  database.write(records: records, collection: 'events_copy')

  # procedural API
  # records = database.records(collection: 'events')
  # processed_records = app.process(records: records, process: Passthrough)
  # database.write(records: processed_records, collection: "events_copy")

  # chaining API
  # database.records(collection: "events").
  #   process_with(process: Passthrough).
  #   write_to(resource: database, collection: "events_copy")
end

class Passthrough
  def process(records:)
    records.each do |r|
      r.value = 'hi there'
    end
    records
  end
end

class SomethingElse 
  def process(records:)
    records.each do |r|
      r.value = 'foo'
    end
    records
  end
end