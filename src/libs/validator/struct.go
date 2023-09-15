package validator

func Struct[T any](field T) error {
	if err := validate.Struct(field); err != nil {
		return err
	}
	return nil
}
