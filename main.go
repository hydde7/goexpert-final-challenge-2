// main.go
package main

import (
	"flag"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var doneCount int64
	var successCount int64

	url := flag.String("url", "", "URL do serviço a ser testado")
	totalReq := flag.Int("requests", 0, "Número total de requests")
	concurrency := flag.Int("concurrency", 1, "Número de chamadas simultâneas")
	flag.Parse()

	if *url == "" || *totalReq <= 0 || *concurrency <= 0 {
		log.Fatalf("Uso incorreto: é obrigatório informar --url, --requests (>0) e --concurrency (>0)")
	}

	log.Infof("Iniciando teste de carga: url=%s | total_requests=%d | concurrency=%d",
		*url, *totalReq, *concurrency)

	start := time.Now()
	jobs := make(chan struct{}, *totalReq)
	statusCounts := make(map[int]int)

	wg.Add(*concurrency)
	for i := range *concurrency {
		go func(id int) {
			defer wg.Done()
			for range jobs {
				atomic.AddInt64(&doneCount, 1)

				resp, err := http.Get(*url)
				if err != nil {
					mu.Lock()
					statusCounts[-1]++
					mu.Unlock()
					log.WithFields(log.Fields{
						"worker": id,
						"error":  err,
					}).Warn("falha na requisição HTTP")
					continue
				}
				code := resp.StatusCode
				resp.Body.Close()

				if code == http.StatusOK {
					atomic.AddInt64(&successCount, 1)
				}

				mu.Lock()
				statusCounts[code]++
				mu.Unlock()
			}
		}(i + 1)
	}

	for range *totalReq {
		jobs <- struct{}{}
	}
	close(jobs)

	wg.Wait()
	elapsed := time.Since(start)

	log.Info("======== Relatório de Carga ========")
	log.Infof("URL testada:               %s", *url)
	log.Infof("Tempo total de execução:   %v", elapsed)
	log.Infof("Total de requests feitos:  %d", atomic.LoadInt64(&doneCount))
	log.Infof("Requests bem-sucedidos:    %d (200)", atomic.LoadInt64(&successCount))
	log.Info("Distribuição de códigos HTTP:")
	mu.Lock()
	for code, cnt := range statusCounts {
		label := http.StatusText(code)
		if code == -1 {
			label = "erro de conexão"
		}
		log.Infof("  %s (%d): %d", label, code, cnt)
	}
	mu.Unlock()
}
