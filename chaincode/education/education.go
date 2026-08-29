package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Education struct {
	CertNo       string `json:"certNo"`
	Name         string `json:"name"`
	StudentID    string `json:"studentId"`
	School       string `json:"school"`
	Major        string `json:"major"`
	Degree       string `json:"degree"`
	GraduationDate string `json:"graduationDate"`
	IssueDate    string `json:"issueDate"`
	Status       string `json:"status"`
	AuthorizedViewers []string `json:"authorizedViewers"`
}

type QueryResult struct {
	Key    string `json:"key"`
	Record Education `json:"record"`
}

type HistoryRecord struct {
	TxID      string `json:"txId"`
	Value     Education `json:"value"`
	Timestamp time.Time `json:"timestamp"`
	IsDelete  bool `json:"isDelete"`
}

func (s *SmartContract) AddEducation(ctx contractapi.TransactionContextInterface,
	certNo, name, studentID, school, major, degree, graduationDate string) error {

	exists, err := s.EducationExists(ctx, certNo)
	if err != nil {
		return fmt.Errorf("failed to check existence: %v", err)
	}
	if exists {
		return fmt.Errorf("education record %s already exists", certNo)
	}

	education := Education{
		CertNo:         certNo,
		Name:           name,
		StudentID:      studentID,
		School:         school,
		Major:          major,
		Degree:         degree,
		GraduationDate: graduationDate,
		IssueDate:      time.Now().Format("2006-01-02"),
		Status:         "active",
		AuthorizedViewers: []string{},
	}

	educationJSON, err := json.Marshal(education)
	if err != nil {
		return fmt.Errorf("failed to marshal education: %v", err)
	}

	err = ctx.GetStub().PutState(certNo, educationJSON)
	if err != nil {
		return fmt.Errorf("failed to put state: %v", err)
	}

	ctx.GetStub().SetEvent("AddEducation", educationJSON)
	return nil
}

func (s *SmartContract) UpdateEducation(ctx contractapi.TransactionContextInterface,
	certNo, name, studentID, school, major, degree, graduationDate string) error {

	education, err := s.QueryEducationByID(ctx, certNo)
	if err != nil {
		return err
	}

	education.Name = name
	education.StudentID = studentID
	education.School = school
	education.Major = major
	education.Degree = degree
	education.GraduationDate = graduationDate

	educationJSON, err := json.Marshal(education)
	if err != nil {
		return fmt.Errorf("failed to marshal: %v", err)
	}

	return ctx.GetStub().PutState(certNo, educationJSON)
}

func (s *SmartContract) QueryEducationByID(ctx contractapi.TransactionContextInterface,
	certNo string) (*Education, error) {

	educationJSON, err := ctx.GetStub().GetState(certNo)
	if err != nil {
		return nil, fmt.Errorf("failed to read from state: %v", err)
	}
	if educationJSON == nil {
		return nil, fmt.Errorf("education record %s does not exist", certNo)
	}

	var education Education
	err = json.Unmarshal(educationJSON, &education)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %v", err)
	}

	return &education, nil
}

func (s *SmartContract) VerifyByCertNoAndName(ctx contractapi.TransactionContextInterface,
	certNo, name string) (bool, error) {

	education, err := s.QueryEducationByID(ctx, certNo)
	if err != nil {
		return false, err
	}

	if education.Name != name {
		return false, nil
	}
	if education.Status != "active" {
		return false, nil
	}

	return true, nil
}

func (s *SmartContract) GetHistoryByID(ctx contractapi.TransactionContextInterface,
	certNo string) ([]HistoryRecord, error) {

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(certNo)
	if err != nil {
		return nil, fmt.Errorf("failed to get history: %v", err)
	}
	defer resultsIterator.Close()

	var records []HistoryRecord
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to iterate: %v", err)
		}

		var education Education
		if !response.IsDelete {
			err = json.Unmarshal(response.Value, &education)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal history: %v", err)
			}
		}

		records = append(records, HistoryRecord{
			TxID:      response.TxId,
			Value:     education,
			Timestamp: time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)),
			IsDelete:  response.IsDelete,
		})
	}

	return records, nil
}

func (s *SmartContract) AuthorizeViewer(ctx contractapi.TransactionContextInterface,
	certNo, viewerID string) error {

	education, err := s.QueryEducationByID(ctx, certNo)
	if err != nil {
		return err
	}

	for _, v := range education.AuthorizedViewers {
		if v == viewerID {
			return fmt.Errorf("viewer %s already authorized", viewerID)
		}
	}

	education.AuthorizedViewers = append(education.AuthorizedViewers, viewerID)

	educationJSON, err := json.Marshal(education)
	if err != nil {
		return fmt.Errorf("failed to marshal: %v", err)
	}

	return ctx.GetStub().PutState(certNo, educationJSON)
}

func (s *SmartContract) QueryAllEducation(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {

	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, fmt.Errorf("failed to get range: %v", err)
	}
	defer resultsIterator.Close()

	var results []QueryResult
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to iterate: %v", err)
		}

		var education Education
		err = json.Unmarshal(queryResponse.Value, &education)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal: %v", err)
		}

		results = append(results, QueryResult{
			Key:    queryResponse.Key,
			Record: education,
		})
	}

	return results, nil
}

func (s *SmartContract) EducationExists(ctx contractapi.TransactionContextInterface,
	certNo string) (bool, error) {

	educationJSON, err := ctx.GetStub().GetState(certNo)
	if err != nil {
		return false, fmt.Errorf("failed to read state: %v", err)
	}
	return educationJSON != nil, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		fmt.Printf("Error creating chaincode: %v", err)
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting chaincode: %v", err)
	}
}
