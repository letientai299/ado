package models

// IdentityRef represents a reference to an identity.
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-requests/get-pull-requests-by-project
type IdentityRef struct {
	// This field contains zero or more interesting links about the graph subject.
	// These links may be invoked to obtain additional relationships or more
	// detailed information about this graph subject.
	Links *ReferenceLinks `json:"_links,omitempty"`
	// The descriptor is the primary way to reference the graph subject while the
	// system is running.
	// This field will uniquely identify the same graph subject across both Accounts
	// and Organizations.
	Descriptor string `json:"descriptor,omitempty"`
	// Deprecated - Can be retrieved by querying the Graph user referenced in the
	// "self" entry of the IdentityRef "_links" dictionary
	DirectoryAlias string `json:"directoryAlias,omitempty"`
	// This is the non-unique display name of the graph subject.
	// To change this field, you must alter its value in the source provider.
	DisplayName string `json:"displayName,omitempty"`
	Id          string `json:"id,omitempty"`
	ImageUrl    string `json:"imageUrl,omitempty"`
	// True if we are unable to identify this identity, so we use a fallback
	// identity instead (one with the same name but a different primary key)
	Inactive bool `json:"inactive,omitempty"`
	// True if the identity is a group.
	IsAadIdentity bool `json:"isAadIdentity,omitempty"`
	// True if the identity is a group.
	IsContainer    bool `json:"isContainer,omitempty"`
	IsExternalUser bool `json:"isExternalUser,omitempty"`
	// Meta-type of the graph subject, e.g., User, Group, Scope, etc.
	SubjectKind string `json:"subjectKind,omitempty"`
	// Used to help identify the source of the graph subject, e.g., VSTS, AAD, MSA, etc.
	UniqueName string `json:"uniqueName,omitempty"`
	Url        string `json:"url,omitempty"`
}

// IdentityRefWithVote represents an identity with a vote on a pull request.
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-requests/get-pull-requests-by-project
type IdentityRefWithVote struct {
	// Identity reference.
	IdentityRef
	// Indicates if this reviewer is required.
	IsRequired bool `json:"isRequired,omitempty"`
	// Vote on a pull request:
	// 10 - approved 5 - approved with suggestions 0 - no vote -5 - waiting for
	// author -10 - rejected
	Vote int `json:"vote,omitempty"`
	// Groups or teams that this reviewer is a member of.
	VotedFor []IdentityRef `json:"votedFor,omitempty"`
}

type Identity struct {
	Id                  string `json:"id"`
	Descriptor          string `json:"descriptor"`
	SubjectDescriptor   string `json:"subjectDescriptor"`
	ProviderDisplayName string `json:"providerDisplayName"`
	IsActive            bool   `json:"isActive"`
	Account             string `json:"properties.Account.$value"`
}
