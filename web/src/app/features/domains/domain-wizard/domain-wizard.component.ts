import { Component, EventEmitter, inject, Output, signal, computed } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { FormsModule } from '@angular/forms';

interface Provider {
  name: string;
  type: string;
}

interface Recommendation {
  type: string;
  title: string;
  message: string;
  action?: string;
}

interface DNSAnalysis {
  domain: string;
  has_mx: boolean;
  mx_records: string[];
  detected_provider?: Provider;
  has_spf: boolean;
  existing_spf?: string;
  merged_spf: string;
  has_dmarc: boolean;
  existing_dmarc?: string;
  has_existing_dkim: boolean;
  recommendations: Recommendation[];
}

@Component({
  selector: 'app-domain-wizard',
  standalone: true,
  imports: [FormsModule],
  templateUrl: './domain-wizard.component.html',
  styleUrl: './domain-wizard.component.scss',
})
export class DomainWizardComponent {
  private readonly http = inject(HttpClient);
  private readonly router = inject(Router);

  @Output() cancelled = new EventEmitter<void>();

  // Step 1: Enter domain
  readonly step = signal<'input' | 'analysis' | 'confirm'>('input');
  readonly domainName = signal('');
  readonly isAnalyzing = signal(false);
  readonly analyzeError = signal<string | null>(null);

  // Step 2: Analysis results
  readonly analysis = signal<DNSAnalysis | null>(null);
  readonly enableInbound = signal(false);

  // Step 3: Creating
  readonly isCreating = signal(false);
  readonly createError = signal<string | null>(null);

  // Derived
  readonly hasProviderConflict = computed(() => {
    const a = this.analysis();
    return a?.detected_provider != null || a?.has_mx === true;
  });

  readonly providerName = computed(() => {
    return this.analysis()?.detected_provider?.name ?? 'an existing email provider';
  });

  readonly suggestedSubdomain = computed(() => {
    const name = this.domainName();
    if (!name) return '';
    // If already a subdomain, don't suggest another
    if (name.split('.').length > 2) return '';
    return 'mail.' + name;
  });

  analyze(): void {
    const name = this.domainName().trim().toLowerCase();
    if (!name) return;

    this.isAnalyzing.set(true);
    this.analyzeError.set(null);

    this.http
      .post<{ data: DNSAnalysis }>('/v1/domains/analyze', { name })
      .subscribe({
        next: (res) => {
          this.analysis.set(res.data);
          this.isAnalyzing.set(false);
          this.step.set('analysis');

          // Auto-disable inbound if provider detected
          if (res.data.detected_provider || res.data.has_mx) {
            this.enableInbound.set(false);
          }
        },
        error: (err) => {
          this.isAnalyzing.set(false);
          this.analyzeError.set(err.error?.error?.message ?? 'Failed to analyze domain. Please check the domain name.');
        },
      });
  }

  useSuggestedSubdomain(): void {
    const sub = this.suggestedSubdomain();
    if (sub) {
      this.domainName.set(sub);
      this.step.set('input');
      this.analysis.set(null);
    }
  }

  proceedToConfirm(): void {
    this.step.set('confirm');
  }

  backToInput(): void {
    this.step.set('input');
    this.analysis.set(null);
  }

  backToAnalysis(): void {
    this.step.set('analysis');
  }

  createDomain(): void {
    const a = this.analysis();
    if (!a) return;

    this.isCreating.set(true);
    this.createError.set(null);

    this.http
      .post<{ data: { domain: { id: string } } }>('/v1/domains', {
        name: a.domain,
        enable_inbound: this.enableInbound(),
        skip_dmarc: a.has_dmarc,
        merged_spf: a.has_spf ? a.merged_spf : '',
      })
      .subscribe({
        next: (res) => {
          this.isCreating.set(false);
          this.router.navigate(['/domains', res.data.domain.id]);
        },
        error: (err) => {
          this.isCreating.set(false);
          this.createError.set(err.error?.error?.message ?? 'Failed to create domain.');
        },
      });
  }

  cancel(): void {
    this.cancelled.emit();
  }

  recIcon(type: string): string {
    switch (type) {
      case 'warning': return 'warning';
      case 'success': return 'check_circle';
      case 'info': return 'info';
      default: return 'info';
    }
  }

  recClass(type: string): string {
    return 'rec rec--' + type;
  }
}
