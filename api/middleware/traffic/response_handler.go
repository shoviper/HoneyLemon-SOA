package traffic

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

func ConvertXMLToJSON(xmlData []byte) ([]byte, error) {
	decoder := xml.NewDecoder(bytes.NewReader(xmlData))
	result := make(map[string]interface{})
	stack := []map[string]interface{}{result}
	elementNames := []string{}

	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		switch token := token.(type) {
		case xml.StartElement:
			newMap := make(map[string]interface{})
			currentMap := stack[len(stack)-1]
			elementName := token.Name.Local
			if currentMap != nil {
				if existing, ok := currentMap[elementName]; ok {
					if existingArray, ok := existing.([]interface{}); ok {
						// Append to existing array
						currentMap[elementName] = append(existingArray, newMap)
					} else {
						// Create an array with existing and new elements
						currentMap[elementName] = []interface{}{existing, newMap}
					}
				} else {
					currentMap[elementName] = newMap
				}
			}
			stack = append(stack, newMap)
			elementNames = append(elementNames, elementName)

		case xml.EndElement:
			if len(stack) > 1 {
				stack = stack[:len(stack)-1]
				elementNames = elementNames[:len(elementNames)-1]
			}

		case xml.CharData:
			if len(stack) > 0 {
				currentMap := stack[len(stack)-1]
				text := string(token)
				if text = strings.TrimSpace(text); text != "" {
					if existing, ok := currentMap["content"]; ok {
						currentMap["content"] = existing.(string) + text
					} else {
						currentMap["content"] = text
					}
				}
			}
		}
	}

	// Marshal the map into JSON
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

// ExtractBody extracts the Body part and replaces "content" with its value.
func ExtractBody(jsonData []byte, responseTag string) ([]byte, error) {
	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return nil, err
	}

	body, ok := result["Envelope"].(map[string]interface{})["Body"].(map[string]interface{})
	if !ok {
		return nil, nil
	}

	response, ok := body[responseTag]
	if !ok {
		return nil, fmt.Errorf("expected key '%s' not found", responseTag)
	}

	// Handle multiple elements
	if elements, ok := response.([]interface{}); ok {
		var cleanedElement []interface{}
		for _, element := range elements {
			cleanedElement = append(cleanedElement, cleanContent(element))
		}
		// Marshal the result into JSON
		return json.MarshalIndent(cleanedElement, "", "  ")
	}

	// Replace "content" keys with their values
	cleanedResponse := cleanContent(response)

	// Marshal the result into JSON
	responseJson, err := json.MarshalIndent(cleanedResponse, "", "  ")
	if err != nil {
		return nil, err
	}

	return responseJson, nil
}

// cleanContent replaces "content" keys with their values in the given data.
func cleanContent(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		cleaned := make(map[string]interface{})
		for key, value := range v {
			if key == "content" {
				return value
			}
			cleaned[key] = cleanContent(value)
		}
		return cleaned
	case []interface{}:
		for i, item := range v {
			v[i] = cleanContent(item)
		}
		return v
	default:
		return v
	}
}
