// Package testing holds tools for testing documentation code.
package golden

import (
	"testing"
	"time"
)

func Test_withinToleranceTime(t *testing.T) {
	type args struct {
		a         time.Time
		b         time.Time
		tolerance time.Duration
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "within tolerance b > a",
			args: args{
				a:         time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				b:         time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC),
				tolerance: 10 * time.Second,
			},
			want: true,
		},
		{
			name: "within tolerance a > b",
			args: args{
				a:         time.Date(2020, 1, 1, 0, 0, 4, 0, time.UTC),
				b:         time.Date(2020, 1, 1, 0, 0, 3, 0, time.UTC),
				tolerance: 10 * time.Second,
			},
			want: true,
		},
		{
			name: "within tolerance not ok a > b",
			args: args{
				a:         time.Date(2020, 1, 1, 0, 10, 4, 0, time.UTC),
				b:         time.Date(2020, 1, 1, 0, 0, 3, 0, time.UTC),
				tolerance: 10 * time.Second,
			},
			want: false,
		},
		{
			name: "within tolerance not ok b > a",
			args: args{
				a:         time.Date(2020, 1, 1, 0, 0, 4, 0, time.UTC),
				b:         time.Date(2020, 1, 1, 0, 2, 3, 0, time.UTC),
				tolerance: 10 * time.Second,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := withinToleranceTime(tt.args.a, tt.args.b, tt.args.tolerance); got != tt.want {
				t.Errorf("withinToleranceTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_withinTolerance(t *testing.T) {
	type args[T int | float64 | time.Duration] struct {
		a         T
		b         T
		tolerance T
	}
	tests := []struct {
		name string
		args args[int]
		want bool
	}{
		{
			name: "within tolerance int b > a",
			args: args[int]{
				a:         1,
				b:         2,
				tolerance: 10,
			},
			want: true,
		},
		{
			name: "within tolerance int a > b",
			args: args[int]{
				a:         4,
				b:         2,
				tolerance: 10,
			},
			want: true,
		},
		{
			name: "within tolerance not ok int b > a",
			args: args[int]{
				a:         1,
				b:         5,
				tolerance: 2,
			},
			want: false,
		},
		{
			name: "within tolerance not ok int a > b",
			args: args[int]{
				a:         8,
				b:         2,
				tolerance: 3,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := withinTolerance(tt.args.a, tt.args.b, tt.args.tolerance); got != tt.want {
				t.Errorf("withinTolerance() = %v, want %v", got, tt.want)
			}
		})
	}

	testsDuration := []struct {
		name string
		args args[time.Duration]
		want bool
	}{
		{
			name: "within tolerance int b > a",
			args: args[time.Duration]{
				a:         1 * time.Second,
				b:         2 * time.Second,
				tolerance: 10 * time.Second,
			},
			want: true,
		},
		{
			name: "within tolerance int a > b",
			args: args[time.Duration]{
				a:         4 * time.Second,
				b:         2 * time.Second,
				tolerance: 10 * time.Minute,
			},
			want: true,
		},
		{
			name: "within tolerance not ok int b > a",
			args: args[time.Duration]{
				a:         1 * time.Second,
				b:         2 * time.Minute,
				tolerance: 10 * time.Second,
			},
			want: false,
		},
		{
			name: "within tolerance not ok int a > b",
			args: args[time.Duration]{
				a:         1 * time.Minute,
				b:         2 * time.Second,
				tolerance: 10 * time.Nanosecond,
			},
			want: false,
		},
	}
	for _, tt := range testsDuration {
		t.Run(tt.name, func(t *testing.T) {
			if got := withinTolerance(tt.args.a, tt.args.b, tt.args.tolerance); got != tt.want {
				t.Errorf("withinTolerance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_valuesAreEqual(t *testing.T) {
	type args struct {
		config   Config
		key      string
		output   any
		expected any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "int ok",
			args: args{
				config:   Config{},
				key:      "test",
				output:   1,
				expected: 1,
			},
			wantErr: false,
		},
		{
			name: "int not ok",
			args: args{
				config:   Config{},
				key:      "test",
				output:   1,
				expected: 2,
			},
			wantErr: true,
		},
		{
			name: "int ok within threshold",
			args: args{
				config: Config{
					Thresholds: Tresholds{
						Int: 4,
					},
				},
				key:      "test",
				output:   1,
				expected: 2,
			},
			wantErr: false,
		},
		{
			name: "float ok",
			args: args{
				config:   Config{},
				key:      "test",
				output:   1.234,
				expected: 1.234,
			},
			wantErr: false,
		},
		{
			name: "float not ok",
			args: args{
				config:   Config{},
				key:      "test",
				output:   1.234,
				expected: 1.2345,
			},
			wantErr: true,
		},
		{
			name: "float ok within threshold",
			args: args{
				config: Config{
					Thresholds: Tresholds{
						Float: 0.1,
					},
				},
				key:      "test",
				output:   1.234,
				expected: 1.2345,
			},
			wantErr: false,
		},
		{
			name: "float ok within custom threshold",
			args: args{
				config: Config{
					Thresholds: Tresholds{
						Float: 0.1,
						CustomThresholds: CustomThresholds{
							Float: map[string]float64{
								"test": 0.25,
							},
						},
					},
				},
				key:      "test",
				output:   1.2,
				expected: 1.0,
			},
			wantErr: false,
		},
		{
			name: "bool ok",
			args: args{
				config:   Config{},
				key:      "test",
				output:   true,
				expected: true,
			},
			wantErr: false,
		},
		{
			name: "bool not ok",
			args: args{
				config:   Config{},
				key:      "test",
				output:   true,
				expected: false,
			},
			wantErr: true,
		},
		{
			name: "string ok",
			args: args{
				config:   Config{},
				key:      "test",
				output:   "foo",
				expected: "foo",
			},
			wantErr: false,
		},
		{
			name: "string not ok",
			args: args{
				config:   Config{},
				key:      "test",
				output:   "foo",
				expected: "bar",
			},
			wantErr: true,
		},
		{
			name: "time.Time ok",
			args: args{
				config:   Config{},
				key:      "test",
				output:   "2023-01-01T00:00:00Z",
				expected: "2023-01-01T00:00:00Z",
			},
			wantErr: false,
		},
		{
			name: "time.Time not ok",
			args: args{
				config:   Config{},
				key:      "test",
				output:   "2023-01-01T00:00:00Z",
				expected: "2023-01-01T12:20:00Z",
			},
			wantErr: true,
		},
		{
			name: "time.Time ok with threshold",
			args: args{
				config: Config{
					Thresholds: Tresholds{
						Time: 16 * time.Hour,
					},
				},
				key:      "test",
				output:   "2023-01-01T00:00:00Z",
				expected: "2023-01-01T12:20:00Z",
			},
			wantErr: false,
		},
		{
			name: "time.Duration ok",
			args: args{
				config:   Config{},
				key:      "test",
				output:   "456ms",
				expected: "456ms",
			},
			wantErr: false,
		},
		{
			name: "time.Duration not ok",
			args: args{
				config:   Config{},
				key:      "test",
				output:   "456ms",
				expected: "789ms",
			},
			wantErr: true,
		},
		{
			name: "time.Duration ok with threshold",
			args: args{
				config: Config{
					Thresholds: Tresholds{
						Duration: 500 * time.Millisecond,
					},
				},
				key:      "test",
				output:   "456ms",
				expected: "789ms",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := valuesAreEqual(tt.args.config, tt.args.key, tt.args.output, tt.args.expected); (err != nil) != tt.wantErr {
				t.Errorf("valuesAreEqual() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
