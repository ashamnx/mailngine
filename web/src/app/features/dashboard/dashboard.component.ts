import { Component, inject } from '@angular/core';
import { AuthService } from '../../core/auth/auth.service';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  templateUrl: './dashboard.component.html',
  styleUrl: './dashboard.component.scss',
})
export class DashboardComponent {
  protected readonly authService = inject(AuthService);

  readonly stats = [
    { icon: 'mail', title: 'Emails Sent', value: '0' },
    { icon: 'check_circle', title: 'Delivered', value: '--' },
    { icon: 'error', title: 'Bounce Rate', value: '--' },
    { icon: 'dns', title: 'Active Domains', value: '0' },
  ];
}
