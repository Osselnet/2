package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
	"log"
)
var actions = []string{"logged in", "logged out", "created record", "deleted record", "updated account"}

type logItem struct {
	action    string
	timestamp time.Time
}

type User struct {
	id    int
	email string
	logs  []logItem
}

func (u User) getActivityInfo() string {
	output := fmt.Sprintf("UID: %d; Email: %s;\nActivity Log:\n", u.id, u.email)
	for index, item := range u.logs {
		output += fmt.Sprintf("%d. [%s] at %s\n", index, item.action, item.timestamp.Format(time.RFC3339))
	}

	return output
}

func main() {

	rand.Seed(time.Now().Unix())

	startTime := time.Now()

	const jobsCount, workwerCount = 100, 5
	jobs := make(chan int, jobsCount)
	users := make(chan User, jobsCount)

	for i := 0; i < workwerCount; i++ {
		go worker(jobs, users)
	}

	for i := 0; i < jobsCount; i++ {
		jobs <- i + 1
	}
	close(jobs)

	for user := range users {
		saveUserInfo(user)
	}

	fmt.Printf("DONE! Time Elapsed: %.2f seconds\n", time.Since(startTime).Seconds())
}

func worker(jobs <-chan int, users chan<- User){
	for range jobs {
		users <- generateUser(<-jobs)
	}
}

func saveUserInfo(user User) {
	fmt.Printf("WRITING FILE FOR UID %d\n", user.id)

	filename := fmt.Sprintf("users/uid%d.txt", user.id)
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}

	file.WriteString(user.getActivityInfo())
	time.Sleep(time.Second)
}

// func generateUsers(count int) []User {
// 	users := make([]User, count)

// 	for i := 0; i < count; i++ {
// 		users[i] = User{
// 			id:    i + 1,
// 			email: fmt.Sprintf("user%d@company.com", i+1),
// 			logs:  generateLogs(rand.Intn(1000)),
// 		}
// 		fmt.Printf("generated user %d\n", i+1)
// 		time.Sleep(time.Millisecond * 100)
// 	}

// 	return users
// }

func generateUser(id int) User {
		user := User{
			id:    id,
			email: fmt.Sprintf("user%d@company.com", id),
			logs:  generateLogs(rand.Intn(1000)),
		}
		fmt.Printf("generated user %d\n", id)
		time.Sleep(time.Millisecond * 100)
		return user
}

func generateLogs(count int) []logItem {
	logs := make([]logItem, count)

	for i := 0; i < count; i++ {
		logs[i] = logItem{
			action:    actions[rand.Intn(len(actions)-1)],
			timestamp: time.Now(),
		}
	}

	return logs
}