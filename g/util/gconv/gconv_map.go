// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import (
    "github.com/gogf/gf/g/internal/empty"
    "github.com/gogf/gf/g/text/gstr"
    "reflect"
    "strings"
)

// Map converts any variable <i> to map[string]interface{}.
// If the parameter <i> is not a map type, then the conversion will fail and returns nil.
// If <i> is a struct object, the second parameter noTagCheck means that
// the json tag is not detected,
// otherwise the json tag will be used as the map key name.
func Map(value interface{}, noTagCheck...bool) map[string]interface{} {
    if value == nil {
        return nil
    }
    if r, ok := value.(map[string]interface{}); ok {
        return r
    } else {
        // Only assert the common combination type of maps, and finally use reflection.
        m := make(map[string]interface{})
        switch value.(type) {
            case map[interface{}]interface{}:
                for k, v := range value.(map[interface{}]interface{}) {
                    m[String(k)] = v
                }
            case map[interface{}]string:
                for k, v := range value.(map[interface{}]string) {
                    m[String(k)] = v
                }
            case map[interface{}]int:
                for k, v := range value.(map[interface{}]int) {
                    m[String(k)] = v
                }
            case map[interface{}]uint:
                for k, v := range value.(map[interface{}]uint) {
                    m[String(k)] = v
                }
            case map[interface{}]float32:
                for k, v := range value.(map[interface{}]float32) {
                    m[String(k)] = v
                }
            case map[interface{}]float64:
                for k, v := range value.(map[interface{}]float64) {
                    m[String(k)] = v
                }

            case map[string]bool:
                for k, v := range value.(map[string]bool) {
                    m[k] = v
                }
            case map[string]int:
                for k, v := range value.(map[string]int) {
                    m[k] = v
                }
            case map[string]uint:
                for k, v := range value.(map[string]uint) {
                    m[k] = v
                }
            case map[string]float32:
                for k, v := range value.(map[string]float32) {
                    m[k] = v
                }
            case map[string]float64:
                for k, v := range value.(map[string]float64) {
                    m[k] = v
                }

            case map[int]interface{}:
                for k, v := range value.(map[int]interface{}) {
                    m[String(k)] = v
                }
            case map[int]string:
                for k, v := range value.(map[int]string) {
                    m[String(k)] = v
                }
            case map[uint]string:
                for k, v := range value.(map[uint]string) {
                    m[String(k)] = v
                }
            // Not a common type, use reflection
            default:
                rv   := reflect.ValueOf(value)
                kind := rv.Kind()
                // If it is a pointer, we should find its real data type.
                if kind == reflect.Ptr {
                    rv   = rv.Elem()
                    kind = rv.Kind()
                }
                switch kind {
                    case reflect.Map:
                        ks := rv.MapKeys()
                        for _, k := range ks {
                            m[String(k.Interface())] = rv.MapIndex(k).Interface()
                        }
                    case reflect.Struct:
                        rt   := rv.Type()
                        name := ""
                        for i := 0; i < rv.NumField(); i++ {
                            // Only convert the public attributes.
                            fieldName := rt.Field(i).Name
                            if !gstr.IsLetterUpper(fieldName[0]) {
                                continue
                            }
                            name = ""
                            // Tag check, supporting "gconv" and "json" tag,
                            // "gconv" has the high priority to use.
                            if len(noTagCheck) == 0 || !noTagCheck[0] {
                                tag := rt.Field(i).Tag
                                if name = tag.Get("gconv"); name == "" {
                                    name = tag.Get("json")
                                }
                            }
                            if name == "" {
                                name = strings.TrimSpace(fieldName)
                            } else {
                                // Support json tag feature: -, omitempty
                                name = strings.TrimSpace(name)
                                if name == "-" {
                                    continue
                                }
                                array := strings.Split(name, ",")
                                if len(array) > 1 {
                                    switch strings.TrimSpace(array[1]) {
                                        case "omitempty":
                                            if empty.IsEmpty(rv.Field(i).Interface()) {
                                                continue
                                            } else {
                                                name = strings.TrimSpace(array[0])
                                            }
                                        default:
                                            name = strings.TrimSpace(array[0])
                                    }
                                }
                            }
                            m[name] = rv.Field(i).Interface()
                        }
                    default:
                        return nil
                }
        }
        return m
    }
}
