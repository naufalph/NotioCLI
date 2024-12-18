package ui

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	m "tdlst/models"
	"tdlst/pkg/utils"

	"github.com/olekukonko/tablewriter"
)

func PrintListNew(listType string, tasks []m.Task) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Description", "Created", "Due", "Status"})
	table.SetAutoWrapText(false)

	for _, task := range tasks {
		table.Append(addRowBasedOnStatus(ConvertToTaskVO(task)))
	}
	table.Render()
}

func PrintLine(line string) {
	fmt.Println(line)
}

func boldify(text string) string {
	if text == "" {
		return ""
	} else {
		return fmt.Sprintf("\033[1m%v\033[0m", text)
	}
}

func strikeThrough(text string) string {
	if text == "" {
		return ""
	} else {
		return fmt.Sprintf("\033[9m%v\033[0m", text)
	}
}

type TaskVO struct {
	ID          string
	Description string
	CreatedAt   string
	DueDate     string
	Status      string
}

func ConvertToTaskVO(task m.Task) TaskVO {
	return TaskVO{
		ID:          strconv.FormatUint(uint64(task.ID), 6),
		Description: task.Description,
		CreatedAt:   task.CreatedAt.Format("02 Jan 2006"),
		DueDate:     task.DueDate.Format("02 Jan 2006"),
		Status:      task.Status,
	}
}
func addRowBasedOnStatus(taskVO TaskVO) []string {
	resultArray := []string{
		boldify(taskVO.ID),
		taskVO.Description,
		taskVO.CreatedAt,
		taskVO.DueDate,
		taskVO.Status,
	}
	if strings.Compare(utils.StatusDone, taskVO.Status) == 0 {
		for i := 1; i < len(resultArray); i++ {
			resultArray[i] = strikeThrough(resultArray[i])
		}
	}
	return resultArray
}
