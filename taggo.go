package taggo

import (
	"encoding/json"
)

type Discriminator[ValueT any] interface {
	GetType() (ValueT, error)
}

type Discriminated[ValueT any, DiscriminatorT Discriminator[ValueT]] struct {
	Value ValueT
}

func (d *Discriminated[ValueT, DiscriminatorT]) UnmarshalJSON(bytes []byte) error {
	var discriminator DiscriminatorT
	err := json.Unmarshal(bytes, &discriminator)
	if err != nil {
		return err
	}

	d.Value, err = discriminator.GetType()
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, &d.Value)
}
