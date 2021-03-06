package lib

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"cloud.google.com/go/storage"
)

// Sets the name for the new bucket.
const bucketName = "gocity-cache"

type Storage interface {
	Get(projectName string) (bool, []byte, error)
	Save(projectName string, content []byte) error
	Delete(projectName string) error
}

type GCS struct {
	ctx    context.Context
	client *storage.Client
}

func NewGCS(ctx context.Context) (Storage, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	return &GCS{
		ctx:    ctx,
		client: client,
	}, nil
}

func getObjectName(name string) string {
	return fmt.Sprintf("%s.json", name)
}

func (g *GCS) Get(projectName string) (bool, []byte, error) {
	object := g.client.Bucket(bucketName).Object(getObjectName(projectName))
	if object == nil {
		return false, nil, nil
	}

	reader, err := object.NewReader(g.ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			return false, nil, nil
		}

		return false, nil, err
	}
	defer reader.Close()

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return false, nil, err
	}

	return true, data, nil
}

func (g *GCS) Save(projectName string, content []byte) error {
	client, err := storage.NewClient(g.ctx)
	if err != nil {
		return err
	}

	buffer := bytes.NewBuffer(content)
	wc := client.Bucket(bucketName).Object(getObjectName(projectName)).NewWriter(g.ctx)
	if _, err = io.Copy(wc, buffer); err != nil {
		return err
	}

	return wc.Close()
}

func (g *GCS) Delete(projectName string) error {
	return nil
}
