package builtin

import (
	"testing"
)

func TestVerifiedDealWeightMultiplierAt(t *testing.T) {
	startEpoch := int64(1000)
	transitionEpochs := int64(100)

	tests := []struct {
		name     string
		epoch    int64
		start    int64
		duration int64
		want     int64
	}{
		{"no transition (duration=0)", 5000, 1000, 0, 100},
		{"before start", 500, startEpoch, transitionEpochs, 100},
		{"at start", startEpoch, startEpoch, transitionEpochs, 100},
		{"1% through", startEpoch + 1, startEpoch, transitionEpochs, 99},
		{"10% through", startEpoch + 10, startEpoch, transitionEpochs, 91},
		{"25% through", startEpoch + 25, startEpoch, transitionEpochs, 77},
		{"50% through", startEpoch + 50, startEpoch, transitionEpochs, 55},
		{"75% through", startEpoch + 75, startEpoch, transitionEpochs, 32},
		{"99% through", startEpoch + 99, startEpoch, transitionEpochs, 10},
		{"at end", startEpoch + transitionEpochs, startEpoch, transitionEpochs, 10},
		{"after end", startEpoch + 200, startEpoch, transitionEpochs, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := VerifiedDealWeightMultiplierAt(tt.epoch, tt.start, tt.duration)
			if got.Int64() != tt.want {
				t.Errorf("VerifiedDealWeightMultiplierAt(%d, %d, %d) = %d, want %d",
					tt.epoch, tt.start, tt.duration, got.Int64(), tt.want)
			}
		})
	}

	// Verify monotonicity: multiplier should never increase as epoch advances
	prev := int64(100)
	for epoch := startEpoch; epoch <= startEpoch+transitionEpochs+10; epoch++ {
		val := VerifiedDealWeightMultiplierAt(epoch, startEpoch, transitionEpochs).Int64()
		if val > prev {
			t.Errorf("monotonicity violation at epoch %d: %d > %d", epoch, val, prev)
		}
		prev = val
	}
}
