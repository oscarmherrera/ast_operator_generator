package processors

import "go.uber.org/zap"

func FieldListHandler(fieldsMap map[string]interface{}) {
	var fields []interface{}
	if fieldsMap["List"] != nil {
		fields = fieldsMap["List"].([]interface{})
	} else {
		logger.Debug("FieldListHandler->List is nil")
		return
	}
	for _, field := range fields {
		fieldMap := field.(map[string]interface{})
		fieldType := fieldMap["Type"].(map[string]interface{})
		fieldNames := fieldMap["Names"].([]interface{})
		nameList := []string{}
		if fieldNames != nil {
			for _, fieldName := range fieldNames {
				nameList = append(nameList, fieldName.(map[string]interface{})["Name"].(string))
			}
		}
		switch fieldType["NodeType"].(string) {
		case "Ident":
			paramNameList, paramTypeList := IdentTypeHandler(fieldMap, fieldType)
			logger.Debug("FieldListHandler->Ident", zap.Any("names", paramNameList), zap.Any("types", paramTypeList))
		case "SelectorExpr":
			SelectorExprHandler(fieldType)
		case "StarExpr":
			StarExprHandler(fieldType, fieldMap)
		case "ArrayType":
			ArrayTypeHandler(fieldType, fieldMap)
		case "MapType":
			MapTypeHandler(fieldType, fieldMap)
		case "FuncType":
			logger.Debug("FieldListHandler->FuncType", zap.Any("names", nameList))
		case "InterfaceType":
			logger.Debug("FieldListHandler->InterfaceType", zap.Any("names", nameList))
		case "Ellipsis":
			logger.Debug("FieldListHandler->Ellipsis", zap.Any("names", nameList))
		case "ChanType":
			logger.Debug("FieldListHandler->ChanType", zap.Any("names", nameList))
		case "StructType":
			//FieldListHandler(fieldType)
			logger.Debug("FieldListHandler->StructType", zap.Any("names", nameList))
		default:
			logger.Debug("FieldListHandler Unknown", zap.Any("type", fieldType["NodeType"]))
		}
	}
}
