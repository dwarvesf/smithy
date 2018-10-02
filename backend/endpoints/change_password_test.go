package endpoints

import "testing"

func Test_checkPassword(t *testing.T) {
	type args struct {
		pwd string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: TooShort,
			want: TooShort,
		},
		{
			name: VeryWeak,
			args: args{
				"abc",
			},
			want: VeryWeak,
		},
		{
			name: Weak,
			args: args{
				"abcDeEa",
			},
			want: Weak,
		},
		{
			name: Good,
			args: args{
				"abcDeEa@",
			},
			want: Good,
		},
		{
			name: Strong,
			args: args{
				"abcDeEa@1",
			},
			want: Strong,
		},
		{
			name: VeryStrong,
			args: args{
				"abcDeEa@12!",
			},
			want: VeryStrong,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkPassword(tt.args.pwd); got != tt.want {
				t.Errorf("checkPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
