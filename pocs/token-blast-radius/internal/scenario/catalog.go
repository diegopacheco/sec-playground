package scenario

import (
	"time"

	"token-blast-radius/internal/domain"
	"token-blast-radius/internal/token"
)

func Build(tokens *token.Service) ([]domain.Scenario, error) {
	expires := time.Now().UTC().Add(24 * time.Hour).Unix()
	scenarios := []domain.Scenario{
		contained(expires),
		overDelegated(expires),
		supportSession(expires),
	}
	for index := range scenarios {
		value, err := tokens.Issue(scenarios[index].Claims)
		if err != nil {
			return nil, err
		}
		scenarios[index].Token = value
	}
	return scenarios, nil
}

func contained(expires int64) domain.Scenario {
	return domain.Scenario{
		ID: "01-contained", Name: "Contained analyst", Summary: "A read-only agent reaches two approved reports through one team boundary.",
		Question: "Can this agent read anything beyond the incident workspace?",
		Claims:   domain.Claims{Issuer: "identity.local", Subject: "agent:analyst", Actor: "user:maya", Audience: "knowledge-api", Scopes: []string{"read"}, MaximumDepth: 3, ExpiresAt: expires, ID: "tok-contained-7a91"},
		Nodes: []domain.Node{
			node("user:maya", "Maya Chen", "human", "internal", "Human principal accountable for the agent"),
			node("agent:analyst", "Analysis agent", "agent", "internal", "Non-human identity named by the token subject"),
			node("workspace:incidents", "Incident workspace", "workspace", "confidential", "Security team relationship boundary"),
			node("document:timeline", "Incident timeline", "document", "confidential", "Approved chronology for active investigations"),
			node("document:runbook", "Response runbook", "document", "confidential", "Approved operational response procedures"),
			node("vault:tokens", "Support token vault", "vault", "critical", "Live support credentials and session material"),
			node("service:tenant-admin", "Tenant administration", "service", "critical", "Privileged tenant control surface"),
		},
		Relationships: []domain.Relationship{
			edge("agent:analyst", "workspace:incidents", "read", "read", "agent is assigned to the incident workspace"),
			edge("workspace:incidents", "document:timeline", "read", "read", "workspace readers inherit timeline access"),
			edge("workspace:incidents", "document:runbook", "read", "read", "workspace readers inherit runbook access"),
			edge("workspace:incidents", "vault:tokens", "export", "export", "vault export requires a stronger scope"),
			edge("workspace:incidents", "service:tenant-admin", "admin", "admin", "tenant control requires administrative scope"),
		},
	}
}

func overDelegated(expires int64) domain.Scenario {
	return domain.Scenario{
		ID: "02-over-delegated", Name: "Over-delegated agent", Summary: "Broad scopes and deep delegation connect an assistant to financial, identity, and credential systems.",
		Question: "How far can one helpful-looking agent credential travel?",
		Claims:   domain.Claims{Issuer: "identity.local", Subject: "agent:ops-copilot", Actor: "user:leon", Audience: "operations-api", Scopes: []string{"read", "export", "admin", "delegate"}, MaximumDepth: 6, ExpiresAt: expires, ID: "tok-wide-f042"},
		Nodes: []domain.Node{
			node("user:leon", "Leon Brooks", "human", "internal", "Human principal accountable for the agent"),
			node("agent:ops-copilot", "Operations copilot", "agent", "internal", "Agent holding the presented credential"),
			node("workspace:operations", "Operations workspace", "workspace", "internal", "Shared operational authorization boundary"),
			node("database:customer", "Customer database", "database", "critical", "Customer identity and profile records"),
			node("service:billing", "Billing service", "service", "critical", "Invoices, refunds, and payment state"),
			node("vault:production", "Production vault", "vault", "critical", "Runtime secrets for production systems"),
			node("service:identity-admin", "Identity administration", "service", "critical", "User, factor, and session administration"),
			node("agent:automation", "Automation agent", "agent", "internal", "Downstream agent trusted by operations"),
			node("document:quarterly", "Quarterly report", "document", "confidential", "Internal financial performance report"),
			node("account:vendor", "Vendor account", "account", "confidential", "External procurement account"),
		},
		Relationships: []domain.Relationship{
			edge("agent:ops-copilot", "workspace:operations", "read", "read", "agent belongs to the operations workspace"),
			edge("workspace:operations", "database:customer", "read", "read", "workspace role inherits customer lookup"),
			edge("workspace:operations", "service:billing", "admin", "admin", "operations administrators manage billing"),
			edge("workspace:operations", "vault:production", "export", "export", "operations can export runtime credentials"),
			edge("workspace:operations", "service:identity-admin", "admin", "admin", "operations administrators manage identity"),
			edge("agent:ops-copilot", "agent:automation", "delegate", "delegate", "primary agent can delegate background work"),
			edge("agent:automation", "document:quarterly", "read", "read", "automation agent reads financial reports"),
			edge("agent:automation", "account:vendor", "admin", "admin", "automation agent manages procurement"),
			edge("service:identity-admin", "vault:production", "export", "export", "identity administration can recover connector secrets"),
		},
	}
}

func supportSession(expires int64) domain.Scenario {
	return domain.Scenario{
		ID: "03-support-session", Name: "Captured support session", Summary: "A support identity can cross tenant boundaries through reset and impersonation relationships.",
		Question: "What can an attacker reach before this session is revoked?",
		Claims:   domain.Claims{Issuer: "identity.local", Subject: "agent:support-assistant", Actor: "user:support-engineer", Audience: "support-api", Scopes: []string{"read", "reset", "impersonate", "export"}, MaximumDepth: 5, ExpiresAt: expires, ID: "tok-support-b811"},
		Nodes: []domain.Node{
			node("user:support-engineer", "Support engineer", "human", "internal", "Human principal linked to the support session"),
			node("agent:support-assistant", "Support assistant", "agent", "internal", "Agent operating with the captured session"),
			node("service:case-system", "Case system", "service", "confidential", "Customer cases and uploaded diagnostics"),
			node("document:har", "Uploaded HAR", "document", "critical", "Diagnostic archive containing session material"),
			node("account:tenant-a", "Tenant A administrator", "account", "critical", "Administrative identity in Tenant A"),
			node("account:tenant-b", "Tenant B administrator", "account", "critical", "Administrative identity in Tenant B"),
			node("service:mfa", "Factor reset service", "service", "critical", "Privileged factor reset workflow"),
			node("database:directory", "Customer directory", "database", "confidential", "Cross-tenant customer contacts"),
		},
		Relationships: []domain.Relationship{
			edge("agent:support-assistant", "service:case-system", "read", "read", "support assistant can read assigned cases"),
			edge("service:case-system", "document:har", "export", "export", "case readers can retrieve uploaded diagnostics"),
			edge("service:case-system", "account:tenant-a", "impersonate", "impersonate", "support workflow can enter Tenant A"),
			edge("service:case-system", "account:tenant-b", "impersonate", "impersonate", "support workflow can enter Tenant B"),
			edge("account:tenant-a", "service:mfa", "reset", "reset", "tenant administrator can reset factors"),
			edge("account:tenant-b", "service:mfa", "reset", "reset", "tenant administrator can reset factors"),
			edge("service:case-system", "database:directory", "read", "read", "support can search the customer directory"),
		},
	}
}

func node(id string, label string, kind string, sensitivity string, description string) domain.Node {
	return domain.Node{ID: id, Label: label, Type: kind, Sensitivity: sensitivity, Description: description}
}

func edge(from string, to string, permission string, scope string, rationale string) domain.Relationship {
	return domain.Relationship{From: from, To: to, Permission: permission, RequiredScope: scope, Rationale: rationale}
}
