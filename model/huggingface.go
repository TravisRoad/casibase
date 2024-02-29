// Copyright 2023 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/casibase/casibase/proxy"
	"github.com/henomis/lingoose/llm/huggingface"
)

type HuggingFaceModelProvider struct {
	subType     string
	secretKey   string
	temperature float32
}

func NewHuggingFaceModelProvider(subType string, secretKey string, temperature float32) (*HuggingFaceModelProvider, error) {
	return &HuggingFaceModelProvider{subType: subType, secretKey: secretKey, temperature: temperature}, nil
}

func (p *HuggingFaceModelProvider) GetPricing() (string, string) {
	return "Free", "Free"
}

func (p *HuggingFaceModelProvider) caculatePrice(mr *ModelResult) {
	mr.TotalPrice = 0.0
}

func (p *HuggingFaceModelProvider) QueryText(question string, writer io.Writer, history []*RawMessage, prompt string, knowledgeMessages []*RawMessage) (*ModelResult, error) {
	client := huggingface.New(p.subType, p.temperature, false).WithToken(p.secretKey).WithHTTPClient(proxy.ProxyHttpClient).WithMode(huggingface.ModeTextGeneration)

	ctx := context.Background()

	resp, err := client.Completion(ctx, question)
	if err != nil {
		return nil, err
	}

	resp = strings.Split(resp, "\n")[0]
	fmt.Println(resp)

	_, writeErr := fmt.Fprint(writer, resp)
	if writeErr != nil {
		return nil, writeErr
	}

	// get token count and price
	mr := new(ModelResult)
	promptTokenCount, err := GetTokenSize(p.subType, question)
	if err != nil {
		return nil, err
	}
	mr.PromptTokenCount = promptTokenCount
	responseTokenCount, err := GetTokenSize(p.subType, resp)
	if err != nil {
		return nil, err
	}
	mr.ResponseTokenCount = responseTokenCount
	mr.TotalTokenCount = promptTokenCount + responseTokenCount
	p.caculatePrice(mr)

	return mr, nil
}
