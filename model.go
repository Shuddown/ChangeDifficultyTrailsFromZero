package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	zstdFilepath string
	data         []byte
	choices      []Difficulty
	cursor       int
	selected     Difficulty
	success      bool
	err          error
}

func initialModel(zstdFilepath string) model {
	return model{
		zstdFilepath: zstdFilepath,
		choices:      []Difficulty{Easy, Normal, Hard, Nightmare},
		cursor:       0,
		selected:     None,
		success:      false,
		err:          nil,
	}
}

func (m model) zstdCheck() tea.Msg {
	if t, err := isZstd(m.zstdFilepath); !t {
		return errMsg{err}
	}
	return nil
}

func (m model) Init() tea.Cmd {
	return m.zstdCheck
}

func (m model) decompressFile() tea.Msg {
	data, err := os.ReadFile(m.zstdFilepath)
	if err != nil {
		return errMsg{err}
	}
	decompressedData, err := Decompress(data)
	if err != nil {
		return errMsg{err}
	}
	m.data = decompressedData
	return m.validateSaveFile()
}

func (m model) validateSaveFile() tea.Msg {
	if len(m.data) != SaveFileSize {
		return errMsg{fmt.Errorf("the savefile's filesize (%d) doesn't match the expected file size (%d)", len(m.data), SaveFileSize)}
	}
	return m.updateDifficulty()
}

func (m model) updateDifficulty() tea.Msg {
	difficultyByte := byte(m.selected)
	m.data[DifficultyOffset] = difficultyByte
	return m.updateChecksum()
}

func (m model) updateChecksum() tea.Msg {
	checksum := calculateTrailsFromZeroChecksum(m.data[:len(m.data)-ChecksumSize])
	copy(m.data[len(m.data)-ChecksumSize:], checksum)
	return m.compressToFile()
}

func writeFileSafe(name string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.Write(data); err != nil {
		return err
	}

	return f.Sync()
}

func (m model) compressToFile() tea.Msg {
	newSaveData := Compress(m.data)
	err := writeFileSafe(m.zstdFilepath, newSaveData, 0o644)
	if err != nil {
		return errMsg{err}
	}
	return SUCCESS
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			m.cursor = max(m.cursor-1, 0)

		case "down", "j":
			m.cursor = min(m.cursor+1, len(m.choices)-1)

		case "enter", " ":
			m.selected = m.choices[m.cursor]
			return m, m.decompressFile
		}

	case errMsg:
		m.err = msg
		return m, tea.Quit

	case statusMsg:
		m.success = true
		return m, tea.Quit
	}
	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nError: %v\n\n", m.err)
	}

	if m.success {
		return fmt.Sprintf("\nSuccessfully changed difficulty to %s.\n\n", m.selected.String())
	}

	if m.selected != None {
		return fmt.Sprintf("\nYou selected %s. Updating save file...\n\n", m.selected.String())
	}

	var s strings.Builder
	s.WriteString("What difficulty should your save be?\n\n")

	for i, difficulty := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		fmt.Fprintf(&s, "%s %s\n", cursor, difficulty.String())
	}

	s.WriteString("\nPress q to quit.\n")
	return s.String()
}

func calculateTrailsFromZeroChecksum(data []byte) []byte {
	var sum32Bit uint32 = 0
	var notSum32Bit uint32 = 0

	for i := 0; i < len(data); i += 4 {
		num := binary.LittleEndian.Uint32(data[i : i+4])
		sum32Bit += num
		notSum32Bit += ^num
	}

	checksum := make([]byte, 0, ChecksumSize)
	checksum = binary.LittleEndian.AppendUint32(checksum, sum32Bit)
	checksum = binary.LittleEndian.AppendUint32(checksum, notSum32Bit)

	return checksum
}
