package services

import (
	"context"
	"fmt"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type SMSService struct {
	client     *twilio.RestClient
	fromNumber string
}

func NewSMSService(accountSID, authToken, fromNumber string) *SMSService {
	if accountSID == "" || authToken == "" {
		return nil // SMS service disabled
	}

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSID,
		Password: authToken,
	})

	return &SMSService{
		client:     client,
		fromNumber: fromNumber,
	}
}

func (s *SMSService) SendSMS(ctx context.Context, toNumber, message string) error {
	if s == nil || s.client == nil {
		return nil // SMS service not configured
	}

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(toNumber)
	params.SetFrom(s.fromNumber)
	params.SetBody(message)

	resp, err := s.client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}

	if resp.Sid != nil {
		return nil
	}

	return fmt.Errorf("SMS send failed: no SID returned")
}

func (s *SMSService) SendApplicationUpdateSMS(ctx context.Context, toNumber string, applicationID uint, status string) error {
	if s == nil {
		return nil
	}

	message := fmt.Sprintf("Your visa application #%d status has been updated to: %s. Check your account for details.", applicationID, status)
	return s.SendSMS(ctx, toNumber, message)
}
