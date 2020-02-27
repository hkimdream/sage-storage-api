package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v6"
)

func listByBucket(w http.ResponseWriter, r *http.Request) {
	buckets, err := minioClient.ListBuckets()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, bucket := range buckets {
		fmt.Fprintf(w, "Bucket: %v\n", bucket)
	}
	w.WriteHeader(http.StatusOK)
}

func getObjectFromBucket(w http.ResponseWriter, r *http.Request) {
	//Parse
	pathParams := mux.Vars(r)
	bucketName := pathParams["bucket"]
	objectName := pathParams["object"]
	//Check if Bucket exist
	exists, errBucketExists := minioClient.BucketExists(bucketName)
	if !exists && errBucketExists != nil {
		log.Printf("Bucket: %s does not exist.", bucketName)
		log.Fatalln(errBucketExists)
	}
	object, err := minioClient.GetObject(bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}
	localFile, err := os.Create(filePath + objectName)
	if err != nil {
		fmt.Println(err)
		return
	}
	if _, err = io.Copy(localFile, object); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintf(w, "Download - Bucket: %v and Object: %v\n", bucketName, objectName)
	fmt.Fprintf(w, "Destination: %v\n", filePath+objectName)
	w.WriteHeader(http.StatusOK)
}
