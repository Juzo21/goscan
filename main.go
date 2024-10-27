package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

type Config struct {
	UserAgent       string
	Cookie          string
	ClientCert      string
	OutputFile      string
	ProxyURL        string
	Extensions      string
	Delay           int
	IgnoreCode      int
	AuthUser        string
	AuthPassword    string
	CaseInsensitive bool
	SilentMode      bool
}

func AsciiGolmonHelp() {
	asciiArt := []string{
		"    .-----.",
		"   .' -   - '.",
		"  /  .-. .-.  \\",
		"  |  | | | |  |",
		"   \\ \\o/ \\o/ /",
		"  _/    ^    \\_",
		" | \\  '---'  / |",
		" / /`--. .--`\\ \\",
		"/ /'---` `---'\\ \\",
		"'.__.       .__.'",
		"    `|     |`",
		"     |     \\",
		"     \\      '--.",
		"      '.        `\\",
		"        `'---.   |",
		"             ,__) /",
		"            `..'  ",
	}

	for _, line := range asciiArt {
		c := color.New(color.FgCyan)
		c.Println(line)
	}
}

func printHelp() {
	c := color.New(color.FgCyan)
	d := color.New(color.FgRed)
	c.Println("Options disponibles :")
	c.Println("-url <URL>         : URL de base à scanner")
	c.Println("-wordlist <file>   : Fichier contenant la wordlist [ Si t'en as pas va les prendre ici c'est bien : https://github.com/v0re/dirb/blob/master/wordlists/]")
	c.Println("-a <agent>         : Agent utilisateur personnalisé")
	c.Println("-c <cookie>        : Cookie pour les requêtes HTTP")
	c.Println("-E <certificat>    : Chemin vers le certificat client  [Inutile]")
	c.Println("-o <file>          : Enregistrer les résultats dans un fichier")
	c.Println("-p <proxy>         : URL du proxy : Exemple : http://10.10.10.10:8080/")
	c.Println("-X <extensions>    : Extensions à ajouter aux mots")
	c.Println("-z <delay>         : Délai entre les requêtes en millisecondes")
	c.Println("-N <ignoreCode>    : Ignorer les réponses avec ce code HTTP")
	c.Println("-u <user>          : Nom d'utilisateur pour l'authentification HTTP")
	c.Println("-P <password>      : Mot de passe pour l'authentification HTTP")
	c.Println("-i ------------->  : Recherche insensible à la casse")
	c.Println("-S ------------->  : Mode silencieux (ca t'affiche rien c'est discret zehma)")
	c.Println("-h ------------->  : Afficher le menu d'aide")

	d.Println("[Utilisation : go run main.go -url <URL> -wordlist <file>]")
}

func loadlist(filename string, caseInsensitive bool) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := scanner.Text()
		if caseInsensitive {
			word = strings.ToLower(word)
		}
		words = append(words, word)
	}
	return words, scanner.Err()
}

func urlscan(baseURL, word string, config Config, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()

	baseURL = strings.TrimRight(baseURL, "/")

	urlToTest := fmt.Sprintf("%s/%s", baseURL, word)
	if config.Extensions != "" {
		urlToTest += config.Extensions
	}

	req, err := http.NewRequest("GET", urlToTest, nil)
	if err != nil {
		return
	}

	if config.UserAgent != "" {
		req.Header.Set("User-Agent", config.UserAgent)
	}
	if config.Cookie != "" {
		req.Header.Set("Cookie", config.Cookie)
	}
	if config.AuthUser != "" && config.AuthPassword != "" {
		req.SetBasicAuth(config.AuthUser, config.AuthPassword)
	}

	client := &http.Client{}
	if config.ClientCert != "" {
		cert, err := tls.LoadX509KeyPair(config.ClientCert, config.ClientCert)
		if err != nil {
			fmt.Println("Erreur de chargement du certificat:", err)
			return
		}
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{Certificates: []tls.Certificate{cert}},
		}
	}
	if config.ProxyURL != "" {
		proxyURL, _ := url.Parse(config.ProxyURL)
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == config.IgnoreCode {
		return
	}

	if !config.SilentMode {
		c := color.New(color.FgRed)
		c.Printf("URL trouvée : %s (Code %d)\n", urlToTest, resp.StatusCode)
	}
}

func main() {

	AsciiGolmonHelp()

	baseURL := flag.String("url", "", "URL de base à scanner")
	wordlistFile := flag.String("wordlist", "", "Fichier contenant la wordlist")
	userAgent := flag.String("a", "", "Agent utilisateur personnalisé")
	cookie := flag.String("c", "", "Cookie pour les requêtes HTTP")
	clientCert := flag.String("E", "", "Chemin vers le certificat client [Inutile]")
	outputFile := flag.String("o", "", "Enregistrer les résultats dans un fichier")
	proxy := flag.String("p", "", "URL du proxy : Exemple : http://10.10.10.10:8080/")
	extensions := flag.String("X", "", "Extensions à ajouter aux mots")
	delay := flag.Int("z", 0, "Délai entre les requêtes en ms")
	ignoreCode := flag.Int("N", 0, "Ignorer les réponses avec ce code HTTP")
	authUser := flag.String("u", "", "Nom d'utilisateur pour l'authentification HTTP")
	authPassword := flag.String("P", "", "Mot de passe pour l'authentification HTTP")
	caseInsensitive := flag.Bool("i", false, "Recherche insensible à la casse")
	silentMode := flag.Bool("S", false, "Mode silencieux (ca t'affiche rien c'est discret zehma)")
	help := flag.Bool("h", false, "Afficher le menu d'aide")

	flag.Parse()

	if *help || (*baseURL == "" && *wordlistFile == "") {
		printHelp()
		return
	}

	c := color.New(color.FgCyan).Add(color.Underline)
	if *baseURL == "" || *wordlistFile == "" {
		c.Println("Utilisation : go run main.go -url <URL> -wordlist <file>")
		return
	}

	config := Config{
		UserAgent:       *userAgent,
		Cookie:          *cookie,
		ClientCert:      *clientCert,
		OutputFile:      *outputFile,
		ProxyURL:        *proxy,
		Extensions:      *extensions,
		Delay:           *delay,
		IgnoreCode:      *ignoreCode,
		AuthUser:        *authUser,
		AuthPassword:    *authPassword,
		CaseInsensitive: *caseInsensitive,
		SilentMode:      *silentMode,
	}

	words, err := loadlist(*wordlistFile, *caseInsensitive)
	if err != nil {
		fmt.Printf("Erreur lors du chargement de la liste de mots: %v\n", err)
		return
	}

	results := make(chan string)
	var wg sync.WaitGroup

	for _, word := range words {
		wg.Add(1)
		go urlscan(*baseURL, word, config, &wg, results)
		if config.Delay > 0 {
			time.Sleep(time.Duration(config.Delay) * time.Millisecond)
		}
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var file *os.File
	if config.OutputFile != "" {
		file, err = os.Create(config.OutputFile)
		if err != nil {
			fmt.Println("Erreur de création du fichier :", err)
			return
		}
		defer file.Close()
	}

	for result := range results {
		if config.OutputFile != "" {
			file.WriteString(result + "\n")
		} else {
			fmt.Println("URL trouvée :", result)
		}
	}
}
