package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

type Follower struct {
	Username   string `json:"username"`
	FollowedOn string `json:"followed_on"`
}

type Subscriber struct {
	Username      string `json:"username"`
	User          string `json:"user"`
	AmountCents   int    `json:"amount_cents"`
	AmountDollars int    `json:"amount_dollars"`
	SubscribedOn  string `json:"subscribed_on"`
}

type Livestream struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	CreatedOn   string `json:"created_on"`
	IsLive      bool   `json:"is_live"`
	StreamKey   string `json:"stream_key"`
	Likes       int    `json:"likes"`
	Dislikes    int    `json:"dislikes"`
	WatchingNow int    `json:"watching_now"`
}

type ResponseData struct {
	Now           int64       `json:"now"`
	Type          string      `json:"type"`
	UserID        string      `json:"user_id"`
	ChannelID     interface{} `json:"channel_id"`
	Since         interface{} `json:"since"`
	MaxNumResults int         `json:"max_num_results"`
	Followers     struct {
		NumFollowers      int        `json:"num_followers"`
		NumFollowersTotal int        `json:"num_followers_total"`
		LatestFollower    Follower   `json:"latest_follower"`
		RecentFollowers   []Follower `json:"recent_followers"`
	} `json:"followers"`
	Subscribers struct {
		NumSubscribers    int          `json:"num_subscribers"`
		LatestSubscriber  Subscriber   `json:"latest_subscriber"`
		RecentSubscribers []Subscriber `json:"recent_subscribers"`
	} `json:"subscribers"`
	Livestreams []Livestream `json:"livestreams"`
}

func fetchRumbleData() (*ResponseData, error) {
	url := "https://rumble.com/-livestream-api/get-data?key=0me1clTKb6jODPcPvZTqj6SA_5cervHSucsmVpUEJqOIj3g3PpsGNQJWfGVp7wJ8wN2CwTutNfoc33tzqXtLVA"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var responseData ResponseData
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	return &responseData, nil
}

func main() {
	// Load templates using Fiber's HTML engine
	engine := html.New("./views", ".html")

	// Initialize Fiber with the template engine
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Serve static files (like CSS)
	app.Static("/static", "./static")

	// Render the landing page with dynamic content
	app.Get("/", func(c *fiber.Ctx) error {
		// Define title
		title := "My Landing Page"

		// Load body content from file
		bodyContent, err := os.ReadFile("./content/bodycontent.html")
		if err != nil {
			log.Fatalf("Error reading body content file: %v", err)
		}

		// Render the page with dynamic content
		return c.Render("index", fiber.Map{
			"Title":       title,
			"BodyContent": template.HTML(bodyContent), // Prevent escaping HTML content
		})
	})

	// Create a new endpoint to fetch Rumble data
	app.Get("/rumble", func(c *fiber.Ctx) error {
		responseData, err := fetchRumbleData()
		if err != nil {
			return c.Status(500).SendString(fmt.Sprintf("Failed to fetch Rumble data: %v", err))
		}

		return c.JSON(responseData)
	})

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the server
	log.Printf("Server started on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
