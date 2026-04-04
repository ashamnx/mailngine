import { Component, inject, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { ActivatedRoute, Router } from '@angular/router';
import { FormsModule } from '@angular/forms';

interface DnsRecord {
  id: string;
  record_type: string;
  host: string;
  value: string;
  purpose: string;
  status: string;
}

interface DomainInfo {
  id: string;
  name: string;
  status: string;
  open_tracking: boolean;
  click_tracking: boolean;
  created_at: string;
}

interface DomainDetailResponse {
  domain: DomainInfo;
  dns_records: DnsRecord[];
}

@Component({
  selector: 'app-domain-detail',
  standalone: true,
  imports: [FormsModule],
  templateUrl: './domain-detail.component.html',
  styleUrl: './domain-detail.component.scss',
})
export class DomainDetailComponent {
  private readonly http = inject(HttpClient);
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);

  readonly domainInfo = signal<DomainInfo | null>(null);
  readonly dnsRecords = signal<DnsRecord[]>([]);
  readonly isLoading = signal(true);
  readonly isVerifying = signal(false);
  readonly isSavingSettings = signal(false);
  readonly isConfiguringDns = signal(false);
  readonly cfConfigSuccess = signal(false);
  readonly showDeleteConfirm = signal(false);
  readonly isDeleting = signal(false);
  readonly copiedKey = signal<string | null>(null);
  readonly expandedIndex = signal<number | null>(null);
  readonly domainConnectSupported = signal(false);
  readonly domainConnectProvider = signal('');
  readonly domainConnectURL = signal('');
  readonly justConnected = signal(false);

  readonly openTracking = signal(false);
  readonly clickTracking = signal(false);

  private domainId = '';

  constructor() {
    this.domainId = this.route.snapshot.params['id'];
    this.loadDomain();
  }

  loadDomain(): void {
    this.isLoading.set(true);
    this.http
      .get<{ data: DomainDetailResponse }>(`/v1/domains/${this.domainId}`)
      .subscribe({
        next: (res) => {
          this.domainInfo.set(res.data.domain);
          this.dnsRecords.set(res.data.dns_records || []);
          this.openTracking.set(res.data.domain.open_tracking);
          this.clickTracking.set(res.data.domain.click_tracking);
          this.isLoading.set(false);
          this.checkDomainConnect();
        },
        error: () => this.isLoading.set(false),
      });
  }

  private checkDomainConnect(): void {
    this.http
      .get<{ data: { supported: boolean; provider?: string; redirect_url?: string } }>(
        `/v1/domains/${this.domainId}/connect-url`
      )
      .subscribe({
        next: (res) => {
          this.domainConnectSupported.set(res.data.supported);
          this.domainConnectProvider.set(res.data.provider ?? '');
          this.domainConnectURL.set(res.data.redirect_url ?? '');
        },
        error: () => {},
      });

    // Check if we just returned from Domain Connect
    if (window.location.search.includes('connected=true')) {
      this.justConnected.set(true);
      this.verifyDns();
      // Clean up the URL
      window.history.replaceState({}, '', window.location.pathname);
    }
  }

  autoConfigure(): void {
    const url = this.domainConnectURL();
    if (url) {
      window.location.href = url;
    }
  }

  goBack(): void {
    this.router.navigate(['/domains']);
  }

  verifyDns(): void {
    this.isVerifying.set(true);
    this.http
      .post<{ data: any }>(`/v1/domains/${this.domainId}/verify`, {})
      .subscribe({
        next: () => {
          this.isVerifying.set(false);
          this.loadDomain(); // Refresh all data
        },
        error: () => this.isVerifying.set(false),
      });
  }

  autoConfigureDns(): void {
    this.isConfiguringDns.set(true);
    this.cfConfigSuccess.set(false);
    this.http
      .post<{ data: any }>(`/v1/domains/${this.domainId}/auto-dns`, {})
      .subscribe({
        next: () => {
          this.isConfiguringDns.set(false);
          this.cfConfigSuccess.set(true);
          this.loadDomain(); // Refresh records
        },
        error: () => this.isConfiguringDns.set(false),
      });
  }

  saveSettings(): void {
    this.isSavingSettings.set(true);
    this.http
      .patch<{ data: any }>(`/v1/domains/${this.domainId}`, {
        open_tracking: this.openTracking(),
        click_tracking: this.clickTracking(),
      })
      .subscribe({
        next: () => {
          this.isSavingSettings.set(false);
          this.loadDomain();
        },
        error: () => this.isSavingSettings.set(false),
      });
  }

  confirmDelete(): void {
    this.showDeleteConfirm.set(true);
  }

  cancelDelete(): void {
    this.showDeleteConfirm.set(false);
  }

  deleteDomain(): void {
    this.isDeleting.set(true);
    this.http.delete(`/v1/domains/${this.domainId}`).subscribe({
      next: () => {
        this.isDeleting.set(false);
        this.router.navigate(['/domains']);
      },
      error: () => this.isDeleting.set(false),
    });
  }

  copyValue(value: string, key: string): void {
    navigator.clipboard.writeText(value).then(() => {
      this.copiedKey.set(key);
      setTimeout(() => this.copiedKey.set(null), 2000);
    });
  }

  toggleExpand(index: number): void {
    this.expandedIndex.update((current) => (current === index ? null : index));
  }

  statusBadgeClass(status: string): string {
    switch (status) {
      case 'verified': return 'badge badge--success';
      case 'pending': return 'badge badge--warning';
      case 'failed': return 'badge badge--error';
      default: return 'badge badge--neutral';
    }
  }

  recordStatusIcon(status: string): string {
    switch (status) {
      case 'verified': return 'check_circle';
      case 'pending': return 'schedule';
      case 'failed': return 'cancel';
      default: return 'help';
    }
  }

  recordStatusClass(status: string): string {
    switch (status) {
      case 'verified': return 'record-status record-status--verified';
      case 'pending': return 'record-status record-status--pending';
      case 'failed': return 'record-status record-status--failed';
      default: return 'record-status';
    }
  }
}
