package main

import "errors"

// Flag ... cliの入力値を格納
type Flag struct {
	ecrProfile  string
	ecsProfiles []string
	region      string
	keep        int
}

func (f Flag) validate() error {
	if f.ecrProfile == "" {
		return errors.New("--ecr-profile option is required")
	}
	if len(f.ecsProfiles) == 0 {
		return errors.New("--ecs-profiles option is required")
	}
	if f.region == "" {
		return errors.New("-r or --region option is required")
	}

	return nil
}
