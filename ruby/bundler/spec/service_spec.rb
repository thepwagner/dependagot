require 'spec_helper'
require_relative '../lib/service'

RSpec.describe Updater do
  it 'should return response' do
    u = Updater.new
    res = u.files(nil, nil)
    expect(res.paths).to contain_exactly("ruby", "server")
  end
end