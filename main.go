package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/adamwoolhether/tictactoe/api"
	"github.com/adamwoolhether/tictactoe/game"
)

func main() {
	g := game.New()

	a := api.New(g)

	fmt.Println("Game server started:::")
	if err := http.ListenAndServe(":3000", a); err != nil {
		fmt.Fprintf(os.Stderr, "server err: %v", err)
	}
}
