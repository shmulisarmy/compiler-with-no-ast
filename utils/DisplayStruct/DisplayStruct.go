package displayStruct

import (
	"fmt"
	"os"
	"reflect"
	"slices"
	"strings"
)

// ANSI color codes
const (
	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
	colorGray    = "\033[90m"
)

func DisplayStruct(v any) string {
	var sb strings.Builder
	display(reflect.ValueOf(v), 0, &sb)
	return sb.String()
}

func Print(v any) {
	if !slices.Contains(os.Args, "logging-off") {
		fmt.Println(DisplayStruct(v))
	}
}

func display(v reflect.Value, indent int, sb *strings.Builder) {
	if !v.IsValid() {
		sb.WriteString(colorGray + "nil" + colorReset)
		return
	}

	// Dereference interface
	for v.Kind() == reflect.Interface && !v.IsNil() {
		v = v.Elem()
	}

	if !v.IsValid() {
		sb.WriteString(colorGray + "nil" + colorReset)
		return
	}

	pad := strings.Repeat("  ", indent)

	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			sb.WriteString(colorGray + "nil" + colorReset)
		} else {
			sb.WriteString(colorGray + "&" + colorReset)
			display(v.Elem(), indent, sb)
		}

	case reflect.Struct:
		typeName := v.Type().Name()
		if typeName == "" {
			typeName = "struct"
		}
		sb.WriteString(colorYellow + typeName + colorReset + " {\n")
		for i := 0; i < v.NumField(); i++ {
			f := v.Type().Field(i)
			field := v.Field(i)

			sb.WriteString(pad + "  " + colorCyan + f.Name + colorReset + ": ")

			if !field.CanInterface() {
				sb.WriteString(colorGray + "<unexported>" + colorReset)
			} else {
				display(field, indent+1, sb)
			}
			sb.WriteString("\n")
		}
		sb.WriteString(pad + "}")

	case reflect.Slice:
		if v.IsNil() {
			sb.WriteString(colorGray + "nil" + colorReset)
		} else {
			sb.WriteString(colorMagenta + "[" + colorReset + "\n")
			for i := 0; i < v.Len(); i++ {
				sb.WriteString(pad + "  ")
				display(v.Index(i), indent+1, sb)
				sb.WriteString("\n")
			}
			sb.WriteString(pad + colorMagenta + "]" + colorReset)
		}

	case reflect.Array:
		sb.WriteString(colorMagenta + "[" + colorReset + "\n")
		for i := 0; i < v.Len(); i++ {
			sb.WriteString(pad + "  ")
			display(v.Index(i), indent+1, sb)
			sb.WriteString("\n")
		}
		sb.WriteString(pad + colorMagenta + "]" + colorReset)

	case reflect.Map:
		if v.IsNil() {
			sb.WriteString(colorGray + "nil" + colorReset)
		} else {
			sb.WriteString(colorMagenta + "{" + colorReset + "\n")
			iter := v.MapRange()
			for iter.Next() {
				sb.WriteString(pad + "  ")
				display(iter.Key(), indent+1, sb)
				sb.WriteString(": ")
				display(iter.Value(), indent+1, sb)
				sb.WriteString("\n")
			}
			sb.WriteString(pad + colorMagenta + "}" + colorReset)
		}

	case reflect.String:
		sb.WriteString(colorGreen + fmt.Sprintf("%q", v.String()) + colorReset)

	case reflect.Bool:
		sb.WriteString(colorBlue + fmt.Sprintf("%v", v.Bool()) + colorReset)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		sb.WriteString(colorBlue + fmt.Sprintf("%d", v.Int()) + colorReset)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		sb.WriteString(colorBlue + fmt.Sprintf("%d", v.Uint()) + colorReset)

	case reflect.Float32, reflect.Float64:
		sb.WriteString(colorBlue + fmt.Sprintf("%v", v.Float()) + colorReset)

	default:
		if v.CanInterface() {
			sb.WriteString(fmt.Sprintf("%v", v.Interface()))
		} else {
			sb.WriteString(colorGray + "<unexported>" + colorReset)
		}
	}
}
