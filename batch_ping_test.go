package ping

import (
	"testing"
)

func TestNewBatchPinger(t *testing.T) {
	type args struct {
		addrs      []string
		privileged bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "newBatchPinger",
			args: args{
				addrs: []string{"baidu.com"},
                privileged: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewBatchPinger(tt.args.addrs, tt.args.privileged)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBatchPinger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNewBatchPinger_multiAddr(t *testing.T) {
	type args struct {
		addrs      []string
		privileged bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "multi_addr",
			args: args{
				addrs: []string{
					"39.156.69.1",
					"39.156.69.2",
					"39.156.69.3",
					"39.156.69.4",
					"39.156.69.5",
					"39.156.69.6",
					"39.156.69.7",
					"39.156.69.8",
					"39.156.69.9",
					"39.156.69.10",
					"39.156.69.11",
					"39.156.69.12",
					"39.156.69.13",
					"39.156.69.14",
					"39.156.69.15",
					"39.156.69.16",
					"39.156.69.17",
					"39.156.69.18",
					"39.156.69.19",
					"39.156.69.20",
					"39.156.69.21",
					"39.156.69.22",
					"39.156.69.23",
					"39.156.69.24",
					"39.156.69.25",
					"39.156.69.26",
					"39.156.69.27",
					"39.156.69.28",
					"39.156.69.29",
					"39.156.69.30",
					"39.156.69.31",
					"39.156.69.32",
					"39.156.69.33",
					"39.156.69.34",
					"39.156.69.35",
					"39.156.69.36",
					"39.156.69.37",
					"39.156.69.38",
					"39.156.69.39",
					"39.156.69.40",
					"39.156.69.41",
					"39.156.69.42",
					"39.156.69.43",
					"39.156.69.44",
					"39.156.69.45",
					"39.156.69.46",
					"39.156.69.47",
					"39.156.69.48",
					"39.156.69.49",
					"39.156.69.50"},
                    privileged: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			batchPinger, err := NewBatchPinger(tt.args.addrs, tt.args.privileged)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBatchPinger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			err = batchPinger.Run()
			if (err != nil) != tt.wantErr {
				t.Errorf("multi ping  error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNewBatchPinger_ipv6(t *testing.T) {
	type args struct {
		addrs      []string
		privileged bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "newBatchPinger",
			args: args{
				addrs: []string{"2400:da00:2::29"},
                privileged: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bp, err := NewBatchPinger(tt.args.addrs, tt.args.privileged)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBatchPinger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			err = bp.Run()
			if (err != nil) != tt.wantErr {
				t.Errorf("ping ipv6 error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
