package pgrst

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func setOutput(input []byte, output any) (err error) {

	// if either parameters are null, return immediately
	if nil == input || nil == output {
		return nil
	}
	//fmt.Println(string(input))

	//fmt.Println(reflect.ValueOf(output).Elem().Kind())

	oKind := reflect.ValueOf(output).Elem()

	switch oKind.Kind() {
	case reflect.String:
		inputStr := string(input)
		if inputStr != "null" {
			oKind.Set(reflect.ValueOf(inputStr))
		}

	case reflect.Slice:
		err = json.Unmarshal(input, output)
		if nil != err {
			return fmt.Errorf("error unmarshaling input to slice output: %w", err)
		}

	case reflect.Struct, reflect.Map:
		var rMap ([]map[string]any)

		//fmt.Println(string(input))
		err = json.Unmarshal(input, &rMap)
		if nil != err {
			return fmt.Errorf("error unmarshaling input to map: %w", err)
		}

		// if there isn't EXACTLY one entry in the result, we return immediately
		if len(rMap) != 1 {
			return
		}

		bytes, err := json.Marshal(rMap[0])
		if nil != err {
			return fmt.Errorf("error marshaling map: %w", err)
		}

		/*
			to unmarshal slices, write a custom json unmarshal to convert the
			string value to a specific slice type in the struct.
		*/

		//fmt.Println(string(bytes))
		err = json.Unmarshal(bytes, output)
		if nil != err {
			return fmt.Errorf("error unmarshaling map to output: %w", err)
		}

	}

	return nil
}
