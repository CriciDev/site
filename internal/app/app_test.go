package app

import "testing"

func TestDisplayURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		addr string
		want string
	}{
		{
			name: "port only uses localhost",
			addr: ":8080",
			want: "http://localhost:8080",
		},
		{
			name: "ipv4 any uses localhost",
			addr: "0.0.0.0:8080",
			want: "http://localhost:8080",
		},
		{
			name: "ipv6 any uses localhost",
			addr: "[::]:8080",
			want: "http://localhost:8080",
		},
		{
			name: "ipv4 localhost stays unchanged",
			addr: "127.0.0.1:3000",
			want: "http://127.0.0.1:3000",
		},
		{
			name: "hostname stays unchanged",
			addr: "localhost:8080",
			want: "http://localhost:8080",
		},
		{
			name: "invalid address falls back to raw value",
			addr: "localhost",
			want: "http://localhost",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := displayURL(tc.addr)
			if got != tc.want {
				t.Fatalf("displayURL(%q) = %q, want %q", tc.addr, got, tc.want)
			}
		})
	}
}
