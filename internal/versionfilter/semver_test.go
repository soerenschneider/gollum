package versionfilter

import (
	"log"
	"testing"

	"github.com/Masterminds/semver/v3"
)

func mustConstraint(constraint string) *semver.Constraints {
	c, err := semver.NewConstraint(constraint)
	if err != nil {
		log.Fatal(err)
	}
	return c
}

func TestSemverReleaseFilter_Matches(t *testing.T) {
	type fields struct {
		constraint *semver.Constraints
	}
	type args struct {
		tag string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "",
			fields: fields{
				constraint: mustConstraint(">= 1.0.0"),
			},
			args: args{
				tag: "v1.0.1",
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &SemverReleaseFilter{
				constraint: tt.fields.constraint,
			}
			got, err := f.Matches(tt.args.tag)
			if (err != nil) != tt.wantErr {
				t.Errorf("Matches() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Matches() got = %v, want %v", got, tt.want)
			}
		})
	}
}
