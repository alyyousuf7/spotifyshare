// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Connect to a spotifyshare server",
	Run: func(cmd *cobra.Command, args []string) {
		d := websocket.Dialer{}
		conn, _, err := d.Dial("ws://localhost:8080/ws", nil)
		if err != nil {
			log.Error(err)
			return
		}

		go func() {
			for {
				t, msg, err := conn.ReadMessage()
				if err != nil {
					log.Error(err)
					break
				}
				log.Infof("[%d] %s", t, msg)
			}
		}()

		conn.SetCloseHandler(func(code int, text string) error {
			log.Debugf("Closed [%s]: %s", code, text)
			return nil
		})

		quitCh := make(chan os.Signal)
		signal.Notify(quitCh, os.Interrupt)
		select {
		case <-quitCh:
			fmt.Println("Shutting down")
			if err := conn.Close(); err != nil {
				log.Error(err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(clientCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clientCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clientCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
