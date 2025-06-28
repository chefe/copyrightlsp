package analysis

import "testing"

func TestMatchesTemplateLine(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		line     string
		template string
		want     bool
	}{
		{
			name:     "empty line and empty template",
			line:     "",
			template: "",
			want:     true,
		},
		{
			name:     "different line and template",
			line:     "foo",
			template: "bar",
			want:     false,
		},
		{
			name:     "matches template with year placeholder",
			line:     "# Copyright (C) 2024 AUTHOR",
			template: "# Copyright (C) {year} AUTHOR",
			want:     true,
		},
		{
			name:     "year is not a number",
			line:     "# Copyright (C) asdf AUTHOR",
			template: "# Copyright (C) {year} AUTHOR",
			want:     false,
		},
		{
			name:     "line is longer then template",
			line:     "# Copyright (C) 2024 AND MORE",
			template: "# Copyright (C) {year}",
			want:     false,
		},
		{
			name:     "line is shorter then template",
			line:     "# Copyright (C) 2024",
			template: "# Copyright (C) {year} AUTHOR",
			want:     false,
		},
		{
			name:     "template with multiple year placeholder",
			line:     "# Copyright 2020 (C) 2022 AUTHOR 2024",
			template: "# Copyright {year} (C) {year} AUTHOR {year}",
			want:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := matchesTemplateLine(test.line, test.template)
			if got != test.want {
				t.Fatalf("expected: %t, got: %t", test.want, got)
			}
		})
	}
}

func TestContainsTemplateLines(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		lines    []string
		template []string
		want     bool
	}{
		{
			name:     "nil for lines and template",
			lines:    nil,
			template: nil,
			want:     false,
		},
		{
			name:     "empty slice for line and template",
			lines:    []string{},
			template: []string{},
			want:     false,
		},
		{
			name:     "single line template with one matching line",
			lines:    []string{"# Copyright (C) 2024 AUTHOR"},
			template: []string{"# Copyright (C) {year} AUTHOR"},
			want:     true,
		},
		{
			name:     "single line template with second line matching",
			lines:    []string{"#!/bin/sh", "# Copyright (C) 2024 AUTHOR"},
			template: []string{"# Copyright (C) {year} AUTHOR"},
			want:     true,
		},
		{
			name:     "single line template with no line matching",
			lines:    []string{"#!/bin/sh", "echo 'test'"},
			template: []string{"# Copyright (C) {year} AUTHOR"},
			want:     false,
		},
		{
			name:     "single line template with no lines",
			lines:    []string{},
			template: []string{"# Copyright (C) {year} AUTHOR"},
			want:     false,
		},
		{
			name:     "multi line template with match starting at line zero",
			lines:    []string{"/*", " * Copyright (C) 2024 AUTHOR", " */"},
			template: []string{"/*", " * Copyright (C) {year} AUTHOR", " */"},
			want:     true,
		},
		{
			name:     "multi line template with match starting at another line",
			lines:    []string{"// Another comment", "", "/*", " * Copyright (C) 2024 AUTHOR", " */"},
			template: []string{"/*", " * Copyright (C) {year} AUTHOR", " */"},
			want:     true,
		},
		{
			name:     "multi line template with only some lines matching",
			lines:    []string{"/*", " * Copyright (C) 2024 AUTHOR"},
			template: []string{"/*", " * Copyright (C) {year} AUTHOR", " */"},
			want:     false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := containsTemplateLines(test.lines, test.template)
			if got != test.want {
				t.Fatalf("expected: %t, got: %t", test.want, got)
			}
		})
	}
}

func TestContainsCopyrightString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		content     string
		template    []string
		searchRange uint8
		want        bool
	}{
		{
			name:        "empty string and nil template",
			content:     "",
			template:    nil,
			searchRange: 0,
			want:        false,
		},
		{
			name:        "empty string and empty template",
			content:     "",
			template:    []string{},
			searchRange: 0,
			want:        false,
		},
		{
			name:        "match on first line",
			content:     "# Copyright (C) 2024 AUTHOR",
			template:    []string{"# Copyright (C) {year} AUTHOR"},
			searchRange: 0,
			want:        true,
		},
		{
			name:        "match on second line",
			content:     "#!/bin/sh\n# Copyright (C) 2024 AUTHOR",
			template:    []string{"# Copyright (C) {year} AUTHOR"},
			searchRange: 1,
			want:        true,
		},
		{
			name:        "multiline match starting on last line of search range",
			content:     "\n\n/*\n * Copyright (C) 2024 AUTHOR\n */",
			template:    []string{"/*", " * Copyright (C) {year} AUTHOR", " */"},
			searchRange: 2,
			want:        true,
		},
		{
			name:        "match starting outside of search range",
			content:     "\n\n\n# Copyright (C) 2024 AUTHOR",
			template:    []string{"# Copyright (C) {year} AUTHOR"},
			searchRange: 2,
			want:        false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := ContainsCopyrightString(test.content, test.template, test.searchRange)
			if got != test.want {
				t.Fatalf("expected: %t, got: %t", test.want, got)
			}
		})
	}
}
