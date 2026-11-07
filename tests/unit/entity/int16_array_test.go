package entity_test

import (
	"reflect"
	"testing"

	"github.com/rescuebyte/rescuebyte-be/domain/entity"
)

func TestInt16ArrayScan(t *testing.T) {
	tests := []struct {
		name    string
		src     any
		want    entity.Int16Array
		wantErr bool
	}{
		{"string", "{1,2,3,4,5,6}", entity.Int16Array{1, 2, 3, 4, 5, 6}, false},
		{"bytes", []byte("{1,7}"), entity.Int16Array{1, 7}, false},
		{"empty", "{}", entity.Int16Array{}, false},
		{"nil", nil, entity.Int16Array{}, false},
		{"invalid element", "{1,a}", nil, true},
		{"unsupported type", 42, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got entity.Int16Array
			err := got.Scan(tt.src)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Scan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt16ArrayValue(t *testing.T) {
	tests := []struct {
		name string
		in   entity.Int16Array
		want string
	}{
		{"values", entity.Int16Array{1, 2, 3}, "{1,2,3}"},
		{"empty", entity.Int16Array{}, "{}"},
		{"nil", nil, "{}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.in.Value()
			if err != nil {
				t.Fatalf("Value() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}
