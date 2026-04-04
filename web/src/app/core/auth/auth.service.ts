import { Injectable, computed, inject, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { Observable, tap } from 'rxjs';
import { MeResponse, Organization, OrgListItem, User } from './auth.models';

@Injectable({ providedIn: 'root' })
export class AuthService {
  private readonly TOKEN_KEY = 'hm_token';
  private readonly http = inject(HttpClient);
  private readonly router = inject(Router);

  private readonly _user = signal<User | null>(null);
  private readonly _organization = signal<Organization | null>(null);
  private readonly _role = signal<string>('');
  private readonly _organizations = signal<OrgListItem[]>([]);
  private readonly _isLoading = signal<boolean>(true);

  readonly user = this._user.asReadonly();
  readonly organization = this._organization.asReadonly();
  readonly role = this._role.asReadonly();
  readonly organizations = this._organizations.asReadonly();
  readonly isAuthenticated = computed(() => !!this._user());
  readonly isLoading = this._isLoading.asReadonly();

  get token(): string | null {
    return localStorage.getItem(this.TOKEN_KEY);
  }

  loginWithGoogle(): void {
    window.location.href = '/v1/auth/google';
  }

  handleCallback(token: string): void {
    localStorage.setItem(this.TOKEN_KEY, token);
    this.loadMe().subscribe({
      next: () => this.router.navigate(['/dashboard']),
      error: () => {
        this.clearState();
        this.router.navigate(['/login']);
      },
    });
  }

  loadMe(): Observable<{ data: MeResponse }> {
    return this.http.get<{ data: MeResponse }>('/v1/auth/me').pipe(
      tap((res) => {
        const me = res.data;
        this._user.set(me.user);
        this._organization.set(me.organization);
        this._role.set(me.role);
        this._organizations.set(me.organizations);
        this._isLoading.set(false);
      }),
    );
  }

  logout(): void {
    this.http.post('/v1/auth/logout', {}).subscribe({
      complete: () => {
        this.clearState();
        this.router.navigate(['/login']);
      },
      error: () => {
        this.clearState();
        this.router.navigate(['/login']);
      },
    });
  }

  private clearState(): void {
    localStorage.removeItem(this.TOKEN_KEY);
    this._user.set(null);
    this._organization.set(null);
    this._role.set('');
    this._organizations.set([]);
    this._isLoading.set(false);
  }
}
