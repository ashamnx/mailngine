import { Component, inject, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { DatePipe } from '@angular/common';
import { DomainWizardComponent } from './domain-wizard/domain-wizard.component';

interface Domain {
  id: string;
  name: string;
  status: string;
  created_at: string;
}

@Component({
  selector: 'app-domains',
  standalone: true,
  imports: [FormsModule, DatePipe, DomainWizardComponent],
  templateUrl: './domains.component.html',
  styleUrl: './domains.component.scss',
})
export class DomainsComponent {
  private readonly http = inject(HttpClient);
  private readonly router = inject(Router);

  readonly domains = signal<Domain[]>([]);
  readonly isLoading = signal(true);
  readonly showAddForm = signal(false);

  constructor() {
    this.loadDomains();
  }

  loadDomains(): void {
    this.isLoading.set(true);
    this.http
      .get<{ data: Domain[] }>('/v1/domains')
      .subscribe({
        next: (res) => {
          this.domains.set(res.data);
          this.isLoading.set(false);
        },
        error: () => this.isLoading.set(false),
      });
  }

  openAddForm(): void {
    this.showAddForm.set(true);
  }

  closeAddForm(): void {
    this.showAddForm.set(false);
  }

  navigateToDomain(id: string): void {
    this.router.navigate(['/domains', id]);
  }

  statusBadgeClass(status: string): string {
    switch (status) {
      case 'verified':
        return 'badge badge--success';
      case 'pending':
        return 'badge badge--warning';
      case 'failed':
        return 'badge badge--error';
      default:
        return 'badge badge--neutral';
    }
  }
}
