# NotifAPI

[![Deploy to Production](https://github.com/coma-toast/NotifAPI/actions/workflows/deploy.yml/badge.svg)](https://github.com/coma-toast/NotifAPI/actions/workflows/deploy.yml)  
A simple API to receive data and send push notifications to multiple devices via [Pusher](https://pusher.com)

## Features

-   Send messages from all your various systems and services to all of your devices via API call
-   Logging of all alerts to a SQLite DB

## Installation

### Server

-   Sign up for a [Pusher](https://pusher.com) account
-   Clone this repo
-   `go build`
-   Run manually
    -   `./notifapi`
-   Install as a service
    -   Copy `notifapi.sh` somewhere (home folder?)
    -   Modify and copy `notifapi.service` to `/etc/systemd/system/`
    -   `systemctl enable notifapi.service` and `systemctl start notifapi.service`
-   On first run, a blank config file will be generated
-   Modify `config.yaml` and re-run
-   Navigate to 127.0.0.1:<port> to install the Service Worker in Chrome to receive notifications

### Client

Example:

```go
package notifications

import (
	"fmt"

	"github.com/coma-toast/notifapi/pkg/client"
	"github.com/coma-toast/notifapi/pkg/notification"
)

type Notification struct {
	Target string
}

func (n *Notification) SendMessage(title, body string) error {
	client := client.Client{Target: n.Target}
	message := notification.Message{
		Interests: []string{"hello"},
		Title:     title,
		Body:      body,
		Source:    "Source Application Name", // what application are you sending this from
	}

	response, err := client.SendMessage(message)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(response.Status)

	return nil
}
```

## Send a notification

Send a POST request with the following JSON body:

```json
{
    "interests": ["hello"],
    "title": "Test message",
    "body": "This is a test from my testing platform",
    "link": "",
    "source": "My-Laptop",
    "metadata": {
        "server": "macbook-pro"
    }
}
```

## Roadmap

-   Authentication (JWT)
-   Front End to view alert logs (React, or the JS framework du jour)
-   Better documentation
-   Better testing
-   Docker images


