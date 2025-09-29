package model

import (
	"encoding/base64"
	"encoding/binary"
	"math"
)

// Embedding is a special format of data representation that can be easily utilized by machine
// learning models and algorithms. The embedding is an information dense representation of the
// semantic meaning of a piece of text. Each embedding is a vector of floating point numbers,
// such that the distance between two embeddings in the vector space is correlated with semantic similarity
// between two inputs in the original format. For example, if two texts are similar,
// then their vector representations should also be similar.
type Embedding struct {
	Object    string    `json:"object"`
	Embedding []float32 `json:"embedding"`
	Index     int       `json:"index"`
}

// EmbeddingResponse is the response from a Create embeddings request.
type EmbeddingResponse struct {
	ID      string      `json:"id"`
	Created int         `json:"created"`
	Object  string      `json:"object"`
	Data    []Embedding `json:"data"`
	Model   string      `json:"model"`
	Usage   Usage       `json:"usage"`

	HttpHeader
}

type base64String string

func (b base64String) Decode() ([]float32, error) {
	decodedData, err := base64.StdEncoding.DecodeString(string(b))
	if err != nil {
		return nil, err
	}

	const sizeOfFloat32 = 4
	floats := make([]float32, len(decodedData)/sizeOfFloat32)
	for i := 0; i < len(floats); i++ {
		floats[i] = math.Float32frombits(binary.LittleEndian.Uint32(decodedData[i*4 : (i+1)*4]))
	}

	return floats, nil
}

// Base64Embedding is a container for base64 encoded embeddings.
type Base64Embedding struct {
	Object    string       `json:"object"`
	Embedding base64String `json:"embedding"`
	Index     int          `json:"index"`
}

// EmbeddingResponseBase64 is the response from a Create embeddings request with base64 encoding format.
type EmbeddingResponseBase64 struct {
	Object string            `json:"object"`
	Data   []Base64Embedding `json:"data"`
	Model  string            `json:"model"`
	Usage  Usage             `json:"usage"`

	HttpHeader
}

// ToEmbeddingResponse converts an embeddingResponseBase64 to an EmbeddingResponse.
func (r *EmbeddingResponseBase64) ToEmbeddingResponse() (EmbeddingResponse, error) {
	data := make([]Embedding, len(r.Data))

	for i, base64Embedding := range r.Data {
		embedding, err := base64Embedding.Embedding.Decode()
		if err != nil {
			return EmbeddingResponse{}, err
		}

		data[i] = Embedding{
			Object:    base64Embedding.Object,
			Embedding: embedding,
			Index:     base64Embedding.Index,
		}
	}

	return EmbeddingResponse{
		Object: r.Object,
		Model:  r.Model,
		Data:   data,
		Usage:  r.Usage,
	}, nil
}

type EmbeddingRequestConverter interface {
	// Needs to be of type EmbeddingRequestStrings or EmbeddingRequestTokens
	Convert() EmbeddingRequest
}

// EmbeddingEncodingFormat is the format of the embeddings data.
// Currently, only "float" and "base64" are supported, however, "base64" is not officially documented.
// If not specified will use "float".
type EmbeddingEncodingFormat string

const (
	EmbeddingEncodingFormatFloat  EmbeddingEncodingFormat = "float"
	EmbeddingEncodingFormatBase64 EmbeddingEncodingFormat = "base64"
)

type EmbeddingRequest struct {
	Input          interface{}             `json:"input"`
	Model          string                  `json:"model"`
	User           string                  `json:"user"`
	EncodingFormat EmbeddingEncodingFormat `json:"encoding_format,omitempty"`
	// Dimensions The number of dimensions the resulting output embeddings should have.
	// Only supported in text-embedding-3 and later models.
	Dimensions int `json:"dimensions,omitempty"`
}

func (r EmbeddingRequest) Convert() EmbeddingRequest {
	return r
}

// EmbeddingRequestStrings is the input to a create embeddings request with a slice of strings.
type EmbeddingRequestStrings struct {
	// Input is a slice of strings for which you want to generate an Embedding vector.
	// Each input must not exceed 8192 tokens in length.
	// OpenAPI suggests replacing newlines (\n) in your input with a single space, as they
	// have observed inferior results when newlines are present.
	// E.g.
	//	"The food was delicious and the waiter..."
	Input []string `json:"input"`
	// ID of the model to use. You can use the List models API to see all of your available models,
	// or see our Model overview for descriptions of them.
	Model string `json:"model"`
	// A unique identifier representing your end-user, which will help to monitor and detect abuse.
	User string `json:"user"`
	// EmbeddingEncodingFormat is the format of the embeddings data.
	// Currently, only "float" and "base64" are supported, however, "base64" is not officially documented.
	// If not specified will use "float".
	EncodingFormat EmbeddingEncodingFormat `json:"encoding_format,omitempty"`
	// Dimensions The number of dimensions the resulting output embeddings should have.
	// Only supported in text-embedding-3 and later models.
	Dimensions int `json:"dimensions,omitempty"`
}

func (r EmbeddingRequestStrings) Convert() EmbeddingRequest {
	return EmbeddingRequest{
		Input:          r.Input,
		Model:          r.Model,
		User:           r.User,
		EncodingFormat: r.EncodingFormat,
		Dimensions:     r.Dimensions,
	}
}

type EmbeddingRequestTokens struct {
	// Input is a slice of slices of ints ([][]int) for which you want to generate an Embedding vector.
	// Each input must not exceed 8192 tokens in length.
	// OpenAPI suggests replacing newlines (\n) in your input with a single space, as they
	// have observed inferior results when newlines are present.
	// E.g.
	//	"The food was delicious and the waiter..."
	Input [][]int `json:"input"`
	// ID of the model to use. You can use the List models API to see all of your available models,
	// or see our Model overview for descriptions of them.
	Model string `json:"model"`
	// A unique identifier representing your end-user, which will help to monitor and detect abuse.
	User string `json:"user"`
	// EmbeddingEncodingFormat is the format of the embeddings data.
	// Currently, only "float" and "base64" are supported, however, "base64" is not officially documented.
	// If not specified will use "float".
	EncodingFormat EmbeddingEncodingFormat `json:"encoding_format,omitempty"`
	// Dimensions The number of dimensions the resulting output embeddings should have.
	// Only supported in text-embedding-3 and later models.
	Dimensions int `json:"dimensions,omitempty"`
}

func (r EmbeddingRequestTokens) Convert() EmbeddingRequest {
	return EmbeddingRequest{
		Input:          r.Input,
		Model:          r.Model,
		User:           r.User,
		EncodingFormat: r.EncodingFormat,
		Dimensions:     r.Dimensions,
	}
}
