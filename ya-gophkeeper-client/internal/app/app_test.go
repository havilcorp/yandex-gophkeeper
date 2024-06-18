package app

import (
	"testing"
)

func Test_getCert(t *testing.T) {
	type args struct {
		ca  string
		crt string
		key string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "good",
			args: args{
				ca:  "../../tls/ca.crt",
				crt: "../../tls/server.crt",
				key: "../../tls/server.key",
			},
			wantErr: false,
		},
		{
			name: "err1",
			args: args{
				ca:  "../../tls/notfound.crt",
				crt: "../../tls/server.crt",
				key: "../../tls/server.key",
			},
			wantErr: true,
		},
		{
			name: "err2",
			args: args{
				ca:  "../../tls/ca.crt",
				crt: "../../tls/notfound.crt",
				key: "../../tls/server.key",
			},
			wantErr: true,
		},
		{
			name: "err3",
			args: args{
				ca:  "../../tls/ca.crt",
				crt: "../../tls/server.crt",
				key: "../../tls/notfound.key",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := getCert(tt.args.ca, tt.args.crt, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("getCert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
