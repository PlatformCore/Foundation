package validator

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Rule func(path string, value reflect.Value) error

type ValidationError struct {
	Path string
	Err  error
}

func (e ValidationError) Error() string { return e.Path + ": " + e.Err.Error() }

type Errors []ValidationError

func (e Errors) Error() string {
	parts := make([]string, 0, len(e))
	for _, err := range e {
		parts = append(parts, err.Error())
	}
	return strings.Join(parts, "; ")
}

func Validate(v any) error {
	if v == nil {
		return errors.New("config validator: nil config")
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return errors.New("config validator: nil pointer")
		}
		rv = rv.Elem()
	}
	errs := walk("", rv)
	if len(errs) > 0 {
		return errs
	}
	return nil
}

func walk(path string, v reflect.Value) Errors {
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil
	}
	t := v.Type()
	var errs Errors
	for i := 0; i < v.NumField(); i++ {
		sf := t.Field(i)
		if sf.PkgPath != "" {
			continue
		}
		name := sf.Name
		if tag := sf.Tag.Get("yaml"); tag != "" {
			name = strings.Split(tag, ",")[0]
		}
		if tag := sf.Tag.Get("json"); tag != "" && tag != "-" {
			name = strings.Split(tag, ",")[0]
		}
		childPath := name
		if path != "" {
			childPath = path + "." + name
		}
		fv := v.Field(i)
		if err := apply(childPath, fv, sf.Tag.Get("validate")); err != nil {
			if ve, ok := err.(Errors); ok {
				errs = append(errs, ve...)
			} else {
				errs = append(errs, ValidationError{Path: childPath, Err: err})
			}
		}
		errs = append(errs, walk(childPath, fv)...)
	}
	return errs
}

func apply(path string, v reflect.Value, tag string) error {
	if tag == "" {
		return nil
	}
	var errs Errors
	for _, part := range strings.Split(tag, ",") {
		if part == "" {
			continue
		}
		name, arg, _ := strings.Cut(part, "=")
		if err := evalRule(name, arg, v); err != nil {
			errs = append(errs, ValidationError{Path: path, Err: err})
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

func evalRule(name, arg string, v reflect.Value) error {
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			if name == "required" {
				return errors.New("is required")
			}
			return nil
		}
		v = v.Elem()
	}
	switch name {
	case "required":
		if isZero(v) {
			return errors.New("is required")
		}
	case "min":
		n, _ := strconv.ParseFloat(arg, 64)
		if number(v) < n {
			return fmt.Errorf("must be >= %v", arg)
		}
	case "max":
		n, _ := strconv.ParseFloat(arg, 64)
		if number(v) > n {
			return fmt.Errorf("must be <= %v", arg)
		}
	case "duration":
		if v.Kind() == reflect.String {
			if _, err := time.ParseDuration(v.String()); err != nil {
				return err
			}
		}
	case "url":
		if _, err := url.ParseRequestURI(v.String()); err != nil {
			return err
		}
	case "hostport":
		if _, _, err := net.SplitHostPort(v.String()); err != nil {
			return err
		}
	case "file":
		if _, err := os.Stat(v.String()); err != nil {
			return err
		}
	case "regexp":
		if _, err := regexp.Compile(v.String()); err != nil {
			return err
		}
	}
	return nil
}

func isZero(v reflect.Value) bool {
	return !v.IsValid() || reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}

func number(v reflect.Value) float64 {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(v.Uint())
	case reflect.Float32, reflect.Float64:
		return v.Float()
	case reflect.String:
		return float64(len(v.String()))
	case reflect.Slice, reflect.Array, reflect.Map:
		return float64(v.Len())
	}
	return 0
}
