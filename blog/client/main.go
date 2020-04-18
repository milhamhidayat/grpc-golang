package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc"

	bpb "blog/pb"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	c := bpb.NewBlogServiceClient(cc)

	// fmt.Println("creating the blog")
	// blog := bpb.Blog{
	// 	AuthorId: "John",
	// 	Title:    "Doe",
	// 	Content:  "Content of the first blog",
	// }

	/**
	 * Create Blog
	 */
	// res, err := c.CreateBlog(context.Background(), &bpb.CreateBlogRequest{Blog: &blog})
	// if err != nil {
	// 	log.Fatalf("unexpected err: %v", err)
	// }
	// fmt.Printf("blog has been created: %v\n", res)

	/**
	 * Read Blog
	 */
	// fmt.Println("reading a blog")
	// res2, err2 := c.ReadBlog(context.Background(), &bpb.ReadBlogRequest{
	// 	BlogId: "5e98de361079d5101bfb6cc3",
	// })
	// if err2 != nil {
	// 	log.Fatalf("error happened while reading: %v", err2)
	// }
	// fmt.Println("blog response:", res2)

	/**
	 * Update blog
	 */
	// newBlog := &bpb.Blog{
	// 	Id:       "5e98de361079d5101bfb6cc3",
	// 	AuthorId: "changedAuthorID",
	// 	Title:    "My first blog (edited)",
	// 	Content:  "Content of the first blog, with some awesome additions",
	// }
	// res3, err3 := c.UpdateBlog(context.Background(), &bpb.UpdateBlogRequest{Blog: newBlog})
	// if err3 != nil {
	// 	log.Fatalf("error happened while updating; %v", err3)
	// }
	// fmt.Printf("blog was updated: %v\n", res3)

	/**
	 *
	 * Delete blog
	 */
	// res4, err4 := c.DeleteBlog(context.Background(), &bpb.DeleteBlogRequest{BlogId: "5e98de361079d5101bfb6cc3"})
	// if err4 != nil {
	// 	log.Fatalf("error happened while deleteing: %v", err4)
	// }
	// fmt.Printf("blog was deleted: %v\n", res4)

	/**
	 * List blogs
	 */
	stream, err := c.ListBlog(context.Background(), &bpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("error while calling listBlog RPC: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("something happened: %v", err)
		}
		fmt.Println(res.GetBlog())
	}
}
