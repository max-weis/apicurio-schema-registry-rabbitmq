package validate

import "github.com/linkedin/goavro"

type Validator struct {
	schema string
}

func NewValidator(schema string) Validator {
	return Validator{schema: schema}
}

func (v Validator) Validate(obj any) (bool, error) {
	codec, err := goavro.NewCodec(v.schema)
	if err != nil {
		return false, err
	}

	if _, err = codec.BinaryFromNative(nil, obj); err != nil {
		return false, err
	}

	return true, nil
}
