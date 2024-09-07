package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	configs, err := parseConfig("config.yaml")
	if err != nil {
		log.Fatalf("falha ao fazer parsing do arquivo de configuração: %v", err)
	}

	var wg sync.WaitGroup
	for _, route := range configs.Routes {
		fmt.Printf("Gateway: %s\n", route.Name)
		fmt.Printf("Rota: %s\n", route.Protocol)
		fmt.Printf("Porta: %d\n", route.Port)
		wg.Add(1)
		switch route.Protocol {
		case "http":
			go configureHttp(&route, &wg)
		case "tcp":
			go configureTcp(&route, &wg)
		case "udp":
			go configureUdp(&route, &wg)
		default:
			log.Fatalf("protocolo inválido: %s", route.Protocol)
		}
	}
	wg.Wait()
}

func configureHttp(route *Route, waitGroup *sync.WaitGroup) {
	// panic("unimplemented")
}

func configureTcp(route *Route, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	var portTcp = strconv.Itoa(route.Port)
	listener, err := net.Listen("tcp", ":"+portTcp)
	if err != nil {
		log.Fatalf("Erro ao iniciar servidor TCP na porta %s: %v", portTcp, err)
	}
	defer listener.Close()

	log.Printf("Servidor TCP rodando na porta %s", portTcp)

	var wg sync.WaitGroup
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Erro ao aceitar conexão TCP: %v", err)
			continue
		}

		clientAddr := conn.RemoteAddr().String()
		log.Printf("Cliente TCP conectado: %s", clientAddr)

		wg.Add(1)
		go handleTcpConnection(conn, &wg, clientAddr, route)
	}
	wg.Wait()
}

func configureUdp(route *Route, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	var port = route.Port

	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP("0.0.0.0"),
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Fatalf("Erro ao iniciar servidor UDP na porta %d: %v", port, err)
	}
	defer conn.Close()

	log.Printf("Servidor UDP rodando na porta %d", port)

	var wg sync.WaitGroup
	for {
		buffer := make([]byte, 65536)
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Erro ao ler dados UDP: %v", err)
			continue
		}

		var data string

		if route.Entry.Compressed {
			decompressedData, err := decompressGzip(buffer[:n])
			if err != nil {
				log.Printf("Erro ao descomprimir os dados: %v", err)
				continue
			}
			data = strings.TrimSpace(string(decompressedData))
		} else {
			data = strings.TrimSpace(string(buffer[:n]))
		}

		wg.Add(1)
		go handleUDPConnection([]byte(data), clientAddr, route, &wg)
	}
	wg.Wait()
}

func handleUDPConnection(data []byte, clientAddr *net.UDPAddr, route *Route, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	log.Printf("Cliente UDP %s enviou %d bytes", clientAddr.String(), len(data))

	// Verifica se o JSON está completo e processa
	if isCompleteMessage(data, route.Entry.ContentType) {
		processOutput(data, route)
	} else {
		log.Printf("Mensagem UDP incompleta, %q", string(data))
	}
}

func handleTcpConnection(conn net.Conn, waitGroup *sync.WaitGroup, clientAddr string, route *Route) {
	defer conn.Close()
	defer waitGroup.Done()

	reader := bufio.NewReader(conn)
	var buffer bytes.Buffer

	for {
		var buff = make([]byte, 65732)
		n, err := reader.Read(buff)
		if err != nil {
			if err == io.EOF {
				log.Printf("Cliente TCP %s desconectado", clientAddr)
			} else {
				log.Printf("Erro ao ler dados do cliente TCP %s: %v", clientAddr, err)
			}
			break
		}

		var data = strings.TrimSpace(string(buff[:n]))

		if route.Entry.Compressed {
			tempBuff, err := decompressGzip([]byte(data))
			if err != nil {
				log.Printf("Erro ao descomprimir dados %v", err)
			}
			data = string(tempBuff)
		}

		data = strings.TrimSpace(string(buff[:n]))
		data = strings.TrimSuffix(data, "\x00")

		log.Printf("Cliente TCP %s enviou %d bytes", clientAddr, n)

		buffer.Write([]byte(data))

		if isCompleteMessage(buffer.Bytes(), route.Entry.ContentType) {
			processRequest(buffer.Bytes(), route)
			buffer.Reset()
		} else {
			log.Printf("Mensagem TCP incompleta, %q", string(data))
		}
	}

}

func processRequest(data []byte, route *Route) {
	modifiedData, err := applyTransformations(data, route.Transform, route.Entry.ContentType)
	if err != nil {
		log.Printf("Erro ao aplicar transformação de dados %v", err)
	}

	processOutput(modifiedData, route)
}

func processOutput(modifiedData []byte, route *Route) {
	switch route.Output.Protocol {
	case "tcp":
		handleTcpOutput(modifiedData, route)
	case "udp":
		handleUdpOutput(modifiedData, route)
	case "http", "https":
		handleHttpOutput(modifiedData, route)
	}
}

func handleHttpOutput(modifiedData []byte, route *Route) {
	client := &http.Client{
		Timeout: time.Duration(route.Output.Timeout) * time.Second,
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s://%s:%d%s",
		route.Output.Protocol,
		route.Output.Host,
		route.Output.Port,
		route.Output.Path), bytes.NewBuffer(modifiedData))
	if err != nil {
		log.Printf("Erro ao criar requisição: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range route.Output.Headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("erro ao enviar requisição: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		log.Printf("resposta inválida do output: %s", resp.Status)
	}

	log.Println("Dados enviados para o destino com sucesso")
}

func handleUdpOutput(modifiedData []byte, route *Route) {
	// panic("unimplemented")
}

func handleTcpOutput(modifiedData []byte, route *Route) {
	// panic("unimplemented")
}

func isCompleteMessage(data []byte, contentType string) bool {
	switch contentType {
	case "json":
		return bytes.HasSuffix(data, []byte("}"))
	case "text":
		return bytes.HasSuffix(data, []byte("\n"))
	default:
		return true
	}
}
