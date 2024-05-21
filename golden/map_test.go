package golden

import (
	"reflect"
	"testing"
)

func Test_flatten(t *testing.T) {
	type args struct {
		nested map[string]any
	}
	tests := []struct {
		name string
		args args
		want map[string]any
	}{
		{
			name: "flat",
			args: args{
				nested: map[string]any{
					"a": "foo",
					"b": 2,
					"c": true,
				},
			},
			want: map[string]any{
				".a": "foo",
				".b": 2,
				".c": true,
			},
		},
		{
			name: "flat with nil",
			args: args{
				nested: map[string]any{
					"a": "foo",
					"b": nil,
					"c": true,
				},
			},
			want: map[string]any{
				".a": "foo",
				".b": nil,
				".c": true,
			},
		},
		{
			name: "slice",
			args: args{
				nested: map[string]any{
					"a": "foo",
					"b": []any{
						"bar",
						2,
					},
				},
			},
			want: map[string]any{
				".a":    "foo",
				".b[0]": "bar",
				".b[1]": 2,
			},
		},
		{
			name: "nested map",
			args: args{
				nested: map[string]any{
					"a": "foo",
					"b": map[string]any{
						"c": "bar",
						"d": 2,
					},
				},
			},
			want: map[string]any{
				".a":   "foo",
				".b.c": "bar",
				".b.d": 2,
			},
		},
		{
			name: "slice with nested maps",
			args: args{
				nested: map[string]any{
					"a": "foo",
					"b": []any{
						map[string]any{
							"c": "bar",
							"d": 2,
						},
						map[string]any{
							"c": "baz",
							"d": 3,
						},
					},
				},
			},
			want: map[string]any{
				".a":      "foo",
				".b[0].c": "bar",
				".b[0].d": 2,
				".b[1].c": "baz",
				".b[1].d": 3,
			},
		},
		{
			name: "slice with nested maps with nested slice",
			args: args{
				nested: map[string]any{
					"a": "foo",
					"b": []any{
						map[string]any{
							"c": "bar",
							"d": []any{
								2,
								true,
							},
						},
						map[string]any{
							"c": "baz",
							"d": []any{
								3,
								false,
							},
						},
					},
				},
			},
			want: map[string]any{
				".a":         "foo",
				".b[0].c":    "bar",
				".b[0].d[0]": 2,
				".b[0].d[1]": true,
				".b[1].c":    "baz",
				".b[1].d[0]": 3,
				".b[1].d[1]": false,
			},
		},
		{
			name: "slice with nested maps with nested slice with nested map",
			args: args{
				nested: map[string]any{
					"a": "foo",
					"b": []any{
						map[string]any{
							"c": "bar",
							"d": []any{
								map[string]any{
									"e": 2,
								},
								true,
							},
						},
						map[string]any{
							"c": "baz",
							"d": []any{
								map[string]any{
									"e": 3,
								},
								false,
							},
						},
					},
				},
			},
			want: map[string]any{
				".a":           "foo",
				".b[0].c":      "bar",
				".b[0].d[0].e": 2,
				".b[0].d[1]":   true,
				".b[1].c":      "baz",
				".b[1].d[0].e": 3,
				".b[1].d[1]":   false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := flatten(tt.args.nested); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("flatten() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
					".a": "foo",
					".b": 2,
				},
				transientFields: []TransientField{{Key: ".b"}},
			},
			want: map[string]any{
				".a": "foo",
				".b": 123,
			},
		},
		{
			name: "map with transient float",
			args: args{
				original: map[string]any{
					".a": "foo",
					".b": 1.2,
				},
				transientFields: []TransientField{{Key: ".b"}},
			},
			want: map[string]any{
				".a": "foo",
				".b": 0.123,
			},
		},
		{
			name: "map with transient time",
			args: args{
				original: map[string]any{
					".a": "foo",
					".b": "2023-05-04T19:52:53Z",
				},
				transientFields: []TransientField{{Key: ".b"}},
			},
			want: map[string]any{
				".a": "foo",
				".b": "2023-01-01T00:00:00Z",
			},
		},
		{
			name: "map with transient time duration",
			args: args{
				original: map[string]any{
					".a": "foo",
					".b": "666ms",
				},
				transientFields: []TransientField{{Key: ".b"}},
			},
			want: map[string]any{
				".a": "foo",
				".b": "123ms",
			},
		},
		{
			name: "map with transient string",
			args: args{
				original: map[string]any{
					".a": "foo",
					".b": "bar",
				},
				transientFields: []TransientField{{Key: ".b"}},
			},
			want: map[string]any{
				".a": "foo",
				".b": "text",
			},
		},
		{
			name: "map with transient bool",
			args: args{
				original: map[string]any{
					".a": "foo",
					".b": true,
				},
				transientFields: []TransientField{{Key: ".b"}},
			},
			want: map[string]any{
				".a": "foo",
				".b": true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := replaceTransient(tt.args.original, tt.args.transientFields...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("replaceTransient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nest(t *testing.T) {
	type args struct {
		flattened map[string]any
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{
			name: "flat",
			args: args{
				flattened: map[string]any{
					".a": "foo",
					".b": 2,
					".c": true,
				},
			},
			want: map[string]any{
				"a": "foo",
				"b": 2,
				"c": true,
			},
		},
		{
			name: "flat with nil",
			args: args{
				flattened: map[string]any{
					".a": "foo",
					".b": nil,
					".c": true,
				},
			},
			want: map[string]any{
				"a": "foo",
				"b": nil,
				"c": true,
			},
		},
		{
			name: "slice",
			args: args{
				flattened: map[string]any{
					".a":    "foo",
					".b[0]": "bar",
					".b[1]": 2,
				},
			},
			want: map[string]any{
				"a": "foo",
				"b": []any{
					"bar",
					2,
				},
			},
		},
		{
			name: "nested map",
			args: args{
				flattened: map[string]any{
					".a":   "foo",
					".b.c": "bar",
					".b.d": 2,
				},
			},
			want: map[string]any{
				"a": "foo",
				"b": map[string]any{
					"c": "bar",
					"d": 2,
				},
			},
		},
		{
			name: "slice with nested maps",
			args: args{
				flattened: map[string]any{
					".a":      "foo",
					".b[0].c": "bar",
					".b[0].d": 2,
					".b[1].c": "baz",
					".b[1].d": 3,
				},
			},
			want: map[string]any{
				"a": "foo",
				"b": []any{
					map[string]any{
						"c": "bar",
						"d": 2,
					},
					map[string]any{
						"c": "baz",
						"d": 3,
					},
				},
			},
		},
		{
			name: "slice with nested maps with nested slice",
			args: args{
				flattened: map[string]any{
					".a":         "foo",
					".b[0].c":    "bar",
					".b[0].d[0]": 2,
					".b[0].d[1]": true,
					".b[1].c":    "baz",
					".b[1].d[0]": 3,
					".b[1].d[1]": false,
				},
			},
			want: map[string]any{
				"a": "foo",
				"b": []any{
					map[string]any{
						"c": "bar",
						"d": []any{
							2,
							true,
						},
					},
					map[string]any{
						"c": "baz",
						"d": []any{
							3,
							false,
						},
					},
				},
			},
		},
		{
			name: "slice with nested maps with nested slice with nested map",
			args: args{
				flattened: map[string]any{
					".a":           "foo",
					".b[0].c":      "bar",
					".b[0].d[0].e": 2,
					".b[0].d[1]":   true,
					".b[1].c":      "baz",
					".b[1].d[0].e": 3,
					".b[1].d[1]":   false,
				},
			},
			want: map[string]any{
				"a": "foo",
				"b": []any{
					map[string]any{
						"c": "bar",
						"d": []any{
							map[string]any{
								"e": 2,
							},
							true,
						},
					},
					map[string]any{
						"c": "baz",
						"d": []any{
							map[string]any{
								"e": 3,
							},
							false,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := nest(tt.args.flattened); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("nest() = %v, want %v", got, tt.want)
			}
		})
	}
}
