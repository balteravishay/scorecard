// Copyright 2024 OpenSSF Scorecard Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"gopkg.in/yaml.v3"

	sce "github.com/ossf/scorecard/v5/errors"
)

var (
	ErrInvalidCheck  = errors.New("check is not valid")
	ErrInvalidReason = errors.New("reason is not valid")
)

// Config contains configurations defined by maintainers.
type Config struct {
	Annotations []Annotation `yaml:"annotations"`
}

// parseFile takes the scorecard.yml file content and returns a `Config`.
func parseFile(c *Config, content []byte) error {
	unmarshalErr := yaml.Unmarshal(content, c)
	if unmarshalErr != nil {
		return sce.WithMessage(sce.ErrScorecardInternal, unmarshalErr.Error())
	}

	return nil
}

func isValidCheck(check string, checks []string) bool {
	for _, validCheck := range checks {
		if strings.EqualFold(check, validCheck) {
			return true
		}
	}
	return false
}

func validate(c Config, checks []string) error {
	for _, annotation := range c.Annotations {
		for _, check := range annotation.Checks {
			if !isValidCheck(check, checks) {
				return fmt.Errorf("%w: %s", ErrInvalidCheck, check)
			}
		}
		for _, reasonGroup := range annotation.Reasons {
			if !IsValidReason(reasonGroup.Reason) {
				return fmt.Errorf("%w: %s", ErrInvalidReason, reasonGroup.Reason)
			}
		}
	}
	return nil
}

// Parse reads the configuration file from the repo, stored in scorecard.yml, and returns a `Config`.
func Parse(r io.Reader, checks []string) (Config, error) {
	c := Config{}
	// Find scorecard.yml file in the repository's root
	content, err := io.ReadAll(r)
	if err != nil {
		return Config{}, fmt.Errorf("fail to read configuration file: %w", err)
	}

	err = parseFile(&c, content)
	if err != nil {
		return Config{}, fmt.Errorf("fail to parse configuration file: %w", err)
	}

	err = validate(c, checks)
	if err != nil {
		return Config{}, fmt.Errorf("configuration file is not valid: %w", err)
	}

	// Return configuration
	return c, nil
}
