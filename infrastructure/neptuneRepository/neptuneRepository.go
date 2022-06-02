package neptuneRepository

import (
	"encoding/json"
	"eng-1004_nonomuy/domain/model"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/schwartzmx/gremtune"
	"log"
	"strings"
	"time"
)

const layout = "2006-01-02 15:04:05.999999999 -0700 MST"

type neptuneRepository struct {
	readClient  *gremtune.Client
	writeClient *gremtune.Client
}

type graphData struct {
	Type  string          `json:"@type"`
	Value []vertexWrapper `json:"@value"`
}

type vertexWrapper struct {
	Type  string `json:"@type"`
	Value vertex `json:"@value"`
}
type vertex struct {
	Id         string                `json:"id"`
	Label      string                `json:"label"`
	Properties map[string][]property `json:"properties"`
}
type property struct {
	Value valueWrapper `json:"@value"`
}
type valueWrapper struct {
	Value string `json:"value"`
}

func NewNeptuneRepository() *neptuneRepository {
	errs := make(chan error)
	go func(chan error) {
		err := <-errs
		log.Fatal("Lost connection to the database: " + err.Error())
	}(errs)

	readDialer := gremtune.NewDialer("wss://{fixme}.neptune.amazonaws.com:8182")
	readClient, err := gremtune.Dial(readDialer, errs)
	if err != nil {
		fmt.Printf("Dial: %s", err)
	}

	writeDialer := gremtune.NewDialer("wss://{fixme}.neptune.amazonaws.com:8182")
	writeClient, err := gremtune.Dial(writeDialer, errs)
	if err != nil {
		fmt.Printf("Dial: %s", err)
	}

	return &neptuneRepository{&readClient, &writeClient}
}

func (n *neptuneRepository) Tweet(tweetFrom string, message string) {

	newTweet := model.Post{
		Id: tweetFrom, Message: message, PostedAt: time.Now()}

	res, err := n.writeClient.Execute(
		fmt.Sprintf("g.addV('post').property('UserId', '%s').property('Message', '%s').property('PostedAt', '%s')",
			newTweet.Id, newTweet.Message, newTweet.PostedAt),
	)
	if err != nil {
		fmt.Printf("Execute Error: %s", err)
		return
	}
	j, err := json.Marshal(res[0].Result.Data)
	if err != nil {
		fmt.Printf("Marshal Error: %s", err)
		return
	}
	fmt.Printf("MarshalResult: %s", j)
}

func (n *neptuneRepository) GetPosts(userId string) []model.Post {
	res, err := n.readClient.Execute(
		fmt.Sprintf("g.V().has('post', 'UserId', '%s')", userId),
	)
	if err != nil {
		fmt.Printf("Execute: %s", err)
		return nil
	}

	var graphData graphData
	err = json.Unmarshal(res[0].Result.Data, &graphData)
	if err != nil {
		fmt.Printf("Unmarshal error : %s", err)
		return nil
	}

	var results []model.Post
	for _, partialRes := range graphData.Value {
		partialPost := partialRes.Value
		fmt.Printf("partialPost : %s", partialPost)

		var postedAt time.Time
		if len(partialPost.Properties["PostedAt"]) > 0 {
			postedAt, err = time.Parse(layout, strings.Split(partialPost.Properties["PostedAt"][0].Value.Value, " m=")[0])
			if err != nil {
				fmt.Printf("parse error : %s", err)
				return nil
			}
		}
		var userId string
		if len(partialPost.Properties["UserId"]) > 0 {
			userId = partialPost.Properties["UserId"][0].Value.Value
		}

		var message string
		if len(partialPost.Properties["Message"]) > 0 {
			message = partialPost.Properties["Message"][0].Value.Value
		}

		post := model.Post{
			Id:       userId,
			Message:  message,
			PostedAt: postedAt,
		}
		results = append(results, post)
	}

	return results
}

func (n *neptuneRepository) PostUser(name string) {
	userId, _ := n.createNewUserId()

	newUser := model.User{
		Id: userId, Name: name}

	res, err := n.writeClient.Execute(
		fmt.Sprintf("g.addV('user').property('UserId', '%s').property('Name', '%s')",
			newUser.Id, newUser.Name),
	)
	if err != nil {
		fmt.Printf("Execute Error: %s", err)
		return
	}
	j, err := json.Marshal(res[0].Result.Data)
	if err != nil {
		fmt.Printf("Marshal Error: %s", err)
		return
	}
	fmt.Printf("MarshalResult: %s", j)
}

func (n *neptuneRepository) createNewUserId() (string, error) {
	for {
		proposal, _ := uuid.NewV4()

		res, err := n.readClient.Execute(
			fmt.Sprintf("g.V().has('user', 'UserId', '%s')", proposal.String()),
		)
		if err != nil {
			fmt.Printf("Execute: %s", err)
			return "", err
		}

		var graphData graphData
		err = json.Unmarshal(res[0].Result.Data, &graphData)
		if err != nil {
			fmt.Printf("Unmarshal error : %s", err)
			return "", err
		}
		if len(graphData.Value) == 0 {
			return proposal.String(), nil
		}
	}

}

func (n *neptuneRepository) GetUsers() []model.User {
	res, err := n.readClient.Execute(
		fmt.Sprintf("g.V().hasLabel('user')"),
	)
	if err != nil {
		fmt.Printf("Execute: %s", err)
		return []model.User{}
	}

	var graphData graphData
	err = json.Unmarshal(res[0].Result.Data, &graphData)
	if err != nil {
		fmt.Printf("Unmarshal error : %s", err)
		return []model.User{}
	}

	return n.convert(graphData)

}

func (n *neptuneRepository) FollowUser(followerId string, followeeId string) {
	res, err := n.writeClient.Execute(
		fmt.Sprintf("g.V().has('user', 'UserId', '%s').as('followee').V().has('user', 'UserId', '%s').addE('follow').to('followee')",
			followeeId, followerId),
	)
	if err != nil {
		fmt.Printf("Execute Error: %s", err)
		return
	}
	j, err := json.Marshal(res[0].Result.Data)
	if err != nil {
		fmt.Printf("Marshal Error: %s", err)
		return
	}
	fmt.Printf("MarshalResult: %s", j)

}

func (n *neptuneRepository) GetFollowees(followerId string) []model.User {
	res, err := n.readClient.Execute(
		fmt.Sprintf("g.V().has('user', 'UserId', '%s').out('follow')", followerId),
	)
	if err != nil {
		fmt.Printf("Execute: %s", err)
		return []model.User{}
	}

	var graphData graphData
	err = json.Unmarshal(res[0].Result.Data, &graphData)
	if err != nil {
		fmt.Printf("Unmarshal error : %s", err)
		return []model.User{}
	}

	return n.convert(graphData)
}

func (n *neptuneRepository) convert(graphData graphData) []model.User {
	var results []model.User
	for _, partialRes := range graphData.Value {
		partialPost := partialRes.Value
		fmt.Printf("partialPost : %s", partialPost)

		var userId string
		if len(partialPost.Properties["UserId"]) > 0 {
			userId = partialPost.Properties["UserId"][0].Value.Value
		}

		var name string
		if len(partialPost.Properties["Name"]) > 0 {
			name = partialPost.Properties["Name"][0].Value.Value
		}

		user := model.User{
			Id:   userId,
			Name: name,
		}
		results = append(results, user)
	}
	return results
}
