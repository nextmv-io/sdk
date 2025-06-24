package golden

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_replaceTransient(t *testing.T) {
	type args struct {
		original        map[string]any
		transientFields []TransientField
	}
	tests := []struct {
		name string
		args args
		want map[string]any
	}{
		{
			name: "map with transient int",
			args: args{
				original: map[string]any{
					"$.a": "foo",
					"$.b": 2,
				},
				transientFields: []TransientField{{Key: ".b"}},
			},
			want: map[string]any{
				"$.a": "foo",
				"$.b": 123,
			},
		},
		{
			name: "map with transient float",
			args: args{
				original: map[string]any{
					"$.a": "foo",
					"$.b": 1.2,
				},
				transientFields: []TransientField{{Key: ".b"}},
			},
			want: map[string]any{
				"$.a": "foo",
				"$.b": 0.123,
			},
		},
		{
			name: "map with transient time",
			args: args{
				original: map[string]any{
					"$.a": "foo",
					"$.b": "2023-05-04T19:52:53Z",
				},
				transientFields: []TransientField{{Key: ".b"}},
			},
			want: map[string]any{
				"$.a": "foo",
				"$.b": "2023-01-01T00:00:00Z",
			},
		},
		{
			name: "map with transient time duration",
			args: args{
				original: map[string]any{
					"$.a": "foo",
					"$.b": "666ms",
				},
				transientFields: []TransientField{{Key: ".b"}},
			},
			want: map[string]any{
				"$.a": "foo",
				"$.b": "123ms",
			},
		},
		{
			name: "map with transient string",
			args: args{
				original: map[string]any{
					"$.a": "foo",
					"$.b": "bar",
				},
				transientFields: []TransientField{{Key: ".b"}},
			},
			want: map[string]any{
				"$.a": "foo",
				"$.b": "text",
			},
		},
		{
			name: "map with transient bool",
			args: args{
				original: map[string]any{
					"$.a": "foo",
					"$.b": true,
				},
				transientFields: []TransientField{{Key: ".b"}},
			},
			want: map[string]any{
				"$.a": "foo",
				"$.b": true,
			},
		},
		{
			name: "map with array",
			args: args{
				original: map[string]any{
					"$.a[0].b": "foo",
					"$.a[0].c": 1.2,
					"$.a[1].b": "bar",
					"$.a[1].c": 3.4,
				},
				transientFields: []TransientField{
					{Key: "$.a[].b", Replacement: "text"},
				},
			},
			want: map[string]any{
				"$.a[0].b": "text",
				"$.a[0].c": 1.2,
				"$.a[1].b": "text",
				"$.a[1].c": 3.4,
			},
		},
		{
			name: "map with replaced parent key",
			args: args{
				original: map[string]any{
					"$.a.b": "foo",
					"$.a.c": 1.2,
				},
				transientFields: []TransientField{
					{Key: ".a", Replacement: map[string]int{"foo": 123}},
				},
			},
			want: map[string]any{
				"$.a": map[string]int{
					"foo": 123,
				},
			},
		},
		{
			name: "map with multiple replaced parent keys",
			args: args{
				original: map[string]any{
					"$.a[0].foo.b": 123,
					"$.a[0].foo.c": 123,
					"$.a[1].foo.b": 456,
					"$.a[1].foo.c": 456,
				},
				transientFields: []TransientField{
					{Key: ".a[].foo", Replacement: "replaced"},
				},
			},
			want: map[string]any{
				"$.a[0].foo": "replaced",
				"$.a[1].foo": "replaced",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := replaceTransient("test/path/to-file.golden", tt.args.original, tt.args.transientFields...)
			if err != nil {
				t.Errorf("replaceTransient() error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("replaceTransient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_round(t *testing.T) {
	tests := []struct {
		num       float64
		want      float64
		precision int
	}{
		{
			num:       1.23456789,
			want:      1.235,
			precision: 3,
		},
		{
			num:       1.234567891234567891234,
			want:      1.23456789123456789123,
			precision: 20,
		},
		{
			num:       1.234567891234567891234,
			want:      1,
			precision: 0,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("round %d", tt.precision), func(t *testing.T) {
			if got := round(tt.num, tt.precision); got != tt.want {
				t.Errorf("round() = %v, want %v", got, tt.want)
			}
		})
	}
}
