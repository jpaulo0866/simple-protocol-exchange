package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
)

func decompressGzip(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("falha ao criar o leitor gzip: %w", err)
	}
	defer reader.Close()

	var decompressedData bytes.Buffer
	_, err = io.Copy(&decompressedData, reader)
	if err != nil {
		return nil, fmt.Errorf("erro ao descomprimir os dados: %w", err)
	}

	return decompressedData.Bytes(), nil
}

func applyTransformations(data []byte, transform Transform, contentType string) ([]byte, error) {
	var logData GenericData

	switch contentType {
	case "json":
		err := json.Unmarshal(data, &logData)
		if err != nil {
			return nil, fmt.Errorf("falha ao fazer parsing do JSON: %v", err)
		}
	case "text":
		return data, nil
	}

	// Remapeando os campos
	for _, remap := range transform.Remap {
		remapField(logData, remap.Source, remap.Target, remap.PreserveSource)
	}

	// Adicionando campos estáticos
	for key, value := range transform.StaticFields {
		logData[key] = value
	}

	// Removendo campos
	for _, field := range transform.RemoveFields {
		delete(logData, field)
	}

	// Gerando JSON final
	modifiedData, err := json.Marshal(logData)
	if err != nil {
		return nil, fmt.Errorf("falha ao gerar JSON transformado: %v", err)
	}

	return modifiedData, nil
}

func remapField(logData GenericData, sourceField, targetField string, preserveSource bool) {
	if value, ok := logData[sourceField]; ok {
		logData[targetField] = value

		if !preserveSource {
			delete(logData, sourceField)
		}
	}
}
