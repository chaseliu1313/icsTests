package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool
var fNames []string

func formatICalTime(t time.Time) string {
	return t.UTC().Format("20060102T150405Z")
}

func connectToKSMDB() (*pgxpool.Pool, error) {
	var connectStr = "postgres://goServer:goServer123@localhost:5431/postgres?sslmode=disable"
	pool, err := pgxpool.New(context.Background(), connectStr)
	log.Println("connecting to database: ", connectStr)
	if err != nil {
		//log.Println(os.Stderr)
		log.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
		pool.Close()
		return nil, err
	}

	rows, err := pool.Query(context.Background(), "SELECT first_name FROM guest")

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var firstName string
		err := rows.Scan(&firstName)
		if err != nil {
			log.Fatal(err)
		}
		fNames = append(fNames, firstName)
	}

	Pool = pool
	return pool, nil
}

func getAPTType() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomInt := r.Intn(6) + 1 // generates 1 to 12
	return fmt.Sprintf("type-%d", randomInt)
}

func getLastRandomLName() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	letterCode := 'A' + r.Intn(26) // generates a rune between 'A' and 'Z'
	letter := string(rune(letterCode))

	return letter + "."
}

func getRandomFName() string {
	if len(fNames) == 0 {
		log.Println("first names not loaded")
		return ""
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomInt := r.Intn(len(fNames))
	return fNames[randomInt]
}

func generateFile(fileName string) {
	location := time.Now().Location()
	startOfDay := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 8, 0, 0, 0, location)
	endOfDay := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 17, 0, 0, 0, location)

	interval := 15 * time.Minute

	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("BEGIN:VCALENDAR\n")
	f.WriteString("VERSION:2.0\n")
	f.WriteString("CALSCALE:GREGORIAN\n")

	for current := startOfDay; current.Before(endOfDay); current = current.Add(interval) {
		summry := getRandomFName() + " " + getLastRandomLName() + " (" + getAPTType() + ")"
		end := current.Add(interval)
		f.WriteString("BEGIN:VEVENT\n")
		f.WriteString(fmt.Sprintf("UID:%d@kwaka.ca\n", current.UnixNano()))
		f.WriteString(fmt.Sprintf("DTSTAMP:%s\n", formatICalTime(time.Now())))
		f.WriteString(fmt.Sprintf("DTSTART:%s\n", formatICalTime(current)))
		f.WriteString(fmt.Sprintf("DTEND:%s\n", formatICalTime(end)))
		f.WriteString(fmt.Sprintf("SUMMARY:%s\n", summry))
		f.WriteString("STATUS:CONFIRMED\n")
		f.WriteString("DESCRIPTION:15-minute slot\n")
		f.WriteString("END:VEVENT\n")
	}

	f.WriteString("END:VCALENDAR\n")

	fmt.Println("ICS file generated as", fileName)
}

func main() {

	var fileNames = []string{"link1.ics", "link2.ics", "link3.ics"}
	connectToKSMDB()

	for _, fName := range fileNames {

		generateFile(fName)

	}

}
