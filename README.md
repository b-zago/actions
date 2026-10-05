# Actions and workflows

My GitHub Actions and reusable workflows.

## Actions

- [build-push-ghcr](./build-push-ghcr/action.yml) - simply builds and pushes multiplatform docker image to a ghcr repo.
- [k3s-bump](./k3s-bump/action.yml) - bumps image tag in my [k3s-cluster](https://github.com/b-zago/k3s-cluster) repo.
- [test-build](./test-build/action.yml) - just builds an image to see if it even succeeds. Mainly used for PRs to test.

## Workflows

- [build-deploy-ghcr](./.github/workflows/build-deploy-ghcr.yml) - builds and pushes multiplatform docker image and bumps image tag in [k3s-cluster](https://github.com/b-zago/k3s-cluster) repo.
- [build-push-ecr-lambda](./.github/workflows/build-push-ecr-lambda.yml) - builds and pushes arm64 docker image to ECR repository and updates the lambda function with the latest tag.
