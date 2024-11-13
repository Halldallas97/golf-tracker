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
			scores := loadData(&player)
			player.Scores = append(player.Scores, scores...)
			handicap := calculateHandicap(player.Scores)
			fmt.Println("Your handicap score is: ", handicap)

		case "3":
			scores := loadData(&player)
			player.Scores = append(player.Scores, scores...)
			displayScores(&player)

		case "Q", "q":
			saveData(&player)
			fmt.Println("Exiting... Data saved.")
			return
		}

	}
}

func displayScores(player *Player) {
	fmt.Printf("\nScore data for %s \n", player.Name)
	for i, score := range player.Scores {
		fmt.Printf("%d. Score: %d, Course: %s, Date: %s\n", i+1, score.Score, score.Course, score.Date)
	}
}

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

func loadData(player *Player) []ScoreData {
	filename := player.Name + ".csv"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("No file found at this time.")
		return nil
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading the CSV file:", err)
		return nil
	}

	var loadedScores []ScoreData
	for i, record := range records {
		//skiping the header
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
			Course: record[1],
			Date:   record[2],
		}
		loadedScores = append(loadedScores, scoreData)
	}

	fmt.Println("Data loaded successfully.")
	return loadedScores
}

func saveData(player *Player) {
	file, err := os.OpenFile(player.Name+".csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

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

	for _, score := range player.Scores {
		record := []string{
			strconv.Itoa(score.Score),
			score.Course,
			score.Date,
		}
		if err := writer.Write(record); err != nil {
			fmt.Println("Error writing record to CSV:", err)
			return
		}
	}

	fmt.Println("Scores appended to file titled " + player.Name + ".csv in this directory!")
}
