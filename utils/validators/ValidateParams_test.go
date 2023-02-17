package validators

import (
	"fiber-api/utils/types"
	"fiber-api/utils/types/requests"
	"reflect"
	"testing"
)

func generateOutput(outputs []types.ValidationError) []*types.ValidationError {
	var errors []*types.ValidationError
	for _, output := range outputs {
		errors = append(errors, &output)
	}
	return errors
}

func TestIsValidParams(t *testing.T) {
	type args struct {
		params interface{}
	}
	tests := []struct {
		name string
		args args
		want []*types.ValidationError
	}{
		{"valid email and password", args{params: requests.TeacherRegistrationRequest{Email: "test@gmail.com", Password: "Passw0rd!", Age: 20, Name: "Teacher1"}}, nil},
		{"invalid email", args{params: requests.TeacherRegistrationRequest{Email: "testgmail.com", Password: "Passw0rd!", Age: 20, Name: "Teacher1"}}, generateOutput([]types.ValidationError{{Field: "Email", Tag: "email", Value: ""}})},
		{"invalid password (len < 8)", args{params: requests.TeacherRegistrationRequest{Email: "test@gmail.com", Password: "Pw0!", Age: 20, Name: "Teacher1"}}, generateOutput([]types.ValidationError{{Field: "Password", Tag: "min", Value: "8"}})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidParams(tt.args.params); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IsValidParams() = %v, want %v", got, tt.want)
			}
		})
	}
}
