import { Component, inject, signal, computed } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { DecimalPipe } from '@angular/common';

interface Plan {
  name: string;
  monthly_limit: number;
}

interface Usage {
  emails_sent: number;
  monthly_limit: number;
}

interface UsageHistoryEntry {
  month: string;
  sent: number;
  received: number;
  api_calls: number;
}

interface PlanTier {
  name: string;
  price: string;
  emails: string;
  domains: string;
  features: string[];
}

@Component({
  selector: 'app-billing',
  standalone: true,
  imports: [DecimalPipe],
  templateUrl: './billing.component.html',
  styleUrl: './billing.component.scss',
})
export class BillingComponent {
  private readonly http = inject(HttpClient);

  readonly plan = signal<Plan | null>(null);
  readonly usage = signal<Usage | null>(null);
  readonly usageHistory = signal<UsageHistoryEntry[]>([]);
  readonly isLoading = signal(true);

  readonly usagePercent = computed(() => {
    const u = this.usage();
    if (!u || u.monthly_limit === 0) return 0;
    return Math.min(Math.round((u.emails_sent / u.monthly_limit) * 100), 100);
  });

  readonly planTiers: PlanTier[] = [
    {
      name: 'Free',
      price: '$0',
      emails: '100/mo',
      domains: '1 domain',
      features: ['Basic analytics', 'Email templates', 'Community support'],
    },
    {
      name: 'Starter',
      price: '$20',
      emails: '10,000/mo',
      domains: '5 domains',
      features: ['Advanced analytics', 'Custom templates', 'Email support', 'Webhooks'],
    },
    {
      name: 'Pro',
      price: '$80',
      emails: '100,000/mo',
      domains: 'Unlimited',
      features: ['Full analytics', 'Priority support', 'Dedicated IP', 'Team management', 'Audit logs'],
    },
    {
      name: 'Enterprise',
      price: 'Custom',
      emails: 'Custom',
      domains: 'Unlimited',
      features: ['Everything in Pro', 'Custom SLA', 'Dedicated account manager', 'SSO/SAML', 'Custom integrations'],
    },
  ];

  constructor() {
    this.loadData();
  }

  loadData(): void {
    this.isLoading.set(true);

    this.http.get<{ data: Plan }>('/v1/billing/plan').subscribe({
      next: (res) => this.plan.set(res.data),
      error: () => {},
    });

    this.http.get<{ data: Usage }>('/v1/billing/usage').subscribe({
      next: (res) => {
        this.usage.set(res.data);
        this.isLoading.set(false);
      },
      error: () => this.isLoading.set(false),
    });

    this.http.get<{ data: UsageHistoryEntry[] }>('/v1/billing/usage/history').subscribe({
      next: (res) => this.usageHistory.set(res.data),
      error: () => {},
    });
  }

  isCurrentPlan(tierName: string): boolean {
    return this.plan()?.name?.toLowerCase() === tierName.toLowerCase();
  }

  usageBarClass(): string {
    const pct = this.usagePercent();
    if (pct >= 90) return 'usage-fill usage-fill--critical';
    if (pct >= 75) return 'usage-fill usage-fill--warning';
    return 'usage-fill';
  }
}
