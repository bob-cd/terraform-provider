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
    provider = "resource-git"
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
      store = "artifact-local"
    }
  }
}
