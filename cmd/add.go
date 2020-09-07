package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/GGP1/kure/crypt"
	"github.com/GGP1/kure/db"
	"github.com/GGP1/kure/entry"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	custom bool
	length uint16
	format []uint
	secure bool

	addCmd = &cobra.Command{
		Use:   "add [-c custom] [-l length] [-f format] [-S secure]",
		Short: "Adds a new entry to the database",
		Run: func(cmd *cobra.Command, args []string) {
			var (
				password string
				entropy  float64
				err      error
			)

			// Take entry input from the user
			title, username, password, url, notes, expiration := entryInput()

			if !custom {
				levels := make(map[uint]struct{})

				for _, v := range format {
					levels[v] = struct{}{}
				}
				password, entropy, err = entry.GeneratePassword(length, levels)
				if err != nil {
					fmt.Println("error:", err)
					return
				}
			}

			if secure {
				pwd, err := passInput()
				if err != nil {
					fmt.Println("error:", err)
					return
				}

				encryptedPwd, err := crypt.Encrypt([]byte(password), pwd)
				if err != nil {
					fmt.Println("error:", err)
					return
				}
				password = string(encryptedPwd)
			}

			title, expiration, err = formatFields(title, expiration)
			if err != nil {
				fmt.Println("error:", err)
				return
			}

			entry := entry.New(title, username, password, url, notes, expiration, secure)

			if err := db.CreateEntry(entry); err != nil {
				fmt.Println("error:", err)
				return
			}

			fmt.Printf("\nSucessfully created the entry.\nBits of entropy: %.2f", entropy)
		},
	}
)

func init() {
	RootCmd.AddCommand(addCmd)
	addCmd.Flags().BoolVarP(&custom, "custom", "c", false, "custom password")
	addCmd.Flags().Uint16VarP(&length, "length", "l", 1, "password length")
	addCmd.Flags().UintSliceVarP(&format, "format", "f", []uint{1, 2, 3, 4}, "password format")
	addCmd.Flags().BoolVarP(&secure, "secure", "S", false, "security mode")

	addCmd.MarkFlagRequired("title")
	if !custom {
		addCmd.MarkFlagRequired("length")
	}
}

func entryInput() (title, username, password, url, notes, expiration string) {
	scanner := bufio.NewScanner(os.Stdin)

	title = scan(scanner, "Title", title)
	username = scan(scanner, "Username", username)
	url = scan(scanner, "URL", url)
	notes = scan(scanner, "Notes", notes)
	expiration = scan(scanner, "Expiration", expiration)

	if custom {
		password = scan(scanner, "Password", password)
	}

	return title, username, password, url, notes, expiration
}

func scan(scanner *bufio.Scanner, field string, value string) string {
	fmt.Printf("%s: ", field)
	scanner.Scan()
	value = scanner.Text()

	return value
}

func formatFields(title, expiration string) (string, string, error) {
	t := strings.ToLower(title)

	if expiration == "0s" || expiration == "0" || expiration == "" {
		expiration = "Never"
	} else {
		expTime, err := time.ParseDuration(expiration)
		if err != nil {
			return "", "", errors.Wrap(err, "duration parse")
		}
		// Add duration and format
		expiration = time.Now().Add(expTime).Format(time.RFC3339)
	}

	return t, expiration, nil
}
