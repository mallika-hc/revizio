package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
)

type Entry struct {
	entry_type     string
	request_id     string
	request_token  string
	namespace_path string
	path           string
	mount_type     string
	token_type     string
	token_ttl      float64
	operation      string
	error_text     string
	remote_address string
	error_present  bool
	token_creation bool
}

func main() {
	summarize := flag.Bool("summary", false, "Summarized output instead of CSV.")
	errors := flag.Bool("errors", false, "Output errors.")
	tokens := flag.Bool("tokens", false, "Output token creation metadata.")
	verbose := flag.Bool("verbose", false, "More detail.")
	buffersize := flag.Int("buffersize", 262144, "Adjust line buffer size max.")
	flag.Parse()

	if !*tokens && !*errors {
		flag.PrintDefaults()
		os.Exit(127)
	}

	_, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	entries := []Entry{}

	scanner := bufio.NewScanner(os.Stdin)

	buf := make([]byte, *buffersize)
	scanner.Buffer(buf, *buffersize)

	for scanner.Scan() {
		var t = scanner.Text()

		jsonMap := make(map[string](interface{}))
		err := json.Unmarshal([]byte(t), &jsonMap)

		if err == nil {
			if jsonMap["type"] == "response" {
				e := handleResponse(jsonMap)
				entries = append(entries, e)
				if !*summarize {
					if *tokens && e.token_creation {
						if *verbose {
							printTokenCreationVerbose(e)
						} else {
							printTokenCreation(e)
						}
					}
					if *errors && e.error_present {
						if *verbose {
							printErrorVerbose(e)
						} else {
							printError(e)
						}
					}

				}
			}

		}

		if err := scanner.Err(); err != nil {
			log.Println(err)
		}

	}

	if *summarize {
		printSummary(entries)
	}

}

func printSummary(entries []Entry) {
	unique_tokens := mapset.NewSet[string]()
	unique_namespaces := mapset.NewSet[string]()
	unique_namespacepathspaths := mapset.NewSet[string]()
	total_batch_tokens := 0
	total_service_tokens := 0
	for _, e := range entries {
		unique_tokens.Add(e.request_token)
		unique_namespaces.Add(e.namespace_path)
		unique_namespacepathspaths.Add(e.namespace_path + e.path)

		if e.token_type == "batch" {
			total_batch_tokens += 1
		} else if e.token_type == "service" {
			total_service_tokens += 1
		}
	}

	fmt.Printf("Total Number of Token Creations:    %d\n", len(entries))
	fmt.Printf("Unique Tokens:                      %d\n", unique_tokens.Cardinality())
	fmt.Printf("Unique Namespaces:                  %d\n", unique_namespaces.Cardinality())
	fmt.Printf("Unique Namespace/Path Combinations: %d\n", unique_namespacepathspaths.Cardinality())
	fmt.Printf("Total Number of Batch Tokens:       %d\n", total_batch_tokens)
	fmt.Printf("Total Number of Service Tokens:     %d\n", total_service_tokens)
}

func printError(e Entry) {
	fmt.Printf("\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\"\n", e.entry_type, e.namespace_path, e.path, e.mount_type, e.token_type, e.error_text)
}

func printErrorVerbose(e Entry) {
	fmt.Printf("\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\"\n", e.entry_type, e.request_id, e.remote_address, e.request_token, e.namespace_path, e.path, e.mount_type, e.token_type, e.error_text)
}

func printTokenCreation(e Entry) {
	fmt.Printf("\"%s\",\"%s\",\"%s\",\"%s\",\"%s\"\n", e.entry_type, e.namespace_path, e.path, e.mount_type, e.token_type)
}

func printTokenCreationVerbose(e Entry) {
	fmt.Printf("\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",%d\n", e.entry_type, e.request_id, e.remote_address, e.request_token, e.namespace_path, e.path, e.mount_type, e.token_type, int(e.token_ttl))
}

func handleResponse(line map[string](interface{})) Entry {
	e := Entry{
		entry_type:     "",
		request_id:     "",
		request_token:  "",
		namespace_path: "",
		path:           "",
		mount_type:     "",
		token_type:     "",
		operation:      "",
		error_text:     "",
		remote_address: "",
		token_creation: false,
		error_present:  false,
	}

	if request, okay := line["request"].(map[string]interface{}); okay {
		if operation, okay := request["operation"]; okay {
			e.operation = operation.(string)
		}

		if path, okay := request["path"]; okay {
			e.path = path.(string)
		} else {
			e.path = "<no_path>"
		}

		if request_id, okay := request["id"]; okay {
			e.request_id = request_id.(string)
		}

		if mount_type, okay := request["mount_type"]; okay {
			e.mount_type = mount_type.(string)
		} else {
			e.mount_type = "<no_mount_type>"
		}

		if request_token, okay := request["client_token"]; okay {
			e.request_token = request_token.(string)
		}

		if remote_address, okay := request["remote_address"]; okay {
			e.remote_address = remote_address.(string)
		}

		namespace := request["namespace"].(map[string]interface{})
		if namespace_path, okay := namespace["path"]; okay {
			e.namespace_path = namespace_path.(string)
		} else {
			e.namespace_path = "<root>"
		}
	}

	if error_text, okay := line["error"].(string); okay {
		e.error_present = true
		e.entry_type = "<error>"
		error_text_nonewlines := strings.ReplaceAll(error_text, "\n", "\\n")
		e.error_text = strings.ReplaceAll(error_text_nonewlines, "\t", "\\t")
	}

	if response, okay := line["response"].(map[string]interface{}); okay {

		if auth, okay := response["auth"].(map[string]interface{}); okay {
			if token_type, okay := auth["token_type"]; okay {
				e.token_type = token_type.(string)

				if e.operation == "update" {
					e.token_creation = true
					e.entry_type = "<token_creation>"
					if request_token, okay := auth["client_token"]; okay {
						e.request_token = request_token.(string)
					}

					if token_ttl, okay := auth["token_ttl"].(float64); okay {
						e.token_ttl = token_ttl
					}

				}
			}

		}
	}
	return e

}
