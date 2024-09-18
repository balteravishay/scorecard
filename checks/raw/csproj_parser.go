package raw

import (
	"encoding/xml"
)

type RestoreLockedModePropertyGroup struct {
	XMLName           xml.Name `xml:"PropertyGroup"`
	RestoreLockedMode bool     `xml:"RestoreLockedMode"`
}

type CSProjProject struct {
	XMLName        xml.Name                         `xml:"Project"`
	PropertyGroups []RestoreLockedModePropertyGroup `xml:"PropertyGroup"`
}

func isRestoreLockedModeEnabled(content []byte) (error, bool) {
	var project CSProjProject

	err := xml.Unmarshal(content, &project)
	if err != nil {
		return errInvalidCsProjFile, false
	}

	for _, propertyGroup := range project.PropertyGroups {
		if propertyGroup.RestoreLockedMode {
			return nil, true
		}
	}

	return nil, false
}
