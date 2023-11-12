package processors

import (
	"go.uber.org/zap"
	"strings"
)

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
					for _, name := range nameList {
						paramNameList = append(paramNameList, name.(map[string]interface{})["Name"].(string))
					}
					switch typeMap["NodeType"].(string) {
					case "Ident":
						paramNameList, paramTypeList = IdentTypeHandler(typeMap, param)
						logger.Debug("FuncDeclHandler->Ident")
					case "SelectorExpr":
						SelectorExprHandler(typeMap)
						logger.Debug("FuncDeclHandler->SelectorExpr")
					case "StarExpr":
						StarExprHandler(typeMap, param)
						logger.Debug("FuncDeclHandler->StarExpr")
					case "ArrayType":
						ArrayTypeHandler(typeMap, param)
						logger.Debug("FuncDeclHandler->ArrayType")
					case "MapType":
						MapTypeHandler(typeMap, param)
						logger.Debug("FuncDeclHandler->MapType")
					case "FuncType":
						logger.Debug("FuncDeclHandler->FuncType", zap.Any("name", param["Names"]))
					case "InterfaceType":
						logger.Debug("FuncDeclHandler->InterfaceType", zap.Any("name", param["Names"]))
					case "Ellipsis":
						EllipsisTypeHandler(typeMap, param)
						logger.Debug("FuncDeclHandler->Ellipsis")
					case "ChanType":
						logger.Debug("FuncDeclHandler->ChanType", zap.Any("name", param["Names"]))
					case "StructType":
						StructTypeHandler(typeMap)
						logger.Debug("FuncDeclHandler->StructType")
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

func GenDeclHandler(declMap map[string]interface{}) {

	logger.Debug("GenDecl", zap.Any("name", declMap["Tok"]))
	spec := declMap["Specs"].([]interface{})
	for _, s := range spec {
		sMap := s.(map[string]interface{})
		sType := sMap["NodeType"].(string)
		switch sType {
		case "ImportSpec":
			if sMap["Name"] == nil {
				path := sMap["Path"].(map[string]interface{})
				value := strings.ReplaceAll(path["Value"].(string), "\"", "")
				logger.Debug("ImportSpec", zap.Any("name", value))
			}
		case "TypeSpec":
			for _, s := range spec {
				sMap := s.(map[string]interface{})
				sType := sMap["NodeType"].(string)
				switch sType {
				case "ImportSpec":
					if sMap["Name"] == nil {
						path := sMap["Path"].(map[string]interface{})
						value := strings.ReplaceAll(path["Value"].(string), "\"", "")
						logger.Debug("TypeSpec->ImportSpec", zap.Any("name", value))
					}
				case "TypeSpec":
					if sMap["Name"] != nil {
						name := sMap["Name"].(map[string]interface{})
						typeMap := sMap["Type"].(map[string]interface{})
						if typeMap["Name"] != nil {
							logger.Debug("TypeSpec->TypeSpec", zap.String("type", typeMap["Name"].(string)), zap.Any("name", name["Name"]))
						} else {
							// we are dealing with a struct or an interface
							switch typeMap["NodeType"].(string) {
							case "StructType":
								StructTypeHandler(typeMap)
							case "InterfaceType":
								logger.Debug("TypeSpec->TypeSpec", zap.String("type", "interface"), zap.Any("name", name["Name"]))
							default:
								logger.Debug("TypeSpec->TypeSpec Unknown", zap.Any("type", typeMap))
							}
						}
					} else {
						logger.Debug("Uknown TypeSpec", zap.Any("name", sMap))
					}
				case "ValueSpec":
					ValueSpecItemHandler(sMap)

				default:
					logger.Debug("Unknown", zap.Any("name", sMap["Name"]))
				}
			}
			logger.Debug("TypeSpec", zap.Any("name", sMap["Name"]))
		case "ValueSpec":
			ValueSpecHandler(spec)

		default:
			logger.Debug("Unknown", zap.Any("name", sMap["Name"]))
		}
	}
}
