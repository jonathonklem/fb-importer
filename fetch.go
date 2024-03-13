package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/adlio/trello"
	"github.com/go-sql-driver/mysql"
)

func fetch() {
	data := getFacebookFeed()

	cfg := mysql.Config{
		User:                 os.Getenv("MYSQL_USER"),
		Passwd:               os.Getenv("MYSQL_PASSWORD"),
		Net:                  "tcp",
		Addr:                 os.Getenv("MYSQL_ADDR") + ":3306",
		DBName:               os.Getenv("MYSQL_DATABASE"),
		AllowNativePasswords: true,
	}

	var db *sql.DB

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	pingErr := db.Ping()
	if pingErr != nil {
		fmt.Printf("Error: %s\n", pingErr)
		return
	}

	var unimportable []string
	var unupdatable []string

	current := time.Now()
	stringDate := current.Format("01-02-2006")
	stringDate = stringDate[:5] + ` ` + strconv.Itoa(current.Hour())

	for _, element := range data.Data {
		// source_id = 5090 (youarelinked facebook)
		// url = element.attachments.data[0].url
		// description = element.attachments.data[0].description
		// meta_description = element.attachments.data[0].description
		// thumbnail = element.attachments.data[0].media.src
		// user_id = 36 == fb_import

		featured := 0
		show_in_banner := 0

		var exists bool
		err = db.QueryRow("SELECT IF (COUNT(*), 'true', 'false') FROM fb_ids WHERE value = ?", element.ID).Scan(&exists)
		if err != nil {
			fmt.Printf("Error checking duplicate")
			return
		}

		if exists {
			fmt.Println("Skipping already imported: " + element.ID)
			continue
		} else {
			db.Exec("INSERT INTO fb_ids (value) VALUES (?)", element.ID)
		}

		var descriptionToInsert string
		if element.Attachments.Data[0].Description == "" {
			descriptionToInsert = strings.Title(element.Message)
		} else {
			// use text/cases to title case the description

			descriptionToInsert = strings.Title(element.Attachments.Data[0].Description)
		}

		_, insertErr := db.Exec("INSERT INTO resources (is_featured, show_in_banner, resource_type_id, URL, description, meta_description, thumbnail, user_id, facebook, is_dead, fb_id) VALUES (?, ?, 4, ?, ?, ?, ?, 36, ?, 1, ?)", featured, show_in_banner, element.Attachments.Data[0].URL, descriptionToInsert, element.Attachments.Data[0].Description, element.Attachments.Data[0].Media.Image.Src, element.Message, element.ID)

		if insertErr != nil {
			fmt.Printf("Error: %s", insertErr)
			return
		}

		// mark it as posted
		message := element.Message + ` SAVED ` + stringDate
		if !strings.Contains(element.Message, "SAVED") {
			if updatePost(element.ID, message) != "200 OK" {
				response := updatePost(element.ID, message)
				if response != "200 OK" {
					unupdatable = append(unupdatable, element.ID+" "+response+" - "+element.PermalinkURL)
				}
			}
		}
	}

	runTrello(unupdatable, unimportable, stringDate)
}

func runTrello(uu []string, ui []string, stringDate string) {
	client := trello.NewClient(os.Getenv("TRELLO_KEY"), os.Getenv("TRELLO_TOKEN"))

	list, _ := client.GetList(os.Getenv("TRELLO_LIST"), trello.Defaults())
	if len(ui) > 0 {
		list.AddCard(&trello.Card{Name: "Unable to Import " + stringDate, Desc: "The following were not able to be imported due to permissions: " + strings.Join(ui, ", ")}, trello.Defaults())
	}

	if len(uu) > 0 {
		list.AddCard(&trello.Card{Name: "Unable to Update " + stringDate, Desc: "The following were imported successfully but could not be marked as such: " + strings.Join(uu, ", ")}, trello.Defaults())
	}

}

func updatePost(postId string, newMessage string) string {
	url := "https://graph.facebook.com/v18.0/" + postId + "?message=" + url.QueryEscape(newMessage)

	token := getToken()
	bearer := "Bearer " + token

	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	return resp.Status
}

func getFacebookFeed() FacebookPosts {
	url := os.Getenv("FB_FEED_URL")

	token := getToken()
	bearer := "Bearer " + token

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n", err)
		os.Exit(1)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading body.\n", err)
		os.Exit(1)
	}

	var data FacebookPosts

	err = json.Unmarshal([]byte(body), &data)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	return data
}
