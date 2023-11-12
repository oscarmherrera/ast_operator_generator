package processors

import (
	"go.uber.org/zap"
	"strings"
)

func MapTypeHandler(fieldType map[string]interface{}, fieldMap map[string]interface{}) {

	fieldNames := fieldMap["Names"].([]interface{})
	nameList := []string{}
	if fieldNames != nil {
		for _, fieldName := range fieldNames {
			nameList = append(nameList, fieldName.(map[string]interface{})["Name"].(string))
		}
	}
	var valuesValue string
	switch fieldType["NodeType"].(string) {
	case "Ident":
		logger.Debug("TypeSpec->TypeSpec->StructType->MapType->Ident", zap.String("type", fieldType["Name"].(string)), zap.Any("names", nameList))
	case "SelectorExpr":
		SelectorExprHandler(fieldType)
	case "StarExpr":
		StarExprHandler(fieldType, fieldMap)

	case "ArrayType":
		ArrayTypeHandler(fieldType, fieldMap)
	case "MapType":
		keys := fieldType["Key"].(map[string]interface{})
		values := fieldType["Value"].(map[string]interface{})
		switch values["NodeType"].(string) {
		case "Ident":
			valuesValue = values["Name"].(string)
		case "SelectorExpr":
			SelectorExprHandler(values)
		case "StarExpr":
			StarExprHandler(values, fieldMap)
		case "ArrayType":
			ArrayTypeHandler(values, fieldMap)
		case "MapType":
			keysValue := keys["Name"].(string)
			if values == nil {
				values = fieldType["Value"].(map[string]interface{})
			}
			valuesValue = values["Name"].(string)
			logger.Debug("TypeSpec->TypeSpec->StructType->MapType", zap.String("key type", keysValue), zap.String("value type", valuesValue), zap.Any("names", nameList))
		}
	default:
		logger.Debug("MapType Unknown", zap.Any("type", fieldType["NodeType"]))
	}

}

func SelectorExprHandler(fieldType map[string]interface{}) {
	//xExp := fieldType["X"].(map[string]interface{})
	//selExp := fieldType["Sel"].(map[string]interface{})
	//
	//xName := xExp["Name"].(string)
	//selName := selExp["Name"].(string)
	starExprType := fieldType["X"].(map[string]interface{})
	switch starExprType["NodeType"].(string) {
	case "Ident":
		logger.Debug("SelectorExpr->Ident", zap.String("type", starExprType["Name"].(string)))
	case "SelectorExpr":
		xExp := fieldType["X"].(map[string]interface{})
		selExp := fieldType["Sel"].(map[string]interface{})

		xName := xExp["Name"].(string)
		selName := selExp["Name"].(string)
		logger.Debug("SelectorExpr", zap.String("type", xName), zap.String("Selector", selName))
	default:
		logger.Debug("SelectorExpr->Unknown", zap.String("type", starExprType["Name"].(string)))
	}

}

func ArrayTypeHandler(values map[string]interface{}, fieldMap map[string]interface{}) {
	arrayType := values["Elt"].(map[string]interface{})
	fieldNames := fieldMap["Names"].([]interface{})
	nameList := []string{}
	if fieldNames != nil {
		for _, fieldName := range fieldNames {
			nameList = append(nameList, fieldName.(map[string]interface{})["Name"].(string))
		}
	}
	logger.Debug("ArrayType", zap.Any("type", arrayType["Name"]), zap.Any("names", nameList))

}

func StarExprHandler(values map[string]interface{}, fieldMap map[string]interface{}) {
	starExprType := values["X"].(map[string]interface{})
	nameList := []string{}
	if fieldMap["Names"] != nil {
		fieldNames := fieldMap["Names"].([]interface{})
		for _, fieldName := range fieldNames {
			nameList = append(nameList, fieldName.(map[string]interface{})["Name"].(string))

		}
	}

	switch starExprType["NodeType"].(string) {
	case "Ident":
		logger.Debug("StarExpr->Ident", zap.String("type", starExprType["Name"].(string)), zap.Any("name", fieldMap["Names"]))
	case "SelectorExpr":
		SelectorExprHandler(starExprType)
	default:
		logger.Debug("StarExpr->Unknown", zap.String("type", starExprType["NodeType"].(string)), zap.Any("name", fieldMap["Names"]))
	}
}

func ValueSpecHandler(spec []interface{}) {
	for _, s := range spec {
		sMap := s.(map[string]interface{})
		sType := sMap["NodeType"].(string)
		switch sType {
		case "ImportSpec":
			if sMap["Name"] == nil {
				path := sMap["Path"].(map[string]interface{})
				value := strings.ReplaceAll(path["Value"].(string), "\"", "")
				logger.Debug("ValueSpec->ImportSpec", zap.Any("name", value))
			}
		case "TypeSpec":
			logger.Debug("TypeSpec", zap.Any("name", sMap["Name"]))
		case "ValueSpec":
			ValueSpecItemHandler(sMap)
		default:
			logger.Debug("Unknown", zap.Any("name", sMap["Name"]))
		}
	}
}

func ValueSpecItemHandler(sMap map[string]interface{}) {
	values := []interface{}{}
	typeMap := map[string]interface{}{}

	if sMap["Names"] != nil {
		names := sMap["Names"].([]interface{})
		if sMap["Values"] != nil {
			values = sMap["Values"].([]interface{})
		}

		if sMap["Type"] != nil {
			typeMap = sMap["Type"].(map[string]interface{})
			switch typeMap["NodeType"].(string) {
			case "Ident":
				for i, name := range names {
					nameMap := name.(map[string]interface{})
					if values != nil && len(values) > 0 {
						valueMap := values[i].(map[string]interface{})
						logger.Debug("ValueSpecItem->Ident", zap.Any("name", nameMap["Name"]), zap.Any("value", valueMap["Value"]), zap.Any("value type", valueMap["Kind"]), zap.Any("type", typeMap["Name"]))
					}
				}
			case "SelectorExpr":
				SelectorExprHandler(typeMap)
			case "StarExpr":
				StarExprHandler(typeMap, sMap)
			case "ArrayType":
				ArrayTypeHandler(typeMap, sMap)
			case "MapType":
				MapTypeHandler(typeMap, sMap)
			case "FuncType":
				logger.Debug("ValueSpecItem->Name", zap.Any("name", names), zap.Any("value", values), zap.Any("type", typeMap["Name"]))
			case "InterfaceType":
				logger.Debug("ValueSpecItem->Name", zap.Any("name", names), zap.Any("value", values), zap.Any("type", typeMap["Name"]))
			case "Ellipsis":
				logger.Debug("ValueSpecItem->Name", zap.Any("name", names), zap.Any("value", values), zap.Any("type", typeMap["Name"]))
			case "ChanType":
				logger.Debug("ValueSpecItem->Name", zap.Any("name", names), zap.Any("value", values), zap.Any("type", typeMap["Name"]))
			case "StructType":
				logger.Debug("ValueSpecItem->Name", zap.Any("name", names), zap.Any("value", values), zap.Any("type", typeMap["Name"]))
			default:
				logger.Debug("ValueSpecItem->Name", zap.Any("name", names), zap.Any("value", values), zap.Any("type", typeMap["Name"]))
			}
			//logger.Debug("ValueSpecItem->Name", zap.Any("name", names), zap.Any("value", values), zap.Any("type", typeMap["Name"]))
			return
		} else {
			for i, name := range names {
				nameMap := name.(map[string]interface{})
				if values != nil && len(values) > 0 {
					valueMap := values[i].(map[string]interface{})
					logger.Debug("ValueSpecItem->Name", zap.Any("name", nameMap["Name"]), zap.Any("value", valueMap["Value"]), zap.Any("type", valueMap["Kind"]))
				} else {
					logger.Debug("ValueSpecItem->Name", zap.Any("name", nameMap["Name"]), zap.Any("value", "nil"), zap.Any("type", "nil"))
				}
			}
		}

	} else {
		logger.Fatal("unexpected ValueSpecItem", zap.Any("name", sMap["Name"]))
	}
}
