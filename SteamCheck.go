package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/pkg/browser"
)

var (
	err error
)

func main() {

	var (
		file     *os.File
		reader   *bufio.Scanner
		key      string
		answer   string
		filename string
		start    time.Time
		res      []byte
	)

	fmt.Println("Enter the name of the text file containing your queries: ")

	_, _ = fmt.Scanln(&filename)

	file, err = os.Open(filename)

	if err != nil {
		fmt.Printf("Couldn't open file %v", err)

		os.Exit(-1)
	}

	defer file.Close()

	fmt.Print("Do you have a Steam API Key? [Y/N]: ")

	_, _ = fmt.Scanln(&answer)

	if answer == "n" || answer == "N" {

		err = browser.OpenURL("https://steamcommunity.com/dev/apikey")

		if err != nil {

			fmt.Printf("Couldn't open URL %v", err)
		}

		fmt.Println("Press enter after you've generated your key...")

		_, _ = fmt.Scanln()
	}

	fmt.Print("Enter your Steam API Key: ")

	_, _ = fmt.Scanln(&key)

	fmt.Println("Running, please wait...")

	reader = bufio.NewScanner(file)

	output, _ := os.Create("output.txt")

	defer output.Close()

	start = time.Now().UTC()

	counted := 0

	for reader.Scan() {

		res = readBytes("http://api.steampowered.com/ISteamUser/ResolveVanityURL/v0001/?key="+key+"&vanityurl=", reader.Text())

		counted++

		if bytes.Contains(res, []byte("No match")) {

			fmt.Printf("%v is available! Writing %v to output.txt...", reader.Text(), reader.Text())

			_, err = fmt.Fprintf(output, reader.Text()+"\n")

			if err != nil {
				fmt.Println("Couldn't write " + reader.Text() + " to output.txt")
			}
		} else {

			fmt.Printf("\n%v is unavailable", reader.Text())
		}
	}

	fmt.Printf("Made %v requests.\n", strconv.Itoa(counted))

	fmt.Println("Keep track of the number of requests you make!")

	fmt.Println("A single key can only make 100,000 requests every 24 hours!")

	fmt.Printf("Done. Took %v.\n", time.Since(start))

	fmt.Println("Press enter to exit...")

	_, _ = fmt.Scanln()
}

func readBytes(url string, line string) []byte {

	res, _ := http.Get(url + line)

	read, _ := ioutil.ReadAll(res.Body)

	err = res.Body.Close()

	if err != nil {
		fmt.Println(err)
	}

	return read
}
