package main

import (
	"compress/flate"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mholt/archiver"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	if r.Method == "POST" {
		reader, err := r.MultipartReader()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//copy each part to destination.
		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}

			//if part.FileName() is empty, skip this iteration.
			if part.FileName() == "" {
				continue
			}
			runtime.GOMAXPROCS(runtime.NumCPU())

			// open the uploaded file

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			buff := make([]byte, 512) // why 512 bytes ? see http://golang.org/pkg/net/http/#DetectContentType
			_, err = part.Read(buff)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			filetype := http.DetectContentType(buff)

			fmt.Println(filetype)

			switch filetype {
			case "image/jpeg", "image/jpg", "image/png":
				fmt.Println(filetype)
			default:
				fmt.Println("unknown file type uploaded")
				fmt.Fprint(w, "Unknown file type. Please upload a png, jpg, jpeg, or csv.")
				continue
			}
			dst, err := os.Create("uploaded-images/" + strings.ToLower(part.FileName()))
			defer dst.Close()

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if _, err := io.Copy(dst, part); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
}
func ReturnFile(writer http.ResponseWriter, req *http.Request) {
	z := archiver.Zip{
		CompressionLevel:       flate.DefaultCompression,
		MkdirAll:               true,
		SelectiveCompression:   true,
		ContinueOnError:        false,
		OverwriteExisting:      false,
		ImplicitTopLevelFolder: false,
	}
	files := []string{
		"uploaded-images",
	}

	err := z.Archive(files, "dataset.zip")
	if err != nil {
		log.Fatal(err)
	}
	writer.Header().Set("Content-Disposition", "attachment; filename=dataset.zip")
	writer.Header().Set("Content-type", "application/zip")
	http.ServeFile(writer, req, "dataset.zip")
	//delete file from server once it has been served
	defer os.Remove("dataset.zip")
}
func setupRoutes() {
	router := mux.NewRouter()
	router.HandleFunc("/upload", uploadFile)
	router.HandleFunc("/download/zip", ReturnFile)
	http.ListenAndServe(":8080", router)
}

func main() {
	fmt.Println("Hello World")
	setupRoutes()
}