package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type Result struct {
	Acao  string `json:"acao"`
	Valor int    `json:"valor"`
}

var geminiModel *genai.GenerativeModel
var lastHash string

func Init() {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Println("⚠ GEMINI_API_KEY não setado - rodando sem IA, só log")
		return
	}
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}
	geminiModel = client.GenerativeModel("gemini-1.5-flash")
	// geminiModel.GenerationConfig = genai.GenerationConfig{
	// 	ResponseMIMEType: "application/json",
	// }

	InitHistory()
	log.Println("Gemini inicializado")
}

func ProcessImage(img []byte) Result {
	hash := fmt.Sprintf("%d", len(img))
	if hash == lastHash {
		return Result{Acao: "NADA", Valor: 0}
	}
	lastHash = hash

	if geminiModel == nil {
		Broadcast("[" + time.Now().Format("15:04:05") + "] Sem IA")
		return Result{Acao: "NADA", Valor: 0}
	}

	ctx := context.Background()
	prompt := `Analise esta imagem de partida de poker.
Decida: PAGAR, PASSAR, CORRER, SUBIR.
Avalie banca minha e adversários, probabilidade de vitória.
Responda SOMENTE JSON válido: {"acao":"SUBIR","valor":50}
Se acao não for SUBIR, valor = 0. Acao sempre em MAIUSCULO.`

	resp, err := geminiModel.GenerateContent(ctx,
		genai.Text(prompt),
		genai.ImageData("png", img),
	)
	if err != nil {
		log.Println(err)
		return Result{Acao: "NADA", Valor: 0}
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return Result{Acao: "NADA", Valor: 0}
	}

	// PEGA TEXTO CORRETO
	var txt string
	switch p := resp.Candidates[0].Content.Parts[0].(type) {
	case genai.Text:
		txt = string(p)
	default:
		txt = fmt.Sprintf("%v", p)
	}

	// LIMPA ```json... ```
	txt = strings.TrimSpace(txt)
	txt = strings.TrimPrefix(txt, "```json")
	txt = strings.TrimPrefix(txt, "```")
	txt = strings.TrimSuffix(txt, "```")
	txt = strings.TrimSpace(txt)

	log.Println("Gemini RAW:", txt)
	Broadcast("[" + time.Now().Format("15:04:05") + "] " + txt)
	SaveHistory(txt)

	var res Result
	if err := json.Unmarshal([]byte(txt), &res); err != nil {
		log.Println("ERRO JSON:", err, "txt:", txt)
		return Result{Acao: "NADA", Valor: 0}
	}
	// normaliza
	res.Acao = strings.ToUpper(res.Acao)
	return res
}

// pra tua janela lateral
func ProcessImageSync(img []byte) Result {
	res := ProcessImage(img)
	// pra usar na UI walk
	return res
}
