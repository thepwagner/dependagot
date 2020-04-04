# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: dependabot/v1/dependabot.proto

require 'google/protobuf'

Google::Protobuf::DescriptorPool.generated_pool.build do
  add_file("dependabot/v1/dependabot.proto", :syntax => :proto3) do
    add_message "dependabot.v1.FilesRequest" do
      map :files, :string, :bytes, 1
    end
    add_message "dependabot.v1.FilesResponse" do
      repeated :required_paths, :string, 1
      repeated :optional_paths, :string, 2
    end
    add_message "dependabot.v1.ListDependenciesRequest" do
    end
    add_message "dependabot.v1.Dependency" do
      optional :package, :string, 1
      optional :version, :string, 2
    end
    add_message "dependabot.v1.ListDependenciesResponse" do
      repeated :dependencies, :message, 1, "dependabot.v1.Dependency"
    end
  end
end

module Dependabot
  module V1
    FilesRequest = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("dependabot.v1.FilesRequest").msgclass
    FilesResponse = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("dependabot.v1.FilesResponse").msgclass
    ListDependenciesRequest = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("dependabot.v1.ListDependenciesRequest").msgclass
    Dependency = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("dependabot.v1.Dependency").msgclass
    ListDependenciesResponse = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("dependabot.v1.ListDependenciesResponse").msgclass
  end
end
