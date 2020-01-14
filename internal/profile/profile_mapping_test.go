package profile

import (
	"reflect"
	"testing"

	"github.com/reviz0r/golang-layout/internal/profile/models"
	"github.com/reviz0r/golang-layout/pkg/profile"
)

func Test_userFromProto(t *testing.T) {
	type args struct {
		in *profile.User
	}
	tests := []struct {
		name string
		args args
		want *models.User
	}{
		{
			name: "user is nil",
			args: args{in: nil},
			want: nil,
		},
		{
			name: "user is not nil",
			args: args{in: &profile.User{Name: "user", Email: "user@example.com"}},
			want: &models.User{Name: "user", Email: "user@example.com"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := userFromProto(tt.args.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userFromProto() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userToProto(t *testing.T) {
	type args struct {
		in *models.User
	}
	tests := []struct {
		name string
		args args
		want *profile.User
	}{
		{
			name: "user is nil",
			args: args{in: nil},
			want: nil,
		},
		{
			name: "user is not nil",
			args: args{in: &models.User{Name: "user", Email: "user@example.com"}},
			want: &profile.User{Name: "user", Email: "user@example.com"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := userToProto(tt.args.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userToProto() = %v, want %v", got, tt.want)
			}
		})
	}
}
