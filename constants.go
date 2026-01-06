package main

import (
	"log"
)

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type statusMsg int

const (
	SUCCESS statusMsg = 0
)

type errMsg struct{ err error }

// For messages that contain errors it's often handy to also implement the
// error interface on the message.
func (e errMsg) Error() string { return e.err.Error() }

const (
	ChecksumSize     = 8
	SaveFileSize     = 155624
	DifficultyOffset = 0x0001EF91
)

type Difficulty byte

const (
	Normal Difficulty = iota
	Hard
	Nightmare
	Easy
	None
)

var difficultyName = map[Difficulty]string{
	Normal:    "Normal",
	Hard:      "Hard",
	Nightmare: "Nightmare",
	Easy:      "Easy",
	None:      "None",
}

func (d Difficulty) String() string {
	return difficultyName[d]
}
