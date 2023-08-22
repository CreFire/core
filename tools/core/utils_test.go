package core

import "testing"

func TestCapitalizeFirstLetter(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{s: "aAa"},
			want: "AAa",
		},
		{
			name: "2",
			args: args{s: "Q"},
			want: "Q",
		},
		{
			name: "3",
			args: args{s: "1nm"},
			want: "1nm",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CapitalizeFirstLetter(tt.args.s); got != tt.want {
				t.Errorf("CapitalizeFirstLetter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReverseString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "3",
			args: args{s: "1nm"},
			want: "1nm",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReverseString(tt.args.s); got != tt.want {
				t.Errorf("ReverseString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat64bits(t *testing.T) {
	type args struct {
		f float64
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{
			name: "3",
			args: args{f: 1.00},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Float64bits(tt.args.f); got != tt.want {
				t.Errorf("Float64bits() = %v, want %v", got, tt.want)
			}
		})
	}
}
