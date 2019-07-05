package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/dpapathanasiou/go-recaptcha"
	"github.com/joho/godotenv"
)

var chain string
var recaptchaSecretKey string
var amountFaucet string
var amountSteak string
var key string
var pass string
var node string
var publicUrl string

type claim_struct struct {
	Address  string
	Response string
}

func getEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		fmt.Println(key, "=", value)
		return value
	} else {
		log.Fatal("Error loading environment variable: ", key)
		return ""
	}
}

func main() {
	err := godotenv.Load(".env.local", ".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	chain = getEnv("FAUCET_CHAIN")
	recaptchaSecretKey = getEnv("FAUCET_RECAPTCHA_SECRET_KEY")
	amountFaucet = getEnv("FAUCET_AMOUNT_FAUCET")
	amountSteak = getEnv("FAUCET_AMOUNT_STEAK")
	key = getEnv("FAUCET_KEY")
	pass = getEnv("FAUCET_PASS")
	node = getEnv("FAUCET_NODE")
	publicUrl = getEnv("FAUCET_PUBLIC_URL")

	recaptcha.Init(recaptchaSecretKey)

	http.HandleFunc("/claim", getCoinsHandler)

	if err := http.ListenAndServe(publicUrl, nil); err != nil {
		log.Fatal("failed to start server", err)
	}
}

func executeCmd(command string, writes ...string) {
	cmd, wc, _ := goExecute(command)

	for _, write := range writes {
		wc.Write([]byte(write + "\n"))
	}
	cmd.Wait()
}

func goExecute(command string) (cmd *exec.Cmd, pipeIn io.WriteCloser, pipeOut io.ReadCloser) {
	cmd = getCmd(command)
	pipeIn, _ = cmd.StdinPipe()
	pipeOut, _ = cmd.StdoutPipe()
	go cmd.Start()
	time.Sleep(time.Second)
	return cmd, pipeIn, pipeOut
}

func getCmd(command string) *exec.Cmd {
	// split command into command and args

	var cmd *exec.Cmd

	cmd = exec.Command("/bin/sh", "shell.sh")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", string(out))

	return cmd
}

func getCoinsHandler(w http.ResponseWriter, request *http.Request) {
	var claim claim_struct

	// decode JSON response from front end
	decoder := json.NewDecoder(request.Body)
	decoderErr := decoder.Decode(&claim)
	if decoderErr != nil {
		panic(decoderErr)
	}

	var cmd *exec.Cmd
	cmd = exec.Command("/bin/sh", "shell.sh", claim.Address)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", string(out))

	return

	// // make sure address is bech32
	// readableAddress, decodedAddress, decodeErr := bech32.DecodeAndConvert(claim.Address)
	// if decodeErr != nil {
	// 	panic(decodeErr)
	// }
	// // re-encode the address in bech32
	// encodedAddress, encodeErr := bech32.ConvertAndEncode(readableAddress, decodedAddress)
	// if encodeErr != nil {
	// 	panic(encodeErr)
	// }

	// // make sure captcha is valid
	// clientIP := realip.FromRequest(request)
	// captchaResponse := claim.Response
	// captchaPassed, captchaErr := recaptcha.Confirm(clientIP, captchaResponse)
	// if captchaErr != nil {
	// 	panic(captchaErr)
	// }

	// // send the coins!
	// if captchaPassed || true {
	// 	sendFaucet := fmt.Sprintf(
	// 		"gaiacli send --to=%v --name=%v --chain-id=%v --amount=%v",
	// 		encodedAddress, key, chain, amountFaucet)
	// 	fmt.Println(time.Now().UTC().Format(time.RFC3339), encodedAddress, "[1]")
	// 	executeCmd(sendFaucet, pass)

	// 	time.Sleep(5 * time.Second)

	// 	sendSteak := fmt.Sprintf(
	// 		"gaiacli send --to=%v --name=%v --chain-id=%v --amount=%v",
	// 		encodedAddress, key, chain, amountSteak)
	// 	fmt.Println(time.Now().UTC().Format(time.RFC3339), encodedAddress, "[2]")
	// 	executeCmd(sendSteak, pass)
	// }

	// return
}
