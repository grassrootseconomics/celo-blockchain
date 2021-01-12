package uptime

import (
	"testing"

	"github.com/ethereum/go-ethereum/consensus/istanbul"
)

func TestEpochSizeIsConsistentWithSkippedBlock(t *testing.T) {
	if istanbul.MinEpochSize <= BlocksToSkipAtEpochEnd {
		t.Fatalf("Constant MinEpochSize MUST BE greater than BlocksToSkipAtEpochEnd (%d, %d) ", istanbul.MinEpochSize, BlocksToSkipAtEpochEnd)
	}
}

func TestMonitoringWindow(t *testing.T) {
	type args struct {
		epochNumber        uint64
		epochSize          uint64
		lookbackWindowSize uint64
	}
	tests := []struct {
		name      string
		args      args
		wantStart uint64
		wantEnd   uint64
	}{
		{"monitoringWindow on first epoch", args{1, 10, 2}, 2, 8},
		{"monitoringWindow on second epoch", args{2, 10, 2}, 12, 18},
		{"lookback window too big", args{1, 10, 10}, 10, 8},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := MonitoringWindow(tt.args.epochNumber, tt.args.epochSize, tt.args.lookbackWindowSize)
			if w.Start != tt.wantStart {
				t.Errorf("MonitoringWindow() got = %v, wantStart %v", w.Start, tt.wantStart)
			}
			if w.End != tt.wantEnd {
				t.Errorf("MonitoringWindow() got1 = %v, wantEnd %v", w.End, tt.wantEnd)
			}
		})
	}
}

func TestComputeLookbackWindow(t *testing.T) {
	constant := func(value uint64) func() (uint64, error) {
		return func() (uint64, error) { return value, nil }
	}

	type args struct {
		epochSize             uint64
		defaultLookbackWindow uint64
		isDonut               bool
		getLookbackWindow     func() (uint64, error)
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{"returns default if Donut is not active", args{100, 20, false, constant(24)}, 20},
		{"returns safe minimun if configured is below", args{100, 20, true, constant(10)}, 12},
		{"returns safe maximum if configured is above", args{1000, 20, true, constant(800)}, 720},
		{"returns epochSize-2 if configured is above", args{100, 20, true, constant(99)}, 98},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ComputeLookbackWindow(tt.args.epochSize, tt.args.defaultLookbackWindow, tt.args.isDonut, tt.args.getLookbackWindow)
			if got != tt.want {
				t.Errorf("ComputeLookbackWindow() = %v, want %v", got, tt.want)
			}
		})
	}
}
