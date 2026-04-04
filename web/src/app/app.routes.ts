import { Routes } from '@angular/router';
import { authGuard } from './core/guards/auth.guard';

export const routes: Routes = [
  {
    path: 'login',
    loadComponent: () =>
      import('./features/auth/login/login.component').then(
        (m) => m.LoginComponent,
      ),
  },
  {
    path: 'auth/callback',
    loadComponent: () =>
      import('./features/auth/callback/callback.component').then(
        (m) => m.CallbackComponent,
      ),
  },
  {
    path: '',
    loadComponent: () =>
      import('./layouts/main-layout/main-layout.component').then(
        (m) => m.MainLayoutComponent,
      ),
    canActivate: [authGuard],
    children: [
      { path: '', redirectTo: 'dashboard', pathMatch: 'full' },
      {
        path: 'dashboard',
        loadComponent: () =>
          import('./features/dashboard/dashboard.component').then(
            (m) => m.DashboardComponent,
          ),
      },
      {
        path: 'api-keys',
        loadComponent: () =>
          import('./features/api-keys/api-keys.component').then(
            (m) => m.ApiKeysComponent,
          ),
      },
      {
        path: 'inbox',
        loadComponent: () =>
          import('./features/inbox/inbox.component').then(
            (m) => m.InboxComponent,
          ),
      },
      {
        path: 'emails',
        loadComponent: () =>
          import('./features/emails/emails.component').then(
            (m) => m.EmailsComponent,
          ),
      },
      {
        path: 'emails/:id',
        loadComponent: () =>
          import(
            './features/emails/email-detail/email-detail.component'
          ).then((m) => m.EmailDetailComponent),
      },
      {
        path: 'domains',
        loadComponent: () =>
          import('./features/domains/domains.component').then(
            (m) => m.DomainsComponent,
          ),
      },
      {
        path: 'domains/:id',
        loadComponent: () =>
          import(
            './features/domains/domain-detail/domain-detail.component'
          ).then((m) => m.DomainDetailComponent),
      },
      {
        path: 'suppression',
        loadComponent: () =>
          import('./features/suppression/suppression.component').then(
            (m) => m.SuppressionComponent,
          ),
      },
      {
        path: 'webhooks',
        loadComponent: () =>
          import('./features/webhooks/webhooks.component').then(
            (m) => m.WebhooksComponent,
          ),
      },
      {
        path: 'webhooks/:id',
        loadComponent: () =>
          import(
            './features/webhooks/webhook-detail/webhook-detail.component'
          ).then((m) => m.WebhookDetailComponent),
      },
      {
        path: 'templates',
        loadComponent: () =>
          import('./features/templates/templates.component').then(
            (m) => m.TemplatesComponent,
          ),
      },
      {
        path: 'templates/:id',
        loadComponent: () =>
          import(
            './features/templates/template-detail/template-detail.component'
          ).then((m) => m.TemplateDetailComponent),
      },
      {
        path: 'analytics',
        loadComponent: () =>
          import('./features/analytics/analytics.component').then(
            (m) => m.AnalyticsComponent,
          ),
      },
      {
        path: 'billing',
        loadComponent: () =>
          import('./features/billing/billing.component').then(
            (m) => m.BillingComponent,
          ),
      },
      {
        path: 'team',
        loadComponent: () =>
          import('./features/team/team.component').then(
            (m) => m.TeamComponent,
          ),
      },
      {
        path: 'audit-logs',
        loadComponent: () =>
          import('./features/audit-logs/audit-logs.component').then(
            (m) => m.AuditLogsComponent,
          ),
      },
      {
        path: 'settings',
        loadComponent: () =>
          import('./features/settings/settings.component').then(
            (m) => m.SettingsComponent,
          ),
      },
    ],
  },
  { path: '**', redirectTo: '/' },
];
