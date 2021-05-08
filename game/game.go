package game

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"sync"
)

var initialSeq []int

func init() {
	initialSeq = make([]int, 101)
	for i := 0; i < 101; i++ {
		initialSeq[i] = i
	}
}

type Session struct {
	playersMap map[string]*Player
	candidates []int
	mutex sync.Mutex
}

type Player struct {
	user *discordgo.User
	number int
}

func NewSession() *Session {
	candidates := make([]int, 101)
	copy(candidates, initialSeq)
	return &Session{
		playersMap: make(map[string]*Player),
		candidates: candidates,
		mutex: sync.Mutex{},
	}
}

func (s *Session) JoinUser(user *discordgo.User) int {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	player, ok := s.playersMap[user.ID]
	if ok {
		return player.number
	}
	i := rand.Intn(len(s.candidates))
	num := s.candidates[i]
	s.candidates = append(s.candidates[:i], s.candidates[i + 1:]...)
	player = &Player{
		user:   user,
		number: num,
	}
	s.playersMap[user.ID] = player
	return num
}

func (s *Session) GetPlayerNumber(discordUserId string) (int, error) {
	player, ok := s.playersMap[discordUserId]
	if ok {
		return player.number, nil
	}
	return 0, errors.New("player not found")
}
