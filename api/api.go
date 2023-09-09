package api

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/adamwoolhether/tictactoe/game"
	"github.com/adamwoolhether/tictactoe/polling"
)

//go:embed home.html
var content embed.FS

var count int

type HomeParam struct {
	Move
}

// GamePlay represents the active game being played.
type GamePlay struct {
	Board      [game.BoardSize][game.BoardSize]*game.Player
	PlayerTurn game.Player
	Winner     *game.Player
	Counter    int
}

type Move struct {
	GamePlay
	CreatedAt int64
}

type Handler struct {
	*chi.Mux
	game   game.Game
	pool   *Pool
	queue  *pubsub.Queue[Move]
	pubsub *pubsub.PubSub
}

func New(game game.Game) *Handler {
	h := &Handler{
		Mux:    chi.NewMux(),
		game:   game,
		pool:   NewPool(2),
		queue:  pubsub.NewQueue[Move](9),
		pubsub: pubsub.NewPubSub(),
	}

	h.Use(middleware.Logger)
	h.Use(middleware.Timeout(30 * time.Second))

	h.Get("/", h.Home)
	h.Route("/web", func(r chi.Router) {
		r.Get("/updates", h.Updates)
		r.Get("/player/{playerID}/move", h.Move)
		r.Delete("/reset", h.Reset)
	})

	return h
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	var move Move

	canPlay, amtPlayers := h.pool.CanPlay()
	if canPlay {
		gamePlay := h.createGamePlay(amtPlayers)
		move = h.enqueueAndPublish(gamePlay)
	} else {
		gamePlay := h.createGamePlay(-1)
		move = Move{
			GamePlay:  gamePlay,
			CreatedAt: time.Now().Unix(),
		}
	}

	renderHomePage(w, HomeParam{move})
}

func (h *Handler) Updates(w http.ResponseWriter, r *http.Request) {
	getMovesAfter := func(ts int64) []Move {
		moves := h.queue.Copy()
		index := sort.Search(len(moves), func(i int) bool {
			return moves[i].CreatedAt > ts
		})

		return moves[index:]
	}

	lastUpdate, _ := strconv.ParseInt(r.URL.Query().Get("lastUpdate"), 10, 64)
	moves := getMovesAfter(lastUpdate)

	if len(moves) == 0 {
		ch, closefn := h.pubsub.Subscribe()
		defer closefn()

		select {
		case <-ch:
			moves = getMovesAfter(lastUpdate)
		case <-r.Context().Done():
			w.WriteHeader(http.StatusRequestTimeout)
			w.Write([]byte("requet timed out"))

			return
		}
	}

	// Update client if any moves
	renderHomePage(w, HomeParam{moves[0]})
}

func (h *Handler) Reset(w http.ResponseWriter, r *http.Request) {
	h.game.Start()
	h.pool = NewPool(2)
	count = 0
	gamePlay := h.createGamePlay(count)
	move := h.enqueueAndPublish(gamePlay)

	renderHomePage(w, HomeParam{move})
}

func (h *Handler) Move(w http.ResponseWriter, r *http.Request) {
	player := chi.URLParam(r, "playerID")
	row, _ := strconv.Atoi(r.URL.Query().Get("row"))
	col, _ := strconv.Atoi(r.URL.Query().Get("col"))
	count++

	if err := h.game.Move(game.Player(player), row, col); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	gamePlay := h.createGamePlay(count)
	move := h.enqueueAndPublish(gamePlay)

	renderHomePage(w, HomeParam{move})
}

func (h *Handler) createGamePlay(c int) GamePlay {
	g := GamePlay{
		Board:      h.game.GetBoard(),
		PlayerTurn: h.game.GetTurn(),
		Winner:     h.game.GetWinner(),
		Counter:    c,
	}

	return g
}

func (h *Handler) enqueueAndPublish(gameplay GamePlay) Move {
	move := Move{GamePlay: gameplay, CreatedAt: time.Now().Unix()}
	h.queue.Enqueue(move)
	h.pubsub.Publish()

	return move
}

func renderHomePage(w io.Writer, p HomeParam) {
	home := template.Must(template.ParseFS(content, "home.html"))

	if err := home.Execute(w, p); err != nil {
		fmt.Println("err rendering page: %w", err)
	}
}
