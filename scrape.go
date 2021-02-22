package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

// Problem struct stores basic data of each problem
type Problem struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	URL        string `json:"url"`
	Difficulty string `json:"difficulty"`
}

// [TODO] figure out why separating tasks into own func breaks the action func
// yoinkCode is the opposite of yeetCode ... 🙃🥲
// func yoinkCode(nodes []*cdp.Node) chromedp.Tasks {
// 	fmt.Println("preparing to smol yoink..")
// 	return chromedp.Tasks{
// 		chromedp.Navigate(`https://leetcode.com/problemset/all/`),
//              // [TODO] fill in the tasks
// 	}
// }

// writeToFile creates output.json if the file doesn't exist, or appends if it does. Example from https://golang.org/pkg/os/#OpenFile
func writeToFile(json []byte) {
	filename := `output.json`
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := file.Write(json); err != nil {
		file.Close() // ignore error; Write error takes precedence
		log.Fatal(err)
	}
	if err := file.Close(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	// set chromedp options
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		// Set the headless flag to false to display the browser window
		chromedp.Flag("headless", true),
	)

	// "ExecAllocator is an Allocator which starts new browser processes on the host machine"
	fmt.Println("Creating new instance..")
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// Start new chrome instance / context
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	// Set timeout because reasons
	ctx, cancel = context.WithTimeout(ctx, 1*time.Minute)

	// Create arrays of cdproto type Node
	var rows, titles, links []*cdp.Node

	selector := `.reactable-data`
	titleSelector := `td[label='Title']`
	linkSelector := titleSelector + `> div > a`

	if err := chromedp.Run(ctx,
		chromedp.Navigate(`https://leetcode.com/problemset/all/`),

		// Starting smol..
		// chromedp.Navigate(`https://leetcode.com/problemset/all/?search=power%20of&difficulty=Easy`),

		// Wait until table data is visible
		chromedp.WaitVisible(".reactable-data", chromedp.ByQuery),

		// Scroll to bottom of the table loads (where you can find the show all & pagination)
		chromedp.ScrollIntoView(`.reactable-pagination`, chromedp.ByQuery),

		// Start By Getting Rows
		chromedp.Nodes(selector+`> tr`, &rows, chromedp.ByQueryAll),

		// Sidenote: Selector `td[label='Title'] > div > a:only-child` gets all the free/public problem links
		chromedp.Nodes(titleSelector, &titles, chromedp.ByQueryAll),
		chromedp.Nodes(linkSelector, &links, chromedp.ByQueryAll),
	); err != nil {
		fmt.Printf("something ducked up: %s", err)
		log.Fatal(err)
	}

	log.Printf("Found %d titles", len(rows))

	dataSelector := selector + `> tr:nth-child(%d) > td:nth-child(%d)`

	var problems []*Problem
	var num, title, difficulty, link string

	for i := 0; i < len(rows); i++ {

		numSelector := fmt.Sprintf(dataSelector, i+1, 2)
		difficultySelector := fmt.Sprintf(dataSelector+` > span`, i+1, 6)

		link = links[i].AttributeValue(`href`)
		// Gets Title
		title = titles[i].AttributeValue(`value`)

		if err := chromedp.Run(ctx,
			// Gets problem number
			chromedp.Text(numSelector, &num, chromedp.ByQuery),
			// Gets Difficulty
			chromedp.Text(difficultySelector, &difficulty, chromedp.ByQuery),
		); err != nil {
			log.Fatal(err)
		}

		problem := &Problem{
			ID:         num,
			Title:      title,
			URL:        link,
			Difficulty: difficulty,
		}
		// Append problem to array of problems
		problems = append(problems, problem)

		// Make the JSON pretty
		pJSON, err := json.MarshalIndent(problem, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Added %s", string(pJSON))
	}

	// Make the JSON behave
	jsonOutput, err := json.MarshalIndent(problems, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	writeToFile(jsonOutput)
}
