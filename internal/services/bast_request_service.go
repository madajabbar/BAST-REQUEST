package services

import (
	"bast-request/internal/models"
	"bast-request/internal/repositories"
	"errors"
	"fmt"
	"strings"
	"time"
)

type BastRequestService struct {
	repo          *repositories.BastRequestRepository
	formatService *BastFormatService
	seqService    *BastSequenceService
}

func NewBastRequestService(repo *repositories.BastRequestRepository, formatService *BastFormatService, seqService *BastSequenceService) *BastRequestService {
	return &BastRequestService{
		repo:          repo,
		formatService: formatService,
		seqService:    seqService,
	}
}

func (s *BastRequestService) GetAllRequests(customerID, projectID, status string) ([]models.BastRequest, error) {
	return s.repo.FindAll(customerID, projectID, status)
}

func (s *BastRequestService) GetRequestByID(id string) (models.BastRequest, error) {
	return s.repo.FindByID(id)
}

func (s *BastRequestService) CreateRequest(req *models.BastRequest) error {
	// 1. Get Format
	format, err := s.formatService.GetFormatByID(req.FormatID.String())
	if err != nil {
		return errors.New("invalid format ID")
	}

	// 2. Generate BAST Number
	if req.TipeNomor == "Internal" {
		now := time.Now()
		year := now.Year()
		month := int(now.Month())

		// Generate Next Sequence
		nextNum, err := s.seqService.GenerateNextNumber(format.FormatID.String(), year, month)
		if err != nil {
			return fmt.Errorf("failed to generate running number: %v", err)
		}

		// Replace pattern
		// Example pattern: BAST/{YYYY}/{MM}/{SEQ}
		bastNum := format.FormatPattern
		bastNum = strings.ReplaceAll(bastNum, "{YYYY}", fmt.Sprintf("%04d", year))
		bastNum = strings.ReplaceAll(bastNum, "{MM}", fmt.Sprintf("%02d", month))
		bastNum = strings.ReplaceAll(bastNum, "{SEQ}", fmt.Sprintf("%04d", nextNum)) // 4 digit padding, maybe need config

		req.BastNumber = bastNum
	} else if req.TipeNomor == "PO" {
		if req.PoNumber == "" {
			return errors.New("po_number is required for TipeNomor PO")
		}
		req.BastNumber = req.PoNumber
	} else {
		return errors.New("invalid tipe_nomor, must be Internal or PO")
	}

	req.Status = "Active"

	return s.repo.Create(req)
}

func (s *BastRequestService) UpdateStatus(id string, status string) error {
	validStatuses := map[string]bool{"Active": true, "Used": true, "Void": true}
	if !validStatuses[status] {
		return errors.New("invalid status")
	}
	return s.repo.UpdateStatus(id, status)
}
