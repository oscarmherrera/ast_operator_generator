package ast_json

import (
	"encoding/json"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"strings"
)

type Options struct {
	WithPositions  bool
	WithComments   bool
	WithReferences bool
	WithImports    bool
}

// SourceToJSONWithContent converts the given Go source code to JSON and writes it to the given output file.
// @param input: input file path
// @param path: path of the file
// @param output: output file path
// @param indent: indentation string
// @param options: options for converting the file to JSON
func SourceToJSONWithContent(input *string, path, output string, indent string, options Options) error {
	// Create a new marshaller with the given options
	marshaller := NewMarshaller(options)

	// Set the mode based on options
	mode := parser.AllErrors
	if options.WithComments {
		mode |= parser.ParseComments
	}

	// Parse the file using the marshaller
	tree, err := parser.ParseFile(marshaller.FileSet(), path, strings.NewReader(*input), mode)
	if err != nil {
		return err
	}

	// Marshal the file to a node
	node := marshaller.MarshalFile(tree)

	// Create the output file
	outFile, err := os.Create(output)
	if err != nil {
		return err
	}

	// Create a JSON encoder with the specified indent
	encoder := json.NewEncoder(outFile)
	encoder.SetIndent("", indent)

	// Encode the node to JSON and write it to the output file
	err = encoder.Encode(node)
	if err != nil {
		return err
	}

	// Close the output file
	err = outFile.Close()
	if err != nil {
		return err
	}

	return nil
}

// SourceToJSON converts the given Go source code to JSON and writes it to the given output file.
// @param input: input file path
// @param output: output file path
// @param indent: indentation string
// @param options: options for converting the file to JSON
func SourceToJSON(input, output string, indent string, options Options) error {
	// Create a new marshaller with the given options
	marshaller := NewMarshaller(options)

	// Set the mode based on options
	mode := parser.AllErrors
	if options.WithComments {
		mode |= parser.ParseComments
	}

	// Parse the file using the marshaller
	tree, err := parser.ParseFile(marshaller.FileSet(), input, nil, mode)
	if err != nil {
		return err
	}

	// Marshal the file to a node
	node := marshaller.MarshalFile(tree)

	// Create the output file
	outFile, err := os.Create(output)
	if err != nil {
		return err
	}

	// Create a JSON encoder with the specified indent
	encoder := json.NewEncoder(outFile)
	encoder.SetIndent("", indent)

	// Encode the node to JSON and write it to the output file
	err = encoder.Encode(node)
	if err != nil {
		return err
	}

	// Close the output file
	err = outFile.Close()
	if err != nil {
		return err
	}

	return nil
}

// JSONToSource converts the given JSON file to Go source code and writes it to the given output file.
// @param input: input file path
// @param output: output file path
// @param options: options for converting the file to JSON
func JSONToSource(input, output string, options Options) error {
	// Open the input file
	inFile, err := os.Open(input)
	if err != nil {
		return err
	}

	// Decode the JSON into a FileNode
	var node FileNode
	decoder := json.NewDecoder(inFile)
	err = decoder.Decode(&node)
	if err != nil {
		return err
	}

	// Close the input file
	err = inFile.Close()
	if err != nil {
		return err
	}

	// Create a new unmarshaller with the given options
	unmarshaler := NewUnmarshaller(options)

	// Unmarshal the FileNode to a tree
	tree := unmarshaler.UnmarshalFileNode(&node)

	// Create the output file
	outFile, err := os.Create(output)
	if err != nil {
		return err
	}

	// Print the tree to the output file
	err = printer.Fprint(outFile, unmarshaler.FileSet(), tree)
	if err != nil {
		return err
	}

	// Close the output file
	err = outFile.Close()
	if err != nil {
		return err
	}

	return nil
}

// Loop converts the given Go source code to JSON and writes it to the given output file.
// @param input: input file path
// @param output: output file path
// @param comments: whether to include comments in the output
func Loop(input, output string, comments bool) error {
	// Set the mode based on comments flag
	mode := parser.SkipObjectResolution
	if comments {
		mode |= parser.ParseComments
	}

	// Create a new file set
	fs := token.NewFileSet()

	// Parse the file using the file set and mode
	tree, err := parser.ParseFile(fs, input, nil, mode)
	if err != nil {
		return err
	}

	// Create the output file
	outFile, err := os.Create(output)
	if err != nil {
		return err
	}

	// Print the tree to the output file
	err = printer.Fprint(outFile, fs, tree)
	if err != nil {
		return err
	}

	// Close the output file
	err = outFile.Close()
	if err != nil {
		return err
	}

	return nil
}
