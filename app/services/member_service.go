package services

import (
	"backend/app/models"
	"backend/app/repositories"
)

type MemberService struct {
	Repo *repositories.MembersRepository
}

func (ms *MemberService) GetMembers() ([]models.Members, error) {
	return ms.Repo.GetMembers()
}

func (ms *MemberService) GetMemberById(id string) (*models.Members, error) {
	return ms.Repo.GetMembersById(id)
}

func (ms *MemberService) CreateMember(member *models.Members) error {
	return ms.Repo.CreateMembers(member)
}

func (ms *MemberService) UpdateMember(member *models.Members) error {
	return ms.Repo.UpdateMembers(member)
}
