package redis_repos_test

import (
	"chat-service/internal/domain"
	"chat-service/internal/infra/memory"
	"chat-service/internal/infra/memory/redis_repos"
	"testing"
)

func TestRedisChatsCache_SetUsersChats(t *testing.T) {
	redisClient := memory.NewRedisClient()

	tests := []struct {
		name string 
		redisClient *memory.RedisClient
		user_uid string
		chats    []domain.Chat
		wantErr  bool
	}{
		{
			name:        "SetUsersChats",
			redisClient: redisClient,
			user_uid:    "test-user",
			chats: []domain.Chat{
				{
					Title: "hello",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := redis_repos.NewRedisChatsCache(tt.redisClient)
			gotErr := r.SetUsersChats(tt.user_uid, tt.chats)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("SetUsersChats() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("SetUsersChats() succeeded unexpectedly")
			}
		})
	}
}

func TestRedisChatsCache_GetUsersChats(t *testing.T) {
	redisClient := memory.NewRedisClient()

	tests := []struct {
		name string
		redisClient *memory.RedisClient
		user_uid string
		limit    int
		offset   int
		want     []domain.Chat
		wantErr  bool
	}{
		{
			name:        "GetUsersChats",
			redisClient: redisClient,
			user_uid:    "test-user",
			limit:       10,
			offset:      0,
			want:        []domain.Chat{{Title: "hello"}},
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := redis_repos.NewRedisChatsCache(tt.redisClient)
			got, gotErr := r.GetUsersChats(tt.user_uid, tt.limit, tt.offset)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetUsersChats() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetUsersChats() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
            if got[0].Title != tt.want[0].Title {
				t.Errorf("GetUsersChats() = %v, want %v", got, tt.want)
			}
		})
	}
}
