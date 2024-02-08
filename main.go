package main

import (
	"context"
	"fmt"
	"net/http"
    "bytes"
    "regexp"

	"github.com/gorilla/websocket"

    "github.com/bluesky-social/indigo/api/bsky"
	comatproto "github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/events"
	"github.com/bluesky-social/indigo/events/schedulers/sequential"
    "github.com/bluesky-social/indigo/repo"
    "github.com/bluesky-social/indigo/repomgr"
)

func SolveNyan(text string) bool {
    reg := regexp.MustCompile(`にゃん|🐈|nyan|にゃあ|ねこ`)
    return reg.MatchString(text)
}

func main() {
	fmt.Printf("Run...\n")

	url := "wss://bsky.network/xrpc/com.atproto.sync.subscribeRepos"
	con, _, err := websocket.DefaultDialer.Dial(url, http.Header{})

	if err != nil {
		fmt.Printf("Connection error: %v\n", err)
		return
	}
	defer con.Close()


	rsc := &events.RepoStreamCallbacks{
		RepoCommit: func(evt *comatproto.SyncSubscribeRepos_Commit) error {
            ctx := context.Background()
			for _, op := range evt.Ops {
                ek := repomgr.EventKind(op.Action)
                switch ek {
                case repomgr.EvtKindCreateRecord, repomgr.EvtKindUpdateRecord:
                    r, err := repo.ReadRepoFromCar(ctx, bytes.NewReader(evt.Blocks))
                    if err != nil {
                        fmt.Errorf("ReadRepo error: %v", err)
                        return nil
                    }

                    _, rec, err := r.GetRecord(ctx, op.Path)
                    if err != nil {
                        fmt.Errorf("Get record error: %v", err)
                        return nil
                    }

                    if ek != "create" {
						return nil
					}

                    post, ok := rec.(*bsky.FeedPost)
                    if !ok {
						return nil
					}

                    if SolveNyan(post.Text) {
                        fmt.Printf("[Nyanpost]%v %v\n", evt.Repo, post.CreatedAt)
                        fmt.Printf("%v\n\n", post.Text)
                    }

                default:
                }
			}
            return nil
		},
	}

	sched := sequential.NewScheduler("myfirehose", rsc.EventHandler)
	events.HandleRepoStream(context.Background(), con, sched)
}
