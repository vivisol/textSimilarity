package main

import (
	"flag"
	"fmt"
	"github.com/jina-ai/client-go"
	"github.com/jina-ai/client-go/docarray"
	"github.com/jina-ai/client-go/jina"
)

// Create a Document
func getDoc(sentence string) *docarray.DocumentProto {
	return &docarray.DocumentProto{
		Content: &docarray.DocumentProto_Text{
			Text: sentence,
		},
	}
}

// Create a DocumentArray with 3 Documents
func getDocarrays(s1, s2 string) *docarray.DocumentArrayProto {
	var docs []*docarray.DocumentProto
	docs = append(docs, getDoc(s1))
	docs = append(docs, getDoc(s2))

	return &docarray.DocumentArrayProto{
		Docs: docs,
	}
}

// Create DataRequest with a DocumentArray
func getDataRequest(s1, s2 string) *jina.DataRequestProto {
	return &jina.DataRequestProto{
		Data: &jina.DataRequestProto_DataContentProto{
			Documents: &jina.DataRequestProto_DataContentProto_Docs{
				Docs: getDocarrays(s1, s2),
			},
		},
	}
}

// Generate a stream of requests
func generateDataRequests(s1, s2 string) <-chan *jina.DataRequestProto {
	requests := make(chan *jina.DataRequestProto)
	go func() {
		requests <- getDataRequest(s1, s2)
		defer close(requests)
	}()
	return requests
}

func OnDone(resp *jina.DataRequestProto) {

	fmt.Println("resp.Data:", resp.GetData().GetDocs().Docs[0].GetText()) //注意返回的数据是array，所以用了Docs[0]

	fmt.Println("服务调用成功！:")
}

func OnError(resp *jina.DataRequestProto) {
	fmt.Println("Got an error for request", resp)
}

func main() {
	host := flag.String("host", "", "host of the gateway")
	flag.Parse()

	if *host == "" {
		panic("Please pass a host to check the health of")
	}

	s1 := "我认为你很优秀"
	s2 := "我不得不承认你很优秀"
	GRPCClient, err := client.NewGRPCClient(*host)
	if err != nil {
		panic(err)
	}
	request := generateDataRequests(s1, s2)
	GRPCClient.POST(request, OnDone, OnError, nil)

}
