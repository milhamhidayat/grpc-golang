package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	bpb "blog/pb"
)

var collection *mongo.Collection

type server struct{}

type blogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}

func (s *server) CreateBlog(ctx context.Context, req *bpb.CreateBlogRequest) (*bpb.CreateBlogResponse, error) {
	fmt.Println("create blog request")

	blog := req.GetBlog()

	data := blogItem{
		AuthorID: blog.GetAuthorId(),
		Content:  blog.GetTitle(),
		Title:    blog.GetTitle(),
	}

	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("internal error: %v", err))
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("cannot convert to oid: %v", err))
	}

	return &bpb.CreateBlogResponse{
		Blog: &bpb.Blog{
			Id:       oid.Hex(),
			AuthorId: blog.GetAuthorId(),
			Content:  blog.GetContent(),
			Title:    blog.GetTitle(),
		},
	}, nil
}

func (s *server) ReadBlog(ctx context.Context, req *bpb.ReadBlogRequest) (*bpb.ReadBlogResponse, error) {
	blogID := req.GetBlogId()

	oid, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("cannot parse id"),
		)
	}

	data := &blogItem{}
	filter := bson.M{"_id": oid}

	res := collection.FindOne(context.Background(), filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("cannot find blog with specified id: %v", err),
		)
	}

	return &bpb.ReadBlogResponse{
		Blog: dataToBlogPb(data),
	}, nil
}

func (s *server) UpdateBlog(ctx context.Context, req *bpb.UpdateBlogRequest) (*bpb.UpdateBlogResponse, error) {
	fmt.Println("update blog request")

	blog := req.GetBlog()
	oid, err := primitive.ObjectIDFromHex(blog.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("cannot parse id"),
		)
	}

	data := &blogItem{}
	filter := bson.M{"_id": oid}

	res := collection.FindOne(context.Background(), filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("cannot find blog with specified id: %v", err),
		)
	}

	data.AuthorID = blog.GetAuthorId()
	data.Content = blog.GetContent()
	data.Title = blog.GetTitle()

	_, err = collection.ReplaceOne(context.Background(), filter, data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("cannot update object in mongodb: %v", err),
		)
	}

	return &bpb.UpdateBlogResponse{
		Blog: dataToBlogPb(data),
	}, nil
}

func (s *server) DeleteBlog(ctx context.Context, req *bpb.DeleteBlogRequest) (*bpb.DeleteBlogResponse, error) {
	fmt.Println("delete blog request")

	blogID := req.GetBlogId()
	oid, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("cannot parse id"),
		)
	}

	filter := bson.M{"_id": oid}
	res, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("cannot delete object in mongodb: %v", err),
		)
	}

	if res.DeletedCount == 0 {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("cannot find blog in mongodb: %v", err),
		)
	}

	return &bpb.DeleteBlogResponse{BlogId: req.GetBlogId()}, nil
}

func (s *server) ListBlog(req *bpb.ListBlogRequest, stream bpb.BlogService_ListBlogServer) error {
	fmt.Println("list blog request")

	cur, err := collection.Find(context.Background(), primitive.D{{}})
	if err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("unknown internal error: %v", err),
		)
	}

	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		data := &blogItem{}
		err := cur.Decode(data)
		if err != nil {
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("error while decoding data from mongodb: %v", err),
			)
		}
		stream.Send(&bpb.ListBlogResponse{Blog: dataToBlogPb(data)})
	}
	if err := cur.Err(); err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("unknown internal error: %v", err),
		)
	}
	return nil
}

func dataToBlogPb(data *blogItem) *bpb.Blog {
	return &bpb.Blog{
		Id:       data.ID.Hex(),
		AuthorId: data.AuthorID,
		Content:  data.Content,
		Title:    data.Title,
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Blot Service Started")

	fmt.Println("connecting to mongodb")
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017/blog_mongo"))
	if err != nil {
		log.Fatalf("failed to create new mongodb client: %v", err)
	}

	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatalf("failed to connect to mongo db: %v", err)
	}

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	collection = client.Database("blog_mongo").Collection("blog")

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	bpb.RegisterBlogServiceServer(s, &server{})

	reflection.Register(s)

	go func() {
		fmt.Println("Starting Server...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch
	fmt.Println("stopping the server")
	s.Stop()
	fmt.Println("close the listener")
	lis.Close()
	fmt.Println("closing mongo db connection")
	client.Disconnect(context.TODO())
	fmt.Println("end of program")
}
