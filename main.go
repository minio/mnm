package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"path"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/minio/cli"
	minio "github.com/minio/minio-go"
	homedir "github.com/mitchellh/go-homedir"
)

// Endpoint - represents and S3 endpoint.
type Endpoint struct {
	URL       string `json:"url"`
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}

// Config - represents on disk config file.
type Config struct {
	Auth      string     `json:"auth"`
	Endpoints []Endpoint `json:"endpoints"`
}

// Holds all the info needed for routing.
type ming struct {
	hashMap map[string]*minio.Client
	hashes  []string
	address string
	auth    string
}

// Handles PUT requests.
func (m ming) put(w http.ResponseWriter, r *http.Request) {
	auth := r.URL.Query().Get("auth")
	object := mux.Vars(r)["object"]
	randIdx := rand.Intn(len(m.hashes))
	client := m.hashMap[m.hashes[randIdx]]

	if m.auth != "" && m.auth != auth {
		// If auth configured and the request is not authenticated.
		w.WriteHeader(403)
		return
	}

	_, err := client.PutObject("ming", object, r.Body, "application/octet-stream")
	if err != nil {
		w.WriteHeader(500)
		if errResp, ok := err.(minio.ErrorResponse); ok {
			w.Write([]byte(errResp.Message))
		}
		return
	}
	// Return the URL with hash info which will be used to fetch the object during
	// GET request.
	respStr := fmt.Sprintf("%s/ming/%s/%s", m.address, m.hashes[randIdx], object)
	w.Write([]byte(respStr))
}

// Handles GET request.
func (m ming) get(w http.ResponseWriter, r *http.Request) {
	hash := mux.Vars(r)["hash"]
	object := mux.Vars(r)["object"]
	auth := r.URL.Query().Get("auth")

	if m.auth != "" && m.auth != auth {
		// If auth configured and the request is not authenticated.
		w.WriteHeader(403)
		return
	}
	client := m.hashMap[hash]
	reader, err := client.GetObject("ming", object)

	if err != nil {
		w.WriteHeader(500)
		if errResp, ok := err.(minio.ErrorResponse); ok {
			w.Write([]byte(errResp.Message))
		}
		return
	}
	_, err = io.Copy(w, reader)
	if err != nil {
		// In case no data was written, return error.
		w.WriteHeader(500)
	}
}

// Parse command line and start the ming server.
func startMing(ctx *cli.Context) error {
	confPath := path.Join(ctx.GlobalString("config-dir"), "config.json")

	configContents, err := ioutil.ReadFile(confPath)
	if err != nil {
		log.Fatal("Unable to read configuration file: ", err)
	}

	config := Config{}
	err = json.Unmarshal(configContents, &config)
	if err != nil {
		log.Fatal("Unable to parse configuration file: ", err)
	}

	address := ctx.GlobalString("address")
	m := ming{hashMap: make(map[string]*minio.Client), address: address, auth: config.Auth}

	for _, host := range config.Endpoints {
		url, err := url.Parse(host.URL)
		if err != nil {
			log.Fatalf("Unable to parse %s: %s", host.URL, err)
		}
		client, err := minio.New(url.Host, host.AccessKey, host.SecretKey, url.Scheme == "https")
		if err != nil {
			log.Fatal("Minio client init failed: ", err)
		}
		sum := fmt.Sprintf("%x", md5.Sum([]byte(url.Host)))
		m.hashMap[sum] = client
		m.hashes = append(m.hashes, sum)
	}

	r := mux.NewRouter()
	r.Methods("GET").Path("/ming/{hash}/{object:.+}").Queries("auth", "{auth:.*}").HandlerFunc(m.get)
	r.Methods("GET").Path("/ming/{hash}/{object:.+}").HandlerFunc(m.get)
	r.Methods("PUT").Path("/ming/{object:.+}").Queries("auth", "{auth:.*}").HandlerFunc(m.put)
	r.Methods("PUT").Path("/ming/{object:.+}").HandlerFunc(m.put)

	return http.ListenAndServe(address, r)
}

const globalMingConfigDir = ".ming"

// getConfigPath get server config path
func mustGetConfigDir() string {
	homeDir, err := homedir.Dir()
	if err != nil {
		log.Fatal("Unable to get the home directory", err)
	}
	return filepath.Join(homeDir, globalMingConfigDir)
}

func main() {
	app := cli.NewApp()
	app.Usage = "Minio Gateway"
	app.Author = "Minio.io"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "address",
			Value: "localhost:8000",
			Usage: "Local bind address",
		},
		cli.StringFlag{
			Name:  "config-dir, C",
			Value: mustGetConfigDir(),
			Usage: "Path to the configuration directory",
		},
	}
	app.Before = startMing
	app.RunAndExitOnError()
}
