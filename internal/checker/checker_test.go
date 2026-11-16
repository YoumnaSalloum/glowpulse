package checker

import "testing"

func TestStatusFromCode(t *testing.T) {
	tests := []struct {
		code int
		want Status
	}{
		{200, StatusUp},
		{301, StatusUp},
		{404, StatusDegraded},
		{500, StatusDown},
	}

	for _, tt := range tests {
		got := statusFromCode(tt.code)
		if got != tt.want {
			t.Errorf("statusFromCode(%d) = %s, want %s", tt.code, got, tt.want)
		}
	}
}
