/*
Copyright © 2024 ks6088ts

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package aoai

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/spf13/cobra"
)

// chatcompletionCmd represents the chatcompletion command
var chatcompletionCmd = &cobra.Command{
	Use:   "chatcompletion",
	Short: "A command for Azure OpenAI Service Chat Completion",
	Long:  `ref. https://learn.microsoft.com/azure/ai-services/openai/chatgpt-quickstart?tabs=command-line%2Cpython-new&pivots=programming-language-go`,
	Run: func(cmd *cobra.Command, args []string) {
		// get the values of the flags
		modelDeploymentID, err := cmd.Flags().GetString("modelDeploymentID")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}
		azureOpenAIEndpoint, err := cmd.Flags().GetString("azureOpenAIEndpoint")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}
		azureOpenAIKey, err := cmd.Flags().GetString("azureOpenAIKey")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}
		message, err := cmd.Flags().GetString("message")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}

		keyCredential := azcore.NewKeyCredential(azureOpenAIKey)

		client, err := azopenai.NewClientWithKeyCredential(azureOpenAIEndpoint, keyCredential, nil)

		if err != nil {
			log.Printf("ERROR: %s", err)
			return
		}

		messages := []azopenai.ChatRequestMessageClassification{
			&azopenai.ChatRequestSystemMessage{Content: to.Ptr("You are a helpful assistant.")},
			&azopenai.ChatRequestUserMessage{Content: azopenai.NewChatRequestUserMessageContent(message)},
		}

		gotReply := false

		resp, err := client.GetChatCompletions(context.TODO(), azopenai.ChatCompletionsOptions{
			Messages:       messages,
			DeploymentName: &modelDeploymentID,
		}, nil)

		if err != nil {
			log.Printf("ERROR: %s", err)
			return
		}

		for _, choice := range resp.Choices {
			gotReply = true

			if choice.ContentFilterResults != nil {
				fmt.Fprintf(os.Stderr, "Content filter results\n")

				if choice.ContentFilterResults.Error != nil {
					fmt.Fprintf(os.Stderr, "  Error:%v\n", choice.ContentFilterResults.Error)
				}

				fmt.Fprintf(os.Stderr, "  Hate: sev: %v, filtered: %v\n", *choice.ContentFilterResults.Hate.Severity, *choice.ContentFilterResults.Hate.Filtered)
				fmt.Fprintf(os.Stderr, "  SelfHarm: sev: %v, filtered: %v\n", *choice.ContentFilterResults.SelfHarm.Severity, *choice.ContentFilterResults.SelfHarm.Filtered)
				fmt.Fprintf(os.Stderr, "  Sexual: sev: %v, filtered: %v\n", *choice.ContentFilterResults.Sexual.Severity, *choice.ContentFilterResults.Sexual.Filtered)
				fmt.Fprintf(os.Stderr, "  Violence: sev: %v, filtered: %v\n", *choice.ContentFilterResults.Violence.Severity, *choice.ContentFilterResults.Violence.Filtered)
			}

			if choice.Message != nil && choice.Message.Content != nil {
				fmt.Fprintf(os.Stderr, "Content[%d]: %s\n", *choice.Index, *choice.Message.Content)
			}

			if choice.FinishReason != nil {
				// this choice's conversation is complete.
				fmt.Fprintf(os.Stderr, "Finish reason[%d]: %s\n", *choice.Index, *choice.FinishReason)
			}
		}

		if gotReply {
			fmt.Fprintf(os.Stderr, "Received chat completions reply\n")
		}
	},
}

func init() {
	aoaiCmd.AddCommand(chatcompletionCmd)

	chatcompletionCmd.Flags().StringP("modelDeploymentID", "d", "gpt-4o", "Model Deployment ID")
	chatcompletionCmd.Flags().StringP("azureOpenAIEndpoint", "e", "", "Azure OpenAI Endpoint")
	chatcompletionCmd.Flags().StringP("azureOpenAIKey", "k", "", "Azure OpenAI Key")
	chatcompletionCmd.Flags().StringP("message", "m", "Hello, how are you?", "Message")
}
