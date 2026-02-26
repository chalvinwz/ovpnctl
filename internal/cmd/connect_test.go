package cmd

import (
	"strings"
	"testing"
)

func TestReadOTP(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "valid otp", input: "123456\n", want: "123456"},
		{name: "trim spaces", input: "  654321  \n", want: "654321"},
		{name: "empty otp", input: "   \n", wantErr: true},
		{name: "no input", input: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readOTP(strings.NewReader(tt.input))
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("readOTP() = %q, want %q", got, tt.want)
			}
		})
	}
}
