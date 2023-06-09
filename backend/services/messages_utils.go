package services

import (
	"context"
	"errors"

	"github.com/zanz1n/ws-messaging-app/dba"
)

func ParseMsgToSendableData(db *dba.Queries, query []dba.Message) (*[]MessageCreateReturnedData, error) {
	userMap := make(map[string]string)

	finalData := make([]MessageCreateReturnedData, len(query))

	var (
		ctx        = context.Background()
		ok         bool
		userResult dba.User
		userName   string
		i          int
		iO         = 0
		message    dba.Message
		data       MessageCreateReturnedData
		err        error
	)

	for i = len(query) - 1; i >= 0; i-- {
		message = query[i]
		if _, ok = userMap[message.UserId]; !ok {
			userResult, err = db.GetUserById(ctx, message.UserId)

			if err != nil {
				return nil, errors.New("error while fetching user data")
			}

			userMap[message.UserId] = userResult.Username
		}
		userName = userMap[message.UserId]

		data = MessageCreateReturnedData{
			ID:        message.ID,
			CreatedAt: message.CreatedAt,
			UpdatedAt: message.CreatedAt,
			User: UserReturnedOnMessage{
				ID:       message.UserId,
				Username: userName,
			},
		}

		if message.Content.Valid {
			content := message.Content.String
			data.Content = &content
		}

		if message.ImageUrl.Valid {
			imageUrl := message.ImageUrl.String
			data.Image = &imageUrl
		}

		finalData[iO] = data
		iO++
	}

	return &finalData, nil
}
