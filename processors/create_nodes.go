package processors

import (
	"github.com/thedevsaddam/gojsonq"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"os"
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
		case "GenDecl":
			GenDeclHandler(declMap)
		case "ImportSpec":
			logger.Debug("ImportSpec", zap.Any("name", declMap["Name"]))
		case "TypeSpec":
			logger.Debug("TypeSpec", zap.Any("name", declMap["Name"]))
		case "ValueSpec":
			spec := declMap["Specs"].([]interface{})
			ValueSpecHandler(spec)

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
