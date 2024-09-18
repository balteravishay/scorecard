package raw

import (
	"os"
	"testing"

	"github.com/ossf/scorecard/v5/checker"
)

func TestIsValidFixedVersion(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		version string
		isFixed bool
	}{
		{"fixed version", "10.1.1", true},
		{"fixed beta version", "10.1.1-beta", true},
		{"fixed beta patch", "10.1.1-beta.1", true},
		{"fixed version label zzz", "1.0.1-zzz", true},
		{"fixed version RC with label", "1.0.1-rc.10", true},
		{"fixed version RC with label 2", "1.0.1-rc.2", true},
		{"fixed version with label open", "1.0.1-open", true},
		{"fixed version alpha", "1.0.1-alpha2", true},
		{"fixed version RC with label aaa", "1.0.1-aaa", true},
		{"fixed version range", "[1.0]", true},
		{"version as variable", "$(ComponentDetectionPackageVersion)", true},
		{"version range with inclusive min", "[1.0,)", false},
		{"version range with inclusive min without brackets", "1.0", false},
		{"version range with exclusive min", "(1.0,)", false},
		{"version range with inclusive max", "(,1.0]", false},
		{"version range with exclusive max", "[,1.0)", false},
		{"Exact range, inclusive", "[1.0,2.0]", false},
		{"Exact range, exclusive", "(1.0,2.0)", false},
		{"Mixed inclusive minimum and exclusive maximum version", "(1.0,2.0)", false},
		{"invalid", "(1.0)", false},
	}
	for _, tt := range tests {
		tt := tt // Re-initializing variable so it is not changed while executing the closure below
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			isFixed := isValidFixedVersion(tt.version)
			if tt.isFixed != isFixed {
				t.Errorf("expected %v. Got %v", tt.isFixed, isFixed)
			}
		})
	}
}

func TestAnalyseCentralPackageManagementPinned(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                 string
		filename             string
		pinnedDependencies   int
		unpinnedDependencies int
		expectedError        bool
	}{
		{
			name:                 "Pinned dependencies",
			filename:             "./testdata/Directory.Pinned.packages.props",
			pinnedDependencies:   2,
			unpinnedDependencies: 0,
			expectedError:        false,
		},
		{
			name:                 "Pinned multiple dependencies",
			filename:             "./testdata/Directory.PinnedMultipleGroups.packages.props",
			pinnedDependencies:   3,
			unpinnedDependencies: 0,
			expectedError:        false,
		},
		{
			name:                 "Unpinned CPM false",
			filename:             "./testdata/Directory.CPMFalse.packages.props",
			pinnedDependencies:   0,
			unpinnedDependencies: 1,
			expectedError:        false,
		},
		{
			name:                 "Unpinned CPM undeclared",
			filename:             "./testdata/Directory.Undeclared.packages.props",
			pinnedDependencies:   0,
			unpinnedDependencies: 1,
			expectedError:        false,
		},
		{
			name:                 "Unpinned version undeclared",
			filename:             "./testdata/Directory.UndeclaredVersions.packages.props",
			pinnedDependencies:   1,
			unpinnedDependencies: 1,
			expectedError:        false,
		},
		{
			name:                 "Unpinned version range",
			filename:             "./testdata/Directory.UnpinnedVersions.packages.props",
			pinnedDependencies:   1,
			unpinnedDependencies: 1,
			expectedError:        false,
		},
		{
			name:                 "Unpinned version range in second group",
			filename:             "./testdata/Directory.UnpinnedMultipleGroups.packages.props",
			pinnedDependencies:   2,
			unpinnedDependencies: 1,
			expectedError:        false,
		},
	}
	for _, tt := range tests {
		tt := tt // Re-initializing variable so it is not changed while executing the closure below
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var content []byte
			var err error
			content, err = os.ReadFile(tt.filename)
			if err != nil {
				t.Errorf("cannot read file: %v", err)
			}
			var cpmDeps []checker.Dependency
			err = analyseCentralPackageManagementPinned(tt.filename, content, &cpmDeps)
			if err != nil {
				if !tt.expectedError {
					t.Errorf("unexpected error: %v", err)
				}
				return
			}
			pinned, unpinned := 0, 0
			for _, dep := range cpmDeps {
				if *dep.Pinned {
					pinned++
				} else {
					unpinned++
				}
			}
			if pinned != tt.pinnedDependencies {
				t.Errorf("expected %v pinned dependencies. Got %v", tt.pinnedDependencies, pinned)
			}
			if unpinned != tt.unpinnedDependencies {
				t.Errorf("expected %v unpinned dependencies. Got %v", tt.unpinnedDependencies, unpinned)
			}
		})
	}
}
