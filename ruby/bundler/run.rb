require 'grpc'
require_relative 'service/service'

def main
  port = '0.0.0.0:9999'
  s = GRPC::RpcServer.new
  s.add_http2_port(port, :this_port_is_insecure)
  s.handle(Updater.new)
  s.run_till_terminated_or_interrupted([1, 'int', 'SIGQUIT'])
end

main