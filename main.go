package main

import (
	"html/template"
	//"io/ioutil"
	"net/http"
	"os"
	"fmt"
	"io"
	"strconv"
)

func indexHandle(w http.ResponseWriter, r *http.Request) {
	uploadTemplate := template.Must(template.ParseFiles("index.html"))
	uploadTemplate.Execute(w, nil)
}

func uploadHandle(w http.ResponseWriter, r *http.Request) {
	file, head, _ := r.FormFile("file")
	defer file.Close()
	/* avoid to use ReadAll if file is large and more client to download
	bytes, _ := ioutil.ReadAll(file)
	w.Write(bytes)
	ioutil.WriteFile(head.Filename, bytes, os.ModeAppend)
	*/
	f, err := os.OpenFile("./test/" + head.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)
}

func downloadHandle(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	if filename == "" {
		//Get not set, send a 400 bad request
		http.Error(w, "Get 'file' not specified in url.", 400)
		return
	}

	fmt.Println("client request: " + filename)
	f, err := os.Open("./test/" + filename)
	defer f.Close()
	if err != nil {
		fmt.Println(err)

	}

	//File is found, create and send the correct headers

	//Get the Content-Type of the file
	//Create a buffer to store the header of the file in
	fileHeader := make([]byte, 512)
	//Copy the headers into the FileHeader buffer
	f.Read(fileHeader)
	//Get content type of file
	fileContentType := http.DetectContentType(fileHeader)

	//Get the file size
	fileStat, _ := f.Stat()                     //Get info from file
	fileSize := strconv.FormatInt(fileStat.Size(), 10) //Get file size as a string

	//Send the headers
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", fileContentType)
	w.Header().Set("Content-Length", fileSize)
	fmt.Println(fileContentType)
	fmt.Println(fileSize)
	//Send the file
	//We read 512 bytes from the file already so we reset the offset back to 0
	f.Seek(0, 0)
	io.Copy(w, f) //'Copy' the file to the client
	return
}

func main() {
	http.HandleFunc("/", indexHandle)
	http.HandleFunc("/upload", uploadHandle) //127.0.0.1:9090/upload
	http.HandleFunc("/download", downloadHandle) //27.0.0.1:9090/download?file=test.txt
	fmt.Println(http.ListenAndServe(":9090", nil))
}