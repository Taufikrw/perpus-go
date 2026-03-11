package resources

import "belajar-go/models"

type MemberResource struct {
	ID          string       `json:"id"`
	MemberCode  string       `json:"member_code"`
	PhoneNumber string       `json:"phone_number"`
	Address     string       `json:"address"`
	IsApproved  bool         `json:"is_approved"`
	UserID      UserResource `json:"user"`
	DeletedAt   *string      `json:"deleted_at,omitempty"`
}

func FormatMember(member models.Member) MemberResource {
	var deletedAt *string
	if member.DeletedAt.Valid {
		deletedAtStr := member.DeletedAt.Time.String()
		deletedAt = &deletedAtStr
	}

	return MemberResource{
		ID:          member.ID.String(),
		MemberCode:  member.MemberCode,
		PhoneNumber: member.PhoneNumber,
		Address:     member.Address,
		IsApproved:  member.IsApproved,
		DeletedAt:   deletedAt,
		UserID:      FormatUser(member.User),
	}
}

func FormatMembers(members []models.Member) []MemberResource {
	var memberResources []MemberResource
	for _, member := range members {
		memberResources = append(memberResources, FormatMember(member))
	}
	return memberResources
}
