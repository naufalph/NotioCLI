package notion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"tdlst/config"
	"tdlst/db"
	"tdlst/internal/repository"
	m "tdlst/models"
	"tdlst/pkg/applog"
	"tdlst/pkg/utils"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type NotionResponse struct {
	Object          string      `json:"object"`
	Results         []Page      `json:"results"`
	NextCursor      *string     `json:"next_cursor"`
	HasMore         bool        `json:"has_more"`
	Type            string      `json:"type"`
	PageOrDatabase  interface{} `json:"page_or_database"`
	DeveloperSurvey string      `json:"developer_survey"`
	RequestID       string      `json:"request_id"`
}

type Page struct {
	Object         string      `json:"object"`
	ID             string      `json:"id"`
	CreatedTime    string      `json:"created_time"`
	LastEditedTime string      `json:"last_edited_time"`
	CreatedBy      User        `json:"created_by"`
	LastEditedBy   User        `json:"last_edited_by"`
	Cover          interface{} `json:"cover"` // Can be further specified if necessary
	Icon           interface{} `json:"icon"`  // Can be further specified if necessary
	Parent         Parent      `json:"parent"`
	Archived       bool        `json:"archived"`
	InTrash        bool        `json:"in_trash"`
	Properties     Properties  `json:"properties"`
	URL            string      `json:"url"`
	PublicURL      *string     `json:"public_url"` // Nullable field
}

type Property struct {
	ID     string      `json:"id"`
	Type   string      `json:"type"`
	Status *Status     `json:"status,omitempty"`
	Date   *Date       `json:"date,omitempty"`
	Title  []TextBlock `json:"title,omitempty"`
}

type User struct {
	Object string `json:"object"`
	ID     string `json:"id"`
}

type Parent struct {
	DatabaseID string `json:"database_id"`
}

type Properties struct {
	Status      Status      `json:"Status"`
	DueDate     DueDate     `json:"Due Date"`
	ID          UniqueID    `json:"ID"`
	Description Description `json:"Description"`
}

type Status struct {
	ID     string        `json:"id"`
	Type   string        `json:"type"`
	Status *StatusDetail `json:"status"` // Adjusted to use StatusDetail
}

type StatusDetail struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type DueDate struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Date *Date  `json:"date"` // Nullable field
}

type Date struct {
	Start    *string `json:"start"`     // Nullable
	End      *string `json:"end"`       // Nullable
	TimeZone *string `json:"time_zone"` // Nullable
}

type UniqueID struct {
	ID       string          `json:"id"`
	Type     string          `json:"type"`
	UniqueID UniqueIDDetails `json:"unique_id"`
}

type UniqueIDDetails struct {
	Prefix *string `json:"prefix"` // Nullable
	Number int     `json:"number"`
}

type Description struct {
	ID    string  `json:"id"`
	Type  string  `json:"type"`
	Title []Title `json:"title"`
}

type Title struct {
	Type        string      `json:"type"`
	Text        Text        `json:"text"`
	Annotations Annotations `json:"annotations"`
	PlainText   string      `json:"plain_text"`
	Href        *string     `json:"href"` // Nullable field
}

type Text struct {
	Content string  `json:"content"`
	Link    *string `json:"link"` // Nullable field
}

type Annotations struct {
	Bold          bool   `json:"bold"`
	Italic        bool   `json:"italic"`
	Strikethrough bool   `json:"strikethrough"`
	Underline     bool   `json:"underline"`
	Code          bool   `json:"code"`
	Color         string `json:"color"`
}

type TextBlock struct {
	Type        string      `json:"type"`
	Text        TextContent `json:"text"`
	Annotations Annotations `json:"annotations"`
	PlainText   string      `json:"plain_text"`
	Href        *string     `json:"href,omitempty"`
}

type TextContent struct {
	Content string      `json:"content"`
	Link    interface{} `json:"link"`
}

var DB *gorm.DB

func SyncTask(dbTest *gorm.DB) error {
	dbUse := dbTest
	if dbTest == nil {
		dbUse = db.DB
	}

	taskFromNotion, err := GetTaskList()
	if err != nil {
		applog.Error(err, "Error when query to Notion")
		return err
	}
	taskFromDB, err := repository.ReadTaskAll(dbTest)
	if err != nil {
		applog.Error(err, "Error when query to DB")
		return err
	}
	missingFromNotion, missingFromDB, updateFromNotion, updateFromDB := synchronizedSlice(taskFromNotion, taskFromDB)

	if len(missingFromNotion) > 0 {
		for _, task := range missingFromNotion {
			err = AddTask(task)
			if err != nil {
				return err
			}
		}
	}

	if len(missingFromDB) > 0 {
		for _, task := range missingFromDB {
			err = repository.WriteTask(dbUse, task)
			if err != nil {
				return err
			}
		}
	}

	if len(updateFromNotion) > 0 {
		for _, task := range updateFromNotion {
			err = UpdateTask(task)
			if err != nil {
				return err
			}
		}
	}

	if len(updateFromDB) > 0 {
		for _, task := range updateFromDB {
			err = repository.UpdateTask(dbUse, &task, &task)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func synchronizedSlice(taskFromNotion, taskFromDB []m.Task) (missingFromNotion, missingFromDB, updateFromNotion, updateFromDB []m.Task) {
	notionMap := make(map[uint16]m.Task)
	dbMap := make(map[uint16]m.Task)

	for _, task := range taskFromNotion {
		notionMap[task.ID] = task
	}

	for _, task := range taskFromDB {
		dbMap[task.ID] = task
	}

	// Check for missing and updates
	for id, dbTask := range dbMap {
		notionTask, existsInNotion := notionMap[id]

		if !existsInNotion {
			// Task exists in DB but not in Notion
			missingFromNotion = append(missingFromNotion, dbTask)
		} else {
			// Task exists in both, check for updates based on the latest modification
			if tasksNeedUpdate(dbTask, notionTask) {
				if dbTask.UpdatedAt.After(notionTask.UpdatedAt) {
					updateFromNotion = append(updateFromNotion, notionTask)
				} else {
					updateFromDB = append(updateFromDB, dbTask)
				}
			}
		}
	}

	// Check for tasks in Notion but missing in DB
	for id, notionTask := range notionMap {
		_, existsInDB := dbMap[id]
		if !existsInDB {
			missingFromDB = append(missingFromDB, notionTask)
		}
	}

	return missingFromNotion, missingFromDB, updateFromNotion, updateFromDB
}

func tasksNeedUpdate(task1, task2 m.Task) bool {
	return string(task1.Status) != string(task2.Status) || !task1.DueDate.Equal(task2.DueDate) || strings.Compare(task1.NotionId, task2.NotionId) != 0
}

func AddTask(task m.Task) error {
	type ParentAdd struct {
		DatabaseID string `json:"database_id"`
	}
	type StatusDetailAdd struct {
		Name string `json:"name"`
	}
	type StatusAdd struct {
		Status *StatusDetailAdd `json:"status"` // Adjusted to use StatusDetail
	}
	type PropertyAdd struct {
		Status      StatusAdd   `json:"Status"`
		DueDate     DueDate     `json:"Due Date"`
		Description Description `json:"Description"`
	}
	type NewTaskRequest struct {
		Parent     ParentAdd   `json:"parent"`
		Properties PropertyAdd `json:"properties"`
	}

	notionDbId, notionApiKey := getDbIdAndApi()

	dueDateString := task.DueDate.Format(time.DateOnly)
	payload := NewTaskRequest{
		Parent: ParentAdd{
			DatabaseID: notionDbId,
		},
		Properties: PropertyAdd{
			DueDate: DueDate{
				ID:   "BxQ%7D",
				Type: "date",
				Date: &Date{
					Start:    &dueDateString,
					End:      nil,
					TimeZone: nil,
				},
			},
			Description: Description{
				ID:   "title",
				Type: "title",
				Title: []Title{
					{
						Type: "text",
						Text: Text{
							Content: task.Description,
							Link:    nil,
						},
						Annotations: Annotations{
							Bold:          false,
							Italic:        false,
							Strikethrough: false,
							Underline:     false,
							Code:          false,
							Color:         "default",
						},
						PlainText: task.Description,
						Href:      nil,
					},
				},
			},
			Status: StatusAdd{
				Status: &StatusDetailAdd{
					Name: string(task.Status),
				},
			},
		},
	}

	// Serialize the payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to serialize payload: %w", err)
	}

	// Create the HTTP request
	url := "https://api.notion.com/v1/pages/"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Add headers
	req.Header.Set("Authorization", "Bearer "+notionApiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Notion-Version", "2022-06-28")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create task, status: %d, body: %s", resp.StatusCode, string(body))
	}

	fmt.Println("Task created successfully:", string(body))
	return nil
}

func GetTaskList() ([]m.Task, error) {
	notionDbId, notionApiKey := getDbIdAndApi()

	url := fmt.Sprintf("https://api.notion.com/v1/databases/%s/query", notionDbId)
	payload := map[string]interface{}{
		"sorts": []map[string]interface{}{
			{
				"property":  "ID",
				"direction": "ascending",
			},
		},
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %v", notionApiKey)).
		SetHeader("Notion-Version", "2022-06-28").
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post(url)
	if err != nil {
		applog.Error(err, utils.NtnReqError)
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		applog.Error(err, utils.NtnReqError)
		return nil, err
	}

	var notionResponse NotionResponse
	reader := bytes.NewReader(resp.Body())
	err = json.NewDecoder(reader).Decode(&notionResponse)
	if err != nil {
		applog.Error(err, utils.ParseJsonError)
		return nil, err
	}
	var tasks []m.Task
	for i := 0; i < len(notionResponse.Results); i++ {
		prop := notionResponse.Results[i].Properties
		var task m.Task
		task.ID = uint16(prop.ID.UniqueID.Number)

		createdTime, err := time.Parse("2006-01-02T15:04:05.000Z", notionResponse.Results[i].CreatedTime)
		if err == nil {
			task.CreatedAt = createdTime
		}
		if prop.DueDate.Date != nil {
			dueDate, err := time.Parse(time.DateOnly, *prop.DueDate.Date.Start)
			if err == nil {
				task.DueDate = dueDate
			}
		}
		updatedTime, err := time.Parse("2006-01-02T15:04:05.000Z", notionResponse.Results[i].LastEditedTime)
		if err == nil {
			task.UpdatedAt = updatedTime
		}
		task.Status = utils.ParseStatusString(prop.Status.Status.Name)
		task.Description = prop.Description.Title[0].Text.Content
		task.NotionId = notionResponse.Results[i].ID
		tasks = append(tasks, task)
	}

	return tasks, err
}

func UpdateTask(newTask m.Task) error {
	type StatusDetailUpdate struct {
		Name string `json:"name"`
	}
	type StatusUpdate struct {
		Status *StatusDetailUpdate `json:"status"`
	}
	type DueDateDetailUpdate struct {
		Start string `json:"start"`
	}
	type DueDateUpdate struct {
		Date *DueDateDetailUpdate `json:"date"`
	}
	type TitleText struct {
		Content string  `json:"content"`
		Link    *string `json:"link"` // Nullable field
	}
	type TitleUpdate struct {
		Type      string    `json:"type"`
		Text      TitleText `json:"text"`
		PlainText string    `json:"plain_text"`
		Href      *string   `json:"href"` // Nullable field
	}
	type DescriptionUpdate struct {
		Title []TitleUpdate `json:"title"`
	}
	type PropertiesUpdate struct {
		Status      StatusUpdate      `json:"Status"`
		DueDate     DueDateUpdate     `json:"Due Date"`
		Description DescriptionUpdate `json:"Description"`
	}
	type UpdateTaskRequest struct {
		Properties PropertiesUpdate `json:"properties"`
	}
	// Payload Construction
	dueDateString := newTask.DueDate.Format(time.DateOnly)
	payload := UpdateTaskRequest{
		Properties: PropertiesUpdate{
			DueDate: DueDateUpdate{
				Date: &DueDateDetailUpdate{
					Start: dueDateString,
				},
			},
			Description: DescriptionUpdate{
				Title: []TitleUpdate{
					{
						Type: "text",
						Text: TitleText{
							Content: newTask.Description,
							Link:    nil,
						},
						PlainText: newTask.Description,
					},
				},
			},
			Status: StatusUpdate{
				Status: &StatusDetailUpdate{
					Name: string(newTask.Status),
				},
			},
		},
	}

	// Serialize the payload to JSON
	payloadBytes, err := json.Marshal(payload)
	applog.Debug(string(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to serialize payload: %w", err)
	}

	// Create the HTTP request
	url := "https://api.notion.com/v1/pages/" + newTask.NotionId
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	_, notionApiKey := getDbIdAndApi()
	// Add headers
	req.Header.Set("Authorization", "Bearer "+notionApiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Notion-Version", "2022-06-28")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create task, status: %d, body: %s", resp.StatusCode, string(body))
	}

	fmt.Println("Task updated successfully:", string(body))
	return nil
}
func getDbIdAndApi() (string, string) {
	envMap, err := godotenv.Read(config.RealEnv())
	if err != nil {
		applog.Error(err, utils.NtnConnectionError)
		os.Exit(1)
	}

	return envMap["NOTION_DB_ID"], envMap["NOTION_API_KEY"]
}
