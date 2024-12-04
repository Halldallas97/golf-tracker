package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Player struct {
	Name   string
	Scores []ScoreData
}
type ScoreData struct {
	Score  int
	Date   string
	Course string
}

/*
Function to display main menu all user responses once complete will return here.
*/
func main() {
	player := Player{}

	fmt.Println("Enter your name:")
	var name string
	fmt.Scanln(&name)
	player.Name = name

	for {
		fmt.Println("Please choose an option:\n" +
			"1. Add a new golf score\n" +
			"2. calculate your current handicap score\n" +
			"3. view all scores\n" +
			"Type 'Q' or 'q' to quit.")

		var input string
		fmt.Scanln(&input)

		switch input {

		case "1":
			score := getScoreInput()
			player.Scores = append(player.Scores, score)
			saveData(&player)

		case "2":
			loadData(&player)
			handicap := calculateHandicap(player.Scores)
			fmt.Println("Your average golf score is:", handicap)

		case "3":
			loadData(&player)
			displayScores(&player)

		case "Q", "q":
			saveData(&player)
			fmt.Println("Exiting... Data saved.")
			return
		}
	}
}

/*
Simple function to display in the terminal all of the players scores.
*/
func displayScores(player *Player) {
	fmt.Printf("\nScore data for %s \n", player.Name)
	for i, score := range player.Scores {
		fmt.Printf("%d. Score: %d, Course: %s, Date: %s\n", i+1, score.Score, score.Course, score.Date)
	}
}

/*
Simple mathematics function to calculate the average score of all your golf games.
will only work if there are at min of three scores. The avg is calculated from the players top 3 scores.
*/
func calculateHandicap(scoreData []ScoreData) float64 {
	//sort the top 3 scores
	sort.Slice(scoreData, func(i, j int) bool {
		return scoreData[i].Score > scoreData[j].Score
	})
	if len(scoreData) < 3 {
		fmt.Println("Not enough scores to calculate handicap. At least 3 scores are needed.")
		return 0.0
	}

	handicap := float64(scoreData[0].Score+scoreData[1].Score+scoreData[2].Score) / 3
	return handicap
}

/*
User input function to get the data from the user via the terminal,
The data is validated and sent to the savedata method after this function is
called.
*/
func getScoreInput() ScoreData {
	var scoreData ScoreData

	fmt.Println("Enter your score:")
	var score string
	fmt.Scanln(&score)
	scoreInt, err := strconv.Atoi(score)
	if err != nil {
		fmt.Println("Error: could not convert score to int:", err)
	}
	scoreData.Score = scoreInt

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter the course:")
	course, _ := reader.ReadString('\n')
	scoreData.Course = course[:len(course)-1]

	scoreData.Date = getValidDateInput(reader)

	return scoreData
}

/*
To keep dates proper this function will validate
using regular expression the date that the user enters.
*/
func getValidDateInput(reader *bufio.Reader) string {
	var date string
	datePattern := `^\d{4}-\d{2}-\d{2}$`
	re := regexp.MustCompile(datePattern)

	for {
		fmt.Println("Enter date played in yyyy-mm-dd format:")
		date, _ = reader.ReadString('\n')
		date = strings.TrimSpace(date)

		if re.MatchString(date) {
			break
		} else {
			fmt.Println("Invalid date format. Please enter the date in yyyy-mm-dd format.")
		}
	}

	return date
}

/*
The purpose of this function is to identify if the player-name.csv already exists
if the csv exists the function will attempt to read the file. It will update the player
struct with the current values of the csv.
*/
func loadData(player *Player) {
	filename := player.Name + ".csv"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("No file found at this time.")
		player.Scores = nil
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading the CSV file:", err)
		player.Scores = nil
		return
	}

	var loadedScores []ScoreData
	//ignoring headers
	for i, record := range records {
		if i == 0 {
			continue
		}
		if len(record) < 3 {
			continue
		}

		score, err := strconv.Atoi(record[0])
		if err != nil {
			fmt.Println("Error converting score to int:", err)
			continue
		}
		scoreData := ScoreData{
			Score:  score,
			Course: strings.TrimSpace(record[1]),
			Date:   strings.TrimSpace(record[2]),
		}
		loadedScores = append(loadedScores, scoreData)
	}

	player.Scores = loadedScores
	fmt.Println("Data loaded successfully.")
}

/*
This function aims to save the players data
first the function will check to see if the players-name.csv exist
if it does not exist the function will create the file.
once the file is created it will create the headers for the data.
if a file already exists it will not save the headers.
Then the func will save the players data into the csv. and return to the main menu.
*/
func saveData(player *Player) {
	filename := player.Name + ".csv"
	// check existing scores
	existingScores := map[string]bool{}

	file, err := os.Open(filename)
	if err == nil {
		defer file.Close()
		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err == nil {
			//skip the header
			for _, record := range records[1:] {
				if len(record) < 3 {
					continue
				}
				key := record[0] + record[1] + record[2]
				existingScores[key] = true
			}
		}
	}

	// Open the file for appending new scores
	file, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// create headers for new files.
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return
	}
	if fileInfo.Size() == 0 {
		headers := []string{"Score", "Course", "Date"}
		if err := writer.Write(headers); err != nil {
			fmt.Println("Error writing headers to CSV:", err)
			return
		}
	}

	// Write only new scores
	for _, score := range player.Scores {
		record := []string{
			strconv.Itoa(score.Score),
			strings.TrimSpace(score.Course),
			strings.TrimSpace(score.Date),
		}
		key := record[0] + record[1] + record[2]
		if !existingScores[key] {
			if err := writer.Write(record); err != nil {
				fmt.Println("Error writing record to CSV:", err)
				return
			}
			existingScores[key] = true
		}
	}
	fmt.Println("Scores saved to file titled " + player.Name + ".csv in this directory!")
}
