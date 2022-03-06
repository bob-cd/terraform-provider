resource "bob_pipeline" "dev_test" {
  group = "dev"
  name  = "test"
  image = "docker.io/library/golang:alpine"
  vars = {
    GOOS   = "linux"
    GOARCH = "amd64"
  }

  resource {
    name     = "source"
    type     = "external"
    provider = bob_resource_provider.resource_git.name
    params = {
      repo   = "https://github.com/lispyclouds/bob-example"
      branch = "main"
    }
  }

  step {
    needs_resource = "source"
    cmd            = "go test"
  }

  step {
    needs_resource = "source"
    cmd            = "go build -o app"
    produces_artifact {
      name  = "app"
      path  = "app"
      store = bob_artifact_store.artifact_local.name
    }
  }
}
