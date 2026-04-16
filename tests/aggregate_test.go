package tests

import (
	"testing"

	"github.com/carddemo/project/src/domain/account/model"
	"github.com/carddemo/project/src/domain/card/model"
	"github.com/carddemo/project/src/domain/cardpolicy/model"
	"github.com/carddemo/project/src/domain/shared"
	"github.com/carddemo/project/src/domain/batchsettlement/model"
	"github.com/carddemo/project/src/domain/exportjob/model"
	"github.com/carddemo/project/src/domain/report/model"
	"github.com/carddemo/project/src/domain/transaction/model"
	"github.com/carddemo/project/src/domain/userprofile/model"
)

func TestAccountAggregate_Execute(t *testing.T) {
	acc := model.NewAccount("123")
	type UnknownCommand struct{}

	events, err := acc.Execute(UnknownCommand{})

	if err != shared.ErrUnknownCommand {
		t.Errorf("Expected ErrUnknownCommand, got %v", err)
	}
	if events != nil {
		t.Errorf("Expected nil events, got %v", events)
	}
}

func TestUserProfileAggregate_Execute(t *testing.T) {
	prof := model.NewUserProfile("456")
	type UnknownCommand struct{}

	events, err := prof.Execute(UnknownCommand{})

	if err != shared.ErrUnknownCommand {
		t.Errorf("Expected ErrUnknownCommand, got %v", err)
	}
	if events != nil {
		t.Errorf("Expected nil events, got %v", events)
	}
}

func TestCardAggregate_Execute(t *testing.T) {
	card := model.NewCard("789")
	type UnknownCommand struct{}

	events, err := card.Execute(UnknownCommand{})

	if err != shared.ErrUnknownCommand {
		t.Errorf("Expected ErrUnknownCommand, got %v", err)
	}
}

func TestCardPolicyAggregate_Execute(t *testing.T) {
	pol := model.NewCardPolicy("101")
	type UnknownCommand struct{}

	events, err := pol.Execute(UnknownCommand{})

	if err != shared.ErrUnknownCommand {
		t.Errorf("Expected ErrUnknownCommand, got %v", err)
	}
}

func TestTransactionAggregate_Execute(t *testing.T) {
	tx := model.NewTransaction("201")
	type UnknownCommand struct{}

	events, err := tx.Execute(UnknownCommand{})

	if err != shared.ErrUnknownCommand {
		t.Errorf("Expected ErrUnknownCommand, got %v", err)
	}
}

func TestBatchSettlementAggregate_Execute(t *testing.T) {
	batch := model.NewBatchSettlement("301")
	type UnknownCommand struct{}

	events, err := batch.Execute(UnknownCommand{})

	if err != shared.ErrUnknownCommand {
		t.Errorf("Expected ErrUnknownCommand, got %v", err)
	}
}

func TestReportAggregate_Execute(t *testing.T) {
	rep := model.NewReport("401")
	type UnknownCommand struct{}

	events, err := rep.Execute(UnknownCommand{})

	if err != shared.ErrUnknownCommand {
		t.Errorf("Expected ErrUnknownCommand, got %v", err)
	}
}

func TestExportJobAggregate_Execute(t *testing.T) {
	job := model.NewExportJob("501")
	type UnknownCommand struct{}

	events, err := job.Execute(UnknownCommand{})

	if err != shared.ErrUnknownCommand {
		t.Errorf("Expected ErrUnknownCommand, got %v", err)
	}
}
