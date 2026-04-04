package domain

import (
	"context"
	"net"
	"strings"
)

// Provider represents a detected email provider.
type Provider struct {
	Name string `json:"name"`
	Type string `json:"type"` // "workspace" (Google), "exchange" (O365), "hosting", "other"
}

// DNSAnalysis contains the results of analyzing a domain's existing DNS records.
type DNSAnalysis struct {
	Domain           string    `json:"domain"`
	HasMX            bool      `json:"has_mx"`
	MXRecords        []string  `json:"mx_records"`
	DetectedProvider *Provider `json:"detected_provider,omitempty"`
	HasSPF           bool      `json:"has_spf"`
	ExistingSPF      string    `json:"existing_spf,omitempty"`
	MergedSPF        string    `json:"merged_spf"`
	HasDMARC         bool      `json:"has_dmarc"`
	ExistingDMARC    string    `json:"existing_dmarc,omitempty"`
	HasDKIM          bool      `json:"has_existing_dkim"`
	IsCloudflare           bool   `json:"is_cloudflare"`
	Nameservers            []string `json:"nameservers,omitempty"`
	DomainConnectSupported bool   `json:"domain_connect_supported"`
	DomainConnectProvider  string `json:"domain_connect_provider,omitempty"`
	Recommendations  []Recommendation `json:"recommendations"`
}

// Recommendation is a suggested action for the user.
type Recommendation struct {
	Type    string `json:"type"`    // "warning", "info", "success"
	Title   string `json:"title"`
	Message string `json:"message"`
	Action  string `json:"action,omitempty"` // "use_subdomain", "merge_spf", "skip_dmarc", "skip_mx"
}

// AnalyzeDomain performs DNS lookups on a domain to detect existing email configuration.
func AnalyzeDomain(ctx context.Context, name string) (*DNSAnalysis, error) {
	analysis := &DNSAnalysis{
		Domain:          name,
		Recommendations: []Recommendation{},
	}

	// Check MX records
	mxRecords, err := net.LookupMX(name)
	if err == nil && len(mxRecords) > 0 {
		analysis.HasMX = true
		for _, mx := range mxRecords {
			analysis.MXRecords = append(analysis.MXRecords, strings.TrimSuffix(mx.Host, "."))
		}
		analysis.DetectedProvider = detectProvider(analysis.MXRecords)
	}

	// Check existing SPF
	txtRecords, err := net.LookupTXT(name)
	if err == nil {
		for _, txt := range txtRecords {
			if strings.HasPrefix(txt, "v=spf1") {
				analysis.HasSPF = true
				analysis.ExistingSPF = txt
				break
			}
		}
	}

	// Check existing DMARC
	dmarcRecords, err := net.LookupTXT("_dmarc." + name)
	if err == nil {
		for _, txt := range dmarcRecords {
			if strings.HasPrefix(txt, "v=DMARC1") {
				analysis.HasDMARC = true
				analysis.ExistingDMARC = txt
				break
			}
		}
	}

	// Check nameservers (detect Cloudflare)
	nsRecords, err := net.LookupNS(name)
	if err == nil {
		for _, ns := range nsRecords {
			host := strings.TrimSuffix(ns.Host, ".")
			analysis.Nameservers = append(analysis.Nameservers, host)
			if strings.Contains(strings.ToLower(host), "cloudflare") ||
				strings.Contains(strings.ToLower(host), "ns.cloudflare.com") {
				analysis.IsCloudflare = true
			}
		}
	}

	// Check existing DKIM (our selector)
	dkimRecords, err := net.LookupTXT(dkimSelector + "._domainkey." + name)
	if err == nil && len(dkimRecords) > 0 {
		analysis.HasDKIM = true
	}

	// Check Domain Connect support
	if provider, err := DiscoverProvider(ctx, name); err == nil {
		analysis.DomainConnectSupported = true
		analysis.DomainConnectProvider = provider.ProviderName
	} else {
		analysis.DomainConnectSupported = CheckDomainConnectSupport(name)
	}

	// Generate merged SPF
	analysis.MergedSPF = generateMergedSPF(analysis.ExistingSPF)

	// Build recommendations
	analysis.Recommendations = buildRecommendations(analysis)

	return analysis, nil
}

// detectProvider identifies the email provider from MX records.
func detectProvider(mxRecords []string) *Provider {
	for _, mx := range mxRecords {
		mx = strings.ToLower(mx)
		switch {
		case strings.Contains(mx, "google.com") || strings.Contains(mx, "googlemail.com") || strings.Contains(mx, "smtp.google.com"):
			return &Provider{Name: "Google Workspace", Type: "workspace"}
		case strings.Contains(mx, "outlook.com") || strings.Contains(mx, "protection.outlook.com") || strings.Contains(mx, "mail.protection.outlook.com"):
			return &Provider{Name: "Microsoft 365", Type: "exchange"}
		case strings.Contains(mx, "zoho.com") || strings.Contains(mx, "zoho.eu"):
			return &Provider{Name: "Zoho Mail", Type: "workspace"}
		case strings.Contains(mx, "protonmail.ch") || strings.Contains(mx, "protonmail.com"):
			return &Provider{Name: "ProtonMail", Type: "workspace"}
		case strings.Contains(mx, "mxrouting.net") || strings.Contains(mx, "mxroute.com"):
			return &Provider{Name: "MXRoute", Type: "hosting"}
		case strings.Contains(mx, "hover.com"):
			return &Provider{Name: "Hover", Type: "hosting"}
		case strings.Contains(mx, "icloud.com") || strings.Contains(mx, "apple.com"):
			return &Provider{Name: "iCloud Mail", Type: "workspace"}
		case strings.Contains(mx, "yahoodns.net"):
			return &Provider{Name: "Yahoo Mail", Type: "workspace"}
		case strings.Contains(mx, "mailgun.org"):
			return &Provider{Name: "Mailgun", Type: "transactional"}
		case strings.Contains(mx, "sendgrid.net"):
			return &Provider{Name: "SendGrid", Type: "transactional"}
		}
	}
	if len(mxRecords) > 0 {
		return &Provider{Name: "Custom Mail Server", Type: "other"}
	}
	return nil
}

// generateMergedSPF creates the correct SPF record that includes Hello Mail.
func generateMergedSPF(existingSPF string) string {
	helloMailInclude := "include:spf.hellomail.dev"

	if existingSPF == "" {
		return "v=spf1 " + helloMailInclude + " ~all"
	}

	// Already has our include
	if strings.Contains(existingSPF, "spf.hellomail.dev") {
		return existingSPF
	}

	// Insert our include before the mechanism (~all, -all, ?all)
	for _, mechanism := range []string{" ~all", " -all", " ?all", " +all"} {
		if strings.Contains(existingSPF, mechanism) {
			return strings.Replace(existingSPF, mechanism, " "+helloMailInclude+mechanism, 1)
		}
	}

	// No mechanism found — append
	return existingSPF + " " + helloMailInclude + " ~all"
}

// buildRecommendations generates user-facing recommendations based on DNS analysis.
func buildRecommendations(a *DNSAnalysis) []Recommendation {
	var recs []Recommendation

	if a.DetectedProvider != nil {
		recs = append(recs, Recommendation{
			Type:    "warning",
			Title:   a.DetectedProvider.Name + " detected",
			Message: "This domain uses " + a.DetectedProvider.Name + " for email. Adding Hello Mail MX records would break your existing mailboxes. We recommend using a subdomain (e.g., mail." + a.Domain + ") or sending only without inbound.",
			Action:  "use_subdomain",
		})
	} else if a.HasMX {
		recs = append(recs, Recommendation{
			Type:    "warning",
			Title:   "Existing MX records found",
			Message: "This domain has existing MX records. Enabling inbound email would conflict with your current mail setup.",
			Action:  "skip_mx",
		})
	} else {
		recs = append(recs, Recommendation{
			Type:    "success",
			Title:   "No email provider detected",
			Message: "This domain has no existing MX records. You can safely enable inbound email.",
		})
	}

	if a.HasSPF {
		recs = append(recs, Recommendation{
			Type:    "info",
			Title:   "Existing SPF record found",
			Message: "We've generated a merged SPF record that includes both your existing configuration and Hello Mail. Replace your current SPF record with the merged version.",
			Action:  "merge_spf",
		})
	} else {
		recs = append(recs, Recommendation{
			Type:    "success",
			Title:   "No existing SPF record",
			Message: "We'll create a new SPF record for you.",
		})
	}

	if a.HasDMARC {
		recs = append(recs, Recommendation{
			Type:    "info",
			Title:   "Existing DMARC record found",
			Message: "Your domain already has a DMARC policy. No changes needed — Hello Mail will respect your existing DMARC configuration.",
			Action:  "skip_dmarc",
		})
	}

	if a.HasDKIM {
		recs = append(recs, Recommendation{
			Type:    "info",
			Title:   "DKIM selector already exists",
			Message: "A DKIM record with our selector already exists. This may be from a previous setup.",
		})
	}

	return recs
}
