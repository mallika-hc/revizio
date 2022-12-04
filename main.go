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
	Entry_type     string `json:"entry_type"`
	Request_id     string `json:"request_id"`
	Request_token  string `json:"request_token"`
	Namespace_path string `json:"namespace_path"`
	Path           string `json:"path"`
	Mount_type     string `json:"mount_type"`
	Token_type     string `json:"token_type"`
	Token_ttl      float64 `json:"token_ttl"`
	Operation      string `json:"operation"`
	Error_text     string `json:"error_text"`
	Remote_address string `json:"remote_address"`
	Time           string `json:"time"`
	Error_present  bool `json:"error_present"`
	Token_creation bool `json:"token_creation"`
}

func main() {
	summarize := flag.Bool("summary", false, "Summarized output instead of CSV.")
	errors := flag.Bool("errors", false, "Output errors.")
	tokens := flag.Bool("tokens", false, "Output token creation metadata.")
	verbose := flag.Bool("verbose", false, "More detail.")
	buffersize := flag.Int("buffersize", 262144, "Adjust line buffer size max.")
	
	var fieldnames string
    flag.StringVar(&fieldnames, "fields", "", "list of fileds to be displayed")
	
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
					if *tokens && e.Token_creation {
						if *verbose {
							printTokenCreationVerbose(e)
						} else {
							printTokenCreation(e, fieldnames)
						}
					}
					if *errors && e.Error_present {
						if *verbose {
							printErrorVerbose(e)
						} else {
							printError(e, fieldnames)
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
		unique_tokens.Add(e.Request_token)
		unique_namespaces.Add(e.Namespace_path)
		unique_namespacepathspaths.Add(e.Namespace_path + e.Path)

		if e.Token_type == "batch" {
			total_batch_tokens += 1
		} else if e.Token_type == "service" {
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

func printError(e Entry, fields string) {
	if fields == "" {
	fmt.Printf("\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\"\n", e.Entry_type, e.Time, e.Namespace_path, e.Path, e.Mount_type, e.Token_type, e.Error_text) 
	} else {
		printTargetFields(e, fields)
	}
}

func printTargetFields(e Entry, fields string) {

	b, err := json.Marshal(e)
    if err != nil {
        fmt.Println(err)
        return
	}
	jsonMap := make(map[string](interface{}))
	err = json.Unmarshal([]byte(b), &jsonMap)

	fieldnames := strings.Split(fields, ",")

	fmt.Printf("\"%s\",",e.Entry_type)
	
	for i := 0; i < len(fieldnames); i++ {
		if i != len(fieldnames)-1 {
			fmt.Printf("\"%s\",",jsonMap[fieldnames[i]])
		} else {
			fmt.Printf("\"%s\"",jsonMap[fieldnames[i]])
		}
	}
	fmt.Print("\n")
}

func printErrorVerbose(e Entry) {
	fmt.Printf("\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\"\n", e.Entry_type, e.Time, e.Request_id, e.Remote_address, e.Request_token, e.Namespace_path, e.Path, e.Mount_type, e.Token_type, e.Error_text)
}

func printTokenCreation(e Entry, fields string) {
	if fields == "" {
		fmt.Printf("\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\"\n", e.Entry_type, e.Time, e.Namespace_path, e.Path, e.Mount_type, e.Token_type)
	} else {
		printTargetFields(e, fields)
	}
}

func printTokenCreationVerbose(e Entry) {
	fmt.Printf("\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",%d\n", e.Entry_type, e.Time, e.Request_id, e.Remote_address, e.Request_token, e.Namespace_path, e.Path, e.Mount_type, e.Token_type, int(e.Token_ttl))
}

func handleResponse(line map[string](interface{})) Entry {
	e := Entry{
		Time:           "",
		Entry_type:     "",
		Request_id:     "",
		Request_token:  "",
		Namespace_path: "",
		Path:           "",
		Mount_type:     "",
		Token_type:     "",
		Operation:      "",
		Error_text:     "",
		Remote_address: "",
		Token_creation: false,
		Error_present:  false,
	}

	if time, okay := line["time"].(string); okay {
		e.Time = time
	}

	if request, okay := line["request"].(map[string]interface{}); okay {
		if operation, okay := request["operation"]; okay {
			e.Operation = operation.(string)
		}

		if path, okay := request["path"]; okay {
			e.Path = path.(string)
		} else {
			e.Path = "<no_path>"
		}

		if request_id, okay := request["id"]; okay {
			e.Request_id = request_id.(string)
		}

		if mount_type, okay := request["mount_type"]; okay {
			e.Mount_type = mount_type.(string)
		} else {
			e.Mount_type = "<no_mount_type>"
		}

		if request_token, okay := request["client_token"]; okay {
			e.Request_token = request_token.(string)
		}

		if remote_address, okay := request["remote_address"]; okay {
			e.Remote_address = remote_address.(string)
		}

		namespace := request["namespace"].(map[string]interface{})
		if namespace_path, okay := namespace["path"]; okay {
			e.Namespace_path = namespace_path.(string)
		} else {
			e.Namespace_path = "<root>"
		}
	}

	if error_text, okay := line["error"].(string); okay {
		e.Error_present = true
		e.Entry_type = "<error>"
		error_text_nonewlines := strings.ReplaceAll(error_text, "\n", "\\n")
		e.Error_text = strings.ReplaceAll(error_text_nonewlines, "\t", "\\t")
	}

	if response, okay := line["response"].(map[string]interface{}); okay {

		if auth, okay := response["auth"].(map[string]interface{}); okay {
			if token_type, okay := auth["token_type"]; okay {
				e.Token_type = token_type.(string)

				if e.Operation == "update" {
					e.Token_creation = true
					e.Entry_type = "<token_creation>"
					if request_token, okay := auth["client_token"]; okay {
						e.Request_token = request_token.(string)
					}

					if token_ttl, okay := auth["token_ttl"].(float64); okay {
						e.Token_ttl = token_ttl
					}

				}
			}

		}
	}
	return e

}
