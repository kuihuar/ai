package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"testing"

	"golang.org/x/sync/errgroup"
)

var urls = []string{
	"http://www.google.org/",
	"http://www.google.com/",
	"http:///www.somestupidname.com/",
}

func TestWaitGroup(t *testing.T) {
	wg := sync.WaitGroup{}
	results := make([]string, len(urls))

	for index, url := range urls {
		url := url
		index := index
		wg.Add(1)

		go func() {
			defer wg.Done()
			resp, err := http.Get(url)
			if err != nil {
				return
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return
			}
			results[index] = string(body)
		}()
	}
	wg.Wait()
}

func TestErrorGroup(t *testing.T) {
	results := make([]string, len(urls))

	errg := new(errgroup.Group)

	for index, url := range urls {
		url := url
		index := index
		errg.Go(func() error {
			resp, err := http.Get(url)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			results[index] = string(body)
			return nil
		})
	}
	if err := errg.Wait(); err != nil {
		fmt.Println("Failured fetch all URLs")
	}
}

func TestErrorGroupWithCtx(t *testing.T) {
	results := make([]string, len(urls))

	errg, ctx := errgroup.WithContext(context.Background())
	errg.SetLimit(2)
	for index, url := range urls {
		index := index
		url := url

		errg.Go(func() error {
			select {
			case <-ctx.Done():
				return errors.New("task is canceled")
			default:
				resp, err := http.Get(url)
				if err != nil {
					return err
				}
				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return err
				}

				results[index] = string(body)
				return nil
			}
		})
	}

	if err := errg.Wait(); err != nil {
		fmt.Println("Failured fetch all URLs")
	}
}
