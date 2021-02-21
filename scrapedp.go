package main

import (
	"context"
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
	var rows []*cdp.Node
	// var nodes []*cdp.Node
	var titles []*cdp.Node

	// Run action/task list (yoinkcode)
	// if err := chromedp.Run(ctx, yoinkCode(nodes)); err != nil {
	// 	fmt.Printf("something ducked up: %s", err)
	// }

	if err := chromedp.Run(ctx,
		// chromedp.Navigate(`https://leetcode.com/problemset/all/`),

		// Starting smol..
		chromedp.Navigate(`https://leetcode.com/problemset/all/?search=power%20of&difficulty=Easy`),

		// Wait until the bottom of the page loads (where you can find the show all & pagination)
		chromedp.WaitVisible(`.reactable-pagination`, chromedp.ByQuery),

		// Scrolls to footer
		chromedp.ScrollIntoView(`#footer-root`, chromedp.ByID),

		// Start with table rows
		chromedp.Nodes(`.reactable-data > tr`, &rows, chromedp.ByQueryAll),

		// prints out: 2021/02/21 13:32:44 map[label:Title value:Consecutive Characters]
		chromedp.Value(`td[label='Title']`, &titles, chromedp.ByQueryAll),

		// Selector `td[label='Title'] > div > a:only-child` gets all the free/public problem links
		// chromedp.Nodes(`td[label='Title'] > div > a:only-child`, &nodes, chromedp.ByQueryAll),
	); err != nil {
		fmt.Printf("something ducked up: %s", err)
	}

	log.Println("Found rows", len(rows))

	// var ProblemSets []*Problem

	// loops through td nodes
	const childSelector = `.reactable-data > tr:nth-child(%d) > td:nth-child(2)`
	var col string
	for i := 0; i < len(rows); i++ {
		log.Println(nodes[i].AttributeValue("href"))
		if err := chromedp.Run(ctx,
			chromedp.Text(fmt.Sprintf(childSelector, i+1), &col),
		); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("row %d: value = %s\n", i, col)
		// log.Println(titles[i])
	}
}
