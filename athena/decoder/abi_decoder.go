package decoder

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type ABIType struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type StructType struct {
	ABIType
	Members []Member `json:"members"`
}

type InterfaceType struct {
	ABIType
	Items []Function `json:"items"`
}

type ImplType struct {
	ABIType
	InterfaceName string `json:"interface_name"`
}

type EnumType struct {
	ABIType
	Variants []Variant `json:"variants"`
}

type EventType struct {
	ABIType
	Members []Member `json:"members"`
}

type Function struct {
	Type            string   `json:"type"`
	Name            string   `json:"name"`
	Inputs          []Member `json:"inputs"`
	Outputs         []Output `json:"outputs"`
	StateMutability string   `json:"state_mutability"`
}

type Member struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Variant struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Output struct {
	Type string `json:"type"`
}

func TypeToReadableName(typ string) string {
	typeMappings := map[string]string{
		"felt252":       "Felt252",
		"address":       "ContractAddress",
		"u256":          "U256",
		"core::felt252": "Felt252",
		"core::address": "ContractAddress",
		"core::u256":    "U256",
	}

	parts := strings.Split(typ, "::")
	if len(parts) > 0 {
		if name, exists := typeMappings[parts[len(parts)-1]]; exists {
			return name
		}
	}

	return typ
}

func GetParsedAbi(abi_to_decode string) {
	// Check if the input ABI is a valid JSON array
	var abi []map[string]interface{}
	err := json.Unmarshal([]byte(abi_to_decode), &abi)
	if err != nil {
		// If ABI is a string, unmarshal it as a string and try again
		var abiString string
		if err = json.Unmarshal([]byte(abi_to_decode), &abiString); err == nil {
			// Attempt to parse the string as JSON array
			err = json.Unmarshal([]byte(abiString), &abi)
			if err != nil {
				log.Fatalf("Error parsing ABI string: %v", err)
			}
		} else {
			log.Fatalf("Error unmarshalling ABI: %v", err)
		}
	}

	// Process ABI items
	for _, item := range abi {
		itemBytes, err := json.Marshal(item)
		if err != nil {
			log.Printf("Error marshaling item: %v", err)
			continue
		}

		switch item["type"] {
		case "interface":
			var i InterfaceType
			if err := json.Unmarshal(itemBytes, &i); err == nil {
				for _, funcItem := range i.Items {
					inputs := []string{}
					for _, input := range funcItem.Inputs {
						inputType := TypeToReadableName(input.Type)
						inputs = append(inputs, fmt.Sprintf("%s: %s", input.Name, inputType))
					}
					outputs := []string{}
					for _, output := range funcItem.Outputs {
						outputType := TypeToReadableName(output.Type)
						outputs = append(outputs, outputType)
					}
					inputSignature := strings.Join(inputs, ", ")
					outputSignature := strings.Join(outputs, ", ")
					fmt.Printf("Function: %s(%s) -> (%s) [State Mutability: %s]\n", funcItem.Name, inputSignature, outputSignature, funcItem.StateMutability)
				}
			}
		case "event":
			var event EventType
			if err := json.Unmarshal(itemBytes, &event); err == nil {
				parts := strings.Split(event.Name, "::")
				lastWord := parts[len(parts)-1]
				members := []string{}
				for _, member := range event.Members {
					memberType := TypeToReadableName(member.Type)
					members = append(members, fmt.Sprintf("%s: %s", member.Name, memberType))
				}
				membersSignature := strings.Join(members, ", ")
				fmt.Printf("Event: %s(%s)\n", lastWord, membersSignature)
			}
		default:
			continue
		}
	}
}
