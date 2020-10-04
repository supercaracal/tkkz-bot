# frozen_string_literal: true

# Download Slack history archive
# cat general/*.json | jq -r '.[] | .text' > messages.txt
# ruby data_importer.rb messages.txt

require 'uri'
require 'net/http'

File.open(ARGV[0], 'r') do |f|
  loop do
    line = f.readline
    next if line.empty?

    uri = URI('http://127.0.0.1:3000')
    params = { speaker: 'human', message: line }
    uri.query = URI.encode_www_form(params)

    res = Net::HTTP.get_response(uri)
    p line unless res.is_a?(Net::HTTPSuccess)
  end
end
