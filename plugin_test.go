package main

import (
	"reflect"
	"testing"
)

func Test_getEnvironment(t *testing.T) {
	type args struct {
		Environment []string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "normal",
			args: args{
				Environment: []string{"a=b", "c=d=e"},
			},
			want: map[string]string{
				"a": "b",
				"c": "d=e",
			},
		},
		{
			name: "expect",
			args: args{
				Environment: []string{"a=b", "c=d=e", "d"},
			},
			want: map[string]string{
				"a": "b",
				"c": "d=e",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getEnvironment(tt.args.Environment); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getEnvironment() = %v, want %v", got, tt.want)
			}
		})
	}
}
