# frozen_string_literal: true

require 'pstore'
require 'redis'
require 'webrick' # Oops!

PSTORE_DB_FILE_PATH = ENV['PSTORE_DB_FILE_PATH'] || '/var/tmp/reudy/db'
REDIS_URL = ENV['REDIS_URL'] || 'redis://127.0.0.1:6379/0'

pstore = PStore.new(PSTORE_DB_FILE_PATH)
redis = Redis.new(url: REDIS_URL)

redis.flushdb

pstore.transaction do
  pstore.roots.each do |k|
    v = pstore[k]
    case v
    when Array
      v.each { |e| redis.rpush(k, e) }
    when String
      redis.rpush(k, v)
    else
      raise "Not expected type=#{v.class} key=#{k} value=#{v.inspect}"
    end
  end
end
