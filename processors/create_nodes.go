package processors

import (
	"github.com/thedevsaddam/gojsonq"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func ReadFiles(dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	//createDir(dir + "/funcs")

	for _, file := range files {
		fileName := file.Name()
		//err := GetFunctions(dir, fileName, "funcs")
		createNodeMap(dir, fileName)
		if err != nil {
			logger.Fatal(err.Error())
		}

	}
	return nil
}

var packageArray map[string][]*Node = map[string][]*Node{}

type Node struct {
	parent   *Node
	children []*Node
}

func createNodeMap(dir string, filename string) error {

	var packageName map[string]interface{}
	var decls []interface{}

	jq := gojsonq.New().File(dir + "/" + filename)
	p := jq.From("Name").Get()
	packageName = p.(map[string]interface{})
	jq.Reset()
	d := jq.From("Decls").Get()
	decls = d.([]interface{})
	logger.Debug("Package name", zap.Any("name", packageName["Name"]))

	for _, decl := range decls {
		declMap := decl.(map[string]interface{})

		declType := declMap["NodeType"].(string)

		switch declType {
		case "FuncDecl":
			FuncDeclHandler(declMap)

			//var funcName string
			//var params []interface{}
			//
			//if declMap["Name"] != nil {
			//	funcName = declMap["Name"].(map[string]interface{})["Name"].(string)
			//}
			//
			//paramNameList := []string{}
			//paramTypeList := []string{}
			//
			//if declMap["Type"] != nil {
			//	p := declMap["Type"].(map[string]interface{})["Params"]
			//	paramMap := p.(map[string]interface{})
			//	if paramMap["List"] != nil {
			//		params = paramMap["List"].([]interface{})
			//		for _, pm := range params {
			//			var param map[string]interface{}
			//			param = pm.(map[string]interface{})
			//			nameList := param["Names"].([]interface{})
			//			if param["Type"] != nil {
			//				typeMap := param["Type"].(map[string]interface{})
			//				//typeName := typeMap["Name"].(string)
			//				for _, name := range nameList {
			//					paramNameList = append(paramNameList, name.(map[string]interface{})["Name"].(string))
			//					//paramTypeList = append(paramTypeList, typeName)
			//				}
			//				switch typeMap["NodeType"].(string) {
			//				case "Ident":
			//					typeName := typeMap["Name"].(string)
			//					for _, name := range nameList {
			//						paramNameList = append(paramNameList, name.(map[string]interface{})["Name"].(string))
			//						paramTypeList = append(paramTypeList, typeName)
			//					}
			//					logger.Debug("FuncDecl->Ident", zap.Any("name", paramNameList), zap.String("type", typeName))
			//				case "SelectorExpr":
			//					SelectorExprHandler(typeMap)
			//				case "StarExpr":
			//					StarExprHandler(typeMap, param)
			//				case "ArrayType":
			//					ArrayTypeHandler(typeMap, param)
			//				case "MapType":
			//					MapTypeHandler(typeMap, param)
			//				case "FuncType":
			//					logger.Debug("FuncDecl->FuncType", zap.Any("name", param["Names"]))
			//				case "InterfaceType":
			//					logger.Debug("FuncDecl->InterfaceType", zap.Any("name", param["Names"]))
			//				case "Ellipsis":
			//					logger.Debug("FuncDecl->Ellipsis", zap.Any("name", param["Names"]))
			//				case "ChanType":
			//					logger.Debug("FuncDecl->ChanType", zap.Any("name", param["Names"]))
			//				case "StructType":
			//					logger.Debug("FuncDecl->StructType", zap.Any("name", param["Names"]))
			//				default:
			//					logger.Debug("FuncDecl->Unknown", zap.Any("name", param["Names"]))
			//				}
			//
			//			}
			//		}
			//	}
			//}
			//if declMap["Results"] != nil {
			//	p := declMap["Results"].(map[string]interface{})["List"]
			//	paramMap := p.([]interface{})
			//	for _, param := range paramMap {
			//		typeMap := param.(map[string]interface{})
			//		typeName := typeMap["Type"].(map[string]interface{})["Name"].(string)
			//		paramTypeList = append(paramTypeList, typeName)
			//	}
			//}
			//if declMap["Body"] != nil {
			//	logger.Debug("Function body", zap.Any("name", declMap["Body"]))
			//}
			//
			//logger.Debug("FuncDecl", zap.String("function name", funcName), zap.Any("params", paramNameList), zap.Any("types", paramTypeList))

		case "GenDecl":
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
										fieldMap := typeMap["Fields"].(map[string]interface{})
										if fieldMap["List"] != nil {
											fields := fieldMap["List"].([]interface{})
											for _, field := range fields {
												fieldMap := field.(map[string]interface{})
												fieldType := fieldMap["Type"].(map[string]interface{})
												switch fieldType["NodeType"].(string) {
												case "Ident":
													logger.Debug("TypeSpec->TypeSpec->StructType->Ident", zap.String("type", fieldType["Name"].(string)), zap.Any("name", fieldMap["Names"]))
												case "SelectorExpr":
													SelectorExprHandler(fieldType)
												case "StarExpr":
													StarExprHandler(fieldType, fieldMap)
												case "ArrayType":
													ArrayTypeHandler(fieldType, fieldMap)
												case "MapType":
													MapTypeHandler(fieldType, fieldMap)
												case "FuncType":
													logger.Debug("TypeSpec->TypeSpec->StructType->FuncType", zap.Any("name", fieldMap["Names"]))
												case "InterfaceType":
													logger.Debug("TypeSpec->TypeSpec->StructType->InterfaceType", zap.Any("name", fieldMap["Names"]))
												case "Ellipsis":
													logger.Debug("TypeSpec->TypeSpec->StructType->Ellipsis", zap.Any("name", fieldMap["Names"]))
												case "ChanType":
													logger.Debug("TypeSpec->TypeSpec->StructType->ChanType", zap.Any("name", fieldMap["Names"]))
												case "StructType":
													fieldNames := fieldMap["Names"].([]interface{})
													if fieldNames != nil {
														nameList := []string{}
														for _, fieldName := range fieldNames {
															nameList = append(nameList, fieldName.(map[string]interface{})["Name"].(string))
														}
														if fieldType["Name"] != nil {
															logger.Debug("TypeSpec->TypeSpec->StructType->Name", zap.String("type", fieldType["Name"].(string)), zap.Any("names", nameList))
														} else {
															logger.Debug("TypeSpec->TypeSpec->StructType Nil Name", zap.String("type", fieldType["NodeType"].(string)), zap.Any("name", fieldMap["Names"]))
														}
													} else {
														logger.Debug("TypeSpec->TypeSpec->StructType Nil Name", zap.String("type", fieldType["NodeType"].(string)), zap.Any("name", fieldMap["Names"]))
													}
												default:
													logger.Debug("TypeSpec->TypeSpec->StructType Unknown", zap.Any("type", fieldType["NodeType"]))
												}
											}
										}
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
							//if sMap["Name"] == nil {
							//	names := sMap["Names"].([]interface{})
							//	for _, name := range names {
							//		nameMap := name.(map[string]interface{})
							//		logger.Debug("TypeSpec->ValueSpec->Name", zap.Any("name", nameMap["Name"]))
							//	}
							//}
						default:
							logger.Debug("Unknown", zap.Any("name", sMap["Name"]))
						}
					}
					logger.Debug("TypeSpec", zap.Any("name", sMap["Name"]))
				case "ValueSpec":
					ValueSpecHandler(spec)
					//for _, s := range spec {
					//	sMap := s.(map[string]interface{})
					//	sType := sMap["NodeType"].(string)
					//	switch sType {
					//	case "ImportSpec":
					//		if sMap["Name"] == nil {
					//			path := sMap["Path"].(map[string]interface{})
					//			value := strings.ReplaceAll(path["Value"].(string), "\"", "")
					//			logger.Debug("ValueSpec->ImportSpec", zap.Any("name", value))
					//		}
					//	case "TypeSpec":
					//		logger.Debug("TypeSpec", zap.Any("name", sMap["Name"]))
					//	case "ValueSpec":
					//		if sMap["Name"] == nil {
					//			names := sMap["Names"].([]interface{})
					//			for _, name := range names {
					//				nameMap := name.(map[string]interface{})
					//				logger.Debug("ValueSpec->ValueSpec->Name", zap.Any("name", nameMap["Name"]))
					//			}
					//		}
					//	default:
					//		logger.Debug("Unknown", zap.Any("name", sMap["Name"]))
					//	}
					//}
				default:
					logger.Debug("Unknown", zap.Any("name", sMap["Name"]))
				}
			}
		case "ImportSpec":
			logger.Debug("ImportSpec", zap.Any("name", declMap["Name"]))
		case "TypeSpec":
			logger.Debug("TypeSpec", zap.Any("name", declMap["Name"]))
		case "ValueSpec":
			spec := declMap["Specs"].([]interface{})
			ValueSpecHandler(spec)
			//for _, s := range spec {
			//	sMap := s.(map[string]interface{})
			//	sType := sMap["NodeType"].(string)
			//	switch sType {
			//	case "ImportSpec":
			//		if sMap["Name"] == nil {
			//			path := sMap["Path"].(map[string]interface{})
			//			value := strings.ReplaceAll(path["Value"].(string), "\"", "")
			//			logger.Debug("ValueSpec->ImportSpec", zap.Any("name", value))
			//		}
			//	case "TypeSpec":
			//		logger.Debug("TypeSpec", zap.Any("name", sMap["Name"]))
			//	case "ValueSpec":
			//		if sMap["Name"] == nil {
			//			names := sMap["Names"].([]interface{})
			//			for _, name := range names {
			//				nameMap := name.(map[string]interface{})
			//				logger.Debug("ValueSpec->ValueSpec->Name", zap.Any("name", nameMap["Name"]))
			//			}
			//		}
			//	default:
			//		logger.Debug("Unknown", zap.Any("name", sMap["Name"]))
			//	}
			//}
			logger.Debug("ValueSpec", zap.Any("name", declMap["Name"]))
		default:
			logger.Debug("Unknown", zap.Any("name", declMap["Name"]))
		}

	}

	return nil
}

func GetFunctions(dir string, filename string, tempLocation string) error {
	jq := gojsonq.New().File(dir + "/" + filename)

	outFile, err := os.Create(dir + "/" + tempLocation + "/" + filename + "_functions.json")
	if err != nil {
		return err
	}
	defer outFile.Close()

	resjq := jq.From("Decls").Where("NodeType", "=", "FuncDecl")
	if resjq.Count() > 0 {
		resjq.Writer(outFile)
	} else {
		logger.Info("No functions found in file", zap.String("file", filename))
	}

	return nil
}

func createDir(directory string) {
	// Check if the directory already exists
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		// Create the directory recursively
		err := os.MkdirAll(directory, os.ModePerm)
		if err != nil {
			logger.Error("Error creating directory", zap.Error(err))
			return
		}
		logger.Debug("Directory created successfully.")
	} else {
		logger.Debug("Directory already exists.")
	}
}
