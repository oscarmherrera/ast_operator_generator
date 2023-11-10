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

	createDir(dir + "/funcs")

	for _, file := range files {
		fileName := file.Name()
		err := GetFunctions(dir, fileName, "funcs")
		if err != nil {
			logger.Fatal(err.Error())
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

	//resjq :=
	jq.From("Decls").Where("NodeType", "=", "FuncDecl").Writer(outFile)
	//res, err := resjq.GetR()
	//if err != nil {
	//	return err
	//}
	//resString, err := res..String()
	//if err != nil {
	//	return err
	//}
	//_, err = outFile.WriteString(resString)
	//if err != nil {
	//	return err
	//}

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
