package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func encoding(dist, summary string) error {
	if _, err := os.Stat(dist); os.IsNotExist(err) {
		return err
	}
	var GLOBAL_COUNTER int
	var GLOBAL_CONTENT string
	if err := filepath.Walk(dist, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			fmt.Println("Find >>>", path)
			if file, err := os.ReadFile(path); err != nil {
				return err
			} else {
				var content string
				lines := strings.Split(string(file), "\n")
				for _, line := range lines {
					re := regexp.MustCompile("[\u4e00-\u9faf\u3040-\u309f\u30a0-\u30ff]+")
					matches := re.FindAllString(line, -1)
					if len(matches) > 0 {
						GLOBAL_COUNTER++
						GLOBAL_CONTENT += fmt.Sprintf("%s\n", line)
						content += fmt.Sprintf("[NLJB]-%d\n", GLOBAL_COUNTER)
					} else {
						content += line + "\n"
					}
				}
				if err := os.WriteFile(path, []byte(content), 0644); err != nil {
					return err
				}
			}
		}
		return nil
	}); err != nil {
		return err
	}
	if err := os.WriteFile(summary, []byte(GLOBAL_CONTENT), 0644); err != nil {
		return err
	}
	return nil
}

func decoding(summary, dist string) error {
	// 提取结果 ...
	file, err := os.ReadFile(summary)
	if err != nil {
		return err
	}
	var GLOBAL_COUNTER int
	GLOBAL_LINE := make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(string(file)))
	for scanner.Scan() {
		GLOBAL_COUNTER++
		line := scanner.Text()
		GLOBAL_LINE[fmt.Sprintf("%d", GLOBAL_COUNTER)] = line
	}
	// 回写结果 ...
	if _, err := os.Stat(dist); os.IsNotExist(err) {
		return err
	}
	if err := filepath.Walk(dist, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			fmt.Println("Find >>>", path)
			if file, err := os.ReadFile(path); err != nil {
				return err
			} else {
				var content string
				lines := strings.Split(string(file), "\n")
				for _, line := range lines {
					if strings.HasPrefix(line, "[NLJB]") {
						parts := strings.SplitN(line, "-", 2)
						content += GLOBAL_LINE[parts[1]] + "\n"
					} else {
						content += line + "\n"
					}
				}
				if err := os.WriteFile(path, []byte(content), 0644); err != nil {
					return err
				}
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func main() {
	// 使用 Google 翻译 https://translate.google.com/
	dist := "E:\\katydist-jp-bak-xiugai\\dist"
	summary := "E:\\katydist-jp-bak-xiugai\\summary"
	//if err := encoding(dist, summary); err != nil {
	//	log.Fatal(err)
	//	return
	//}
	if err := decoding(summary, dist); err != nil {
		log.Fatal(err)
		return
	}
}
