package server

import (
	"net/http"
	"voldy/pkg/cluster"
	"voldy/pkg/store"
)

type RestServer struct {
	cluster       cluster.Cluster
	storageEngine store.StorageEngine
}

type GetRequest struct {
	Key     string `json:"key"`
	Version string `json:"version"`
}

type GetResponse struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	Version string `json:"version"`
}

func (rs *RestServer) Get(w http.ResponseWriter, r *http.Request) {

}

type PutRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type PutResponse struct {
	Version string `json:"version"`
}

func (rs *RestServer) Put(w http.ResponseWriter, r *http.Request) {

}

type DeleteRequest struct {
	Key     string `json:"key"`
	Version string `json:"version"`
}

type DeleteResponse struct {
	Deleted bool `json:"deleted"`
}

func (rs *RestServer) Delete(w http.ResponseWriter, r *http.Request) {

}
