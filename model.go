package main

import (
	"encoding/binary"
	"fmt"
	"os"

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
	if t, err := isZstd(m.zstdFilepath); t == false {
		return errMsg{err}
	}
	return nil
}

func (m model) Init() tea.Cmd {
	return m.zstdCheck
}

func (m model) DecompressFile() tea.Msg {
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
	if len(m.data) != SAVE_FILE_SIZE {
		return errMsg{fmt.Errorf("The savefile's filesize (%d) doesn't match the expected file size (%d)", len(m.data), SAVE_FILE_SIZE)}
	}
	return m.UpdateDifficulty()
}

func (m model) UpdateDifficulty() tea.Msg {
	difficultyByte := byte(m.selected)
	m.data[DIFFICULTY_OFFSET] = difficultyByte
	return m.updateChecksum()
}

func (m model) updateChecksum() tea.Msg {
	checksum := calculateTrailsFromZeroChecksum(m.data[:len(m.data)-CHECKSUM_SIZE])
	copy(m.data[len(m.data)-CHECKSUM_SIZE:], checksum)
	return m.compressToFile()
}

func (m model) compressToFile() tea.Msg {
	newSaveData := Compress(m.data)
	os.WriteFile(SAVE_FILE_NAME, newSaveData, 0644)
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
			return m, m.DecompressFile
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

	s := "What difficulty should your save be?\n\n"

	for i, difficulty := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %s\n", cursor, difficulty.String())
	}

	s += "\nPress q to quit.\n"
	return s

}

func calculateTrailsFromZeroChecksum(data []byte) []byte {

	var sum32Bit uint32 = 0
	var notSum32Bit uint32 = 0

	for i := 0; i < len(data); i += 4 {
		num := binary.LittleEndian.Uint32(data[i : i+4])
		sum32Bit += num
		notSum32Bit += ^num
	}

	checksum := make([]byte, 0, CHECKSUM_SIZE)
	checksum = binary.LittleEndian.AppendUint32(checksum, sum32Bit)
	checksum = binary.LittleEndian.AppendUint32(checksum, notSum32Bit)

	return checksum
}
