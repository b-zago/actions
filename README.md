# actions

Reusable GitHub Actions and workflows that all my repositories use for CI/CD. Each application's pipeline is a few lines that call into this repo, so building, versioning and deploying work the same way everywhere and are maintained in one place.

```yaml
# The whole CI/CD pipeline of an application repository
on:
  push:
    branches: [main, staging]

jobs:
  deploy:
    uses: b-zago/actions/.github/workflows/box-deployment.yml@v1
    permissions:
      contents: write
      packages: write
    secrets: inherit
```

A push to `staging` builds the image, tags a release and deploys it to my staging Kubernetes cluster. A push to `main` does the same for production.

## Highlights

- **GitOps deployments.** Kubernetes deployments never touch the cluster directly. The workflow updates the image tag in my [infrastructure repo](https://github.com/b-zago/box), and Argo CD rolls the change out. The branch decides the environment.
- **Automatic semantic versioning.** Every build gets a version from the commit message: patch by default, `[MINOR]` or `[MAJOR]` for bigger bumps, and prerelease tags like `v1.4.2-staging.37` on other branches.
- **Multi-architecture images** (`amd64` and `arm64`) with build caching.
- **No static cloud credentials.** The AWS workflow authenticates through GitHub's OIDC provider and assumes an IAM role, so no AWS keys are stored in GitHub.
- **Versioned like public actions.** Releases are tagged `vX.Y.Z` and a floating `v1` tag follows the latest release, so consumers get fixes without changing their pipelines.

## Contents

- [Workflows](#workflows)
  - [box-deployment](#box-deployment)
  - [build-push-ecr-lambda](#build-push-ecr-lambda)
- [Actions](#actions)
  - [build-push-ghcr](#build-push-ghcr)
  - [semver-bump](#semver-bump)
  - [test-build](#test-build)
- [Versioning](#versioning)

## Workflows

Reusable workflows, called with `uses: b-zago/actions/.github/workflows/<name>.yml@v1`.

### box-deployment

Builds and publishes the repository's image, then deploys it by updating the image tag in the matching environment of [box](https://github.com/b-zago/box).

1. Builds a multi-arch image and pushes it to GHCR with a new semantic version ([build-push-ghcr](#build-push-ghcr)).
2. Picks the environment from the branch: `main` deploys to `prod`, any other branch to the environment of the same name.
3. Updates the tag in `k3s/overlays/<env>/workloads/<repository>/kustomization.yaml` in box and commits the change. If the app doesn't exist in that environment, the step is skipped.

| Input       | Default         | Description                                       |
| ----------- | --------------- | ------------------------------------------------- |
| `imageName` | repository name | Image name under `ghcr.io/<owner>/<repository>/`. |

**Requires:** a `PAT_TOKEN` secret with write access to box, `secrets: inherit` in the caller, and `contents: write` plus `packages: write` permissions (see the example above).

### build-push-ecr-lambda

Builds an `arm64` image on a native ARM runner, pushes it to Amazon ECR and deploys it to AWS Lambda, then waits until the function is updated.

- `main` deploys to the function `<repository>-prod`, any other branch to `<repository>-staging`.
- Images are tagged `<env>-<commit sha>`.
- AWS access goes through OIDC: the job assumes the IAM role in the `AWS_ROLE` secret, with no access keys involved.

```yaml
jobs:
  deploy:
    uses: b-zago/actions/.github/workflows/build-push-ecr-lambda.yml@v1
    permissions:
      id-token: write
      contents: read
    secrets: inherit
```

**Requires:** an ECR repository named like the GitHub repository, the two Lambda functions, and an IAM role that trusts GitHub's OIDC provider for this repository (its ARN in the `AWS_ROLE` secret). Region: `eu-central-1`.

## Actions

Composite actions, called with `uses: b-zago/actions/<name>@v1`.

### build-push-ghcr

Bumps the repository's version with [semver-bump](#semver-bump), then builds a multi-arch image (`linux/amd64`, `linux/arm64`) and pushes it to `ghcr.io/<owner>/<repository>/<imageName>`, tagged with the version and `latest`. Uses the GitHub Actions cache for Docker layers.

| Input       | Required | Description                                                               |
| ----------- | -------- | ------------------------------------------------------------------------- |
| `imageName` | yes      | Image name.                                                               |
| `token`     | yes      | Token with `packages: write` and `contents: write` (for pushing the tag). |

| Output | Description                                           |
| ------ | ----------------------------------------------------- |
| `tag`  | The new version without the `v` prefix, e.g. `1.4.2`. |

Used directly by [rikami-operator](https://github.com/b-zago/rikami-operator)'s release pipeline.

### semver-bump

Creates and pushes the next semantic version tag, based on the latest `vX.Y.Z` tag in the repository.

| Commit message contains | Result                     |
| ----------------------- | -------------------------- |
| nothing special         | patch: `v1.4.1` → `v1.4.2` |
| `[MINOR]`               | minor: `v1.4.1` → `v1.5.0` |
| `[MAJOR]`               | major: `v1.4.1` → `v2.0.0` |

The first tag in a repository is `v1.0.0`. On branches other than `main`, the version gets a prerelease suffix with the branch name and run number, e.g. `v1.4.2-staging.37`, so prereleases never become the base for the next stable version.

| Input   | Required | Description                   |
| ------- | -------- | ----------------------------- |
| `token` | yes      | Token with `contents: write`. |

| Output    | Description                    |
| --------- | ------------------------------ |
| `new-tag` | The pushed tag, e.g. `v1.4.2`. |

### test-build

Builds the image without pushing it, to check that it builds. Meant for pull requests.

| Input       | Required | Description                         |
| ----------- | -------- | ----------------------------------- |
| `imageName` | yes      | Image name, used for the local tag. |

## Versioning

This repository is released like a public action:

- Pushing a tag `vX.Y.Z` triggers `build-actions.yml`, which moves the floating major tag (`v1`) to the new release. Consumers pin `@v1` and receive compatible updates automatically.
- The same workflow also supports container-based actions: anything under `custom/<name>/` that changed since the previous release is built, pushed to GHCR and referenced from its `action.yml` with the new version.
