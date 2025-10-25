package common

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type ConditionOperation struct {
	Path           string      `json:"path"`             
	Mode           string      `json:"mode"`             
	Value          interface{} `json:"value"`            
	Invert         bool        `json:"invert"`           
	PassMissingKey bool        `json:"pass_missing_key"` 
}

type ParamOperation struct {
	Path       string               `json:"path"`
	Mode       string               `json:"mode"` 
	Value      interface{}          `json:"value"`
	KeepOrigin bool                 `json:"keep_origin"`
	From       string               `json:"from,omitempty"`
	To         string               `json:"to,omitempty"`
	Conditions []ConditionOperation `json:"conditions,omitempty"` 
	Logic      string               `json:"logic,omitempty"`      
}

func ApplyParamOverride(jsonData []byte, paramOverride map[string]interface{}) ([]byte, error) {
	if len(paramOverride) == 0 {
		return jsonData, nil
	}

	
	if operations, ok := tryParseOperations(paramOverride); ok {
		
		result, err := applyOperations(string(jsonData), operations)
		return []byte(result), err
	}

	
	return applyOperationsLegacy(jsonData, paramOverride)
}

func tryParseOperations(paramOverride map[string]interface{}) ([]ParamOperation, bool) {
	
	if opsValue, exists := paramOverride["operations"]; exists {
		if opsSlice, ok := opsValue.([]interface{}); ok {
			var operations []ParamOperation
			for _, op := range opsSlice {
				if opMap, ok := op.(map[string]interface{}); ok {
					operation := ParamOperation{}

					
					if path, ok := opMap["path"].(string); ok {
						operation.Path = path
					}
					if mode, ok := opMap["mode"].(string); ok {
						operation.Mode = mode
					} else {
						return nil, false 
					}

					
					if value, exists := opMap["value"]; exists {
						operation.Value = value
					}
					if keepOrigin, ok := opMap["keep_origin"].(bool); ok {
						operation.KeepOrigin = keepOrigin
					}
					if from, ok := opMap["from"].(string); ok {
						operation.From = from
					}
					if to, ok := opMap["to"].(string); ok {
						operation.To = to
					}
					if logic, ok := opMap["logic"].(string); ok {
						operation.Logic = logic
					} else {
						operation.Logic = "OR" 
					}

					
					if conditions, exists := opMap["conditions"]; exists {
						if condSlice, ok := conditions.([]interface{}); ok {
							for _, cond := range condSlice {
								if condMap, ok := cond.(map[string]interface{}); ok {
									condition := ConditionOperation{}
									if path, ok := condMap["path"].(string); ok {
										condition.Path = path
									}
									if mode, ok := condMap["mode"].(string); ok {
										condition.Mode = mode
									}
									if value, ok := condMap["value"]; ok {
										condition.Value = value
									}
									if invert, ok := condMap["invert"].(bool); ok {
										condition.Invert = invert
									}
									if passMissingKey, ok := condMap["pass_missing_key"].(bool); ok {
										condition.PassMissingKey = passMissingKey
									}
									operation.Conditions = append(operation.Conditions, condition)
								}
							}
						}
					}

					operations = append(operations, operation)
				} else {
					return nil, false
				}
			}
			return operations, true
		}
	}

	return nil, false
}

func checkConditions(jsonStr string, conditions []ConditionOperation, logic string) (bool, error) {
	if len(conditions) == 0 {
		return true, nil 
	}
	results := make([]bool, len(conditions))
	for i, condition := range conditions {
		result, err := checkSingleCondition(jsonStr, condition)
		if err != nil {
			return false, err
		}
		results[i] = result
	}

	if strings.ToUpper(logic) == "AND" {
		for _, result := range results {
			if !result {
				return false, nil
			}
		}
		return true, nil
	} else {
		for _, result := range results {
			if result {
				return true, nil
			}
		}
		return false, nil
	}
}

func checkSingleCondition(jsonStr string, condition ConditionOperation) (bool, error) {
	
	path := processNegativeIndex(jsonStr, condition.Path)
	value := gjson.Get(jsonStr, path)
	if !value.Exists() {
		if condition.PassMissingKey {
			return true, nil
		}
		return false, nil
	}

	
	targetBytes, err := json.Marshal(condition.Value)
	if err != nil {
		return false, fmt.Errorf("failed to marshal condition value: %v", err)
	}
	targetValue := gjson.ParseBytes(targetBytes)

	result, err := compareGjsonValues(value, targetValue, strings.ToLower(condition.Mode))
	if err != nil {
		return false, fmt.Errorf("comparison failed for path %s: %v", condition.Path, err)
	}

	if condition.Invert {
		result = !result
	}
	return result, nil
}

func processNegativeIndex(jsonStr string, path string) string {
	re := regexp.MustCompile(`\.(-\d+)`)
	matches := re.FindAllStringSubmatch(path, -1)

	if len(matches) == 0 {
		return path
	}

	result := path
	for _, match := range matches {
		negIndex := match[1]
		index, _ := strconv.Atoi(negIndex)

		arrayPath := strings.Split(path, negIndex)[0]
		if strings.HasSuffix(arrayPath, ".") {
			arrayPath = arrayPath[:len(arrayPath)-1]
		}

		array := gjson.Get(jsonStr, arrayPath)
		if array.IsArray() {
			length := len(array.Array())
			actualIndex := length + index
			if actualIndex >= 0 && actualIndex < length {
				result = strings.Replace(result, match[0], "."+strconv.Itoa(actualIndex), 1)
			}
		}
	}

	return result
}


func compareGjsonValues(jsonValue, targetValue gjson.Result, mode string) (bool, error) {
	switch mode {
	case "full":
		return compareEqual(jsonValue, targetValue)
	case "prefix":
		return strings.HasPrefix(jsonValue.String(), targetValue.String()), nil
	case "suffix":
		return strings.HasSuffix(jsonValue.String(), targetValue.String()), nil
	case "contains":
		return strings.Contains(jsonValue.String(), targetValue.String()), nil
	case "gt":
		return compareNumeric(jsonValue, targetValue, "gt")
	case "gte":
		return compareNumeric(jsonValue, targetValue, "gte")
	case "lt":
		return compareNumeric(jsonValue, targetValue, "lt")
	case "lte":
		return compareNumeric(jsonValue, targetValue, "lte")
	default:
		return false, fmt.Errorf("unsupported comparison mode: %s", mode)
	}
}

func compareEqual(jsonValue, targetValue gjson.Result) (bool, error) {
	
	if (jsonValue.Type == gjson.True || jsonValue.Type == gjson.False) &&
		(targetValue.Type == gjson.True || targetValue.Type == gjson.False) {
		return jsonValue.Bool() == targetValue.Bool(), nil
	}

	
	if jsonValue.Type != targetValue.Type {
		return false, fmt.Errorf("compare for different types, got %v and %v", jsonValue.Type, targetValue.Type)
	}

	switch jsonValue.Type {
	case gjson.True, gjson.False:
		return jsonValue.Bool() == targetValue.Bool(), nil
	case gjson.Number:
		return jsonValue.Num == targetValue.Num, nil
	case gjson.String:
		return jsonValue.String() == targetValue.String(), nil
	default:
		return jsonValue.String() == targetValue.String(), nil
	}
}

func compareNumeric(jsonValue, targetValue gjson.Result, operator string) (bool, error) {
	
	if jsonValue.Type != gjson.Number || targetValue.Type != gjson.Number {
		return false, fmt.Errorf("numeric comparison requires both values to be numbers, got %v and %v", jsonValue.Type, targetValue.Type)
	}

	jsonNum := jsonValue.Num
	targetNum := targetValue.Num

	switch operator {
	case "gt":
		return jsonNum > targetNum, nil
	case "gte":
		return jsonNum >= targetNum, nil
	case "lt":
		return jsonNum < targetNum, nil
	case "lte":
		return jsonNum <= targetNum, nil
	default:
		return false, fmt.Errorf("unsupported numeric operator: %s", operator)
	}
}


func applyOperationsLegacy(jsonData []byte, paramOverride map[string]interface{}) ([]byte, error) {
	reqMap := make(map[string]interface{})
	err := json.Unmarshal(jsonData, &reqMap)
	if err != nil {
		return nil, err
	}

	for key, value := range paramOverride {
		reqMap[key] = value
	}

	return json.Marshal(reqMap)
}

func applyOperations(jsonStr string, operations []ParamOperation) (string, error) {
	result := jsonStr
	for _, op := range operations {
		
		ok, err := checkConditions(result, op.Conditions, op.Logic)
		if err != nil {
			return "", err
		}
		if !ok {
			continue 
		}
		
		opPath := processNegativeIndex(result, op.Path)
		opFrom := processNegativeIndex(result, op.From)
		opTo := processNegativeIndex(result, op.To)

		switch op.Mode {
		case "delete":
			result, err = sjson.Delete(result, opPath)
		case "set":
			if op.KeepOrigin && gjson.Get(result, opPath).Exists() {
				continue
			}
			result, err = sjson.Set(result, opPath, op.Value)
		case "move":
			result, err = moveValue(result, opFrom, opTo)
		case "prepend":
			result, err = modifyValue(result, opPath, op.Value, op.KeepOrigin, true)
		case "append":
			result, err = modifyValue(result, opPath, op.Value, op.KeepOrigin, false)
		default:
			return "", fmt.Errorf("unknown operation: %s", op.Mode)
		}
		if err != nil {
			return "", fmt.Errorf("operation %s failed: %v", op.Mode, err)
		}
	}
	return result, nil
}

func moveValue(jsonStr, fromPath, toPath string) (string, error) {
	sourceValue := gjson.Get(jsonStr, fromPath)
	if !sourceValue.Exists() {
		return jsonStr, fmt.Errorf("source path does not exist: %s", fromPath)
	}
	result, err := sjson.Set(jsonStr, toPath, sourceValue.Value())
	if err != nil {
		return "", err
	}
	return sjson.Delete(result, fromPath)
}

func modifyValue(jsonStr, path string, value interface{}, keepOrigin, isPrepend bool) (string, error) {
	current := gjson.Get(jsonStr, path)
	switch {
	case current.IsArray():
		return modifyArray(jsonStr, path, value, isPrepend)
	case current.Type == gjson.String:
		return modifyString(jsonStr, path, value, isPrepend)
	case current.Type == gjson.JSON:
		return mergeObjects(jsonStr, path, value, keepOrigin)
	}
	return jsonStr, fmt.Errorf("operation not supported for type: %v", current.Type)
}

func modifyArray(jsonStr, path string, value interface{}, isPrepend bool) (string, error) {
	current := gjson.Get(jsonStr, path)
	var newArray []interface{}
	
	addValue := func() {
		if arr, ok := value.([]interface{}); ok {
			newArray = append(newArray, arr...)
		} else {
			newArray = append(newArray, value)
		}
	}
	
	addOriginal := func() {
		current.ForEach(func(_, val gjson.Result) bool {
			newArray = append(newArray, val.Value())
			return true
		})
	}
	if isPrepend {
		addValue()
		addOriginal()
	} else {
		addOriginal()
		addValue()
	}
	return sjson.Set(jsonStr, path, newArray)
}

func modifyString(jsonStr, path string, value interface{}, isPrepend bool) (string, error) {
	current := gjson.Get(jsonStr, path)
	valueStr := fmt.Sprintf("%v", value)
	var newStr string
	if isPrepend {
		newStr = valueStr + current.String()
	} else {
		newStr = current.String() + valueStr
	}
	return sjson.Set(jsonStr, path, newStr)
}

func mergeObjects(jsonStr, path string, value interface{}, keepOrigin bool) (string, error) {
	current := gjson.Get(jsonStr, path)
	var currentMap, newMap map[string]interface{}

	
	if err := json.Unmarshal([]byte(current.Raw), &currentMap); err != nil {
		return "", err
	}
	
	switch v := value.(type) {
	case map[string]interface{}:
		newMap = v
	default:
		jsonBytes, _ := json.Marshal(v)
		if err := json.Unmarshal(jsonBytes, &newMap); err != nil {
			return "", err
		}
	}
	
	result := make(map[string]interface{})
	for k, v := range currentMap {
		result[k] = v
	}
	for k, v := range newMap {
		if !keepOrigin || result[k] == nil {
			result[k] = v
		}
	}
	return sjson.Set(jsonStr, path, result)
}
