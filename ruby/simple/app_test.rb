require 'minitest/autorun'
require_relative './app.rb'

class PassthroughTest < Minitest::Test

  def fixture(name)
    file = File.read("./fixtures/#{name}")
    json = JSON.parse(file)

    json.keys.each do |collection_name|
      json[collection_name] = json[collection_name].map do |record|
                                record_helper(record)
                              end
    end

    return json
  end

  def record_helper(hash)
    record = hash.keys.each_with_object({}) do |key, obj|
      obj[key.to_sym] = hash[key]
      if key.to_sym == :value
        obj[key.to_sym] = hash[key].to_json
      end
      if key.to_sym == :timestamp
        obj[key.to_sym] = hash[key].to_i
      end
    end

    pb_record = Io::Meroxa::Funtime::Record.new(record)
    TurbineRb::Record.new(pb_record)
  end

  # When records come into the Passthrough class, they
  # should be the same on the other side. There are some
  # helper functions to use the fixture data
  #
  def test_record_passthrough
    data = fixture("demo.json")
    passthrough = Passthrough.new
    records = passthrough.call(records: data["collection_name"])
    data["collection_name"].each_with_index do |data_record, index|
      [:key, :value, :timestamp].each do |method|
        assert data_record.send(method) == records[index].send(method)
      end
    end
  end

end
