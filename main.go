package tickspot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type TickClient struct {
	tickURL   string
	tickToken string
	userAgent string
	Client    *http.Client
}

func NewTickspotClient(tickProject, tickToken, userAgent string) (*TickClient, error) {
	if tickProject == "" {
		return nil, errors.New("TICK_URL environment variable is not set")
	}
	if tickToken == "" {
		return nil, errors.New("TICK_TOKEN environment variable is not set")
	}
	if userAgent == "" {
		return nil, errors.New("User-Agent is not set")
	}

	tickURL := fmt.Sprintf("https://tickspot.com/%s/api/v2", tickProject)

	client := &http.Client{}

	return &TickClient{
		tickURL:   tickURL,
		tickToken: tickToken,
		userAgent: userAgent,
		Client:    client,
	}, nil
}

func (T *TickClient) sendRequest(method, url string, body []byte) ([]byte, error) {
	headerToken := fmt.Sprintf(`Token token=%s`, T.tickToken)
	r, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	r.Header.Add("Authorization", headerToken)
	r.Header.Add("User-Agent", T.userAgent)
	r.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode >= 300 {
		bodyContent, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(string(bodyContent))
	}

	bodyContent, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 300 {
		return nil, errors.New(string(bodyContent))
	}
	return bodyContent, nil
}

func (T *TickClient) GetTasks(userID int, startDate, endDate string) ([]TickEntry, error) {
	getURL := fmt.Sprintf("%s/users/%d/entries?start_date=%s&end_date=%s", T.tickToken, userID, startDate, endDate)

	bodyContent, err := T.sendRequest("GET", getURL, nil)
	if err != nil {
		return nil, err
	}

	tasks := []TickEntry{}
	errUnmarshal := json.Unmarshal(bodyContent, &tasks)
	if errUnmarshal != nil {
		return nil, errUnmarshal
	}
	return tasks, nil
}

func (T *TickClient) DeleteTask(taskID int) error {
	deleteURL := fmt.Sprintf("%s/entries/%d.json", T.tickURL, taskID)
	_, err := T.sendRequest("DELETE", deleteURL, nil)
	if err != nil {
		return err
	}
	return nil
}

func (T *TickClient) UploadTask(tickEntry TickEntry) error {
	postURL := fmt.Sprintf("%s/entries.json", T.tickURL)
	body, err := json.Marshal(tickEntry)
	if err != nil {
		return err
	}
	_, err = T.sendRequest("POST", postURL, body)
	if err != nil {
		return err
	}

	return nil
}

func (T *TickClient) GetUsers() ([]UsersTick, error) {
	getURL := fmt.Sprintf("%s/users.json", T.tickURL)
	bodyContent, err := T.sendRequest("GET", getURL, nil)
	if err != nil {
		return nil, err
	}
	usersTick := []UsersTick{}
	errUnmarshal := json.Unmarshal(bodyContent, &usersTick)
	if errUnmarshal != nil {
		return nil, errUnmarshal
	}
	return usersTick, nil
}
func (T *TickClient) GetUserByEmail(email string) (*UsersTick, error) {
	getURL := fmt.Sprintf("%s/users.json?email=%s", T.tickURL, email)
	bodyContent, err := T.sendRequest("GET", getURL, nil)
	if err != nil {
		return nil, err
	}

	usersTick := []UsersTick{}
	errUnmarshal := json.Unmarshal(bodyContent, &usersTick)
	if errUnmarshal != nil {
		return nil, errUnmarshal
	}
	if len(usersTick) == 0 {
		return nil, errors.New("user not found")
	}
	return &usersTick[0], nil
}
