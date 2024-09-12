package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// Função para garantir que o caminho do arquivo exista e abrir o arquivo em modo append
func openFileWithAppend(filePath string) (*os.File, error) {
	// Cria o diretório se ele não existir
	err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar diretório: %v", err)
	}

	// Abre o arquivo no modo append, cria o arquivo se ele não existir
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir arquivo: %v", err)
	}

	return file, nil
}
