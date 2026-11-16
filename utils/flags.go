package utils

import (
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/urfave/cli/v2"

	"github.com/ghthor/gotty/v2/pkg/homedir"
)

func GenerateFlags(options ...any) (flags []cli.Flag, mappings map[string]string, err error) {
	mappings = make(map[string]string)

	for _, struct_ := range options {
		v := reflect.ValueOf(struct_).Elem()
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field, structField := v.Field(i), t.Field(i)
			flagName := structField.Tag.Get("flagName")
			if flagName == "" {
				continue
			}
			envName := "GOTTY_" + strings.ToUpper(strings.ReplaceAll(flagName, "-", "_"))
			mappings[flagName] = structField.Name

			var aliases []string
			if flagShortName := structField.Tag.Get("flagSName"); flagShortName != "" {
				aliases = []string{flagShortName}
			}

			flagDescription := structField.Tag.Get("flagDescribe")

			switch field.Kind() {
			case reflect.String:
				flags = append(flags, &cli.StringFlag{
					Name:    flagName,
					Aliases: aliases,
					Value:   field.String(),
					Usage:   flagDescription,
					EnvVars: []string{envName},
				})
			case reflect.Bool:
				flags = append(flags, &cli.BoolFlag{
					Name:    flagName,
					Aliases: aliases,
					Usage:   flagDescription,
					EnvVars: []string{envName},
				})
			case reflect.Int:
				flags = append(flags, &cli.IntFlag{
					Name:    flagName,
					Aliases: aliases,
					Value:   int(field.Int()),
					Usage:   flagDescription,
					EnvVars: []string{envName},
				})
			}
		}
	}

	return
}

func ApplyFlags(
	flags []cli.Flag,
	mappingHint map[string]string,
	c *cli.Context,
	options ...any,
) {
	values := make([]reflect.Value, len(options))
	for i, struct_ := range options {
		values[i] = reflect.ValueOf(struct_).Elem()
	}

	for flagName, fieldName := range mappingHint {
		if !c.IsSet(flagName) {
			continue
		}
		var field reflect.Value
		found := false
		for _, v := range values {
			if f := v.FieldByName(fieldName); f.IsValid() && f.CanSet() {
				field, found = f, true
				break
			}
		}
		if !found {
			continue
		}
		switch field.Kind() {
		case reflect.String:
			field.SetString(c.String(flagName))
		case reflect.Bool:
			field.SetBool(c.Bool(flagName))
		case reflect.Int:
			field.SetInt(int64(c.Int(flagName)))
		}
	}
}

func ApplyConfigFile(filePath string, options ...any) error {
	filePath = homedir.Expand(filePath)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return err
	}

	log.Printf("Loading config file at: %s", filePath)
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCLFile(filePath)
	if diags.HasErrors() {
		return diags
	}

	for _, object := range options {
		if diags := gohcl.DecodeBody(lenientBody{file.Body}, nil, object); diags.HasErrors() {
			return diags
		}
	}

	return nil
}

// lenientBody wraps an hcl.Body so that gohcl.DecodeBody only sees the
// attributes and blocks declared in the target's schema. Undeclared items
// are ignored instead of raising "unsupported argument" diagnostics,
// matching the lenient behavior of hcl v1's hcl.Decode (the same config
// file is decoded into several option structs).
type lenientBody struct {
	hcl.Body
}

func (b lenientBody) Content(schema *hcl.BodySchema) (*hcl.BodyContent, hcl.Diagnostics) {
	content, _, diags := b.Body.PartialContent(schema)
	return content, diags
}
