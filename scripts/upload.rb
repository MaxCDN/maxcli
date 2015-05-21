#!/usr/bin/env ruby

USAGE = %q{
Usage: scripts/upload.rb [AWS KEY] [AWS SECRET]

  Built to upload binaries to AWS S3. Requires ruby gem, "fog".

  Accepts environment variables:
    PRETEND:: do not update/create AWS::S3 files
    WITHOUT_UPDATE:: only create new AWS::S3 files
    WITHOUT_CREATE:: only update existing AWS::S3 files

---
}

def on_error emessage
  puts USAGE
  abort "ERROR: #{emessage}"
end

require 'rubygems'
require 'pp'
begin
  require 'fog'
rescue
  on_error $!
end

on_error 'Missing required arguments.' unless ARGV.length == 2

options = {
  provider: 'AWS',
  aws_access_key_id: ARGV.first.strip,
  aws_secret_access_key: ARGV.last.strip
}

connection = Fog::Storage.new(options)

directory = connection.directories.get('maxcli')
remote_paths = directory.files.map(&:key).freeze

# prefix for absolute path on disk
local_prefix = File.absolute_path(File.join(File.dirname(__FILE__), '..'))

# local search pattern
local_search = File.join(local_prefix, 'max*', 'builds', '**/*')

# local files (not directories)
local_files  = Dir[local_search].reject { |path| File.directory?(path) }

local_files.each do |local_file|
  # this should never happen
  if File.directory?(local_file)
    raise "Expected a file, but got a directory."
  end

  remote_path = local_file.gsub(local_prefix+'/', '').gsub('builds/', '')

  if remote_paths.include?(remote_path)
    puts "Updating:: #{remote_path}"
    puts " with: #{local_file}"

    # update file
    unless ENV['PRETEND'] or ENV['WITHOUT_UPDATE']
      remote_file = directory.files.get(remote_path)
      remote_file.body = File.open(local_file)
      remote_file.public = true # just to be safe
      remote_file.save or abort "FAILED"
      puts " done!"
    else
      puts " but not really!"
    end
  else
    puts "Creating:: #{remote_path}"
    puts " with: #{local_file}"

    # create file
    unless ENV['PRETEND'] or ENV['WITHOUT_CREATE']
      directory.files.create(key: remote_path, body: File.open(local_file), public: true) or abort "FAILED"
      puts " done!"
    else
      puts " but not really!"
    end
  end

end
