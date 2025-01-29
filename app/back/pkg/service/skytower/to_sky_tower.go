package skytower

import (
	"encoding/json"
	"fmt"
	"main/pkg/service/do_request"
	"net/http"
	"net/url"
	"strconv"
)

type RequestToSkyTower struct {
	Token    string `json:"token"`
	Host     string `json:"host"`
	AuthData Data   `json:"authData"`
}

type Data struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func New(host, username, password string) (*RequestToSkyTower, error) {
	token, err := login(username, password, host)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "request_to_sky_tower.New", err)
	}
	return &RequestToSkyTower{
		Host:     host,
		AuthData: Data{Username: username, Password: password},
		Token:    token,
	}, nil
}

func (a *RequestToSkyTower) updateAuth() error {
	var err error
	a.Token, err = login(a.AuthData.Username, a.AuthData.Password, a.Host)
	return err
}

func login(username, password, host string) (string, error) {
	const method = "login"
	type ResponseToken struct {
		Token string `json:"token"`
	}
	var res ResponseToken
	bJson := Data{Username: username, Password: password}
	bodyJson, err := json.Marshal(bJson)
	if err != nil {
		fmt.Printf("cannot do Marshal: %s", err)
		return "", fmt.Errorf("cannot do Marshal %s: %w", method, err)
	}

	var c = do_request.New(host, "/api/auth")
	data, _, err := c.DoRequestPost(method, bodyJson, http.MethodPost, nil)
	if err != nil {
		return "", fmt.Errorf("cannot do %s: %w", method, err)
	}

	if err = json.Unmarshal(data, &res); err != nil {
		return "", fmt.Errorf("error Unmarshal: %s", string(data))
	}

	return res.Token, nil
}

// sendRequest Выполняет отправку запроса в сервис skytower
func (a *RequestToSkyTower) sendRequest(httpMethod, uri string, bJson []byte, query url.Values, apiRes interface{}) error {
	toApi, status, err := a.send(httpMethod, uri, bJson, query)

	if err != nil {
		return fmt.Errorf("%s", err)
	} else if status >= 400 {
		return fmt.Errorf("error: status code = %s", strconv.Itoa(status))
	} else if err = json.Unmarshal(toApi, &apiRes); err != nil {
		return fmt.Errorf("error Unmarshal: %s", string(toApi))
	}
	return nil
}

func (a *RequestToSkyTower) sendAndHandleRequest(method, url string, body []byte, query url.Values, result interface{}) error {
	if err := a.sendRequest(method, url, body, query, result); err != nil {
		return fmt.Errorf("cannot do %s: %w", url, err)
	}

	// Проверяем, есть ли в результате поле Error, Message
	if resp, ok := result.(*struct {
		Error   int    `json:"error"`
		Message string `json:"message"`
	}); ok {
		if resp.Error > 0 {
			return fmt.Errorf("%s: error code = %d, error_msg = %s", url, resp.Error, resp.Message)
		}
	}

	return nil
}

func (a *RequestToSkyTower) send(httpMethod, uri string, bJson []byte, query url.Values) ([]byte, int, error) {
	var err error
	var c = do_request.Client{}
	c = do_request.New(a.Host, "/api")
	c.TokenSet(a.Token)
	var data []byte
	var statusCode int
	if httpMethod == http.MethodGet {
		data, statusCode, err = c.DoRequest(uri, query)
	} else {
		data, statusCode, err = c.DoRequestPost(uri, bJson, httpMethod, query)
	}
	if statusCode == 401 {
		// попытка перавторизоваться
		if err = a.updateAuth(); err != nil {
			return nil, 400, fmt.Errorf("ошибка при авторизации: %v", err)
		} else {
			fmt.Println("токен не валидный - переавторизовываюсь")
			c.TokenSet(a.Token)
			if httpMethod == http.MethodGet {
				data, statusCode, err = c.DoRequest(uri, query)
			} else {
				data, statusCode, err = c.DoRequestPost(uri, bJson, httpMethod, query)
			}
		}
	}
	if err != nil {
		// Todo Сюда логирование fmt.Errorf("cannot do SendRequest %s: %w", uri, err)
		return nil, statusCode, err
	}

	return data, statusCode, nil
}
