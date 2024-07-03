// Copyright 2024 OpenSSF Scorecard Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	pyaml "github.com/ossf/scorecard/v5/internal/probes/yaml"
)

func printField(w io.Writer, name string, value any) {
	// some fields have extra newlines we can get rid of
	if v, ok := value.(string); ok {
		value = strings.TrimSpace(v)
	}

	fmt.Fprint(w, "**", name, "**: ", value, "\n\n")
}

func printProbe(w io.Writer, p *pyaml.Probe) {
	// short, motivation, implementation, outcome, remediation, ecosystem
	fmt.Fprint(w, "\n"+"## "+p.ID+"\n\n")
	printField(w, "Lifecycle", p.Lifecycle)
	printField(w, "Description", p.Short)
	printField(w, "Motivation", p.Motivation)
	printField(w, "Implementation", p.Implementation)
	printField(w, "Outcomes", "\n\n"+strings.Join(p.Outcomes, "\n"))
	// TODO remediation
	// TODO ecosystem
}

func walk(path string, d fs.DirEntry, err error) error {
	// no special handling of errors, we can stop now
	if err != nil {
		return err
	}

	if !strings.EqualFold(filepath.Base(path), "def.yml") {
		return nil
	}
	var probe pyaml.Probe

	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read probe definition: %w", err)
	}
	err = yaml.Unmarshal(content, &probe)
	if err != nil {
		return fmt.Errorf("parse yaml: %w", err)
	}
	printProbe(os.Stdout, &probe)
	return nil
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %s <probe directory>", os.Args[0])
	}
	probeDir := os.Args[1]

	fmt.Fprint(os.Stdout, `<!-- Do not edit this file manually! Edit the individual probe's def.yml instead. -->
# Probe Documentation

This page describes each Scorecard probe in detail, including description, motivation,
and outcomes. If you have ideas for additions or new detection techniques,
please [contribute](../CONTRIBUTING.md)!
`)
	if err := filepath.WalkDir(probeDir, walk); err != nil {
		log.Fatal(err)
	}
}