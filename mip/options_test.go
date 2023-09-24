package mip_test

import (
	"reflect"
	"testing"

	"github.com/nextmv-io/sdk/mip"
)

func TestControlOptions_ToTyped(t *testing.T) {
	type fields struct {
		Bool   string
		Float  string
		Int    string
		String string
	}
	tests := []struct {
		name    string
		fields  fields
		want    *mip.TypedControlOptions
		wantErr bool
	}{
		{
			name: "empty fields",
			fields: fields{
				Bool:   "",
				Float:  "",
				Int:    "",
				String: "",
			},
			want: &mip.TypedControlOptions{
				Bool:   []mip.TypedControlOption[bool]{},
				Float:  []mip.TypedControlOption[float64]{},
				Int:    []mip.TypedControlOption[int]{},
				String: []mip.TypedControlOption[string]{},
			},
			wantErr: false,
		},
		{
			name: "1 param bool",
			fields: fields{
				Bool:   "foo=true",
				Float:  "",
				Int:    "",
				String: "",
			},
			want: &mip.TypedControlOptions{
				Bool: []mip.TypedControlOption[bool]{
					{
						Name:  "foo",
						Value: true,
					},
				},
				Float:  []mip.TypedControlOption[float64]{},
				Int:    []mip.TypedControlOption[int]{},
				String: []mip.TypedControlOption[string]{},
			},
			wantErr: false,
		},
		{
			name: "1 param float",
			fields: fields{
				Bool:   "",
				Float:  "foo=1",
				Int:    "",
				String: "",
			},
			want: &mip.TypedControlOptions{
				Bool: []mip.TypedControlOption[bool]{},
				Float: []mip.TypedControlOption[float64]{
					{
						Name:  "foo",
						Value: 1,
					},
				},
				Int:    []mip.TypedControlOption[int]{},
				String: []mip.TypedControlOption[string]{},
			},
			wantErr: false,
		},
		{
			name: "1 param int",
			fields: fields{
				Bool:   "",
				Float:  "",
				Int:    "foo=1",
				String: "",
			},
			want: &mip.TypedControlOptions{
				Bool:  []mip.TypedControlOption[bool]{},
				Float: []mip.TypedControlOption[float64]{},
				Int: []mip.TypedControlOption[int]{
					{
						Name:  "foo",
						Value: 1,
					},
				},
				String: []mip.TypedControlOption[string]{},
			},
			wantErr: false,
		},
		{
			name: "1 param string",
			fields: fields{
				Bool:   "",
				Float:  "",
				Int:    "",
				String: "foo=bar",
			},
			want: &mip.TypedControlOptions{
				Bool:  []mip.TypedControlOption[bool]{},
				Float: []mip.TypedControlOption[float64]{},
				Int:   []mip.TypedControlOption[int]{},
				String: []mip.TypedControlOption[string]{
					{
						Name:  "foo",
						Value: "bar",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "wrong format bool 1",
			fields: fields{
				Bool:   "foo=bar",
				Float:  "",
				Int:    "",
				String: "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "wrong format bool 2",
			fields: fields{
				Bool:   "foo;true",
				Float:  "",
				Int:    "",
				String: "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "wrong format float 1",
			fields: fields{
				Bool:   "",
				Float:  "foo=bar",
				Int:    "",
				String: "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "wrong format float 2",
			fields: fields{
				Bool:   "",
				Float:  "foo;1.23",
				Int:    "",
				String: "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "wrong format int 1",
			fields: fields{
				Bool:   "",
				Float:  "",
				Int:    "foo=bar",
				String: "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "wrong format int 2",
			fields: fields{
				Bool:   "",
				Float:  "",
				Int:    "foo;1.23",
				String: "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "wrong format string",
			fields: fields{
				Bool:   "",
				Float:  "",
				Int:    "",
				String: "foo;1.23",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "multi params",
			fields: fields{
				Bool:   "foo=true,bar=false",
				Float:  "foo=1,bar=2.34",
				Int:    "foo=1,bar=2",
				String: "foo=bar,roh=baz",
			},
			want: &mip.TypedControlOptions{
				Bool: []mip.TypedControlOption[bool]{
					{
						Name:  "foo",
						Value: true,
					},
					{
						Name:  "bar",
						Value: false,
					},
				},
				Float: []mip.TypedControlOption[float64]{
					{
						Name:  "foo",
						Value: 1.,
					},
					{
						Name:  "bar",
						Value: 2.34,
					},
				},
				Int: []mip.TypedControlOption[int]{
					{
						Name:  "foo",
						Value: 1,
					},
					{
						Name:  "bar",
						Value: 2,
					},
				},
				String: []mip.TypedControlOption[string]{
					{
						Name:  "foo",
						Value: "bar",
					},
					{
						Name:  "roh",
						Value: "baz",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controlOptions := mip.ControlOptions{
				Bool:   tt.fields.Bool,
				Float:  tt.fields.Float,
				Int:    tt.fields.Int,
				String: tt.fields.String,
			}
			got, err := controlOptions.ToTyped()
			if (err != nil) != tt.wantErr {
				t.Errorf("ControlOptions.ToTyped() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ControlOptions.ToTyped() = %v, want %v", got, tt.want)
			}
		})
	}
}
