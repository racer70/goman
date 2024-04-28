package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"com.rwspeh/goman/config"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func main() {

	showState := flag.Bool("show-state", false, "")
	outputFileFlag := flag.String("output-file", "", "")
	envFlag := flag.String("env", "local", "environment")
	methodFlag := flag.String("m", "", "GET|POST")
	appFlag := flag.String("a", "", "application name")
	pathFlag := flag.String("p", "", "path portion of url")
	authFlag := flag.String("auth", "none", "")
	fileFlag := flag.String("f", "", "path to JSON payload")
	contentTypeFlag := flag.String("ct", "application/json", "application/json or application/x-www-form-urlencode")
	flag.Parse()

	path, err := os.Getwd()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Current working directory")
		fmt.Println(path)
	}

	if *showState {
		//TODO: Show State
	}

	cfg, err := config.GetConfig(*appFlag)

	if err != nil {
		fmt.Println("Error processing config file.")
		panic(err)
	}

	result, err := run(*methodFlag, cfg[*envFlag].GetDomain(*appFlag), *appFlag, *pathFlag, *authFlag, *fileFlag, *contentTypeFlag)

	fmt.Println(result)

	if *outputFileFlag != "" {
		os.WriteFile(*outputFileFlag, []byte(result), 0644)
	}

}

func run(
	method string,
	domain string,
	app string,
	path string,
	authType string,
	payload string,
	contentType string) (string, error) {

	if payload != "" {
		payloadBytes, e := os.ReadFile(payload)
		if e != nil {
			fmt.Println("payload not found")
		}
		payload = string(payloadBytes)

	}

	urlParts := []string{}

	if domain != "" {
		urlParts = append(urlParts, domain)
	}
	if app != "" {
		urlParts = append(urlParts, app)
	}

	if path != "" {
		urlParts = append(urlParts, path)
	}

	fullUrl := strings.Join(urlParts, "/")

	formattedMethod := cases.Upper(language.English).String(method)

	fmt.Printf("\r\n\r\nConstructed Url: %s %s\r\n\r\n", formattedMethod, fullUrl)

	req, err := http.NewRequest(formattedMethod, fullUrl, strings.NewReader(payload))

	if err != nil {
		fmt.Print("Error creating new request:")
		panic(err)
	}

	//TODO: add authorization function
	req.Header.Add("Authorization", fmt.Sprintf("%s %s", formatAuth(authType), ""))

	req.Header.Add("Content-type", contentType)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &http.Client{}

	fmt.Println("Running...")

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error running request")
		panic(err)
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		bodyBytes, err := io.ReadAll(resp.Body)

		if err != nil {
			log.Fatal(err)
		}

		//TODO:  POST PROCESS
		var out bytes.Buffer

		err = json.Indent(&out, bodyBytes, "", "  ")

		if err != nil {
			return "", err
		} else {
			return out.String(), nil
		}

	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("\r\n\r\nStatus Code: %v\r\n\r\n", resp.StatusCode)
	fmt.Printf("\r\n\r\nMessage: %v\r\n\r\n", string(respBytes))
	return "", fmt.Errorf("wtf:  %v", resp)
}

func formatAuth(authType string) string {
	if strings.ToLower(authType) != "none" {
		formattedAuth := cases.Title(language.English).String(authType)
		if formattedAuth == "Admin" {
			formattedAuth = "Basic"
		}
		return formattedAuth
	} else {
		return authType
	}
}
