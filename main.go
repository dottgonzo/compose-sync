package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"gopkg.in/yaml.v2"
)

var composePathFromArgs = flag.String("compose", "docker-compose.yml", "docker compose path")
var composeLabelFromArgs = flag.String("label", "test", "docker compose name")

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
func main() {

	flag.Parse()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	cli, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Println(container.ID[:10] + " " + container.Image + " " + container.Labels["com.docker.compose.project"])
	}

	composePath := *composePathFromArgs
	content, content_err := ioutil.ReadFile(composePath)

	if content_err != nil {
		panic(content_err)
	}
	md5 := GetMD5Hash(string(content))
	println(md5)

	// do request with md5

	m := make(map[interface{}]interface{})

	err = yaml.Unmarshal([]byte(content), &m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- m:\n%v\n\n", m)

	for {
		sig := <-interrupt

		switch sig {
		case os.Interrupt:

			println("EXIT: exiting because of interrupt")

			os.Exit(0)
			return
		case syscall.SIGTERM:

			println("EXIT: exiting because of sigterm")

			os.Exit(0)
			return
		}
	}

}
