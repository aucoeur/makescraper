package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

// Problem struct stores basic data of each problem
// type Problem struct {
// 	ID         string `json:"id"`
// 	Title      string `json:"title"`
// 	URL        string `json:"url"`
// 	Difficulty string `json:"difficulty"`
// }

// yoinkCode is the opposite of yeetCode ... 🙃🥲
// [TODO] figure out why separating tasks into own func breaks the action func
// func yoinkCode(nodes []*cdp.Node) chromedp.Tasks {
// 	fmt.Println("preparing to smol yoink..")
// 	return chromedp.Tasks{
// 		chromedp.Navigate(`https://leetcode.com/problemset/all/`),
//              // [TODO] fill in the tasks
// 	}
// }

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

	// ctx, cancel = context.WithTimeout(ctx, 30*time.Second)

	// Create arrays of cdproto type Node
	// var rows []*cdp.Node
	var links []*cdp.Node
	var titles []*cdp.Node

	// Run action/task list (yoinkcode)
	// if err := chromedp.Run(ctx, yoinkCode(nodes)); err != nil {
	// 	fmt.Printf("something ducked up: %s", err)
	// }

	// Title Selector
	selector := `td[label='Title']`
	linkSelector := selector + `> div > a`
	if err := chromedp.Run(ctx,
		// chromedp.Navigate(`https://leetcode.com/problemset/all/`),

		// Starting smol..
		chromedp.Navigate(`https://leetcode.com/problemset/all/?search=power%20of&difficulty=Easy`),

		// Wait until the bottom of the page loads (where you can find the show all & pagination)
		chromedp.WaitVisible(`.reactable-pagination`, chromedp.ByQuery),

		// Scrolls to footer
		chromedp.ScrollIntoView(`#footer-root`, chromedp.ByID),

		// Start with table rows
		// chromedp.Nodes(`.reactable-data > tr`, &rows, chromedp.ByQueryAll),

		// Selector `td[label='Title'] > div > a:only-child` gets all the free/public problem links
		chromedp.Nodes(linkSelector, &links, chromedp.ByQueryAll),

		// Get Titles
		chromedp.Nodes(selector, &titles, chromedp.ByQueryAll),
	); err != nil {
		fmt.Printf("something ducked up: %s", err)
	}

	log.Printf("Found %d titles", len(titles))

	// var ProblemSets []*Problem

	// loops through td nodes
	// siblingSelect := selector + `(%d).previousSibling`
	const childSelector = `.reactable-data > tr:nth-child(%d) > td:nth-child(%d)`

	var num, title, difficulty, link string

	for i := 0; i < len(titles); i++ {
		title = titles[i].AttributeValue(`value`)
		link = links[i].AttributeValue(`href`)

		if err := chromedp.Run(ctx,
			// Gets problem number
			chromedp.Text(fmt.Sprintf(childSelector, i+1, 2), &num, chromedp.ByQuery),
			// Gets Difficulty
			chromedp.Text(fmt.Sprintf(childSelector+`> span`, i+1, 6), &difficulty, chromedp.ByQuery),
		); err != nil {
			log.Fatal(err)
		}
		problem := &Problem{
			ID:         num,
			Title:      title,
			URL:        link,
			Difficulty: difficulty,
		}
		problemJSON, _ := json.MarshalIndent(problem, "", "  ")
		fmt.Println(string(problemJSON))
	}
}
