package pgdb

import (
	"github.com/pgvector/pgvector-go"
)

type JitEmbedding struct {
	Embedding pgvector.Vector `pg:"type:vector(384)"`
	Text      string
	Subject   string
	Problem   string
	Images    []string
}

type TextbookEmbedding struct {
	Embedding pgvector.Vector `pg:"type:vector(384)"`
	Text      string
	Subject   string
	Page      int
}
