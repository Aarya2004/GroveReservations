package sb

import (
	supabase "github.com/supabase-community/supabase-go"
)

func MustNewSupabaseClient(url string, key string) *supabase.Client {
	client, err := supabase.NewClient(url, key, &supabase.ClientOptions{})
	if err != nil {
		panic(err)
	}
	return client
}
