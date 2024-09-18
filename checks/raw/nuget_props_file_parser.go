package raw

import (
	"encoding/xml"
	"fmt"
	"regexp"

	"github.com/ossf/scorecard/v5/checker"
	"github.com/ossf/scorecard/v5/finding"
)

type CPMPropertyGroup struct {
	XMLName                        xml.Name `xml:"PropertyGroup"`
	ManagePackageVersionsCentrally bool     `xml:"ManagePackageVersionsCentrally"`
}

type PackageVersionItemGroup struct {
	XMLName        xml.Name         `xml:"ItemGroup"`
	PackageVersion []packageVersion `xml:"PackageVersion"`
}

type packageVersion struct {
	XMLName xml.Name `xml:"PackageVersion"`
	Version string   `xml:"Version,attr"`
	Include string   `xml:"Include,attr"`
}

type DirectoryPropsProject struct {
	XMLName        xml.Name                  `xml:"Project"`
	PropertyGroups []CPMPropertyGroup        `xml:"PropertyGroup"`
	ItemGroups     []PackageVersionItemGroup `xml:"ItemGroup"`
}

func analyseCentralPackageManagementPinned(path string, content []byte, pdata *[]checker.Dependency) error {
	var project DirectoryPropsProject

	err := xml.Unmarshal(content, &project)
	if err != nil {
		return errInvalidPropsFile
	}
	for _, propertyGroup := range project.PropertyGroups {
		if propertyGroup.ManagePackageVersionsCentrally {
			dependency := checker.Dependency{
				Location: &checker.File{
					Path:      path,
					Type:      finding.FileTypeSource,
					Offset:    1,
					EndOffset: 1,
					Snippet:   "<ManagePackageVersionsCentrally>true</ManagePackageVersionsCentrally>",
				},
				Pinned: asBoolPointer(true),
				Type:   checker.DependencyUseTypeNugetCommand,
			}
			*pdata = append(*pdata, dependency)
		}
	}
	if len(*pdata) == 0 {
		dependency := checker.Dependency{
			Location: &checker.File{
				Path:      path,
				Type:      finding.FileTypeSource,
				Offset:    1,
				EndOffset: 1,
				Snippet:   "<ManagePackageVersionsCentrally>true</ManagePackageVersionsCentrally>",
			},
			Pinned: asBoolPointer(false),
			Type:   checker.DependencyUseTypeNugetCommand,
		}
		*pdata = append(*pdata, dependency)
		return nil
	}
	for _, itemGroup := range project.ItemGroups {
		for _, packageVersion := range itemGroup.PackageVersion {
			pinned := isValidFixedVersion(packageVersion.Version)
			dependency := checker.Dependency{
				Location: &checker.File{
					Path:      path,
					Type:      finding.FileTypeSource,
					Offset:    1,
					EndOffset: 1,
					Snippet: fmt.Sprintf("<PackageVersion Include=\"%s\" Version=\"%s\" />",
						packageVersion.Include, packageVersion.Version),
				},
				Pinned: asBoolPointer(pinned),
				Type:   checker.DependencyUseTypeNugetCommand,
			}
			*pdata = append(*pdata, dependency)
		}
	}
	return nil
}

// isValidFixedVersion checks if the version string is a valid, fixed version.
// more on version numbers here: https://learn.microsoft.com/en-us/nuget/concepts/package-versioning?tabs=semver20sort
func isValidFixedVersion(version string) bool {
	// Define the regular expression for a valid fixed version
	// ^ asserts the start of the string
	// (\d+)\.(\d+)\.(\d+) matches major.minor.patch version (e.g., 1.0.1)
	// (-[a-zA-Z]+(\.\d+)?|[a-zA-Z]+\d+)? matches optional pre-release tag
	//   -[a-zA-Z]+(\.\d+)? handles pre-release versions with dots (e.g., -beta.12, -rc.10)
	//   -[a-zA-Z]+\d+ handles versions like -alpha2 and -alpha10
	// \[\d+\.\d+\] or \[\d+\.\d+\.\d+\] matches cases like [1.0] and [1.0.1]
	// \$\(ComponentDetectionPackageVersion\) matches special case like $(ComponentDetectionPackageVersion)
	pattern := `^(\d+)\.(\d+)\.(\d+)(-[a-zA-Z]+(\.\d+)?|-[a-zA-Z]+\d*)?$|^\[\d+\.\d+\]$|^\[\d+\.\d+\.\d+\]$|^\$\(.+\)$`

	// Compile the regex
	re := regexp.MustCompile(pattern)

	// Check if the version matches the regex
	return re.MatchString(version)
}
