package notification

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/bshafiee/bonga/scraping"
)

type Channel interface {
	Notify([]scraping.Result) error
}

type Engine struct {
	seenIDs     map[string]bool //list of IDs we have had before
	channels    []Channel
	seenIDsFile string
}

func NewNotificationEngine(dbPath string, c []Channel) *Engine {
	return &Engine{
		seenIDs:     make(map[string]bool),
		channels:    c,
		seenIDsFile: dbPath,
	}
}

func (e *Engine) Initialize() error {
	inFile, err := os.Open(e.seenIDsFile)
	if err != nil {
		return err
	}
	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		if id := strings.TrimSpace(scanner.Text()); len(id) > 0 {
			e.seenIDs[id] = true
		}
	}
	fmt.Println("loaded ", len(e.seenIDs), " existing IDs")
	return nil
}

func (e *Engine) Notify(results []scraping.Result) error {
	//1) find new ones
	newRes := make([]scraping.Result, 0)
	for _, res := range results {
		if !e.seenIDs[res.ID] {
			newRes = append(newRes, res)
			e.seenIDs[res.ID] = true
		}
	}
	//2) notify
	for _, ch := range e.channels {
		if err := ch.Notify(newRes); err != nil {
			return err
		}
	}
	//3) update seen list
	return e.updateSeenIDs()
}

func (e *Engine) updateSeenIDs() error {
	file, err := os.Create(e.seenIDsFile)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for id := range e.seenIDs {
		fmt.Fprintln(w, id)
	}
	return w.Flush()
}
