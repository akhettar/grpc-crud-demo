/*
 *
 * Copyright 2020 akhettar
 *
 * Licensed under the Apache License, Version MIT (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a simple gRPC client that demonstrates how to use gRPC-Go libraries
// to perform unary, client streaming, server streaming and full duplex RPCs.
//
// It interacts with the route guide service whose definition can be found in routeguide/route_guide.proto.
package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/akhettar/grpc-crud-demo/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
	"log"
	"time"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr         = flag.String("server_addr", "localhost:10000", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.youtube.com", "The server name used to verify the hostname returned by the TLS handshake")
)

func main() {
	var opts []grpc.DialOption
	if *tls {
		if *caFile == "" {
			*caFile = testdata.Path("ca.pem")
		}
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	opts = append(opts, grpc.WithBlock())

	// establish connection
	con, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("failed to connect to the server %v", err)
	}
	defer con.Close()

	// blog client
	client := api.NewBookstoreClient(con)

	// post a blog
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// creating shelves
	themes := []string{"science", "action", "history", "cooking", "sport", "wars", "travel", "astronomy", "religion", "finance", "health", "life"}
	for i := 1; i < 10; i++ {
		shelf := &api.Shelf{Theme: themes[i]}
		pub := api.CreateShelfRequest{Shelf: shelf}
		shelf, err := client.CreateShelf(ctx, &pub)
		if err != nil {
			log.Printf("error received for creating shelf %v", err)
		} else {
			log.Printf("shelf successfully created with id %d", shelf.Id)
		}
	}

	// creating books
	for i := 1; i < 10; i++ {
		b := &api.Book{Title: fmt.Sprintf("Title_%d", i), Author: "Ayache"}
		pub := api.CreateBookRequest{Book: b, Shelf: int64(i)}
		book, err := client.CreateBook(ctx, &pub)
		if err != nil {
			log.Printf("error received for creating books %v", err)
		} else {
			log.Printf("book successfully created with id %d", book.Id)
		}

	}

	// Querying books
	for i := 0; i < 20; i++ {
		read := api.GetBookRequest{Book: int64(i)}
		response, err := client.GetBook(ctx, &read)

		if err != nil {
			log.Printf("failed to query book with id %d, %v", i, err)
		} else {

			log.Printf("Book found for id %d and content %v", response.Id, response)
		}
	}
}
