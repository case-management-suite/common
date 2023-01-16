package factory

import (
	"fmt"
	"reflect"
)

type FactorySet struct{}

func (FactorySet) implementsSettable() {}

type FactorySettble interface {
	implementsSettable()
}

func ValidateFactorySet[T FactorySettble](f T) error {
	var err error
	t := reflect.TypeOf(f)
	for i := 0; i < t.NumField(); i++ {
		fieldName := t.Field(i).Name

		refFactoriesValue := reflect.ValueOf(f)
		attributeFactory := refFactoriesValue.Field(i)

		if !refFactoriesValue.IsValid() || refFactoriesValue.IsZero() {
			return fmt.Errorf("factory not valid or not initialized: %s", fieldName)
		}

		switch attributeFactory.Type().Kind() {
		case reflect.Struct:
			switch value := refFactoriesValue.Field(i).Interface().(type) {
			case FactorySettble:
				err = ValidateFactorySet(value)
				if err != nil {
					return fmt.Errorf("factory: %s -> %w;", fieldName, err)
				}
			}
		case reflect.Func:
			if attributeFactory.IsNil() {
				return fmt.Errorf("one of the factory functions is not initialized: %s", fieldName)
			}
		default:
		}

	}
	return err
}
