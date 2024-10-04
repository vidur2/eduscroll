package pgdb

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

var Db *pg.DB

func InitVectorDb() error {
	opt, err := pg.ParseURL("postgres://u_mt3fwb0dxsar8yh:fo36cdauwvgwbnm@02f7e6f1-1adb-4347-835a-02c74fcccb0e.db.cloud.postgresml.org:6432/pgml_mfrodnifmwwnd0j")

	if err != nil {
		return err
	}

	Db = pg.Connect(opt)

	_, err = Db.Exec("DROP TABLE IF EXISTS jit_embedding")
	if err != nil {
		panic(err)
	}

	_, err = Db.Exec("DROP TABLE IF EXISTS textbook_embedding")
	if err != nil {
		panic(err)
	}

	_, err = Db.Exec("CREATE EXTENSION IF NOT EXISTS vector")
	if err != nil {
		panic(err)
	}

	if err != nil {
		return err
	}

	createSchema(Db)
	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}

	return nil
}

func createSchema(db *pg.DB) error {
	models := []interface{}{
		(*JitEmbedding)(nil),
		(*TextbookEmbedding)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp: false,
		})
		if err != nil {
			return err
		}
	}

	_, err := db.Exec("CREATE INDEX ON jit_embeddings USING hnsw (embedding vector_l2_ops)")
	if err != nil {
		return err
	}
	_, err = db.Exec("CREATE INDEX ON textbook_embeddings USING hnsw (embedding vector_l2_ops)")
	if err != nil {
		return err
	}
	return nil
}
