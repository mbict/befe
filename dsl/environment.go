package dsl

import (
	"os"
	"strconv"
)

func FromEnv(name string) string {
	return FromEnvWithDefault(name, "")
}

func FromEnvInt(name string) int {
	return FromEnvWithDefaultInt(name, 0)
}

func FromEnvWithDefault(name string, defaultValue string) string {
	value, found := os.LookupEnv(name)
	if found == false || value == "" {
		return defaultValue
	}
	return value
}

func FromEnvWithDefaultInt(name string, defaultValue int) int {
	value, found := os.LookupEnv(name)
	if found == false || value == "" {
		return defaultValue
	}

	val, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}
	return val
}

/*func FromEnvWithDefault[T int | string](name string, defaultValue T) T {

	value, found := os.LookupEnv(name)
	if found == false || value == "" {
		return defaultValue
	}

	res := new(T)
	v := reflect.ValueOf(res)
	switch v.Elem().Kind() {

	case reflect.Int:
		val, err := strconv.Atoi(value)
		if err != nil {
			panic(err)
		}
		v.Elem().SetInt(int64(val))
	case reflect.String:
		v.Elem().SetString(value)
	}

	return *res
}*/
