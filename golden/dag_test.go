// Package testing holds tools for testing documentation code.
package golden

import (
	"testing"
)

func Test_dagToMermaid(t *testing.T) {
	tests := []struct {
		dagCases []DagTestCase
		name     string
	}{
		{
			dagCases: []DagTestCase{
				{
					Name:  "Case 1",
					Needs: []string{},
				},
				{
					Name:  "Case 2",
					Needs: []string{"Case 1"},
				},
				{
					Name:  "Case 3",
					Needs: []string{"Case 1", "Case 2"},
				},
			},
			name: "valid DAG",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mermaid, err := dagToMermaid(tt.dagCases)
			if err != nil {
				t.Errorf("dagToMermaid() error = %v", err)
				return
			}
			t.Logf("DAG diagram (mermaid):\n%s", mermaid)
		})
	}
}
