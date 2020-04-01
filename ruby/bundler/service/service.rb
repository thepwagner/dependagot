require 'dependabot-ruby-common'

class Updater < ::Dependabot::V1::UpdateService::Service
  def files(req, _call)
    ::Dependabot::V1::FilesResponse.new(paths: ["ruby","server"])
  end
end