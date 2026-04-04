import { Component, inject, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { ActivatedRoute, Router } from '@angular/router';
import { DatePipe } from '@angular/common';

interface EmailDetail {
  id: string;
  from: string;
  to: string;
  subject: string;
  status: string;
  message_id: string;
  created_at: string;
  sent_at: string | null;
  delivered_at: string | null;
}

@Component({
  selector: 'app-email-detail',
  standalone: true,
  imports: [DatePipe],
  templateUrl: './email-detail.component.html',
  styleUrl: './email-detail.component.scss',
})
export class EmailDetailComponent {
  private readonly http = inject(HttpClient);
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);

  readonly email = signal<EmailDetail | null>(null);
  readonly isLoading = signal(true);

  private emailId = '';

  constructor() {
    this.emailId = this.route.snapshot.params['id'];
    this.loadEmail();
  }

  loadEmail(): void {
    this.isLoading.set(true);
    this.http
      .get<{ data: EmailDetail }>(`/v1/emails/${this.emailId}`)
      .subscribe({
        next: (res) => {
          this.email.set(res.data);
          this.isLoading.set(false);
        },
        error: () => this.isLoading.set(false),
      });
  }

  goBack(): void {
    this.router.navigate(['/emails']);
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
