package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type Netflix struct {
	Movie   string `json:"movie"`
	Watched bool   `json:"watched"`
}

func displayMenu() {
	fmt.Println("\nPlease choose a function to run:")
	fmt.Println("1. FetchAllMoviesGet()")
	fmt.Println("2. CreateAMoviePost()")
	fmt.Println("3. MarkAsWatchedPut()")
	fmt.Println("4. DeleteAMovie()")
	fmt.Println("5. DeleteAllMovie()")
	fmt.Println("6. Exit")
}

func main() {
	fmt.Println("Welcome To API Consumer")

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("-----------------------------")
		displayMenu()
		fmt.Println("-----------------------------")
		fmt.Print("Enter your choice: ")

		input, err := reader.ReadString('\n')
		errorHandle(err)
		input = strings.TrimSpace(input)
		choice, err := strconv.Atoi(input)
		errorHandle(err)

		switch choice {
		case 1:
			FetchAllMovies()
		case 2:
			CreateAMovie()
		case 3:
			MarkAsWatched()
		case 4:
			DeleteAMovie()
		case 5:
			DeleteAllMovies()
		case 6:
			fmt.Println("Exiting the program, goodbye!")
			return
		default:
			fmt.Println("Invalid Choice")
		}

	}

}

func FetchAllMovies() {
	const apiUrl = "http://localhost:4000/api/movies"

	response, err := http.Get(apiUrl)
	errorHandle(err)
	defer response.Body.Close()

	fmt.Println("StatusCode: ", response.StatusCode)
	fmt.Println("ContentLength: ", response.ContentLength)

	content, err := io.ReadAll(response.Body)
	errorHandle(err)

	var data interface{}
	err = json.Unmarshal(content, &data)
	errorHandle(err)

	formatedJSON, err := json.MarshalIndent(data, "", "  ")
	errorHandle(err)
	fmt.Println("Server Response:")
	fmt.Println(string(formatedJSON))

}

func CreateAMovie() {
	const apiUrl = "http://localhost:4000/api/movie"

	myMovie := Netflix{
		Movie: "Dhoom", Watched: false,
	}

	finalJSON, err := json.Marshal(myMovie)
	errorHandle(err)
	fmt.Printf("JSON Data ready to send: %s\n", finalJSON)
	response, err := http.Post(apiUrl, "application/json", bytes.NewBuffer(finalJSON))
	errorHandle(err)
	defer response.Body.Close()

	content, err := io.ReadAll(response.Body)
	errorHandle(err)

	var responseString strings.Builder

	byteCount, err := responseString.Write(content)
	errorHandle(err)

	fmt.Println("Response Length: ", byteCount)
	fmt.Println("Server Response: ", responseString.String())
}

func MarkAsWatched() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter Movie ID to Modify: ")
	input, err := reader.ReadString('\n')
	errorHandle(err)
	input = strings.TrimSpace(input)
	partsofUrl := &url.URL{
		Scheme: "http",
		Host:   "localhost:4000",
		Path:   "/api/movie/" + input,
	}

	apiUrl := partsofUrl.String()

	request, err := http.NewRequest(http.MethodPut, apiUrl, nil)
	errorHandle(err)

	request.Header.Set("Content-Type", "application/json")
	//use http Client to send the PUT request
	client := &http.Client{}
	response, err := client.Do(request)
	errorHandle(err)
	defer response.Body.Close()
	content, err := io.ReadAll(response.Body)
	errorHandle(err)

	var responseString strings.Builder

	byteCount, err := responseString.Write(content)
	errorHandle(err)

	fmt.Println("Response Length: ", byteCount)
	fmt.Println("Server Response: ", responseString.String())

}

func DeleteAMovie() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter Movie ID to Delete: ")
	input, err := reader.ReadString('\n')
	errorHandle(err)
	input = strings.TrimSpace(input)
	partsofUrl := &url.URL{
		Scheme: "http",
		Host:   "localhost:4000",
		Path:   "/api/movie/" + input,
	}
	apiUrl := partsofUrl.String()

	request, err := http.NewRequest(http.MethodDelete, apiUrl, nil)
	errorHandle(err)
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	errorHandle(err)

	defer response.Body.Close()
	content, err := io.ReadAll(response.Body)
	errorHandle(err)
	fmt.Println(string(content))
}

func DeleteAllMovies() {
	const apiUrl = "http://localhost:4000/api/deleteallmovie"
	request, err := http.NewRequest(http.MethodDelete, apiUrl, nil)
	errorHandle(err)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	errorHandle(err)
	defer response.Body.Close()

	content, err := io.ReadAll(response.Body)
	errorHandle(err)
	fmt.Println("Server Response: ", string(content))
}
func errorHandle(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}
