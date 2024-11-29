package tui

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type message struct {
	str string
}

type model struct {
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	err         error
	conn        net.Conn
}

func Read(conn net.Conn, msg func(string)) {
	reader := bufio.NewReader(conn)
	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Server closed its connection")
			os.Exit(0)
		}
		str = strings.TrimSpace(str)
		msg(str)
	}
}

func StartTui() {
	conn, err := net.Dial("tcp", "localhost:3000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	m := initialModel(conn)
	p := *tea.NewProgram(m)

	go Read(conn, func(s string) {
		receive := message{
			str: s,
		}
		p.Send(receive)
	})

	if _, err = p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Oof: %v\n", err)
	}
}

func initialModel(conn net.Conn) model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "| "
	ta.CharLimit = 280
	ta.SetWidth(30)
	ta.SetHeight(1)

	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(30, 5)
	vp.SetContent(`Welcome to the chat room!
Type a message and press Enter to send.`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	m := model{
		conn:        conn,
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
	}
	return m
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m *model) WriteMessage(msg string) {
	writer := bufio.NewWriter(m.conn)
	_, err := writer.WriteString(msg + "\n")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	writer.Flush()

	m.messages = append(m.messages, m.senderStyle.Render("You: ")+msg)
	m.viewport.SetContent(strings.Join(m.messages, "\n"))
	m.textarea.Reset()
	m.viewport.GotoBottom()
}
