package main

import (
	"errors"
	"fmt"
	"github.com/pitkve1-wq/internal/database"
	"github.com/pitkve1-wq/internal/config"
	"time"
	"context"
	"github.com/google/uuid"
	"database/sql"
)

type state struct {
	sta *config.Config
	db  *database.Queries
}

type command struct {
	name string
	args []string
}

type commands struct {
	a map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	v, ok := c.a[cmd.name]
	if !ok {
		return errors.New("no command given")
	}
	err := v(s, cmd)
	return err
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.a[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("login must not be empty")
	}
	name := cmd.args[0]
	_, err2 := s.db.GetUser(context.Background(), name)
	if err2 != nil {
		return errors.New("user must not be empty")
	}
	err := s.sta.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Println("username has been set")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("register must not be empty")
	}
	name := cmd.args[0]
	user := database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: name,
	}
	a, err := s.db.CreateUser(context.Background(), user)
	if err != nil {
		return err
	}
	err2 := s.sta.SetUser(cmd.args[0])
	if err2 != nil {
		return err2
	}
	fmt.Println("user created")
	fmt.Println(a.Name)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.ResetUser(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func handlerGetUsers(s *state, cmd command) error {
	use, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range use {
		if user.Name == s.sta.Current_user_name {
			fmt.Printf("* %v (current)\n", user.Name)
		} else {
			fmt.Printf("* %v\n", user.Name)
		}
	}
	return nil
}

func handlerfeeds(s *state, cmd command) error {
feed, err := s.db.GetFeeds(context.Background())
if err != nil {
	return err
}
for _, f := range feed {
fmt.Printf("%v\n",f.Name)
fmt.Printf("%v\n",f.Url)
c, err3 := s.db.Getwhat(context.Background(), f.UserID)
if err3 != nil {
	return err3
}
a, err2 := s.db.GetUser(context.Background(), c)
if err2 != nil {
	return err2
}
fmt.Printf("%v\n",a.Name)
}
return nil
}

func addfeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("usage: %v", cmd.name)
	}
	feed := database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: cmd.args[0],
		Url: cmd.args[1],
		UserID: user.ID,
		}
	
	feeded, err := s.db.CreateFeed(context.Background(), feed)
	if err != nil {
		return err
	}
	feedf := database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: feeded.ID,
	}
	b, err4 := s.db.CreateFeedFollow(context.Background(), feedf)
	if err4 != nil {
		return err4
	}
	fmt.Println(b.FeedNames)
	fmt.Println(feeded)
	return nil
}
func handlerfollow(s *state, cmd command, user database.User) error {
 url := cmd.args[0]
 feed, err3 := s.db.GetUrl(context.Background(), url)
 if err3 != nil {
	return err3
 }
 parm := database.CreateFeedFollowParams{
	ID: uuid.New(),
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
	UserID: user.ID,
	FeedID: feed.ID,
}
a, err := s.db.CreateFeedFollow(context.Background(), parm)
if err != nil {
	return err
}
fmt.Printf("%v\n", a.FeedNames)
fmt.Printf("%v\n", a.UserNames)
return nil
}
func handlerfollowing(s *state, cmd command, user database.User) error {
	a, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}
	for _, f := range a {
		fmt.Printf("%v\n", f.FeedName)
	}
	return nil
}
func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
	current := s.sta.Current_user_name
	user, err := s.db.GetUser(context.Background(), current)
	if err != nil {
		return err
	}
	return handler(s , cmd, user)
}	
}
func unfollow(s *state, cmd command, user database.User) error {
	url := cmd.args[0]
	feed, err2 := s.db.GetUrl(context.Background(), url)
	if err2 != nil {
		return err2
	}
	parm := database.DeleteFeedFollowParams{
	UserID: user.ID,
	FeedID: feed.ID,
	}
	err := s.db.DeleteFeedFollow(context.Background(), parm)
	if err != nil {
		return err
	}
	return nil
}
func scrapeFeeds(s *state) error {
feed, err := s.db.GetNextFeedToFetch(context.Background())
if err != nil {
	return err
}
_, err2 := s.db.MarkFeedFetched(context.Background(), feed.ID)
if err2 != nil {
	return err2
}
rssfeed, err := fetchFeed(context.Background(), feed.Url)
if err != nil {
	return err
}
for _, item := range rssfeed.Channel.Item {
	publishedAt := sql.NullTime{}
parsedTime, err := time.Parse(time.RFC1123Z, item.PubDate)
if err == nil {
    publishedAt = sql.NullTime{
        Time:  parsedTime,
        Valid: true,
    }
}
parm2 := database.CreatePostParams{
	ID: uuid.New(),
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
	Title: item.Title,
	Url: item.Link,
	Description: item.Description,
	PublishedAt: publishedAt,
	FeedID: feed.ID,
}
 _, err7 := s.db.CreatePost(context.Background(), parm2)
 if err7 != nil {
	return err7
 }

}
return nil
}
func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) < 1 {
	 return fmt.Errorf("has to be more than 0 arguments")
	}
	dur, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Collecting feeds every %v", dur)
	ticker := time.NewTicker(dur)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
	return nil
}
