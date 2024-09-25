package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

const (
	monthWidth  = 23
	monthMargin = 2
)

var (
	monthStyle = lipgloss.NewStyle().
			Width(monthWidth).
			MarginRight(monthMargin)

	monthHeadingStyle = lipgloss.NewStyle().Width(monthWidth)

	highlight = lipgloss.NewStyle().
			Background(lipgloss.Color("7")).
			Foreground(lipgloss.Color("18"))

	today = time.Now()
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if len(os.Args) == 1 {
		now := time.Now()
		fmt.Println(renderMonth(now.Year(), int(now.Month()), fmt.Sprintf("%s %d", now.Month(), now.Year())))
		return nil
	}

	year, err := strconv.Atoi(os.Args[1])
	if err != nil {
		return err
	}

	termWidth, _, _ := term.GetSize(int(os.Stdin.Fd()))
	columns := termWidth / (monthWidth + monthMargin)
	if columns > 6 {
		columns = 6
	}

	fmt.Println(renderYear(year, columns))

	return nil
}

func renderYear(year int, columns int) string {
	headingStyle := lipgloss.NewStyle().Width((monthWidth + monthMargin - 1) * columns).Align(lipgloss.Center)
	var res string
	res += headingStyle.Render(strconv.Itoa(year))
	for mn := 1; mn <= 12; {
		var row []string
		for col := 0; col < columns; col++ {
			if mn > 12 {
				break
			}
			month := time.Month(mn)
			row = append(row, renderMonth(year, mn, month.String()))
			mn++
		}
		res = res + "\n" + lipgloss.JoinHorizontal(lipgloss.Top, row...) + "\n"
	}
	return strings.TrimRight(res, " \n")
}

func renderMonth(year, month int, heading string) string {
	const weekHeading = "Wk Mo Tu We Th Fr Sa Su\n"

	res := strings.Builder{}
	res.WriteString(monthHeadingStyle.Align(lipgloss.Center).Render(heading) + "\n")

	date := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	res.WriteString(weekHeading)

	if date.Weekday() != time.Monday {
		_, wk := date.ISOWeek()
		res.WriteString(fmt.Sprintf("%2d", wk))
	}

	offset := int(date.Weekday())
	if date.Weekday() == time.Sunday {
		offset = 7
	}
	for i := 1; i < offset; i++ {
		res.WriteString("   ")
	}

	for {
		if int(date.Month()) != month {
			break
		}
		num := strconv.Itoa(date.Day())
		if isToday(date) {
			num = highlight.Render(strconv.Itoa(date.Day()))
		}
		switch date.Weekday() {
		case 0:
			res.WriteString(fmt.Sprintf(" %2s\n", num))
		case 1:
			year, wk := date.ISOWeek()
			if year != date.Year() {
				wk = 53
			}
			res.WriteString(fmt.Sprintf("%2d %2s", wk, num))
		default:
			res.WriteString(fmt.Sprintf(" %2s", num))
		}
		date = date.Add(time.Hour * 24)
	}
	return monthStyle.Render(res.String())
}

func isToday(date time.Time) bool {
	return date.Year() == today.Year() && date.Month() == today.Month() && date.Day() == today.Day()
}
