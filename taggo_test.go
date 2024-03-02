package taggo_test

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/eduardomosko/taggo"
)

type circle struct {
	Radius int `json:"radius"`
}

type square struct {
	Side int `json:"side"`
}

var UnknownDiscriminatorError = errors.New("unknown discriminator")

type ShapeDiscriminator struct {
	Shape string `json:"shape"`
}


func (sd *ShapeDiscriminator) GetType() (any, error) {
	switch sd.Shape {
	case "square":
		return &square{}, nil
	case "circle":
		return &circle{}, nil
	}
	return nil, UnknownDiscriminatorError
}


func TestUnmarshalShape(t *testing.T) {
	tt := []struct {
		Input    string
		Expected any
		ExpectedErr error
	}{
		{Input: `{"shape":"square"}`, Expected: &square{}},
		{Input: `{"shape":"circle"}`, Expected: &circle{}},
		{Input: `{"shape":"square","side":10}`, Expected: &square{10}},
		{Input: `{"shape":"circle","radius":5}`, Expected: &circle{5}},
		{Input: `{"shape":"circle","side":11,"radius":20}`, Expected: &circle{20}},
		{Input: `{"shape":"square","side":11,"radius":20}`, Expected: &square{11}},
		{Input: `{"shape":"sqare","side":11,"radius":20}`, ExpectedErr: UnknownDiscriminatorError},
		{Input: `{"shpe":"square"}`, ExpectedErr: UnknownDiscriminatorError},
	}

	for _, tc := range tt {
		t.Run(tc.Input, func(t *testing.T) {
			type ShapeUnion = taggo.Discriminated[any, *ShapeDiscriminator]

			var shape ShapeUnion

			err := json.Unmarshal([]byte(tc.Input), &shape)
			if err != nil {
				if tc.ExpectedErr == nil {
					t.Fatalf("unmarshal error: %v", err)
				}
				if !errors.Is(err, UnknownDiscriminatorError) {
					t.Errorf("expected error %#v; got %#v", tc.ExpectedErr, err)
				}
			}

			if !reflect.DeepEqual(tc.Expected, shape.Value) {
				t.Errorf("expected %#v; got %#v", tc.Expected, shape.Value)
			}
		})
	}
}
