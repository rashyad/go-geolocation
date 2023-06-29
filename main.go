package main

import (
	"fmt"
	"io"
	"net/http"
	"errors"
	"os"
	"bytes"
	"encoding/json"
)

func ReadUserIP(r *http.Request) string {
    IPAddress := r.Header.Get("X-Real-Ip")
    if IPAddress == "" {
        IPAddress = r.Header.Get("X-Forwarded-For")
    }
    if IPAddress == "" {
        IPAddress = r.RemoteAddr
    }
    return IPAddress
}


func formatJSON(data []byte) string {
    var out bytes.Buffer
    err := json.Indent(&out, data, "", " ")

    if err != nil {
        fmt.Println(err)
    }

    d := out.Bytes()
    return string(d)
}

func getGeoLocation(userIp string) string{


	apiUrl := "http://ipwho.is/" + userIp

	request, error := http.NewRequest("GET", apiUrl, nil)
	request.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	response, error := client.Do(request)

	if error != nil {
        fmt.Println(error)
    }

    responseBody, error := io.ReadAll(response.Body)

    if error != nil {
        fmt.Println(error)
    }

    formattedData := formatJSON(responseBody)

    // clean up memory after execution
    defer response.Body.Close()

	return formattedData
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, ReadUserIP(r))
	io.WriteString(w, "\n")
	io.WriteString(w, getGeoLocation(ReadUserIP(r)))

}

func main() {
	http.HandleFunc("/", getRoot)

	err := http.ListenAndServe(":3333", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}

}