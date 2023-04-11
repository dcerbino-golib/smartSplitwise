package smartsplitwise

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/aanzolaavila/splitwise.go"
	"github.com/aanzolaavila/splitwise.go/resources"
	"github.com/stretchr/testify/assert"
)

const getFriends200Response = `
{
  "friends": [
    {
      "id": 15,
      "first_name": "Ada",
      "last_name": "Lovelace",
      "email": "ada@example.com",
      "registration_status": "confirmed",
      "picture": {
        "small": "string",
        "medium": "string",
        "large": "string"
      },
      "groups": [
        {
          "group_id": 571,
          "balance": [
            {
              "currency_code": "USD",
              "amount": "414.5"
            }
          ]
        }
      ],
      "balance": [
        {
          "currency_code": "USD",
          "amount": "414.5"
        }
      ],
      "updated_at": "2019-08-24T14:15:22Z"
    },
	{
		"id": 16,
		"first_name": "Pepe",
		"last_name": "ponce",
		"email": "pepe@example.com",
		"registration_status": "confirmed",
		"picture": {
		  "small": "string",
		  "medium": "string",
		  "large": "string"
		},
		"groups": [
		  {
			"group_id": 571,
			"balance": [
			  {
				"currency_code": "USD",
				"amount": "329.5"
			  }
			]
		  }
		],
		"balance": [
		  {
			"currency_code": "USD",
			"amount": "329.5"
		  }
		],
		"updated_at": "2019-08-23T14:15:22Z"
	  }
  ]
}
`

const testUser = `
{
    "user": {
        "id": 123456,
        "first_name": "test",
        "last_name": "test",
        "picture": {
            "small": "https://splitwise.s3.amazonaws.com/uploads/user/avatar/123456/small_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg",
            "medium": "https://splitwise.s3.amazonaws.com/uploads/user/avatar/123456/medium_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg",
            "large": "https://splitwise.s3.amazonaws.com/uploads/user/avatar/123456/large_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg"
        },
        "custom_picture": true,
        "email": "test@test.com",
        "registration_status": "confirmed",
        "force_refresh_at": null,
        "locale": "it",
        "country_code": "AR",
        "date_format": "MM/DD/YYYY",
        "default_currency": "ARS",
        "default_group_id": -1,
        "notifications_read": "2023-04-10T13:28:18Z",
        "notifications_count": 0,
        "notifications": {
            "added_as_friend": true,
            "added_to_group": true,
            "expense_added": false,
            "expense_updated": false,
            "bills": true,
            "payments": true,
            "monthly_summary": true,
            "announcements": true
        }
    }
}
`

const testFriend = `{
    "friend": {
        "id": 12345,
        "first_name": "test",
        "last_name": "test",
        "email": "test@test.com",
        "registration_status": "confirmed",
        "picture": {
            "small": "https://s3.amazonaws.com/splitwise/uploads/user/default_avatars/avatar-orange1-50px.png",
            "medium": "https://s3.amazonaws.com/splitwise/uploads/user/default_avatars/avatar-orange1-100px.png",
            "large": "https://s3.amazonaws.com/splitwise/uploads/user/default_avatars/avatar-orange1-200px.png"
        },
        "balance": [
            {
                "currency_code": "ARS",
                "amount": "5032879.77"
            },
            {
                "currency_code": "USD",
                "amount": "525.0"
            },
            {
                "currency_code": "EUR",
                "amount": "2514.68"
            }
        ],
        "groups": [
            {
                "group_id": 11741221,
                "balance": [
                    {
                        "currency_code": "ARS",
                        "amount": "4983304.52"
                    },
                    {
                        "currency_code": "USD",
                        "amount": "525.0"
                    }
                ]
            },
            {
                "group_id": 11794860,
                "balance": [
                    {
                        "currency_code": "ARS",
                        "amount": "49575.25"
                    }
                ]
            },
            {
                "group_id": 13548002,
                "balance": []
            },
            {
                "group_id": 29044250,
                "balance": [
                    {
                        "currency_code": "EUR",
                        "amount": "2514.68"
                    }
                ]
            },
            {
                "group_id": 0,
                "balance": []
            }
        ],
        "updated_at": "2023-04-10T21:37:48Z"
    }
}`

const testGroup = `{
    "group": {
        "id": 12345,
        "name": "test",
        "created_at": "2019-03-02T02:39:34Z",
        "updated_at": "2023-04-10T21:37:48Z",
        "members": [
            {
                "id": 21623741,
                "first_name": "test1",
                "last_name": "test1",
                "picture": {
                    "small": "https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/small_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg",
                    "medium": "https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/medium_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg",
                    "large": "https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/large_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg"
                },
                "custom_picture": true,
                "email": "test1@test.com",
                "registration_status": "confirmed",
                "balance": [
                    {
                        "currency_code": "ARS",
                        "amount": "4983304.52"
                    },
                    {
                        "currency_code": "USD",
                        "amount": "525.0"
                    }
                ]
            },
            {
                "id": 21679690,
                "first_name": "test2",
                "last_name": "test2",
                "picture": {
                    "small": "https://s3.amazonaws.com/splitwise/uploads/user/default_avatars/avatar-orange1-50px.png",
                    "medium": "https://s3.amazonaws.com/splitwise/uploads/user/default_avatars/avatar-orange1-100px.png",
                    "large": "https://s3.amazonaws.com/splitwise/uploads/user/default_avatars/avatar-orange1-200px.png"
                },
                "custom_picture": false,
                "email": "test2@test.com",
                "registration_status": "confirmed",
                "balance": [
                    {
                        "currency_code": "ARS",
                        "amount": "-4983304.52"
                    },
                    {
                        "currency_code": "USD",
                        "amount": "-525.0"
                    }
                ]
            }
        ],
        "simplify_by_default": false,
        "original_debts": [
            {
                "from": 21679690,
                "to": 21623741,
                "amount": "4983304.52",
                "currency_code": "ARS"
            },
            {
                "from": 21679690,
                "to": 21623741,
                "amount": "525.0",
                "currency_code": "USD"
            }
        ],
        "simplified_debts": [
            {
                "from": 21679690,
                "to": 21623741,
                "amount": "4983304.52",
                "currency_code": "ARS"
            },
            {
                "from": 21679690,
                "to": 21623741,
                "amount": "525.0",
                "currency_code": "USD"
            }
        ],
        "whiteboard": null,
        "group_type": "apartment",
        "invite_link": "https://www.splitwise.com/join/DoD65gHhn9w+cvgzh",
        "group_reminders": null,
        "avatar": {
            "small": "https://splitwise.s3.amazonaws.com/uploads/group/avatar/11741221/small_2a568696-c654-4fca-b761-a5554d9c2c2e.jpeg",
            "medium": "https://splitwise.s3.amazonaws.com/uploads/group/avatar/11741221/medium_2a568696-c654-4fca-b761-a5554d9c2c2e.jpeg",
            "large": "https://splitwise.s3.amazonaws.com/uploads/group/avatar/11741221/large_2a568696-c654-4fca-b761-a5554d9c2c2e.jpeg",
            "xlarge": "https://splitwise.s3.amazonaws.com/uploads/group/avatar/11741221/xlarge_2a568696-c654-4fca-b761-a5554d9c2c2e.jpeg",
            "xxlarge": "https://splitwise.s3.amazonaws.com/uploads/group/avatar/11741221/xxlarge_2a568696-c654-4fca-b761-a5554d9c2c2e.jpeg",
            "original": "https://splitwise.s3.amazonaws.com/uploads/group/avatar/11741221/2a568696-c654-4fca-b761-a5554d9c2c2e.jpeg"
        },
        "tall_avatar": {
            "xlarge": "https://splitwise.s3.amazonaws.com/uploads/group/avatar/11741221/xlarge_2a568696-c654-4fca-b761-a5554d9c2c2e.jpeg",
            "large": "https://splitwise.s3.amazonaws.com/uploads/group/avatar/11741221/large_2a568696-c654-4fca-b761-a5554d9c2c2e.jpeg"
        },
        "custom_avatar": true,
        "cover_photo": {
            "xxlarge": "https://splitwise.s3.amazonaws.com/uploads/group/avatar/11741221/xxlarge_2a568696-c654-4fca-b761-a5554d9c2c2e.jpeg",
            "xlarge": "https://splitwise.s3.amazonaws.com/uploads/group/avatar/11741221/xlarge_2a568696-c654-4fca-b761-a5554d9c2c2e.jpeg"
        }
    }
}`

const testExpence = `{
    "expense": {
        "id": 12345,
        "group_id": 12345,
        "friendship_id": null,
        "expense_bundle_id": null,
        "description": "TelViso",
        "repeats": false,
        "repeat_interval": null,
        "email_reminder": false,
        "email_reminder_in_advance": -1,
        "next_repeat": null,
        "details": null,
        "comments_count": 0,
        "payment": false,
        "creation_method": "equal",
        "transaction_method": "offline",
        "transaction_confirmed": false,
        "transaction_id": null,
        "transaction_status": null,
        "cost": "2048.49",
        "currency_code": "ARS",
        "repayments": [
            {
                "from": 21623741,
                "to": 21702157,
                "amount": "1024.25"
            }
        ],
        "date": "2023-04-11T01:49:27Z",
        "created_at": "2023-04-11T01:50:38Z",
        "created_by": {
            "id": 21702157,
            "first_name": "Test",
            "last_name": "Test",
            "picture": {
                "medium": "https://s3.amazonaws.com/splitwise/uploads/user/default_avatars/avatar-ruby2-100px.png"
            },
            "custom_picture": false
        },
        "updated_at": "2023-04-11T01:51:30Z",
        "updated_by": {
            "id": 21702157,
            "first_name": "Test",
            "last_name": "Test",
            "picture": {
                "medium": "https://s3.amazonaws.com/splitwise/uploads/user/default_avatars/avatar-ruby2-100px.png"
            },
            "custom_picture": false
        },
        "deleted_at": null,
        "deleted_by": null,
        "category": {
            "id": 8,
            "name": "TV/Telefono/Internet"
        },
        "receipt": {
            "large": "https://splitwise.s3.amazonaws.com/uploads/expense/receipt/2292261592/large_6f10dd05-019c-49f0-816a-cb0818aa982b.png",
            "original": "https://splitwise.s3.amazonaws.com/uploads/expense/receipt/2292261592/6f10dd05-019c-49f0-816a-cb0818aa982b.pdf"
        },
        "users": [
            {
                "user": {
                    "id": 21702157,
                    "first_name": "Test",
                    "last_name": "Test",
                    "picture": {
                        "medium": "https://s3.amazonaws.com/splitwise/uploads/user/default_avatars/avatar-ruby2-100px.png"
                    }
                },
                "user_id": 21702157,
                "paid_share": "2048.49",
                "owed_share": "1024.24",
                "net_balance": "1024.25"
            },
            {
                "user": {
                    "id": 21623741,
                    "first_name": "Test2",
                    "last_name": "Test2",
                    "picture": {
                        "medium": "https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/medium_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg"
                    }
                },
                "user_id": 21623741,
                "paid_share": "0.0",
                "owed_share": "1024.25",
                "net_balance": "-1024.25"
            }
        ],
        "comments": []
    }
}`
const unauthorized = `{
    "error": "Invalid API Request: you are not logged in"
}`

type httpClientStub struct {
	DoFunc func(*http.Request) (*http.Response, error)
}

func (c httpClientStub) Do(r *http.Request) (*http.Response, error) {
	if c.DoFunc == nil {
		panic("mocked function is nil")
	}

	return c.DoFunc(r)
}

type logger interface {
	Printf(string, ...interface{})
}

type testLogger struct {
	buf    bytes.Buffer
	logger logger
	once   sync.Once
	T      *testing.T
}

func (l testLogger) Printf(s string, args ...interface{}) {
	l.once.Do(func() {
		tname := l.T.Name()
		prefix := fmt.Sprintf("%s:: ", tname)
		l.logger = log.New(io.Writer(&l.buf), prefix, log.LstdFlags)

		l.T.Cleanup(func() {
			if l.T.Failed() {
				fmt.Print(l.buf.String())
			}
		})
	})

	l.logger.Printf(s, args...)
}

func TestOpen(t *testing.T) {
	assert := assert.New(t)

	token := "test"
	ctx := context.Background()
	log := log.New(os.Stdout, "Splitwise LOG: ", log.Lshortfile)

	result := Open(token, ctx, log)

	assert.Equal(token, result.getClient().Token)
	assert.Equal(ctx, result.getCtx())
	assert.Equal(log, result.getClient().Logger)
}

func TestMainCategoryCache(t *testing.T) {
	assert.Equal(t, true, len(mainCategoryCache) >= 7)
}

func TestCurenciesCache(t *testing.T) {
	assert.Equal(t, true, len(curenciesCache) >= 148)
}

func TestClose(t *testing.T) {
	assert := assert.New(t)

	token := "test"
	ctx := context.Background()
	log := log.New(os.Stdout, "Splitwise LOG: ", log.Lshortfile)

	conn := Open(token, ctx, log)

	executor := conn.GetCurecies()

	assert.Equal(false, executor.isClose())
	executor.Close()
	assert.Equal(true, executor.isClose())
}

func TestGetCategory(t *testing.T) {
	assert := assert.New(t)

	token := "test"
	ctx := context.Background()
	log := log.New(os.Stdout, "Splitwise LOG: ", log.Lshortfile)

	conn := Open(token, ctx, log)

	category, err := conn.GetMainCategory(resources.Identifier(1))
	assert.Equal(nil, err)
	assert.Equal("Utilities", category.Name)

}

func TestGetCategoryNotFound(t *testing.T) {
	token := "test"
	ctx := context.Background()
	log := log.New(os.Stdout, "Splitwise LOG: ", log.Lshortfile)

	conn := Open(token, ctx, log)

	_, err := conn.GetMainCategory(resources.Identifier(0))
	assert.EqualErrorf(t, err, (&ElementNotFound{}).Error(), "Error should be: %v, got: %v", (&ElementNotFound{}).Error(), err)

}

func TestGetCategoryies(t *testing.T) {
	token := "test"
	ctx := context.Background()
	log := log.New(os.Stdout, "Splitwise LOG: ", log.Lshortfile)

	conn := Open(token, ctx, log)

	executor := conn.GetMainCategories()
	cont := 0
	want := 7

	for range executor.GetChan() {
		cont++
	}

	assert.GreaterOrEqual(t, cont, want, "Get categories should return unless %d arg and got it %d", want, cont)
}

func TestGetCurrencies(t *testing.T) {
	token := "test"
	ctx := context.Background()
	log := log.New(os.Stdout, "Splitwise LOG: ", log.Lshortfile)

	conn := Open(token, ctx, log)

	executor := conn.GetCurecies()
	cont := 0
	want := 148

	for range executor.GetChan() {
		cont++
	}

	assert.GreaterOrEqual(t, cont, want, "Get currencies should return unless %d currencies and got it %d", want, cont)
}

func TestGetCurrency(t *testing.T) {
	assert := assert.New(t)

	token := "test"
	ctx := context.Background()
	log := log.New(os.Stdout, "Splitwise LOG: ", log.Lshortfile)

	conn := Open(token, ctx, log)

	currencyCode := "USD"
	currencyUnit := "$"

	currency, err := conn.GetCurency(currencyCode)

	assert.NoError(err, "%s should be present as a currency code", currencyCode)
	assert.Equal(currencyUnit, currency.Unit, "%s currency unit should be %s but got %s", currencyCode, currencyUnit, currency.Unit)
}

func TestGetCurrencyNotFund(t *testing.T) {
	assert := assert.New(t)

	token := "test"
	ctx := context.Background()
	log := log.New(os.Stdout, "Splitwise LOG: ", log.Lshortfile)

	conn := Open(token, ctx, log)

	currencyCode := "US"

	_, err := conn.GetCurency(currencyCode)

	assert.EqualErrorf(err, (&ElementNotFound{}).Error(), "Error should be: %v, got: %v", (&ElementNotFound{}).Error(), err)
}

func TestGetFriends(t *testing.T) {
	doFunc := func(r *http.Request) (*http.Response, error) {
		resposne := http.Response{}
		resposne.Body = io.NopCloser(strings.NewReader(getFriends200Response))
		resposne.Header = make(map[string][]string)
		resposne.Header["Content-Type"] = []string{"application/json", "charset=utf-8"}
		resposne.Status = "200"
		resposne.StatusCode = 200
		return &resposne, nil
	}

	type responseStruct struct {
		Friends []resources.Friend
	}
	wantedRespounce := responseStruct{}

	err := json.Unmarshal([]byte(getFriends200Response), &wantedRespounce)

	if err != nil {
		panic(err)
	}

	conn := getTestConnection(t, doFunc)

	executor := conn.GetFriends()

	cont := 0

	for range executor.GetChan() {
		cont++
	}

	assert.Equal(t, len(wantedRespounce.Friends), cont)
}

func TestCurrentUser(t *testing.T) {
	assert.Nil(t, currentUser)
	doFunc := func(r *http.Request) (*http.Response, error) {
		resposne := http.Response{}
		resposne.Body = io.NopCloser(strings.NewReader(testUser))
		resposne.Header = make(map[string][]string)
		resposne.Header["Content-Type"] = []string{"application/json", "charset=utf-8"}
		resposne.Status = "200"
		resposne.StatusCode = 200
		return &resposne, nil
	}

	type responseStruct struct {
		User resources.User
	}

	wantedRespounce := responseStruct{}
	err := json.Unmarshal([]byte(testUser), &wantedRespounce)

	if err != nil {
		panic(err)
	}

	conn := getTestConnection(t, doFunc)

	user, err := conn.GetCurrentUser()

	if err != nil {
		panic(err)
	}

	assert.Equal(t, wantedRespounce.User, user)
	assert.NotNil(t, currentUser)

	user, err = conn.GetCurrentUser()

	if err != nil {
		panic(err)
	}
	assert.Equal(t, wantedRespounce.User, user)
}
func TestCurrentUserWhenUnathorized(t *testing.T) {
	currentUser = nil
	doFunc := func(r *http.Request) (*http.Response, error) {
		resposne := http.Response{}
		resposne.Body = io.NopCloser(strings.NewReader(unauthorized))
		resposne.Header = make(map[string][]string)
		resposne.Header["Content-Type"] = []string{"application/json", "charset=utf-8"}
		resposne.Status = "401"
		resposne.StatusCode = 401
		return &resposne, nil
	}
	conn := getTestConnection(t, doFunc)

	user, err := conn.GetCurrentUser()
	fmt.Println(user)
	assert.ErrorIs(t, err, splitwise.ErrNotLoggedIn, "Test fail we expected %s but got %s", splitwise.ErrNotLoggedIn, err)

}

func TestGetFriend(t *testing.T) {
	doFunc := func(r *http.Request) (*http.Response, error) {
		resposne := http.Response{}
		resposne.Body = io.NopCloser(strings.NewReader(testFriend))
		resposne.Header = make(map[string][]string)
		resposne.Header["Content-Type"] = []string{"application/json", "charset=utf-8"}
		resposne.Status = "200"
		resposne.StatusCode = 200
		return &resposne, nil
	}

	type responseStruct struct {
		Friend resources.Friend
	}

	wantedRespounce := responseStruct{}
	err := json.Unmarshal([]byte(testFriend), &wantedRespounce)
	assert.NoError(t, err)

	conn := getTestConnection(t, doFunc)
	friend, err := conn.GetFriend(12345)

	assert.NoError(t, err)

	assert.Equal(t, wantedRespounce.Friend, friend)

}

func TestGetExpence(t *testing.T) {
	doFunc := func(r *http.Request) (*http.Response, error) {
		resposne := http.Response{}
		resposne.Body = io.NopCloser(strings.NewReader(testExpence))
		resposne.Header = make(map[string][]string)
		resposne.Header["Content-Type"] = []string{"application/json", "charset=utf-8"}
		resposne.Status = "200"
		resposne.StatusCode = 200
		return &resposne, nil
	}

	type responseStruct struct {
		Expense resources.Expense
	}

	wantedRespounce := responseStruct{}
	err := json.Unmarshal([]byte(testExpence), &wantedRespounce)
	assert.NoError(t, err)

	conn := getTestConnection(t, doFunc)
	expence, err := conn.GetExpense(12345)

	assert.NoError(t, err)

	assert.Equal(t, wantedRespounce.Expense, expence)
}

func TestGetGroup(t *testing.T) {
	doFunc := func(r *http.Request) (*http.Response, error) {
		resposne := http.Response{}
		resposne.Body = io.NopCloser(strings.NewReader(testGroup))
		resposne.Header = make(map[string][]string)
		resposne.Header["Content-Type"] = []string{"application/json", "charset=utf-8"}
		resposne.Status = "200"
		resposne.StatusCode = 200
		return &resposne, nil
	}

	type responseStruct struct {
		Group resources.Group
	}

	wantedRespounce := responseStruct{}
	err := json.Unmarshal([]byte(testGroup), &wantedRespounce)
	assert.NoError(t, err)

	conn := getTestConnection(t, doFunc)
	group, err := conn.GetGroup(12345)

	assert.NoError(t, err)

	assert.Equal(t, wantedRespounce.Group, group)

}

func getTestConnection(t *testing.T, doFunc func(r *http.Request) (*http.Response, error)) SwConnection {
	client := Open("testtoken", context.Background(), log.New(os.Stdout, "Test Splitwise LOG: ", log.Lshortfile))

	bareclient, ok := client.(*swConnectionStruct)

	if !ok {
		panic("unable to convert interface to istance")
	}

	cliststub := httpClientStub{
		DoFunc: doFunc,
	}

	bareclient.client.HttpClient = cliststub
	bareclient.client.Logger = testLogger{
		T: t,
	}

	return bareclient
}
