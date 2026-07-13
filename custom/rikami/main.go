package main

import "github.com/b-zago/rikami/ci"

func main() {
	envs := loadEnvs()

	ci.Run(envs.BASE_GH_TOKEN, envs.BASE_OWNER, envs.BASE_REPO, envs.DOMAIN, envs.URL, envs.APP_NAME)
}
