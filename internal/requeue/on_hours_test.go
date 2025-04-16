package requeue

import (
	"testing"
	"time"
)

type fakeTime struct {
	ret time.Time
}

func (t *fakeTime) Now() time.Time {
	return t.ret
}

func TestOnHoursRequeue_Requeue(t *testing.T) {
	type fields struct {
		fromHour   int
		toHour     int
		timeSource timeSource
	}
	type args struct {
		duration time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   time.Duration
	}{
		{
			name: "within the operating hours",
			fields: fields{
				fromHour:   10,
				toHour:     22,
				timeSource: &fakeTime{ret: time.Date(2025, 05, 30, 15, 0, 0, 0, time.UTC)},
			},
			args: args{
				duration: 30 * time.Minute,
			},
			want: 30 * time.Minute,
		},
		{
			name: "after operating hours",
			fields: fields{
				fromHour:   10,
				toHour:     22,
				timeSource: &fakeTime{ret: time.Date(2025, 05, 30, 21, 0, 0, 0, time.UTC)},
			},
			args: args{
				duration: 90 * time.Minute,
			},
			want: 13 * time.Hour,
		},
		{
			name: "before operating hours",
			fields: fields{
				fromHour:   10,
				toHour:     22,
				timeSource: &fakeTime{ret: time.Date(2025, 05, 30, 9, 0, 0, 0, time.UTC)},
			},
			args: args{
				duration: 30 * time.Minute,
			},
			want: 1 * time.Hour,
		},
		{
			name: "after operating hours, 25h",
			fields: fields{
				fromHour:   10,
				toHour:     22,
				timeSource: &fakeTime{ret: time.Date(2025, 05, 30, 15, 0, 0, 0, time.UTC)},
			},
			args: args{
				duration: 25 * time.Hour,
			},
			want: 25 * time.Hour,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &OnHoursRequeue{
				fromHour:   tt.fields.fromHour,
				toHour:     tt.fields.toHour,
				timeSource: tt.fields.timeSource,
			}
			if got := w.Requeue(tt.args.duration); got != tt.want {
				t.Errorf("Requeue() = %v, want %v", got, tt.want)
			}
		})
	}
}
