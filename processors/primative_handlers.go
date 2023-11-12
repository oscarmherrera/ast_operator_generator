package processors

import "go.uber.org/zap"

func FuncDeclHandler(declMap map[string]interface{}) {
	var funcName string
	var params []interface{}

	if declMap["Name"] != nil {
		funcName = declMap["Name"].(map[string]interface{})["Name"].(string)
	}

	paramNameList := []string{}
	paramTypeList := []string{}

	if declMap["Type"] != nil {
		p := declMap["Type"].(map[string]interface{})["Params"]
		paramMap := p.(map[string]interface{})
		if paramMap["List"] != nil {
			params = paramMap["List"].([]interface{})
			for _, pm := range params {
				var param map[string]interface{}
				param = pm.(map[string]interface{})
				nameList := param["Names"].([]interface{})
				if param["Type"] != nil {
					typeMap := param["Type"].(map[string]interface{})
					//typeName := typeMap["Name"].(string)
					for _, name := range nameList {
						paramNameList = append(paramNameList, name.(map[string]interface{})["Name"].(string))
						//paramTypeList = append(paramTypeList, typeName)
					}
					switch typeMap["NodeType"].(string) {
					case "Ident":
						typeName := typeMap["Name"].(string)
						for _, name := range nameList {
							paramNameList = append(paramNameList, name.(map[string]interface{})["Name"].(string))
							paramTypeList = append(paramTypeList, typeName)
						}
						logger.Debug("FuncDeclHandler->Ident", zap.Any("name", paramNameList), zap.String("type", typeName))
					case "SelectorExpr":
						SelectorExprHandler(typeMap)
					case "StarExpr":
						StarExprHandler(typeMap, param)
					case "ArrayType":
						ArrayTypeHandler(typeMap, param)
					case "MapType":
						MapTypeHandler(typeMap, param)
					case "FuncType":
						logger.Debug("FuncDeclHandler->FuncType", zap.Any("name", param["Names"]))
					case "InterfaceType":
						logger.Debug("FuncDeclHandler->InterfaceType", zap.Any("name", param["Names"]))
					case "Ellipsis":
						logger.Debug("FuncDeclHandler->Ellipsis", zap.Any("name", param["Names"]))
					case "ChanType":
						logger.Debug("FuncDeclHandler->ChanType", zap.Any("name", param["Names"]))
					case "StructType":
						logger.Debug("FuncDeclHandler->StructType", zap.Any("name", param["Names"]))
					default:
						logger.Debug("FuncDeclHandler->Unknown", zap.Any("name", param["Names"]))
					}

				}
			}
		}
	}
	if declMap["Results"] != nil {
		p := declMap["Results"].(map[string]interface{})["List"]
		paramMap := p.([]interface{})
		for _, param := range paramMap {
			typeMap := param.(map[string]interface{})
			typeName := typeMap["Type"].(map[string]interface{})["Name"].(string)
			paramTypeList = append(paramTypeList, typeName)
		}
	}
	if declMap["Body"] != nil {
		logger.Debug("Function body", zap.Any("name", declMap["Body"]))
	}

	logger.Debug("FuncDecl", zap.String("function name", funcName), zap.Any("params", paramNameList), zap.Any("types", paramTypeList))
}
