package main

import (
	"log"
	"os"
	"reflect"
)

type Envs struct {
	BASE_GH_TOKEN,
	BASE_OWNER,
	BASE_REPO,
	DOMAIN,
	URL,
	APP_NAME string
}

func loadEnvs() *Envs {
	envs := Envs{}
	ref := reflect.ValueOf(&envs).Elem()

	for k := range ref.Fields() {
		env, ok := os.LookupEnv(k.Name)
		if !ok {
			log.Fatalf("no env %q is set", k.Name)
		}
		ref.FieldByName(k.Name).SetString(env)
	}

	return &envs
}
