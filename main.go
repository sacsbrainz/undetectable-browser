package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/launcher/flags"
)

var (
	Version     = "dev"
	proxy       = flag.String("proxy", "", "Proxy URL in format [protocol://][username:password@]hostname[:port]")
	quiet       = flag.Bool("quiet", false, "Silence all log output")
	userDir     = flag.String("user-dir", "", "Custom user directory path (defaults to proxy hostname)")
	help        = flag.Bool("help", false, "Show help message")
	showVersion = flag.Bool("version", false, "Show version information")
)

type TimezoneResponse struct {
	Timezone string `json:"timezone"`
}

func getProxyTimezone(proxyIP string) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(fmt.Sprintf("http://ip-api.com/json/%s?fields=timezone", proxyIP))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tzResp TimezoneResponse
	if err := json.Unmarshal(body, &tzResp); err != nil {
		return "", err
	}

	return tzResp.Timezone, nil
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nExample:\n")
	fmt.Fprintf(os.Stderr, "  %s -proxy http://user:pass@localhost:8080\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s -proxy localhost:8080 -user-dir custom_profile\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func monitorBrowser(browser *rod.Browser, done chan struct{}, quiet *bool) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			pages, err := browser.Pages()
			if err != nil {
				if !*quiet {
					log.Println("Error getting pages:", err)
				}
				close(done)
				return
			}

			if len(pages) == 0 {
				if !*quiet {
					log.Println("No pages left open, shutting down...")
				}
				close(done)
				return
			}
		case <-browser.GetContext().Done():
			if !*quiet {
				log.Println("Browser disconnected")
			}
			close(done)
			return
		}
	}
}

func main() {

	flag.Usage = usage
	flag.Parse()

	if *showVersion {
		fmt.Printf("undetectable-browser version %s\n", Version)
		os.Exit(0)
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if *proxy == "" {
		log.Println("Error: proxy is required")
		flag.Usage()
		os.Exit(1)
	}

	proxyURL, err := url.Parse(*proxy)
	if err != nil {
		log.Fatalf("Failed to parse proxy URL: %v", err)
	}

	timezone, err := getProxyTimezone(proxyURL.Hostname())
	if err != nil && !*quiet {
		log.Printf("Warning: Failed to get timezone from proxy IP: %v", err)
		timezone = "UTC"
	}

	userDirectory := *userDir
	if userDirectory == "" {
		userDirectory = proxyURL.Hostname()
	}

	u := launcher.New().
		Env(append(os.Environ(), fmt.Sprintf("TZ=%s", timezone))...).
		NoSandbox(true).
		Headless(false).
		Leakless(true).
		Set("disable-gpu").
		UserDataDir(filepath.Join("data", userDirectory)).
		Set("disable-dev-shm-usage").
		Set("start-maximized").
		Set("no-first-run").
		Set("safebrowsing-disable-download-protection").
		Set("webrtc-ip-handling-policy", "disable_non_proxied_udp").
		Set("enforce-webrtc-ip-permission-check").
		Set("disable-webrtc").
		Set("disable-blink-features", "AutomationControlled,MediaDevices,WebMIDI,AudioWorklet").
		Set("disable-features", "AudioServiceOutOfProcess,IsolateOrigins,site-per-process").
		Set(flags.ProxyServer, proxyURL.Host)

	if chromePath, hasChrome := launcher.LookPath(); hasChrome {
		u.Bin(chromePath)
	}

	controlURL := u.MustLaunch()
	browser := rod.New().
		ControlURL(controlURL).
		MustConnect()

	if !*quiet {
		log.Printf("-> Setting up proxy... %s", *proxy)
		log.Printf("-> Using timezone: %s", timezone)

	}

	username := proxyURL.User.Username()
	password, _ := proxyURL.User.Password()
	if username != "" && password != "" {
		go browser.MustHandleAuth(username, password)()
	}

	browser.MustIgnoreCertErrors(true)

	var wg sync.WaitGroup
	wg.Add(1)

	done := make(chan struct{})

	go func() {
		defer wg.Done()
		monitorBrowser(browser, done, quiet)
	}()

	page := browser.MustPage()

	page.MustNavigate("https://whoer.com")

	<-done

	if err := browser.Close(); err != nil && !*quiet {
		log.Printf("Error closing browser: %v", err)
	}

	wg.Wait()

	if !*quiet {
		log.Println("Program terminated")
	}
	os.Exit(0)
}
