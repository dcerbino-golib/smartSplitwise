package smartsplitwise

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aanzolaavila/splitwise.go"
	"github.com/aanzolaavila/splitwise.go/resources"
	"github.com/stretchr/testify/assert"
)

const testExpences = `{"expenses":[{"id":2123851796,"group_id":11741221,"friendship_id":null,"expense_bundle_id":null,"description":"Jumbo","repeats":false,"repeat_interval":null,"email_reminder":false,"email_reminder_in_advance":-1,"next_repeat":null,"details":null,"comments_count":1,"payment":false,"creation_method":null,"transaction_method":"offline","transaction_confirmed":false,"transaction_id":null,"transaction_status":null,"cost":"1083.92","currency_code":"ARS","repayments":[{"from":21679690,"to":21623741,"amount":"541.96"}],"date":"2023-01-09T14:41:00Z","created_at":"2023-01-11T14:42:02Z","created_by":{"id":21623741,"first_name":"test1","last_name":"test","picture":{"medium":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/medium_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg"},"custom_picture":true},"updated_at":"2023-01-11T17:18:08Z","updated_by":{"id":21623741,"first_name":"test1","last_name":"test","picture":{"medium":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/medium_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg"},"custom_picture":true},"deleted_at":null,"deleted_by":null,"category":{"id":12,"name":"Alimentari"},"receipt":{"large":"https://splitwise.s3.amazonaws.com/uploads/expense/receipt/2123851796/large_f6af2889-3255-429a-b6ba-915d8d4b3376.png","original":"https://splitwise.s3.amazonaws.com/uploads/expense/receipt/2123851796/f6af2889-3255-429a-b6ba-915d8d4b3376.pdf"},"users":[{"user":{"id":21679690,"first_name":"test2","last_name":"test","picture":{"medium":"https://s3.amazonaws.com/splitwise/uploads/user/default_avatars/avatar-orange1-100px.png"}},"user_id":21679690,"paid_share":"0.0","owed_share":"541.96","net_balance":"-541.96"},{"user":{"id":21623741,"first_name":"test1","last_name":"test","picture":{"medium":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/medium_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg"}},"user_id":21623741,"paid_share":"1083.92","owed_share":"541.96","net_balance":"541.96"}]},{"id":2115167389,"group_id":11741221,"friendship_id":null,"expense_bundle_id":null,"description":"Fiambre","repeats":false,"repeat_interval":null,"email_reminder":false,"email_reminder_in_advance":-1,"next_repeat":null,"details":null,"comments_count":0,"payment":false,"creation_method":"equal","transaction_method":"offline","transaction_confirmed":false,"transaction_id":null,"transaction_status":null,"cost":"1185.0","currency_code":"ARS","repayments":[{"from":21679690,"to":21623741,"amount":"592.5"}],"date":"2023-01-06T21:44:18Z","created_at":"2023-01-06T21:44:52Z","created_by":{"id":21623741,"first_name":"test1","last_name":"test","picture":{"medium":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/medium_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg"},"custom_picture":true},"updated_at":"2023-01-06T21:44:52Z","updated_by":null,"deleted_at":null,"deleted_by":null,"category":{"id":12,"name":"Alimentari"},"receipt":{"large":null,"original":null},"users":[{"user":{"id":21623741,"first_name":"test1","last_name":"test","picture":{"medium":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/medium_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg"}},"user_id":21623741,"paid_share":"1185.0","owed_share":"592.5","net_balance":"592.5"},{"user":{"id":21679690,"first_name":"test2","last_name":"test","picture":{"medium":"https://s3.amazonaws.com/splitwise/uploads/user/default_avatars/avatar-orange1-100px.png"}},"user_id":21679690,"paid_share":"0.0","owed_share":"592.5","net_balance":"-592.5"}]},{"id":2115163070,"group_id":11741221,"friendship_id":null,"expense_bundle_id":null,"description":"Jumbo","repeats":false,"repeat_interval":null,"email_reminder":false,"email_reminder_in_advance":-1,"next_repeat":null,"details":null,"comments_count":0,"payment":false,"creation_method":null,"transaction_method":"offline","transaction_confirmed":false,"transaction_id":null,"transaction_status":null,"cost":"23134.09","currency_code":"ARS","repayments":[{"from":21679690,"to":21623741,"amount":"11567.05"}],"date":"2023-01-06T21:40:55Z","created_at":"2023-01-06T21:41:23Z","created_by":{"id":21623741,"first_name":"test1","last_name":"test","picture":{"medium":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/medium_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg"},"custom_picture":true},"updated_at":"2023-01-06T21:41:24Z","updated_by":{"id":21623741,"first_name":"test1","last_name":"test","picture":{"medium":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/medium_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg"},"custom_picture":true},"deleted_at":null,"deleted_by":null,"category":{"id":12,"name":"Alimentari"},"receipt":{"large":"https://splitwise.s3.amazonaws.com/uploads/expense/receipt/2115163070/large_28a40d87-5e70-4942-a9c8-673b8f8e2f40.png","original":"https://splitwise.s3.amazonaws.com/uploads/expense/receipt/2115163070/28a40d87-5e70-4942-a9c8-673b8f8e2f40.pdf"},"users":[{"user":{"id":21623741,"first_name":"test1","last_name":"test","picture":{"medium":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/medium_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg"}},"user_id":21623741,"paid_share":"23134.09","owed_share":"11567.04","net_balance":"11567.05"},{"user":{"id":21679690,"first_name":"test2","last_name":"test","picture":{"medium":"https://s3.amazonaws.com/splitwise/uploads/user/default_avatars/avatar-orange1-100px.png"}},"user_id":21679690,"paid_share":"0.0","owed_share":"11567.05","net_balance":"-11567.05"}]},{"id":2115160067,"group_id":11741221,"friendship_id":null,"expense_bundle_id":null,"description":"Regalo Nati","repeats":false,"repeat_interval":null,"email_reminder":false,"email_reminder_in_advance":-1,"next_repeat":null,"details":null,"comments_count":0,"payment":false,"creation_method":"equal","transaction_method":"offline","transaction_confirmed":false,"transaction_id":null,"transaction_status":null,"cost":"3680.0","currency_code":"ARS","repayments":[{"from":21679690,"to":21623741,"amount":"1840.0"}],"date":"2023-01-06T21:38:07Z","created_at":"2023-01-06T21:39:04Z","created_by":{"id":21623741,"first_name":"test1","last_name":"test","picture":{"medium":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/medium_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg"},"custom_picture":true},"updated_at":"2023-01-06T21:39:04Z","updated_by":null,"deleted_at":null,"deleted_by":null,"category":{"id":42,"name":"Regali"},"receipt":{"large":null,"original":null},"users":[{"user":{"id":21623741,"first_name":"test1","last_name":"test","picture":{"medium":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/medium_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg"}},"user_id":21623741,"paid_share":"3680.0","owed_share":"1840.0","net_balance":"1840.0"},{"user":{"id":21679690,"first_name":"test2","last_name":"test","picture":{"medium":"https://s3.amazonaws.com/splitwise/uploads/user/default_avatars/avatar-orange1-100px.png"}},"user_id":21679690,"paid_share":"0.0","owed_share":"1840.0","net_balance":"-1840.0"}]},{"id":2114356059,"group_id":11741221,"friendship_id":null,"expense_bundle_id":null,"description":"Payment","repeats":false,"repeat_interval":null,"email_reminder":false,"email_reminder_in_advance":-1,"next_repeat":null,"details":null,"comments_count":1,"payment":true,"creation_method":"payment","transaction_method":"offline","transaction_confirmed":false,"transaction_id":null,"transaction_status":null,"cost":"5000.0","currency_code":"ARS","repayments":[{"from":21679690,"to":21623741,"amount":"5000.0"}],"date":"2023-01-06T13:38:32Z","created_at":"2023-01-06T13:38:35Z","created_by":{"id":21623741,"first_name":"test1","last_name":"test","picture":{"medium":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/medium_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg"},"custom_picture":true},"updated_at":"2023-01-06T13:38:35Z","updated_by":null,"deleted_at":null,"deleted_by":null,"category":{"id":18,"name":"Generali"},"receipt":{"large":null,"original":null},"users":[{"user":{"id":21623741,"first_name":"test1","last_name":"test","picture":{"medium":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/medium_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg"}},"user_id":21623741,"paid_share":"5000.0","owed_share":"0.0","net_balance":"5000.0"},{"user":{"id":21679690,"first_name":"test2","last_name":"test","picture":{"medium":"https://s3.amazonaws.com/splitwise/uploads/user/default_avatars/avatar-orange1-100px.png"}},"user_id":21679690,"paid_share":"0.0","owed_share":"5000.0","net_balance":"-5000.0"}]},{"id":2114350422,"group_id":11741221,"friendship_id":null,"expense_bundle_id":null,"description":"Matafuego","repeats":false,"repeat_interval":null,"email_reminder":false,"email_reminder_in_advance":-1,"next_repeat":null,"details":null,"comments_count":0,"payment":false,"creation_method":null,"transaction_method":"offline","transaction_confirmed":false,"transaction_id":null,"transaction_status":null,"cost":"5000.0","currency_code":"ARS","repayments":[{"from":21623741,"to":21679690,"amount":"2500.0"}],"date":"2023-01-06T13:35:16Z","created_at":"2023-01-06T13:35:25Z","created_by":{"id":21679690,"first_name":"test2","last_name":"test","picture":{"medium":"https://s3.amazonaws.com/splitwise/uploads/user/default_avatars/avatar-orange1-100px.png"},"custom_picture":false},"updated_at":"2023-01-06T21:42:22Z","updated_by":{"id":21623741,"first_name":"test1","last_name":"test","picture":{"medium":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/medium_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg"},"custom_picture":true},"deleted_at":null,"deleted_by":null,"category":{"id":15,"name":"Auto"},"receipt":{"large":null,"original":null},"users":[{"user":{"id":21679690,"first_name":"test2","last_name":"test","picture":{"medium":"https://s3.amazonaws.com/splitwise/uploads/user/default_avatars/avatar-orange1-100px.png"}},"user_id":21679690,"paid_share":"5000.0","owed_share":"2500.0","net_balance":"2500.0"},{"user":{"id":21623741,"first_name":"test1","last_name":"test","picture":{"medium":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/medium_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg"}},"user_id":21623741,"paid_share":"0.0","owed_share":"2500.0","net_balance":"-2500.0"}]},{"id":2114349755,"group_id":11741221,"friendship_id":null,"expense_bundle_id":null,"description":"Chocolates keto","repeats":false,"repeat_interval":null,"email_reminder":false,"email_reminder_in_advance":-1,"next_repeat":null,"details":null,"comments_count":0,"payment":false,"creation_method":null,"transaction_method":"offline","transaction_confirmed":false,"transaction_id":null,"transaction_status":null,"cost":"1140.0","currency_code":"ARS","repayments":[{"from":21623741,"to":21679690,"amount":"570.0"}],"date":"2023-01-06T13:34:28Z","created_at":"2023-01-06T13:35:02Z","created_by":{"id":21679690,"first_name":"test2","last_name":"test","picture":{"medium":"https://s3.amazonaws.com/splitwise/uploads/user/default_avatars/avatar-orange1-100px.png"},"custom_picture":false},"updated_at":"2023-01-06T21:42:31Z","updated_by":{"id":21623741,"first_name":"test1","last_name":"test","picture":{"medium":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/medium_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg"},"custom_picture":true},"deleted_at":null,"deleted_by":null,"category":{"id":12,"name":"Alimentari"},"receipt":{"large":null,"original":null},"users":[{"user":{"id":21679690,"first_name":"test2","last_name":"test","picture":{"medium":"https://s3.amazonaws.com/splitwise/uploads/user/default_avatars/avatar-orange1-100px.png"}},"user_id":21679690,"paid_share":"1140.0","owed_share":"570.0","net_balance":"570.0"},{"user":{"id":21623741,"first_name":"test1","last_name":"test","picture":{"medium":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/medium_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg"}},"user_id":21623741,"paid_share":"0.0","owed_share":"570.0","net_balance":"-570.0"}]},{"id":2114348276,"group_id":11741221,"friendship_id":null,"expense_bundle_id":null,"description":"Verduleria","repeats":false,"repeat_interval":null,"email_reminder":false,"email_reminder_in_advance":-1,"next_repeat":null,"details":null,"comments_count":0,"payment":false,"creation_method":null,"transaction_method":"offline","transaction_confirmed":false,"transaction_id":null,"transaction_status":null,"cost":"2500.0","currency_code":"ARS","repayments":[{"from":21679690,"to":21623741,"amount":"1250.0"}],"date":"2023-01-06T13:33:58Z","created_at":"2023-01-06T13:34:10Z","created_by":{"id":21679690,"first_name":"test2","last_name":"test","picture":{"medium":"https://s3.amazonaws.com/splitwise/uploads/user/default_avatars/avatar-orange1-100px.png"},"custom_picture":false},"updated_at":"2023-01-06T21:42:34Z","updated_by":{"id":21623741,"first_name":"test1","last_name":"test","picture":{"medium":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/medium_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg"},"custom_picture":true},"deleted_at":null,"deleted_by":null,"category":{"id":12,"name":"Alimentari"},"receipt":{"large":null,"original":null},"users":[{"user":{"id":21679690,"first_name":"test2","last_name":"test","picture":{"medium":"https://s3.amazonaws.com/splitwise/uploads/user/default_avatars/avatar-orange1-100px.png"}},"user_id":21679690,"paid_share":"0.0","owed_share":"1250.0","net_balance":"-1250.0"},{"user":{"id":21623741,"first_name":"test1","last_name":"test","picture":{"medium":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/medium_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg"}},"user_id":21623741,"paid_share":"2500.0","owed_share":"1250.0","net_balance":"1250.0"}]},{"id":2114347861,"group_id":11741221,"friendship_id":null,"expense_bundle_id":null,"description":"Verduleria ","repeats":false,"repeat_interval":null,"email_reminder":false,"email_reminder_in_advance":-1,"next_repeat":null,"details":null,"comments_count":0,"payment":false,"creation_method":null,"transaction_method":"offline","transaction_confirmed":false,"transaction_id":null,"transaction_status":null,"cost":"2500.0","currency_code":"ARS","repayments":[{"from":21623741,"to":21679690,"amount":"1250.0"}],"date":"2023-01-06T13:33:47Z","created_at":"2023-01-06T13:33:56Z","created_by":{"id":21679690,"first_name":"test2","last_name":"test","picture":{"medium":"https://s3.amazonaws.com/splitwise/uploads/user/default_avatars/avatar-orange1-100px.png"},"custom_picture":false},"updated_at":"2023-01-06T21:42:38Z","updated_by":{"id":21623741,"first_name":"test1","last_name":"test","picture":{"medium":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/medium_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg"},"custom_picture":true},"deleted_at":null,"deleted_by":null,"category":{"id":12,"name":"Alimentari"},"receipt":{"large":null,"original":null},"users":[{"user":{"id":21623741,"first_name":"test1","last_name":"test","picture":{"medium":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/medium_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg"}},"user_id":21623741,"paid_share":"0.0","owed_share":"1250.0","net_balance":"-1250.0"},{"user":{"id":21679690,"first_name":"test2","last_name":"test","picture":{"medium":"https://s3.amazonaws.com/splitwise/uploads/user/default_avatars/avatar-orange1-100px.png"}},"user_id":21679690,"paid_share":"2500.0","owed_share":"1250.0","net_balance":"1250.0"}]},{"id":2114348729,"group_id":11741221,"friendship_id":null,"expense_bundle_id":null,"description":"Consulta pediatra ","repeats":false,"repeat_interval":null,"email_reminder":false,"email_reminder_in_advance":-1,"next_repeat":null,"details":null,"comments_count":0,"payment":false,"creation_method":null,"transaction_method":"offline","transaction_confirmed":false,"transaction_id":null,"transaction_status":null,"cost":"5500.0","currency_code":"ARS","repayments":[{"from":21623741,"to":21679690,"amount":"2750.0"}],"date":"2023-01-05T13:34:00Z","created_at":"2023-01-06T13:34:26Z","created_by":{"id":21679690,"first_name":"test2","last_name":"test","picture":{"medium":"https://s3.amazonaws.com/splitwise/uploads/user/default_avatars/avatar-orange1-100px.png"},"custom_picture":false},"updated_at":"2023-01-06T21:42:57Z","updated_by":{"id":21623741,"first_name":"test1","last_name":"test","picture":{"medium":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/medium_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg"},"custom_picture":true},"deleted_at":null,"deleted_by":null,"category":{"id":43,"name":"Spese mediche"},"receipt":{"large":null,"original":null},"users":[{"user":{"id":21679690,"first_name":"test2","last_name":"test","picture":{"medium":"https://s3.amazonaws.com/splitwise/uploads/user/default_avatars/avatar-orange1-100px.png"}},"user_id":21679690,"paid_share":"5500.0","owed_share":"2750.0","net_balance":"2750.0"},{"user":{"id":21623741,"first_name":"test1","last_name":"test","picture":{"medium":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/medium_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg"}},"user_id":21623741,"paid_share":"0.0","owed_share":"2750.0","net_balance":"-2750.0"}]}]}`

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

const testGroups = `{"groups":[{"id":0,"name":"Spese senza gruppo","created_at":"2019-03-02T01:25:44Z","updated_at":"2023-04-11T17:30:41Z","members":[{"id":21623741,"first_name":"test1","last_name":"test","picture":{"small":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/small_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg","medium":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/medium_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg","large":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/large_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg"},"custom_picture":true,"email":"test@test.com","registration_status":"confirmed","balance":[]}],"simplify_by_default":false,"original_debts":[],"simplified_debts":[],"avatar":{"small":"https://s3.amazonaws.com/splitwise/uploads/group/default_avatars/v2021/avatar-nongroup-50px.png","medium":"https://s3.amazonaws.com/splitwise/uploads/group/default_avatars/v2021/avatar-nongroup-100px.png","large":"https://s3.amazonaws.com/splitwise/uploads/group/default_avatars/v2021/avatar-nongroup-200px.png","xlarge":"https://s3.amazonaws.com/splitwise/uploads/group/default_avatars/v2021/avatar-nongroup-500px.png","xxlarge":"https://s3.amazonaws.com/splitwise/uploads/group/default_avatars/v2021/avatar-nongroup-1000px.png","original":null},"tall_avatar":{"xlarge":"https://s3.amazonaws.com/splitwise/uploads/group/default_tall_avatars/avatar-nongroup-288px.png","large":"https://s3.amazonaws.com/splitwise/uploads/group/default_tall_avatars/avatar-nongroup-192px.png"},"custom_avatar":false,"cover_photo":{"xxlarge":"https://s3.amazonaws.com/splitwise/uploads/group/default_cover_photos/coverphoto-nongroup-1000px.png","xlarge":"https://s3.amazonaws.com/splitwise/uploads/group/default_cover_photos/coverphoto-nongroup-500px.png"}},{"id":11741221,"name":"Familia test test","created_at":"2019-03-02T02:39:34Z","updated_at":"2023-04-10T21:37:48Z","members":[{"id":21623741,"first_name":"test1","last_name":"test","picture":{"small":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/small_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg","medium":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/medium_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg","large":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/large_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg"},"custom_picture":true,"email":"test@test.com","registration_status":"confirmed","balance":[{"currency_code":"ARS","amount":"4983304.52"},{"currency_code":"USD","amount":"525.0"}]},{"id":21679690,"first_name":"test2","last_name":"test","picture":{"small":"https://s3.amazonaws.com/splitwise/uploads/user/default_avatars/avatar-orange1-50px.png","medium":"https://s3.amazonaws.com/splitwise/uploads/user/default_avatars/avatar-orange1-100px.png","large":"https://s3.amazonaws.com/splitwise/uploads/user/default_avatars/avatar-orange1-200px.png"},"custom_picture":false,"email":"test2@test.com","registration_status":"confirmed","balance":[{"currency_code":"ARS","amount":"-4983304.52"},{"currency_code":"USD","amount":"-525.0"}]}],"simplify_by_default":false,"original_debts":[{"from":21679690,"to":21623741,"amount":"4983304.52","currency_code":"ARS"},{"from":21679690,"to":21623741,"amount":"525.0","currency_code":"USD"}],"simplified_debts":[{"from":21679690,"to":21623741,"amount":"4983304.52","currency_code":"ARS"},{"from":21679690,"to":21623741,"amount":"525.0","currency_code":"USD"}],"whiteboard":null,"group_type":"apartment","invite_link":"https://www.splitwise.com/join/DoD65gHhn9w+cvgzh","group_reminders":null,"avatar":{"small":"https://splitwise.s3.amazonaws.com/uploads/group/avatar/11741221/small_2a568696-c654-4fca-b761-a5554d9c2c2e.jpeg","medium":"https://splitwise.s3.amazonaws.com/uploads/group/avatar/11741221/medium_2a568696-c654-4fca-b761-a5554d9c2c2e.jpeg","large":"https://splitwise.s3.amazonaws.com/uploads/group/avatar/11741221/large_2a568696-c654-4fca-b761-a5554d9c2c2e.jpeg","xlarge":"https://splitwise.s3.amazonaws.com/uploads/group/avatar/11741221/xlarge_2a568696-c654-4fca-b761-a5554d9c2c2e.jpeg","xxlarge":"https://splitwise.s3.amazonaws.com/uploads/group/avatar/11741221/xxlarge_2a568696-c654-4fca-b761-a5554d9c2c2e.jpeg","original":"https://splitwise.s3.amazonaws.com/uploads/group/avatar/11741221/2a568696-c654-4fca-b761-a5554d9c2c2e.jpeg"},"tall_avatar":{"xlarge":"https://splitwise.s3.amazonaws.com/uploads/group/avatar/11741221/xlarge_2a568696-c654-4fca-b761-a5554d9c2c2e.jpeg","large":"https://splitwise.s3.amazonaws.com/uploads/group/avatar/11741221/large_2a568696-c654-4fca-b761-a5554d9c2c2e.jpeg"},"custom_avatar":true,"cover_photo":{"xxlarge":"https://splitwise.s3.amazonaws.com/uploads/group/avatar/11741221/xxlarge_2a568696-c654-4fca-b761-a5554d9c2c2e.jpeg","xlarge":"https://splitwise.s3.amazonaws.com/uploads/group/avatar/11741221/xlarge_2a568696-c654-4fca-b761-a5554d9c2c2e.jpeg"}},{"id":11794860,"name":"Construcción","created_at":"2019-03-04T15:04:52Z","updated_at":"2019-08-02T13:38:25Z","members":[{"id":21623741,"first_name":"test1","last_name":"test","picture":{"small":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/small_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg","medium":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/medium_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg","large":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/large_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg"},"custom_picture":true,"email":"test@test.com","registration_status":"confirmed","balance":[{"currency_code":"ARS","amount":"49575.25"}]},{"id":21679690,"first_name":"test2","last_name":"test","picture":{"small":"https://s3.amazonaws.com/splitwise/uploads/user/default_avatars/avatar-orange1-50px.png","medium":"https://s3.amazonaws.com/splitwise/uploads/user/default_avatars/avatar-orange1-100px.png","large":"https://s3.amazonaws.com/splitwise/uploads/user/default_avatars/avatar-orange1-200px.png"},"custom_picture":false,"email":"test2@test.com","registration_status":"confirmed","balance":[{"currency_code":"ARS","amount":"-49575.25"}]}],"simplify_by_default":false,"original_debts":[{"to":21623741,"from":21679690,"amount":"49575.25","currency_code":"ARS"}],"simplified_debts":[{"from":21679690,"to":21623741,"amount":"49575.25","currency_code":"ARS"}],"whiteboard":null,"group_type":"apartment","invite_link":"https://www.splitwise.com/join/jT1o4B9MQux+cvgzh","group_reminders":null,"avatar":{"small":"https://s3.amazonaws.com/splitwise/uploads/group/default_avatars/avatar-teal23-house-50px.png","medium":"https://s3.amazonaws.com/splitwise/uploads/group/default_avatars/avatar-teal23-house-100px.png","large":"https://s3.amazonaws.com/splitwise/uploads/group/default_avatars/avatar-teal23-house-200px.png","xlarge":"https://s3.amazonaws.com/splitwise/uploads/group/default_avatars/avatar-teal23-house-500px.png","xxlarge":"https://s3.amazonaws.com/splitwise/uploads/group/default_avatars/avatar-teal23-house-1000px.png","original":null},"tall_avatar":{"xlarge":"https://s3.amazonaws.com/splitwise/uploads/group/default_tall_avatars/avatar-teal23-house-288px.png","large":"https://s3.amazonaws.com/splitwise/uploads/group/default_tall_avatars/avatar-teal23-house-192px.png"},"custom_avatar":false,"cover_photo":{"xxlarge":"https://s3.amazonaws.com/splitwise/uploads/group/default_cover_photos/coverphoto-teal-1000px.png","xlarge":"https://s3.amazonaws.com/splitwise/uploads/group/default_cover_photos/coverphoto-teal-500px.png"}},{"id":12683611,"name":"Personal","created_at":"2019-04-29T21:20:51Z","updated_at":"2019-07-02T16:03:15Z","members":[{"id":21623741,"first_name":"test1","last_name":"test","picture":{"small":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/small_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg","medium":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/medium_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg","large":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/large_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg"},"custom_picture":true,"email":"test@test.com","registration_status":"confirmed","balance":[]}],"simplify_by_default":false,"original_debts":[],"simplified_debts":[],"whiteboard":null,"group_type":"apartment","invite_link":"https://www.splitwise.com/join/BbCn6pxf1f4+cvgzh","group_reminders":null,"avatar":{"small":"https://s3.amazonaws.com/splitwise/uploads/group/default_avatars/avatar-blue4-house-50px.png","medium":"https://s3.amazonaws.com/splitwise/uploads/group/default_avatars/avatar-blue4-house-100px.png","large":"https://s3.amazonaws.com/splitwise/uploads/group/default_avatars/avatar-blue4-house-200px.png","xlarge":"https://s3.amazonaws.com/splitwise/uploads/group/default_avatars/avatar-blue4-house-500px.png","xxlarge":"https://s3.amazonaws.com/splitwise/uploads/group/default_avatars/avatar-blue4-house-1000px.png","original":null},"tall_avatar":{"xlarge":"https://s3.amazonaws.com/splitwise/uploads/group/default_tall_avatars/avatar-blue4-house-288px.png","large":"https://s3.amazonaws.com/splitwise/uploads/group/default_tall_avatars/avatar-blue4-house-192px.png"},"custom_avatar":false,"cover_photo":{"xxlarge":"https://s3.amazonaws.com/splitwise/uploads/group/default_cover_photos/coverphoto-blue-1000px.png","xlarge":"https://s3.amazonaws.com/splitwise/uploads/group/default_cover_photos/coverphoto-blue-500px.png"}},{"id":13548002,"name":"Bautismo","created_at":"2019-06-17T20:09:48Z","updated_at":"2019-06-26T18:34:46Z","members":[{"id":21679690,"first_name":"test2","last_name":"test","picture":{"small":"https://s3.amazonaws.com/splitwise/uploads/user/default_avatars/avatar-orange1-50px.png","medium":"https://s3.amazonaws.com/splitwise/uploads/user/default_avatars/avatar-orange1-100px.png","large":"https://s3.amazonaws.com/splitwise/uploads/user/default_avatars/avatar-orange1-200px.png"},"custom_picture":false,"email":"test2@test.com","registration_status":"confirmed","balance":[]},{"id":21623741,"first_name":"test1","last_name":"test","picture":{"small":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/small_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg","medium":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/medium_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg","large":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/large_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg"},"custom_picture":true,"email":"test@test.com","registration_status":"confirmed","balance":[]}],"simplify_by_default":false,"original_debts":[],"simplified_debts":[],"whiteboard":null,"group_type":"apartment","invite_link":"https://www.splitwise.com/join/PkT9FtyWTbL+cvgzh","group_reminders":null,"avatar":{"small":"https://s3.amazonaws.com/splitwise/uploads/group/default_avatars/avatar-orange18-house-50px.png","medium":"https://s3.amazonaws.com/splitwise/uploads/group/default_avatars/avatar-orange18-house-100px.png","large":"https://s3.amazonaws.com/splitwise/uploads/group/default_avatars/avatar-orange18-house-200px.png","xlarge":"https://s3.amazonaws.com/splitwise/uploads/group/default_avatars/avatar-orange18-house-500px.png","xxlarge":"https://s3.amazonaws.com/splitwise/uploads/group/default_avatars/avatar-orange18-house-1000px.png","original":null},"tall_avatar":{"xlarge":"https://s3.amazonaws.com/splitwise/uploads/group/default_tall_avatars/avatar-orange18-house-288px.png","large":"https://s3.amazonaws.com/splitwise/uploads/group/default_tall_avatars/avatar-orange18-house-192px.png"},"custom_avatar":false,"cover_photo":{"xxlarge":"https://s3.amazonaws.com/splitwise/uploads/group/default_cover_photos/coverphoto-orange-1000px.png","xlarge":"https://s3.amazonaws.com/splitwise/uploads/group/default_cover_photos/coverphoto-orange-500px.png"}},{"id":19457330,"name":"Fundación ","created_at":"2020-06-21T17:40:26Z","updated_at":"2023-04-11T12:15:53Z","members":[{"id":33433366,"first_name":"test3","last_name":"test","picture":{"small":"https://s3.amazonaws.com/splitwise/uploads/user/default_avatars/avatar-orange11-50px.png","medium":"https://s3.amazonaws.com/splitwise/uploads/user/default_avatars/avatar-orange11-100px.png","large":"https://s3.amazonaws.com/splitwise/uploads/user/default_avatars/avatar-orange11-200px.png"},"custom_picture":false,"email":"test3@test.com","registration_status":"confirmed","balance":[{"currency_code":"ARS","amount":"-194.06"}]},{"id":21623741,"first_name":"test1","last_name":"test","picture":{"small":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/small_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg","medium":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/medium_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg","large":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/large_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg"},"custom_picture":true,"email":"test@test.com","registration_status":"confirmed","balance":[{"currency_code":"ARS","amount":"22918.26"}]},{"id":21702157,"first_name":"Federico","last_name":"test","picture":{"small":"https://s3.amazonaws.com/splitwise/uploads/user/default_avatars/avatar-ruby2-50px.png","medium":"https://s3.amazonaws.com/splitwise/uploads/user/default_avatars/avatar-ruby2-100px.png","large":"https://s3.amazonaws.com/splitwise/uploads/user/default_avatars/avatar-ruby2-200px.png"},"custom_picture":false,"email":"test4@test.com","registration_status":"confirmed","balance":[{"currency_code":"ARS","amount":"-22724.2"}]}],"simplify_by_default":false,"original_debts":[{"from":33433366,"to":21623741,"amount":"16978.44","currency_code":"ARS"},{"from":21702157,"to":21623741,"amount":"5939.82","currency_code":"ARS"},{"currency_code":"ARS","from":21702157,"to":33433366,"amount":"16784.38"}],"simplified_debts":[{"from":21702157,"to":21623741,"amount":"22724.2","currency_code":"ARS"},{"from":33433366,"to":21623741,"amount":"194.06","currency_code":"ARS"}],"whiteboard":null,"group_type":"apartment","invite_link":"https://www.splitwise.com/join/Q41YUqXHzcq+cvgzh","group_reminders":null,"avatar":{"small":"https://splitwise.s3.amazonaws.com/uploads/group/avatar/19457330/small_d5fade38-f5cc-4afa-a6db-5db157f857f2.jpeg","medium":"https://splitwise.s3.amazonaws.com/uploads/group/avatar/19457330/medium_d5fade38-f5cc-4afa-a6db-5db157f857f2.jpeg","large":"https://splitwise.s3.amazonaws.com/uploads/group/avatar/19457330/large_d5fade38-f5cc-4afa-a6db-5db157f857f2.jpeg","xlarge":"https://splitwise.s3.amazonaws.com/uploads/group/avatar/19457330/xlarge_d5fade38-f5cc-4afa-a6db-5db157f857f2.jpeg","xxlarge":"https://splitwise.s3.amazonaws.com/uploads/group/avatar/19457330/xxlarge_d5fade38-f5cc-4afa-a6db-5db157f857f2.jpeg","original":"https://splitwise.s3.amazonaws.com/uploads/group/avatar/19457330/d5fade38-f5cc-4afa-a6db-5db157f857f2.jpeg"},"tall_avatar":{"xlarge":"https://splitwise.s3.amazonaws.com/uploads/group/avatar/19457330/xlarge_d5fade38-f5cc-4afa-a6db-5db157f857f2.jpeg","large":"https://splitwise.s3.amazonaws.com/uploads/group/avatar/19457330/large_d5fade38-f5cc-4afa-a6db-5db157f857f2.jpeg"},"custom_avatar":true,"cover_photo":{"xxlarge":"https://splitwise.s3.amazonaws.com/uploads/group/avatar/19457330/xxlarge_d5fade38-f5cc-4afa-a6db-5db157f857f2.jpeg","xlarge":"https://splitwise.s3.amazonaws.com/uploads/group/avatar/19457330/xlarge_d5fade38-f5cc-4afa-a6db-5db157f857f2.jpeg"}},{"id":29044250,"name":"Viaggio In Italia","created_at":"2021-12-18T08:53:13Z","updated_at":"2022-09-01T15:06:00Z","members":[{"id":21623741,"first_name":"test1","last_name":"test","picture":{"small":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/small_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg","medium":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/medium_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg","large":"https://splitwise.s3.amazonaws.com/uploads/user/avatar/21623741/large_5b87c7d4-d6f5-4ef1-9a7f-0b8fbd0b806d.jpeg"},"custom_picture":true,"email":"test@test.com","registration_status":"confirmed","balance":[{"currency_code":"EUR","amount":"2514.68"}]},{"id":21679690,"first_name":"test2","last_name":"test","picture":{"small":"https://s3.amazonaws.com/splitwise/uploads/user/default_avatars/avatar-orange1-50px.png","medium":"https://s3.amazonaws.com/splitwise/uploads/user/default_avatars/avatar-orange1-100px.png","large":"https://s3.amazonaws.com/splitwise/uploads/user/default_avatars/avatar-orange1-200px.png"},"custom_picture":false,"email":"test2@test.com","registration_status":"confirmed","balance":[{"currency_code":"EUR","amount":"-2514.68"}]}],"simplify_by_default":false,"original_debts":[{"to":21623741,"from":21679690,"amount":"2514.68","currency_code":"EUR"}],"simplified_debts":[{"from":21679690,"to":21623741,"amount":"2514.68","currency_code":"EUR"}],"whiteboard":null,"group_type":"trip","invite_link":"https://www.splitwise.com/join/LxaFazBPRQ8+cvgzh","group_reminders":null,"avatar":{"small":"https://s3.amazonaws.com/splitwise/uploads/group/default_avatars/v2021/avatar-orange37-trip-50px.png","medium":"https://s3.amazonaws.com/splitwise/uploads/group/default_avatars/v2021/avatar-orange37-trip-100px.png","large":"https://s3.amazonaws.com/splitwise/uploads/group/default_avatars/v2021/avatar-orange37-trip-200px.png","xlarge":"https://s3.amazonaws.com/splitwise/uploads/group/default_avatars/v2021/avatar-orange37-trip-500px.png","xxlarge":"https://s3.amazonaws.com/splitwise/uploads/group/default_avatars/v2021/avatar-orange37-trip-1000px.png","original":null},"tall_avatar":{"xlarge":"https://s3.amazonaws.com/splitwise/uploads/group/default_tall_avatars/avatar-orange37-trip-288px.png","large":"https://s3.amazonaws.com/splitwise/uploads/group/default_tall_avatars/avatar-orange37-trip-192px.png"},"custom_avatar":false,"cover_photo":{"xxlarge":"https://s3.amazonaws.com/splitwise/uploads/group/default_cover_photos/coverphoto-orange-1000px.png","xlarge":"https://s3.amazonaws.com/splitwise/uploads/group/default_cover_photos/coverphoto-orange-500px.png"}}]}`

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

	conn := getClientMockedConnection(t, doFunc)

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

	conn := getClientMockedConnection(t, doFunc)

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
	conn := getClientMockedConnection(t, doFunc)

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

	conn := getClientMockedConnection(t, doFunc)
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

	conn := getClientMockedConnection(t, doFunc)
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

	conn := getClientMockedConnection(t, doFunc)
	group, err := conn.GetGroup(12345)

	assert.NoError(t, err)

	assert.Equal(t, wantedRespounce.Group, group)

}

func getIntFromParams(url *url.URL, param string) int {
	urlValues := url.Query()

	tmps, ok := urlValues[param]

	if !ok {
		return 0
	} else if len(tmps) > 0 {
		rtr, err := strconv.Atoi(tmps[0])
		if err != nil {
			panic("unable to parse limit")
		}
		return rtr
	}
	return 0
}

func TestGetExpences(t *testing.T) {
	type responseStruct struct {
		Expenses []resources.Expense
	}
	wantedRespounce := responseStruct{}
	wantedRespounce.Expenses = make([]resources.Expense, 10)
	err := json.Unmarshal([]byte(testExpences), &wantedRespounce)

	assert.NoError(t, err, "error Unmarshaling struct espected nil but received %s", err)

	doFunc := func(r *http.Request) (*http.Response, error) {
		resposne := http.Response{}

		limit := getIntFromParams(r.URL, "limit")
		offset := getIntFromParams(r.URL, "offset")

		if len(wantedRespounce.Expenses) > offset {
			patial := responseStruct{}
			patial.Expenses = wantedRespounce.Expenses[offset : offset+limit]
			content, err := json.Marshal(patial)
			if err != nil {
				panic("error converting parial strunt")
			}
			resposne.Body = io.NopCloser(bytes.NewReader(content))
		} else {
			resposne.Body = io.NopCloser(bytes.NewReader([]byte{}))
		}

		resposne.Header = make(map[string][]string)
		resposne.Header["Content-Type"] = []string{"application/json", "charset=utf-8"}
		resposne.Status = "200"
		resposne.StatusCode = 200
		return &resposne, nil
	}

	assert.NoError(t, err, "error Unmarshaling struct espected nil but received %s", err)

	conn := getClientMockedConnection(t, doFunc)

	params := splitwise.ExpensesParams{}
	params[splitwise.ExpensesDatedAfter] = time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.Local)
	params[splitwise.ExpensesGroupId] = "12345"
	params[splitwise.ExpensesLimit] = 5

	executor := conn.GetExpenses(params)

	count := 0

	for _ = range executor.GetChan() {
		count++
	}
	assert.Equal(t, 10, count)
	assert.True(t, executor.isClose())
}

func TestGetExpencesWithPanic(t *testing.T) {
	type responseStruct struct {
		Expenses []resources.Expense
	}
	wantedRespounce := responseStruct{}
	wantedRespounce.Expenses = make([]resources.Expense, 10)
	err := json.Unmarshal([]byte(testExpences), &wantedRespounce)

	assert.NoError(t, err, "error Unmarshaling struct espected nil but received %s", err)

	doFunc := func(r *http.Request) (*http.Response, error) {
		panic("painic produced in test")
	}

	assert.NoError(t, err, "error Unmarshaling struct espected nil but received %s", err)

	conn := getClientMockedConnection(t, doFunc)

	params := splitwise.ExpensesParams{}
	params[splitwise.ExpensesDatedAfter] = time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.Local)
	params[splitwise.ExpensesGroupId] = "12345"
	params[splitwise.ExpensesLimit] = 5

	executor := conn.GetExpenses(params)

	count := 0

	for _ = range executor.GetChan() {
		count++
	}
	assert.Equal(t, 0, count)
	assert.True(t, executor.isClose())
}

func TestGetgroups(t *testing.T) {
	doFunc := func(r *http.Request) (*http.Response, error) {
		resposne := http.Response{}
		resposne.Body = io.NopCloser(strings.NewReader(testGroups))
		resposne.Header = make(map[string][]string)
		resposne.Header["Content-Type"] = []string{"application/json", "charset=utf-8"}
		resposne.Status = "200"
		resposne.StatusCode = 200
		return &resposne, nil
	}
	type responseStruct struct {
		Groups []resources.Group
	}

	wantedRespounce := responseStruct{}
	wantedRespounce.Groups = make([]resources.Group, 10)
	err := json.Unmarshal([]byte(testGroups), &wantedRespounce)
	assert.NoError(t, err, "error Unmarshaling struct espected nil but received %s", err)
	conn := getClientMockedConnection(t, doFunc)
	executor := conn.GetGroups()

	count := 0

	for _ = range executor.GetChan() {
		count++
	}
	assert.Equal(t, 7, count)
}

func TestGetgroupsWithClose(t *testing.T) {
	doFunc := func(r *http.Request) (*http.Response, error) {
		resposne := http.Response{}
		resposne.Body = io.NopCloser(strings.NewReader(testGroups))
		resposne.Header = make(map[string][]string)
		resposne.Header["Content-Type"] = []string{"application/json", "charset=utf-8"}
		resposne.Status = "200"
		resposne.StatusCode = 200
		return &resposne, nil
	}
	type responseStruct struct {
		Groups []resources.Group
	}

	wantedRespounce := responseStruct{}
	wantedRespounce.Groups = make([]resources.Group, 10)
	err := json.Unmarshal([]byte(testGroups), &wantedRespounce)
	assert.NoError(t, err, "error Unmarshaling struct espected nil but received %s", err)
	conn := getClientMockedConnection(t, doFunc)
	executor := conn.GetGroups()

	count := 0

	for _ = range executor.GetChan() {
		count++
		if count == 1 {
			executor.Close()
			assert.True(t, executor.isClose())
		}
	}
	assert.Equal(t, 1, count)
}

func getClientMockedConnection(t *testing.T, doFunc func(r *http.Request) (*http.Response, error)) SwConnection {
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
