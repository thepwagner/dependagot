require 'rack'
require 'webrick'

require_relative 'lib/service'

def main
  handler = UpdateServiceHandler.new
  service = ::Dependagot::V1::UpdateServiceService.new(handler)
  path_prefix = "/twirp/" + service.full_name
  server = WEBrick::HTTPServer.new(BindAddress: "0.0.0.0", Port: 9999)
  server.mount path_prefix, Rack::Handler::WEBrick, service
  server.start
end

main