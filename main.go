package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	mapset "github.com/deckarep/golang-set/v2"
)

type Entry struct {
	request_id     string
	request_token  string
	namespace_path string
	path           string
	mount_type     string
	token_type     string
	token_ttl      float64
}

func main() {
	summarize := flag.Bool("summary", false, "Summarized output instead of CSV")
	detailed := flag.Bool("detailed", false, "Detailed view (includes token, response token_ttl).")
	flag.Parse()

	_, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	entries := []Entry{}

	scanner := bufio.NewScanner(os.Stdin)

	const maxCapacity int = 262144 // TODO: Tunable?
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		var t = scanner.Text()

		jsonMap := make(map[string](interface{}))
		err := json.Unmarshal([]byte(t), &jsonMap)

		if err == nil {
			if jsonMap["type"] == "response" {
				e := handleResponse(jsonMap)
				if (Entry{} != e) {
					entries = append(entries, e)
					if !*summarize {
						if *detailed {
							printMapDetailed(e)
						} else {
							printMap(e)
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

func printMap(e Entry) {
	fmt.Printf("\"%s\",\"%s\",\"%s\",\"%s\"\n", e.namespace_path, e.path, e.mount_type, e.token_type)
}

func printMapDetailed(e Entry) {
	fmt.Printf("\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",%d\n", e.request_id, e.request_token, e.namespace_path, e.path, e.mount_type, e.token_type, int(e.token_ttl))
}

func handleResponse(line map[string](interface{})) Entry {
	e := Entry{}

	request := line["request"].(map[string]interface{})

	response := line["response"].(map[string]interface{})

	if auth, okay := response["auth"].(map[string]interface{}); okay {
		if token_type, okay := auth["token_type"]; okay {
			e.token_type = token_type.(string)

			if operation, okay := request["operation"]; okay {
				if operation == "update" {

					if request_id, okay := request["id"]; okay {
						e.request_id = request_id.(string)
					}
					if request_token, okay := request["client_token"]; okay {
						e.request_token = request_token.(string)
					} else if request_token, okay := auth["client_token"]; okay {
						e.request_token = request_token.(string)
					}

					if token_ttl, okay := auth["token_ttl"].(float64); okay {
						e.token_ttl = token_ttl
					}

					if path, okay := request["path"]; okay {
						e.path = path.(string)
					} else {
						e.path = "<no_path>"
					}

					if mount_type, okay := request["mount_type"]; okay {
						e.mount_type = mount_type.(string)
					} else {
						e.mount_type = "<no_mount_type>"
					}

					namespace := request["namespace"].(map[string]interface{})
					if namespace_path, okay := namespace["path"]; okay {
						e.namespace_path = namespace_path.(string)
					} else {
						e.namespace_path = "<root>"
					}
				} else {
					return Entry{}
				}
			} else {
				return Entry{}
			}
		} else {
			return Entry{}
		}

	}
	return e

}
