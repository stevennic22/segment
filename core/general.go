package core

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

func CheckError(err error, hardStop bool) bool {
	if err != nil {
		if hardStop {
			log.Fatalf("Fatal error within application: %s", err.Error())
		}
		log.Printf("Recoverable error within application: %s", err.Error())
		return (true)
	}

	return (false)
}

func AllowedOriginCheck(origin string) bool {
	trusted_domain := Config.RetrieveValue("trusted_domain")

	regexStr := fmt.Sprintf("^http(s)*://(.*\\.)*(%s)$", trusted_domain)
	r := regexp.MustCompile(regexStr)

	if r.MatchString(origin) {
		return true
	} else {
		trustExtensions := Config.RetrieveValue("trust_extensions")
		if trustExtensions == "true" && (strings.Contains(origin, "vscode-webview") || strings.Contains(origin, "moz-extension")) {
			return true
		} else {
			return false
		}
	}
}
