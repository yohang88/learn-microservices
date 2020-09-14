package main

import (
	"context"
	gt "github.com/go-kit/kit/transport/grpc"
	reviewpb "github.com/yohang88/learn-microservices/review/proto"
)

type GRPCServer struct {
	getReviewList gt.Handler
}

func (s *GRPCServer) GetReviewList(ctx context.Context, req *reviewpb.GetReviewListRequest) (*reviewpb.GetReviewListResponse, error) {
	_, resp, err := s.getReviewList.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*reviewpb.GetReviewListResponse), nil
}


func MakeGRPCServer(_ context.Context, endpoint Endpoints) reviewpb.ReviewServiceServer {
	return &GRPCServer{
		getReviewList: gt.NewServer(
			endpoint.GetReviewListEndpoint,
			decodeGetReviewListRequest,
			encodeGetReviewListResponse,
		),
	}
}

type GetReviewListRequest struct {
	ProductId string
}

type GetReviewListResponse []ProductReview

// GRPC Request
func decodeGetReviewListRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*reviewpb.GetReviewListRequest)

	return GetReviewListRequest{
		ProductId: req.ProductId,
	}, nil
}

func encodeGetReviewListResponse(ctx context.Context, r interface{}) (interface{}, error) {
	//res := r.(GetReviewListResponse)

	var reviews []*reviewpb.Review
	reviews = append(reviews, &reviewpb.Review{
		Id: "1",
		Content: "Content 1",
	})

	reviews = append(reviews, &reviewpb.Review{
		Id: "2",
		Content: "Content 2",
	})

	reviews = append(reviews, &reviewpb.Review{
		Id: "3",
		Content: "Content 3",
	})

	return &reviewpb.GetReviewListResponse{
		Reviews: reviews,
	}, nil
}