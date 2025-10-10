package request

import (
	"fmt"
	"reflect"
	"strings"
)

// AssignSliceToStructMap 将切片的值一一对应赋值到结构体字段并返回 map[string]interface{}
// structObj结构体对象，sliceObj 切片对象
// 含有 Authorization
func AssignSliceToStructMap(structObj interface{}, sliceObj []interface{}) (map[string]interface{}, error) {
	// 初始化结果 map
	result := make(map[string]interface{})

	// 检查结构体是否为指针
	structVal := reflect.ValueOf(structObj)
	if structVal.Kind() != reflect.Ptr || structVal.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("first parameter must be a pointer to a struct")
	}
	structVal = structVal.Elem()
	structType := structVal.Type()

	// 检查切片是否有效
	sliceVal := reflect.ValueOf(sliceObj)
	if sliceVal.Kind() != reflect.Slice {
		return nil, fmt.Errorf("second parameter must be a slice")
	}

	// 检查切片长度是否与结构体字段数量匹配
	numFields := structVal.NumField()
	if sliceVal.Len() < numFields {
		return nil, fmt.Errorf("slice length (%d) is less than struct field count (%d)", sliceVal.Len(), numFields)
	}

	// 将切片的值赋值给结构体字段
	for i := 0; i < numFields; i++ {
		field := structVal.Field(i)
		fieldName := structType.Field(i).Name
		sliceElement := sliceVal.Index(i)

		// 检查字段是否可设置
		if !field.CanSet() {
			return nil, fmt.Errorf("cannot set field %s", fieldName)
		}

		// 处理 Authorization 字段
		if fieldName == "Authorization" {
			// 尝试将切片元素转换为字符串
			var bearerValue string
			if sliceElement.Kind() == reflect.String {
				bearerValue = "Bearer " + sliceElement.String()
			} else {
				// 尝试将元素转换为字符串（支持常见类型）
				if sliceElement.CanInterface() {
					bearerValue = fmt.Sprintf("Bearer %v", sliceElement.Interface())
				} else {
					return nil, fmt.Errorf("slice element for Authorization must be convertible to string, got %v", sliceElement.Type())
				}
			}

			// 赋值给字段（任意类型支持）
			if field.Type().Kind() == reflect.Interface || field.Type() == reflect.TypeOf("") {
				field.Set(reflect.ValueOf(bearerValue))
				result[fieldName] = bearerValue
			} else {
				return nil, fmt.Errorf("Authorization field must be string or interface{} type, got %v", field.Type())
			}
		} else {
			// 其他字段的赋值
			if field.Type().Kind() == reflect.Interface || sliceElement.Type().AssignableTo(field.Type()) {
				field.Set(sliceElement)
				result[fieldName] = sliceElement.Interface()
			} else {
				return nil, fmt.Errorf("cannot assign slice element type %v to field %s of type %v",
					sliceElement.Type(), fieldName, field.Type())
			}
		}
	}

	return result, nil
}

// 初始化结构体，并且返回map
func InitStructToMap(strct interface{}, values []interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	v := reflect.ValueOf(strct).Elem() // 获取结构体值
	t := v.Type()

	for i := 0; i < v.NumField() && i < len(values); i++ {
		field := v.Field(i)

		// 处理字段可设置情况
		if field.CanSet() {
			val := reflect.ValueOf(values[i])

			// 类型不一致时尝试转换
			if val.Type().ConvertibleTo(field.Type()) {
				field.Set(val.Convert(field.Type()))
			}
		}

		// 优先用 JSON tag 作为 map key，否则用字段名
		tag := t.Field(i).Tag.Get("json")
		if tag == "" {
			tag = t.Field(i).Name
		}
		result[tag] = v.Field(i).Interface()
	}

	return result, nil
}

// StructToMap 将结构体初始化并将切片值映射到 map   // 可以解决嵌套结构体
func StructToMap(structType interface{}, slice []interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// 获取结构体类型
	val := reflect.ValueOf(structType)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("structType must be a struct or struct pointer")
	}

	t := val.Type()
	sliceIndex := 0

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldName := field.Name
		fieldType := field.Type

		// 检查字段是否可导出
		if field.PkgPath != "" {
			fmt.Printf("Skipping unexported field %s\n", fieldName)
			continue
		}

		// 获取 JSON 标签中的字段名
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			jsonTag = fieldName // 如果没有 JSON 标签，使用字段名
		} else {
			// 提取 JSON 标签中的字段名（忽略其他选项，如 ",omitempty"）
			if commaIdx := strings.Index(jsonTag, ","); commaIdx != -1 {
				jsonTag = jsonTag[:commaIdx]
			}
		}

		// 处理嵌套结构体
		if fieldType.Kind() == reflect.Struct {
			if sliceIndex >= len(slice) {
				fmt.Printf("Slice exhausted at nested field %s, assigning nil\n", fieldName)
				result[jsonTag] = nil
				continue
			}
			nestedMap, err := StructToMap(reflect.New(fieldType).Interface(), slice[sliceIndex:])
			if err != nil {
				return nil, fmt.Errorf("error in nested struct %s: %v", fieldName, err)
			}
			result[jsonTag] = nestedMap
			continue
		}

		// 处理基本类型字段
		if sliceIndex >= len(slice) {
			fmt.Printf("Slice exhausted at field %s, assigning nil\n", fieldName)
			result[jsonTag] = nil
			continue
		}

		sliceVal := reflect.ValueOf(slice[sliceIndex])
		if !sliceVal.IsValid() {
			fmt.Printf("Nil slice value at index %d for field %s\n", sliceIndex, fieldName)
			result[jsonTag] = nil
			sliceIndex++
			continue
		}

		// 检查类型兼容性
		if sliceVal.Type().ConvertibleTo(fieldType) {
			result[jsonTag] = sliceVal.Convert(fieldType).Interface()
		} else {
			return nil, fmt.Errorf("cannot assign %v to field %s of type %v", sliceVal.Type(), fieldName, fieldType)
		}
		sliceIndex++
	}

	return result, nil
}

// FlattenMap 将嵌套的 map[string]interface{} 平铺为一层 map，忽略嵌套路径
func FlattenMap(nestedMap map[string]interface{}) map[string]interface{} {
	flatMap := make(map[string]interface{})

	for key, value := range nestedMap {
		// 如果值是嵌套 map，递归平铺
		if nested, ok := value.(map[string]interface{}); ok {
			// 将嵌套 map 的键值对直接合并到 flatMap
			for nestedKey, nestedValue := range FlattenMap(nested) {
				flatMap[nestedKey] = nestedValue // 后覆盖前
				// fmt.Printf("Flattened key %s with value %v\n", nestedKey, nestedValue)
			}
		} else {
			// 直接赋值非 map 值
			flatMap[key] = value
			// fmt.Printf("Flattened key %s with value %v\n", key, value)
		}
	}

	return flatMap
}
