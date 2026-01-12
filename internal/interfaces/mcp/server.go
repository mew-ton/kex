package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/mew-ton/kex/internal/infrastructure/logger"
	"github.com/mew-ton/kex/internal/usecase/retrieve"
	"github.com/mew-ton/kex/internal/usecase/search"
)

// Server handles MCP JSON-RPC requests
type Server struct {
	SearchUC   *search.UseCase
	RetrieveUC *retrieve.UseCase
}

func New(searchUC *search.UseCase, retrieveUC *retrieve.UseCase) *Server {
	return &Server{
		SearchUC:   searchUC,
		RetrieveUC: retrieveUC,
	}
}

// JSON-RPC types
type request struct {
	JSONRPC string           `json:"jsonrpc"`
	Method  string           `json:"method"`
	Params  json.RawMessage  `json:"params"`
	ID      *json.RawMessage `json:"id,omitempty"` // Pointer to handle null/missing
}

type response struct {
	JSONRPC string           `json:"jsonrpc"`
	Result  interface{}      `json:"result,omitempty"`
	Error   *rpcError        `json:"error,omitempty"`
	ID      *json.RawMessage `json:"id"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Serve starts the JSON-RPC loop on Stdio
func (s *Server) Serve() error {
	scanner := bufio.NewScanner(os.Stdin)
	// Increase buffer size if needed, but default is usually fine for messages
	// MCP messages can be large (tool outputs), but requests are usually small.

	for scanner.Scan() {
		line := scanner.Bytes()
		s.handleMessage(line)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read error: %w", err)
	}
	return nil
}

func (s *Server) handleMessage(msg []byte) {
	var req request
	if err := json.Unmarshal(msg, &req); err != nil {
		// Parse error
		fmt.Fprintf(os.Stderr, "failed to parse request: %v\n", err)
		return
	}

	res := response{
		JSONRPC: "2.0",
		ID:      req.ID,
	}

	// Handle methods
	var err *rpcError
	var result interface{}

	logger.Info("[MCP] Request: %s", req.Method)

	switch req.Method {
	case "initialize":
		result = map[string]interface{}{
			"protocolVersion": "2024-11-05", // Latest known or 0.1.0
			"serverInfo": map[string]string{
				"name":    "kex",
				"version": "1.0.0",
			},
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{},
			},
		}
	case "notifications/initialized":
		// No response needed
		return
	case "ping":
		result = map[string]string{}
	case "tools/list":
		result = s.handleListTools()
	case "tools/call":
		result, err = s.handleCallTool(req.Params)
	default:
		// Ignore unknown notifications
		if req.ID == nil {
			return
		}
		err = &rpcError{Code: -32601, Message: "Method not found"}
	}

	if err != nil {
		res.Error = err
	} else {
		res.Result = result
	}

	s.sendResponse(res)
}

func (s *Server) sendResponse(res response) {
	bytes, err := json.Marshal(res)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to marshal response: %v\n", err)
		return
	}
	fmt.Printf("%s\n", bytes)

	status := "Success"
	if res.Error != nil {
		status = fmt.Sprintf("Error (%d: %s)", res.Error.Code, res.Error.Message)
	}
	logger.Info("[MCP] Response Sent: ID=%s, Status=%s", stringifyID(res.ID), status)
}

func stringifyID(id *json.RawMessage) string {
	if id == nil {
		return "null"
	}
	return string(*id)
}

// -- Handlers --

func (s *Server) handleListTools() interface{} {
	return map[string]interface{}{
		"tools": []map[string]interface{}{
			{
				"name":        "search_documents",
				"description": "Search project guidelines using keywords",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"keywords": map[string]interface{}{
							"type": "array",
							"items": map[string]string{
								"type": "string",
							},
							"description": "Keywords related to the coding task",
						},
						"filePath": map[string]interface{}{
							"type":        "string",
							"description": "The path of the file you are working on. Used for scope filtering.",
						},
						"exactScopeMatch": map[string]interface{}{
							"type":        "boolean",
							"description": "If true, treats keywords as exact scope names to match.",
						},
					},
					"required": []string{"keywords"},
				},
			},
			{
				"name":        "read_document",
				"description": "Read the full content of a specific document",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"id": map[string]interface{}{
							"type":        "string",
							"description": "Document ID",
						},
					},
					"required": []string{"id"},
				},
			},
		},
	}
}

func (s *Server) handleCallTool(paramsRaw json.RawMessage) (interface{}, *rpcError) {
	var params struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	if err := json.Unmarshal(paramsRaw, &params); err != nil {
		return nil, &rpcError{Code: -32700, Message: "Invalid params"}
	}

	switch params.Name {
	case "search_documents":
		return s.handleSearchDocuments(params.Arguments)
	case "read_document":
		return s.handleReadDocument(params.Arguments)
	default:
		return nil, &rpcError{Code: -32601, Message: "Tool not found"}
	}
}

func (s *Server) handleSearchDocuments(argsRaw json.RawMessage) (interface{}, *rpcError) {
	var args struct {
		Keywords        []string `json:"keywords"`
		FilePath        string   `json:"filePath"`
		ExactScopeMatch bool     `json:"exactScopeMatch"`
	}
	if err := json.Unmarshal(argsRaw, &args); err != nil {
		return nil, &rpcError{Code: -32700, Message: "Invalid arguments"}
	}

	// Use Search Use Case
	result := s.SearchUC.Execute(args.Keywords, args.FilePath, args.ExactScopeMatch)

	logger.Info("[Tool:search_documents] Query: Keywords=%v, FilePath=%s, Exact=%v", args.Keywords, args.FilePath, args.ExactScopeMatch)

	var foundIDs []string
	for _, doc := range result.Documents {
		foundIDs = append(foundIDs, doc.ID)
	}
	logger.Info("[Tool:search_documents] Result: Found %d documents, IDs=%v", len(result.Documents), foundIDs)

	var content []map[string]interface{}

	if len(result.Documents) == 0 {
		content = append(content, map[string]interface{}{
			"type": "text",
			"text": "No matching documents found.",
		})
	} else {
		text := "Found documents:\n"
		for _, doc := range result.Documents {
			text += fmt.Sprintf("- **%s** (ID: `%s`): %s\n", doc.Title, doc.ID, doc.Description)
		}
		content = append(content, map[string]interface{}{
			"type": "text",
			"text": text,
		})
	}

	return map[string]interface{}{"content": content}, nil
}

func (s *Server) handleReadDocument(argsRaw json.RawMessage) (interface{}, *rpcError) {
	var args struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(argsRaw, &args); err != nil {
		return nil, &rpcError{Code: -32700, Message: "Invalid arguments"}
	}

	logger.Info("[Tool:read_document] ID: %s", args.ID)

	result := s.RetrieveUC.Execute(args.ID)
	if !result.Found {
		logger.Info("[Tool:read_document] Result: Not Found")
		return map[string]interface{}{
			"content": []map[string]interface{}{
				{"type": "text", "text": "Document not found."},
			},
			"isError": true,
		}, nil
	}

	logger.Info("[Tool:read_document] Result: Success (%d bytes)", len(result.Document.Body))

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("# %s\n\n%s", result.Document.Title, result.Document.Body),
			},
		},
	}, nil
}
