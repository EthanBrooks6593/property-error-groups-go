package main

import "testing"

func TestMaintenanceGrouping(t *testing.T) {
	cases := []struct{ name, id, want string }{
		{"same request", "m-17", "maintenance:m-17"}, {"another request", "m-18", "maintenance:m-18"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := ErrorGroup("maintenance", tc.id); got != tc.want {
				t.Fatalf("got %q want %q", got, tc.want)
			}
		})
	}
}
