import { Component, inject, signal, computed } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { DatePipe } from '@angular/common';

interface Email {
  id: string;
  to: string;
  subject: string;
  status: string;
  sent_at: string | null;
  created_at: string;
}

interface PaginationMeta {
  total: number;
  page: number;
  per_page: number;
  total_pages: number;
}

@Component({
  selector: 'app-emails',
  standalone: true,
  imports: [DatePipe],
  templateUrl: './emails.component.html',
  styleUrl: './emails.component.scss',
})
export class EmailsComponent {
  private readonly http = inject(HttpClient);
  private readonly router = inject(Router);

  readonly emails = signal<Email[]>([]);
  readonly meta = signal<PaginationMeta | null>(null);
  readonly isLoading = signal(true);
  readonly currentPage = signal(1);
  readonly perPage = 20;

  readonly showingFrom = computed(() => {
    const m = this.meta();
    if (!m || m.total === 0) return 0;
    return (m.page - 1) * m.per_page + 1;
  });

  readonly showingTo = computed(() => {
    const m = this.meta();
    if (!m) return 0;
    return Math.min(m.page * m.per_page, m.total);
  });

  readonly hasPrev = computed(() => this.currentPage() > 1);

  readonly hasNext = computed(() => {
    const m = this.meta();
    return m ? this.currentPage() < m.total_pages : false;
  });

  constructor() {
    this.loadEmails();
  }

  loadEmails(): void {
    this.isLoading.set(true);
    const page = this.currentPage();
    this.http
      .get<{ data: Email[]; meta: PaginationMeta }>(
        `/v1/emails?page=${page}&per_page=${this.perPage}`,
      )
      .subscribe({
        next: (res) => {
          this.emails.set(res.data);
          this.meta.set(res.meta);
          this.isLoading.set(false);
        },
        error: () => this.isLoading.set(false),
      });
  }

  goToPage(page: number): void {
    this.currentPage.set(page);
    this.loadEmails();
  }

  prevPage(): void {
    if (this.hasPrev()) {
      this.goToPage(this.currentPage() - 1);
    }
  }

  nextPage(): void {
    if (this.hasNext()) {
      this.goToPage(this.currentPage() + 1);
    }
  }

  navigateToEmail(id: string): void {
    this.router.navigate(['/emails', id]);
  }

  statusBadgeClass(status: string): string {
    switch (status) {
      case 'delivered':
        return 'badge badge--success';
      case 'sent':
        return 'badge badge--info';
      case 'queued':
        return 'badge badge--neutral';
      case 'bounced':
      case 'failed':
        return 'badge badge--error';
      default:
        return 'badge badge--neutral';
    }
  }
}
