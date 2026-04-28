package v1

import (
	"context"

	"github.com/boldlogic/portfolio-lens-quik/quik-portfolio/internal/models"
	quikv1 "github.com/boldlogic/portfolio-lens-quik/proto/gen/go/quik/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func currentQuoteToProto(q models.CurrentQuote) *quikv1.CurrentQuote {
	pb := &quikv1.CurrentQuote{
		Ticker:          q.Ticker,
		InstrumentType:  q.InstrumentType,
		PriceCurrency:   q.PriceCurrency,
		AccruedCurrency: q.AccruedCurrency,
		QuoteDate:       timestamppb.New(q.QuoteDate),
	}
	if q.ISIN != nil {
		pb.Isin = q.ISIN
	}
	if q.LastPrice != nil {
		pb.LastPrice = ptr(q.LastPrice.String())
	}
	if q.ClosePrice != nil {
		pb.ClosePrice = ptr(q.ClosePrice.String())
	}
	if q.AccruedInt != nil {
		pb.AccruedInt = ptr(q.AccruedInt.String())
	}
	if q.FaceValue != nil {
		pb.FaceValue = ptr(q.FaceValue.String())
	}
	return pb
}

func (h *Handler) GetCurrentQuotes(ctx context.Context, req *quikv1.GetCurrentQuotesRequest) (*quikv1.GetCurrentQuotesResponse, error) {
	quotes, err := h.service.GetCurrentQuotes(ctx)
	if err != nil {
		h.logger.Error("котировки: чтение", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to get current quotes: %v", err)
	}

	pbQuotes := make([]*quikv1.CurrentQuote, len(quotes))
	for i, q := range quotes {
		pbQuotes[i] = currentQuoteToProto(q)
	}

	return &quikv1.GetCurrentQuotesResponse{Quotes: pbQuotes}, nil
}

func (h *Handler) StreamCurrentQuotes(req *quikv1.GetCurrentQuotesRequest, stream quikv1.LimitsService_StreamCurrentQuotesServer) error {
	quotes, err := h.service.GetCurrentQuotes(stream.Context())
	if err != nil {
		h.logger.Error("котировки: стрим, чтение", zap.Error(err))
		return status.Errorf(codes.Internal, "failed to stream current quotes: %v", err)
	}

	for _, q := range quotes {
		if err := stream.Send(currentQuoteToProto(q)); err != nil {
			h.logger.Error("котировки: стрим, отправка", zap.Error(err))
			return status.Errorf(codes.Internal, "failed to stream current quote: %v", err)
		}
	}

	return nil
}

func (h *Handler) StreamCurrentQuotesForKeys(req *quikv1.GetCurrentQuotesForKeysRequest, stream quikv1.LimitsService_StreamCurrentQuotesForKeysServer) error {
	quotes, err := h.service.GetCurrentQuotesForKeys(stream.Context(), req.GetInstrumentKeys())
	if err != nil {
		h.logger.Error("котировки по ключам: чтение", zap.Error(err))
		return status.Errorf(codes.Internal, "failed to stream current quotes for keys: %v", err)
	}
	for _, q := range quotes {
		if err := stream.Send(currentQuoteToProto(q)); err != nil {
			h.logger.Error("котировки по ключам: стрим, отправка", zap.Error(err))
			return status.Errorf(codes.Internal, "failed to stream current quote for keys: %v", err)
		}
	}
	return nil
}
