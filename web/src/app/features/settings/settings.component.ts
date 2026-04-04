import { Component, signal } from '@angular/core';

@Component({
  selector: 'app-settings',
  standalone: true,
  imports: [],
  templateUrl: './settings.component.html',
  styleUrl: './settings.component.scss',
})
export class SettingsComponent {
  readonly notifyEmails = signal(true);
  readonly notifyUpdates = signal(false);
  readonly notifySecurity = signal(true);

  toggleNotifyEmails(): void {
    this.notifyEmails.update((v) => !v);
  }

  toggleNotifyUpdates(): void {
    this.notifyUpdates.update((v) => !v);
  }

  toggleNotifySecurity(): void {
    this.notifySecurity.update((v) => !v);
  }
}
