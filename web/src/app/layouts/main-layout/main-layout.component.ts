import { Component, inject, signal } from '@angular/core';
import { RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';
import { AuthService } from '../../core/auth/auth.service';

@Component({
  selector: 'app-main-layout',
  standalone: true,
  imports: [RouterOutlet, RouterLink, RouterLinkActive],
  templateUrl: './main-layout.component.html',
  styleUrl: './main-layout.component.scss',
})
export class MainLayoutComponent {
  protected readonly authService = inject(AuthService);
  protected readonly userMenuOpen = signal(false);

  readonly navItems = [
    { path: '/dashboard', icon: 'dashboard', label: 'Dashboard' },
    { path: '/inbox', icon: 'inbox', label: 'Inbox' },
    { path: '/emails', icon: 'mail', label: 'Emails' },
    { path: '/domains', icon: 'dns', label: 'Domains' },
    { path: '/suppression', icon: 'block', label: 'Suppression' },
    { path: '/api-keys', icon: 'key', label: 'API Keys' },
    { path: '/webhooks', icon: 'webhook', label: 'Webhooks' },
    { path: '/templates', icon: 'draft', label: 'Templates' },
    { path: '/analytics', icon: 'analytics', label: 'Analytics' },
    { path: '/team', icon: 'group', label: 'Team' },
    { path: '/audit-logs', icon: 'history', label: 'Audit Logs' },
    { path: '/settings', icon: 'settings', label: 'Settings' },
  ];

  toggleUserMenu(): void {
    this.userMenuOpen.update((open) => !open);
  }

  closeUserMenu(): void {
    this.userMenuOpen.set(false);
  }

  signOut(): void {
    this.authService.logout();
  }
}
