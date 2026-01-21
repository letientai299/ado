package models

// IdentityRef represents a reference to an Azure DevOps identity.
// Identities can be users, groups, or service accounts.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-requests/get#identityref
type IdentityRef struct {
	// Links contains additional navigation links for this identity.
	Links *ReferenceLinks `json:"_links,omitempty"`

	// Descriptor is the primary way to reference the identity while the system runs.
	// This field uniquely identifies the identity across Accounts and Organizations.
	Descriptor string `json:"descriptor,omitempty"`

	// DirectoryAlias is deprecated. Use the Graph API via the "_links" dictionary instead.
	// Previously contained the user's directory alias.
	DirectoryAlias string `json:"directoryAlias,omitempty"`

	// DisplayName is the human-readable name of the identity.
	// This is not unique and is managed by the identity provider.
	DisplayName string `json:"displayName,omitempty"`

	// Id is the unique identifier (GUID) of the identity.
	Id string `json:"id,omitempty"`

	// ImageUrl is the URL of the identity's avatar image.
	ImageUrl string `json:"imageUrl,omitempty"`

	// Inactive indicates if this is a fallback identity.
	// When true, the system couldn't identify the original identity and is using
	// a fallback with the same name but different primary key.
	Inactive bool `json:"inactive,omitempty"`

	// IsAadIdentity indicates if the identity is from Azure Active Directory.
	IsAadIdentity bool `json:"isAadIdentity,omitempty"`

	// IsContainer indicates if the identity is a group or container.
	IsContainer bool `json:"isContainer,omitempty"`

	// IsExternalUser indicates if the identity is an external (guest) user.
	IsExternalUser bool `json:"isExternalUser,omitempty"`

	// SubjectKind indicates the meta-type of the identity.
	// Common values: "user", "group", "scope".
	SubjectKind string `json:"subjectKind,omitempty"`

	// UniqueName helps identify the source of the identity.
	// Often contains the email address or UPN.
	// Common sources: VSTS, AAD, MSA.
	UniqueName string `json:"uniqueName,omitempty"`

	// Url is the REST API URL for this identity.
	Url string `json:"url,omitempty"`
}

// IdentityRefWithVote represents an identity with their vote on a pull request.
// This extends IdentityRef to include reviewer-specific information.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-request-reviewers/list#identityrefwithvote
type IdentityRefWithVote struct {
	// IdentityRef contains the base identity information.
	IdentityRef

	// HasDeclined indicates if the reviewer has declined to review.
	HasDeclined bool `json:"hasDeclined,omitempty"`

	// IsFlagged indicates if the reviewer has been flagged for attention.
	IsFlagged bool `json:"isFlagged,omitempty"`

	// IsRequired indicates if this reviewer is required for PR completion.
	IsRequired bool `json:"isRequired,omitempty"`

	// IsReRequired indicates if re-approval is required after changes.
	// Set when policies require re-approval after source branch updates.
	IsReRequired bool `json:"isReRequired,omitempty"`

	// ReviewerUrl is the URL to this reviewer's review details.
	ReviewerUrl string `json:"reviewerUrl,omitempty"`

	// Vote is the reviewer's current vote on the pull request.
	// 10: approved, 5: approved with suggestions, 0: no vote,
	// -5: waiting for author, -10: rejected.
	Vote int `json:"vote,omitempty"`

	// VotedFor contains groups or teams that this reviewer's vote applies to.
	// When a reviewer is part of a required group, their vote counts for the group.
	VotedFor []IdentityRef `json:"votedFor,omitempty"`
}

// Identity represents the authenticated user's identity information.
// This is returned from the connection data API.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/ims/identities/read-identities
type Identity struct {
	// Id is the unique identifier (GUID) of the identity.
	Id string `json:"id"`

	// Descriptor is the security descriptor for this identity.
	Descriptor string `json:"descriptor"`

	// SubjectDescriptor is the subject descriptor for VSSPS.
	SubjectDescriptor string `json:"subjectDescriptor"`

	// ProviderDisplayName is the display name from the identity provider.
	ProviderDisplayName string `json:"providerDisplayName"`

	// IsActive indicates if the identity is currently active.
	IsActive bool `json:"isActive"`

	// Account is the account name from identity properties.
	Account string `json:"properties.Account.$value"`
}
