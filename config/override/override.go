package override

import (
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type EnvOptions struct {
	Prefix    string
	Separator string
}

func ApplyEnv(v any, opts EnvOptions) error {
	if opts.Separator == "" {
		opts.Separator = "_"
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return nil
	}
	return apply("", rv.Elem(), opts)
}

func apply(path string, v reflect.Value, opts EnvOptions) error {
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil
	}
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		sf := t.Field(i)
		if sf.PkgPath != "" {
			continue
		}
		name := sf.Name
		if tag := sf.Tag.Get("env"); tag != "" {
			name = tag
		}
		envKey := strings.ToUpper(strings.Join(nonEmpty(opts.Prefix, path, name), opts.Separator))
		f := v.Field(i)
		if f.Kind() == reflect.Struct && f.Type() != reflect.TypeOf(time.Duration(0)) {
			if err := apply(strings.Join(nonEmpty(path, name), opts.Separator), f, opts); err != nil {
				return err
			}
			continue
		}
		if raw, ok := os.LookupEnv(envKey); ok {
			if err := set(f, raw); err != nil {
				return err
			}
		}
	}
	return nil
}

func nonEmpty(parts ...string) []string {
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func set(v reflect.Value, raw string) error {
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}
	if !v.CanSet() {
		return nil
	}
	if v.Type() == reflect.TypeOf(time.Duration(0)) {
		d, err := time.ParseDuration(raw)
		if err != nil {
			return err
		}
		v.SetInt(int64(d))
		return nil
	}
	switch v.Kind() {
	case reflect.String:
		v.SetString(raw)
	case reflect.Bool:
		b, err := strconv.ParseBool(raw)
		if err != nil {
			return err
		}
		v.SetBool(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(raw, 10, v.Type().Bits())
		if err != nil {
			return err
		}
		v.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(raw, 10, v.Type().Bits())
		if err != nil {
			return err
		}
		v.SetUint(n)
	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(raw, v.Type().Bits())
		if err != nil {
			return err
		}
		v.SetFloat(n)
	case reflect.Slice:
		parts := strings.Split(raw, ",")
		s := reflect.MakeSlice(v.Type(), len(parts), len(parts))
		for i := range parts {
			if err := set(s.Index(i), strings.TrimSpace(parts[i])); err != nil {
				return err
			}
		}
		v.Set(s)
	}
	return nil
}
