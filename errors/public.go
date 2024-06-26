// Copyright 2021 OpenSSF Scorecard Authors
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

package errors

import (
	"errors"
	"fmt"
)

var (
	// some of these errors didn't follow naming conventions when they were introduced.
	// for backward compatibility reasons, they can't be changed and have nolint directives.

	// ErrScorecardInternal indicates a runtime error in Scorecard code.
	ErrScorecardInternal = errors.New("internal error")
	// ErrRepoUnreachable indicates Scorecard is unable to establish connection with the repository.
	ErrRepoUnreachable = errors.New("repo unreachable")
	// ErrUnsupportedHost indicates the repo's host is unsupported.
	ErrUnsupportedHost = errors.New("unsupported host")
	// ErrInvalidURL indicates the repo's full URL was not passed.
	ErrInvalidURL = errors.New("invalid repo flag")
	// ErrShellParsing indicates there was an error when parsing shell code.
	ErrShellParsing = errors.New("error parsing shell code")
	// ErrJobOSParsing indicates there was an error when detecting a job's operating system.
	ErrJobOSParsing = errors.New("error parsing job operating system")
	// ErrUnsupportedCheck indicates check cannot be run for given request.
	ErrUnsupportedCheck = errors.New("check is not supported for this request")
	// ErrCheckRuntime indicates an individual check had a runtime error.
	ErrCheckRuntime = errors.New("check runtime error")
)

// WithMessage wraps any of the errors listed above.
// For examples, see errors/errors.md.
func WithMessage(e error, msg string) error {
	// Note: Errorf automatically wraps the error when used with `%w`.
	if len(msg) > 0 {
		return fmt.Errorf("%w: %v", e, msg)
	}
	// We still need to use %w to prevent callers from using e == ErrInvalidDockerFile.
	return fmt.Errorf("%w", e)
}

// GetName returns the name of the error.
func GetName(err error) string {
	switch {
	case errors.Is(err, ErrScorecardInternal):
		return "ErrScorecardInternal"
	case errors.Is(err, ErrRepoUnreachable):
		return "ErrRepoUnreachable"
	case errors.Is(err, ErrShellParsing):
		return "ErrShellParsing"
	default:
		return "ErrUnknown"
	}
}
