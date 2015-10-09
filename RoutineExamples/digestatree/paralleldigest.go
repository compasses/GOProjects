//refer to: http://blog.golang.org/pipelines

package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

type result struct {
	path string
	sum  [md5.Size]byte
	err  error
}

func sumFile(done <-chan struct{}, root string) (<-chan result, <-chan error) {
	//for each of regular file, start a goroutine that sums the file and sends the
	//result on c, send the result of the walk on err
	c := make(chan result)
	errc := make(chan error, 1)

	go func() {
		var wg sync.WaitGroup

		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.Mode().IsRegular() {
				return nil
			}

			wg.Add(1)
			go func() {
				data, err := ioutil.ReadFile(path)
				select {
				case c <- result{path, md5.Sum(data), err}:
				case <-done:
				}
				wg.Done()
			}()
			//abort the walk if done is closed
			select {
			case <-done:
				return errors.New("walk canceled")
			default:
				return nil
			}
		})
		//walk has returned, so all calls to wg.Add are done. Start a goroutine to close
		//c once all the sends are done
		go func() {
			wg.Wait()
			close(c)
		}()
		errc <- err
	}()

	return c, errc
}

func MD5ALL(root string) (map[string][md5.Size]byte, error) {

	//MD5All closes the done channel when it returns. it may do so before receiving all the values from c and errc
	done := make(chan struct{})
	defer close(done)

	c, errc := sumFile(done, root)

	m := make(map[string][md5.Size]byte)

	for r := range c {
		if r.err != nil {
			return nil, r.err
		}
		m[r.path] = r.sum
	}

	if err := <-errc; err != nil {
		return nil, err
	}

	return m, nil
}

func main() {
	start := time.Now()

	m, err := MD5ALL(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	elapsed := time.Since(start)
	fmt.Println("time elapsed: ", elapsed)

	var paths []string
	for path := range m {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		fmt.Printf("%x, %s\n", m[path], path)
	}
}
