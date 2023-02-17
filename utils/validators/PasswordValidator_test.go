package validators

import (
	"testing"
)

func Test_IsValidPassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"valid password", args{"Password1!"}, true},
		{"invalid password (no uppercase no special char)", args{"password1"}, false},
		{"invalid password (no lowercase no special char)", args{"PASSWORD1"}, false},
		{"invalid password (no lowercase no uppercase)", args{"23213414132451!"}, false},
		{"invalid password (no lowercase no uppercase no special char)", args{"23213414132451"}, false},
		{"invalid password (no lowercase no uppercase no digit)", args{"!@#$%^&*()_+"}, false},
		{"invalid password (no lowercase no digit no special char)", args{"PASSWORD"}, false},
		{"invalid password (len < 8)", args{"Pas12!"}, false},
		{"invalid password (len > 20)", args{"Password123456789012$4567890"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidPassword(tt.args.password); got != tt.want {
				t.Errorf("IsValidPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
