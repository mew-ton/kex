package retrieve

import (
	"github.com/mew-ton/kex/internal/domain"
)

type UseCase struct {
	Repo domain.DocumentRepository
}

func New(repo domain.DocumentRepository) *UseCase {
	return &UseCase{Repo: repo}
}

type Result struct {
	Document *domain.Document
	Found    bool
}

func (uc *UseCase) Execute(id string) Result {
	doc, ok := uc.Repo.GetByID(id)
	return Result{
		Document: doc,
		Found:    ok,
	}
}
