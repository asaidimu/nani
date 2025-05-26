package ai

import (
	"encoding/xml"
	"fmt"
	"strings"
)

// parseAIResponse takes an XML-like string following the
// <response><think></think><summary></summary><content></content></response>
// structure. It attempts to parse it into a Response and validates
// that all key fields are non-empty.
// It returns the parsed Response and an error if parsing or validation fails.
// parseAIResponse takes an XML-like string following the
// <response><think></think><summary></summary><content></content></response>
// structure. It attempts to parse it into a Response and validates
// that all key fields are non-empty.
// It returns the parsed Response and an error if parsing or validation fails.
// For malformed XML or missing/empty required tags, it returns a Response
// with the input in Content and an error.
func parseAIResponse(responseText string) (Response, error) {
	var aiResponse Response

	// Parse the XML string into the Response struct
	err := xml.Unmarshal([]byte(responseText), &aiResponse)
	if err != nil {
		// Return input text in Content field with an error for malformed XML
		return Response{
			Think:   "No think block",
			Summary: "No summary block",
			Content: responseText,
		}, fmt.Errorf("failed to parse XML: %w", err)
	}

	// Validate that all required fields contain non-empty content
	if strings.TrimSpace(aiResponse.Think) == "" {
		return Response{
			Think:   "No think block",
			Summary: "No summary block",
			Content: responseText,
		}, fmt.Errorf("validation failed: think field is empty or missing")
	}
	if strings.TrimSpace(aiResponse.Summary) == "" {
		return Response{
			Think:   "No think block",
			Summary: "No summary block",
			Content: responseText,
		}, fmt.Errorf("validation failed: summary field is empty or missing")
	}
	if strings.TrimSpace(aiResponse.Content) == "" {
		return Response{
			Think:   "No think block",
			Summary: "No summary block",
			Content: responseText,
		}, fmt.Errorf("validation failed: content field is empty or missing")
	}

	// If all checks pass, return the parsed response
	return aiResponse, nil
}
