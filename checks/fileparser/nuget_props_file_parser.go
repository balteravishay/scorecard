package fileparser

import (
	"encoding/xml"
	"regexp"
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

type NugetPackage struct {
	Name    string
	Version string
	IsFixed bool
}

type CentralPackageManagementConfig struct {
	PackageVersions []NugetPackage
	IsCPMEnabled    bool
}

func GetCentralPackageManagementConfig(path string, content []byte) (CentralPackageManagementConfig, error) {
	var project DirectoryPropsProject

	err := xml.Unmarshal(content, &project)
	if err != nil {
		return CentralPackageManagementConfig{}, errInvalidPropsFile
	}

	cpmConfig := CentralPackageManagementConfig{
		IsCPMEnabled: isCentralPackageManagementEnabled(&project),
	}

	if cpmConfig.IsCPMEnabled {
		cpmConfig.PackageVersions = extractNugetPackages(&project)
	}

	return cpmConfig, nil
}

func isCentralPackageManagementEnabled(project *DirectoryPropsProject) bool {
	for _, propertyGroup := range project.PropertyGroups {
		if propertyGroup.ManagePackageVersionsCentrally {
			return true
		}
	}

	return false
}

func extractNugetPackages(project *DirectoryPropsProject) []NugetPackage {
	var nugetPackages []NugetPackage
	for _, itemGroup := range project.ItemGroups {
		for _, packageVersion := range itemGroup.PackageVersion {
			nugetPackages = append(nugetPackages, NugetPackage{
				Name:    packageVersion.Include,
				Version: packageVersion.Version,
				IsFixed: isValidFixedVersion(packageVersion.Version),
			})
		}
	}
	return nugetPackages
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
