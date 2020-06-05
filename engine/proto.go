package engine

import (
	"dynamic-protobuf-json-service/config"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
	"go.uber.org/zap"
)

const (
	protoFileExtn      = ".proto"
	protoNameSeparator = "."
)

// protoMetadata ...
type protoMetadata struct {
	files                   []string
	mapProtoFileDescriptors map[string]*desc.FileDescriptor
}

var metadata *protoMetadata

// Run ...
func Run(cfg *config.Config) error {

	zap.S().Info("Starting protobuf parsing engine...")

	var err error
	var protoDesc = &protoMetadata{
		files:                   []string{},
		mapProtoFileDescriptors: make(map[string]*desc.FileDescriptor),
	}

	err = protoDesc.populateRelativeFilePaths(cfg.ProtoDir)
	if err != nil {
		return err
	}

	protoDesc.files, err = protoparse.ResolveFilenames(nil, protoDesc.files...)
	if err != nil {
		zap.S().Errorf("Failed to resolve proto file names, reason: %v", err.Error())
		return err
	}

	err = protoDesc.buildFileDescriptors(cfg.ProtoDir)
	if err != nil {
		return err
	}

	metadata = protoDesc
	zap.S().Info("Protobuf parsing engine started!")
	return nil
}

// populateRelativeFilePaths ...
func (p *protoMetadata) populateRelativeFilePaths(directory string) error {

	var err error
	err = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			zap.S().Errorf("Walking %v directory failed, reason: %v", path, err.Error())
			return err
		}

		if !info.IsDir() && strings.Contains(info.Name(), protoFileExtn) {
			p.files = append(p.files, path)
		}

		return nil
	})
	if err != nil {
		zap.S().Errorf("Walking proto directory failed, reason: %v", err.Error())
		return err
	}

	return nil
}

// buildFileDescriptors ...
func (p *protoMetadata) buildFileDescriptors(importDir string) error {

	var protoFile string
	var parser = protoparse.Parser{IncludeSourceCodeInfo: true, ImportPaths: []string{importDir}}
	for _, file := range p.files {
		protoFile = strings.Replace(file, importDir, "", 1)
		fileDescriptor, err := parser.ParseFiles(protoFile)
		if err != nil {
			zap.S().Errorf("Failed to parse proto file %v, reason: %v", protoFile, err)
			return err
		}

		p.mapProtoFileDescriptors[protoFile] = fileDescriptor[0]
	}

	return nil
}

// JSONToProtobuf ...
func JSONToProtobuf(fileName string, msgName string, jsonData []byte) ([]byte, error) {

	desc, ok := metadata.mapProtoFileDescriptors[fileName]
	if !ok {
		zap.S().Errorf("Invalid filename supplied: %v", fileName)
		return nil, fmt.Errorf("Invalid filename supplied: %v", fileName)
	}

	msgDesc := desc.FindMessage(desc.GetPackage() + protoNameSeparator + msgName)
	if msgDesc == nil {
		zap.S().Errorf("Invalid message name supplied: %v", msgName)
		return nil, fmt.Errorf("Invalid message name supplied: %v", msgName)
	}

	msg := dynamic.NewMessage(msgDesc)
	err := msg.UnmarshalJSON(jsonData)
	if err != nil {
		zap.S().Errorf("Failed to unmarshal json data to %v, reason: %v", msgDesc.GetName(), err.Error())
		return nil, fmt.Errorf("Failed to unmarshal json data to %v, reason: %v", msgDesc.GetName(), err.Error())
	}

	protoBytes, err := msg.Marshal()
	if err != nil {
		zap.S().Errorf("Failed to marshal json data to proto binary data, reason: %v", err.Error())
		return nil, fmt.Errorf("Failed to marshal json data to proto binary data, reason: %v", err.Error())
	}

	return protoBytes, nil
}

// ProtobufToJSON ...
func ProtobufToJSON(fileName string, msgName string, protoData []byte) ([]byte, error) {

	desc, ok := metadata.mapProtoFileDescriptors[fileName]
	if !ok {
		zap.S().Errorf("Invalid filename supplied: %v", fileName)
		return nil, fmt.Errorf("Invalid filename supplied: %v", fileName)
	}

	msgDesc := desc.FindMessage(desc.GetPackage() + protoNameSeparator + msgName)
	if msgDesc == nil {
		zap.S().Errorf("Invalid message name supplied: %v", msgName)
		return nil, fmt.Errorf("Invalid message name supplied: %v", msgName)
	}

	msg := dynamic.NewMessage(msgDesc)
	err := msg.Unmarshal(protoData)
	if err != nil {
		zap.S().Errorf("Failed to unmarshal proto data to %v, reason: %v", msgDesc.GetName(), err.Error())
		return nil, fmt.Errorf("Failed to unmarshal proto data to %v, reason: %v", msgDesc.GetName(), err.Error())
	}

	jsonBytes, err := msg.MarshalJSON()
	if err != nil {
		zap.S().Errorf("Failed to marshal proto data to json, reason: %v", err.Error())
		return nil, fmt.Errorf("Failed to marshal proto data to json, reason: %v", err.Error())
	}

	return jsonBytes, nil
}
