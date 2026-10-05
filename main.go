package main
import _ "github.com/lib/pq"
import (
	"fmt"
	"os"
	"github.com/pitkve1-wq/internal/database"
	"github.com/pitkve1-wq/internal/config"
	"database/sql"
)

func main() {
	c, err := config.Read()
	if err != nil {
		os.Exit(1)
	}
	a, err3 := sql.Open("postgres", c.Db_url)
	if err3 != nil {
		fmt.Println(err3)
		os.Exit(1)
	}
	d := database.New(a)
	s := &state{
		sta: &c,
		db: d,
	}
	v := map[string]func(*state, command) error{}
	cmds := commands{
		a: v,
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerGetUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(addfeed))
	cmds.register("feeds", handlerfeeds)
	cmds.register("follow", middlewareLoggedIn(handlerfollow))
	cmds.register("following", middlewareLoggedIn(handlerfollowing))
	cmds.register("unfollow", middlewareLoggedIn(unfollow))
	cmds.register("agg", handlerAgg)
	cmds.register("browse", scrapeFeeds)
	if len(os.Args) < 2 {
		os.Exit(1)
	}
	x := command{
		name: os.Args[1],
		args: os.Args[2:],
	}
	err2 := cmds.run(s, x)
	if err2 != nil {
		fmt.Println(err2)
		os.Exit(1)
	}

}
