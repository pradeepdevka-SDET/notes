package main

import "testing"

func TestValidateTitle(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		wantErr bool
	}{
		{name: "valid title", title: "buy milk", wantErr: false},
		{name: "empty title", title: "", wantErr: true},
		{name: "whitespace only", title: "  ", wantErr: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validateTitile(tc.title)
			if tc.wantErr && err == nil {
				t.Errorf("expected an error for %q, but got nil", tc.title)
			}
			if !tc.wantErr && err != nil {
				t.Errorf("expected no error for %q, but got: %v", tc.title, err)
			}
		})
	}
}
