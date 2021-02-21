package main

import (
	"context"
	"fmt"
	"log"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

// // yoinkCode is the opposite of yeetCode ... 🙃🥲
// func yoinkCode(nodes []*cdp.Node) chromedp.Tasks {
// 	fmt.Println("preparing to smol yoink..")
// 	return chromedp.Tasks{
// 		// chromedp.Navigate(`https://leetcode.com/problemset/all/`),
// 		chromedp.Navigate(`https://leetcode.com/problemset/all/?search=two`),

// 		// Wait until the bottom of the page loads (where you can find the show all & pagination)
// 		chromedp.WaitVisible(`.reactable-pagination`),

// 		// Selector `td[label='Title'] > div > a:only-child` gets all the free/public problems
// 		chromedp.ScrollIntoView(`td[label='Title'] > div > a:only-child`),

// 		// chromedp.Nodes(`td[label='Title'] > div > a:only-child`, &nodes),
// 		chromedp.Nodes(`td[label='Title'] > div > a:only-child`, &nodes, chromedp.ByQuery),
// 	}
// }

func main() {
	// create chrome instance
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		// Set the headless flag to false to display the browser window
		chromedp.Flag("headless", false),
	)

	// Start new chrome instance / context
	fmt.Println("Creating new instance..")
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	// ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Create array of cdproto type Node
	var nodes []*cdp.Node

	// Run action/task list (yoinkcode)
	// err := chromedp.Run(context, yoinkCode(nodes))
	if err := chromedp.Run(ctx,
		// chromedp.Navigate(`https://leetcode.com/problemset/all/`),
		chromedp.Navigate(`https://leetcode.com/problemset/all/?search=two`),

		// Wait until the bottom of the page loads (where you can find the show all & pagination)
		chromedp.WaitVisible(`.reactable-pagination`, chromedp.ByQuery),

		// Selector `td[label='Title'] > div > a:only-child` gets all the free/public problems
		chromedp.ScrollIntoView(`#footer-root`, chromedp.ByID),

		// chromedp.Nodes(`td[label='Title'] > div > a:only-child`, &nodes),
		// chromedp.Nodes(`a`, &nodes),
		chromedp.Nodes(`td[label='Title'] > div > a:only-child`, &nodes, chromedp.ByQueryAll),
	); err != nil {
		fmt.Printf("something ducked up: %s", err)
	}

	// fmt.Println(nodes[0].AttributeValue("href"))
	for i := 0; i < len(nodes); i++ {
		log.Println(nodes[i].AttributeValue("href"))
	}
}
