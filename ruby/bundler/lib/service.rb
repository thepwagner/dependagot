require 'dependabot-ruby-common'

class UpdateServiceHandler

  def files(req, _call)
    st = state
    req.files.each { |fn, data| 
      st[fn] = data
    }

    optional_paths = []
    if !st.has_key?("Gemfile")
      optional_paths.push("Gemfile")
    end
    if !st.has_key?("Gemfile.lock")
      optional_paths.push("Gemfile.lock")
    end

    ::Dependabot::V1::FilesResponse.new(
      optional_paths: optional_paths,
    )
  end

  def list_dependencies(req, _call) 
    ::Dependabot::V1::ListDependenciesResponse.new(
      dependencies: [],
    )
  end

  def update_dependencies(req, _call) 
    ::Dependabot::V1::UpdateDependenciesResponse.new(
      new_files: Hash.new,
    )
  end

  def state
    @state || Hash.new
  end
end