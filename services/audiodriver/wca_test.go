package audiodriver_test

import (
	"context"
	"game-radar/models"
	"game-radar/services/audiodriver"
	"testing"
)

type WCAPeaksReaderWriter interface {
	Read() []float32
	Write([]float32)
	Close()
}

type wcaPeaksReaderWriter struct {
	peaks chan []float32
}

func newWcaPeaksReaderWriter() *wcaPeaksReaderWriter {
	return &wcaPeaksReaderWriter{
		peaks: make(chan []float32, 10),
	}
}

func (wprw *wcaPeaksReaderWriter) Write(peaks []float32) {
	wprw.peaks <- peaks
}

func (wprw *wcaPeaksReaderWriter) Read() []float32 {
	peaks := <-wprw.peaks
	return peaks
}

func (wprw *wcaPeaksReaderWriter) Close() {
	if _, ok := <-wprw.peaks; ok {
		close(wprw.peaks)
	}
}

func TestWCA_LoopbackDevice(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		device models.DeviceSpeakers
	}{
		{
			name: "default",
			device: models.DeviceSpeakers{
				ID:            "{0.0.0.00000000}.{404fcbe5-0c6b-438b-98cf-dbf3d28eae7a}",
				Name:          "Corsair",
				FormatPCM:     32,
				Rate:          48000,
				Channels:      8,
				DefaultPeriod: 100000,
				MinimumPeriod: 30000,
				Latency:       10000000,
				IsDefault:     true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := audiodriver.NewWCA(nil)
			ctx := context.Background()

			wprw := newWcaPeaksReaderWriter()

			go func() {
				gotErr := w.LoopbackDeviceSpeakers(ctx, tt.device, wprw)
				if gotErr != nil {
					wprw.Close()
					// t.Fatalf("loopback err: %v", gotErr)
				}
			}()

			t.Log("started")

			// wait for init
			// time.Sleep(1 * time.Second)
			for range 5 {
				t.Log("peaks", wprw.Read())
			}
		})
	}
}
