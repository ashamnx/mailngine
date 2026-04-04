import { Component, inject, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { FormsModule } from '@angular/forms';
import { DatePipe } from '@angular/common';

interface Org {
  id: string;
  name: string;
}

interface Member {
  id: string;
  name: string;
  email: string;
  role: string;
  avatar_url: string | null;
  joined_at: string;
}

@Component({
  selector: 'app-team',
  standalone: true,
  imports: [FormsModule, DatePipe],
  templateUrl: './team.component.html',
  styleUrl: './team.component.scss',
})
export class TeamComponent {
  private readonly http = inject(HttpClient);

  readonly org = signal<Org | null>(null);
  readonly members = signal<Member[]>([]);
  readonly isLoading = signal(true);

  // Org edit
  readonly editOrgName = signal('');
  readonly isSavingOrg = signal(false);
  readonly orgSaveSuccess = signal(false);

  // Invite
  readonly showInviteForm = signal(false);
  readonly inviteEmail = signal('');
  readonly inviteRole = signal('member');
  readonly isInviting = signal(false);
  readonly inviteError = signal<string | null>(null);
  readonly inviteSuccess = signal(false);

  // Role change
  readonly editingRoleId = signal<string | null>(null);
  readonly editingRoleValue = signal('');

  // Delete confirmation
  readonly deletingMemberId = signal<string | null>(null);

  readonly roles = ['member', 'admin', 'viewer'];

  constructor() {
    this.loadData();
  }

  loadData(): void {
    this.isLoading.set(true);

    this.http.get<{ data: Org }>('/v1/org').subscribe({
      next: (res) => {
        this.org.set(res.data);
        this.editOrgName.set(res.data.name);
      },
      error: () => {},
    });

    this.http.get<{ data: Member[] }>('/v1/org/members').subscribe({
      next: (res) => {
        this.members.set(res.data);
        this.isLoading.set(false);
      },
      error: () => this.isLoading.set(false),
    });
  }

  saveOrgName(): void {
    const name = this.editOrgName().trim();
    if (!name) return;

    this.isSavingOrg.set(true);
    this.orgSaveSuccess.set(false);

    this.http.patch<{ data: Org }>('/v1/org', { name }).subscribe({
      next: (res) => {
        this.org.set(res.data);
        this.isSavingOrg.set(false);
        this.orgSaveSuccess.set(true);
        setTimeout(() => this.orgSaveSuccess.set(false), 3000);
      },
      error: () => this.isSavingOrg.set(false),
    });
  }

  openInviteForm(): void {
    this.showInviteForm.set(true);
    this.inviteEmail.set('');
    this.inviteRole.set('member');
    this.inviteError.set(null);
    this.inviteSuccess.set(false);
  }

  closeInviteForm(): void {
    this.showInviteForm.set(false);
    this.inviteEmail.set('');
    this.inviteError.set(null);
  }

  inviteMember(): void {
    const email = this.inviteEmail().trim();
    if (!email) return;

    this.isInviting.set(true);
    this.inviteError.set(null);
    this.inviteSuccess.set(false);

    this.http
      .post<{ data: Member }>('/v1/org/members/invite', {
        email,
        role: this.inviteRole(),
      })
      .subscribe({
        next: () => {
          this.isInviting.set(false);
          this.inviteSuccess.set(true);
          this.inviteEmail.set('');
          this.loadData();
          setTimeout(() => this.inviteSuccess.set(false), 3000);
        },
        error: (err) => {
          this.isInviting.set(false);
          this.inviteError.set(
            err.error?.message ?? 'Failed to invite member. Please try again.',
          );
        },
      });
  }

  startEditRole(member: Member): void {
    this.editingRoleId.set(member.id);
    this.editingRoleValue.set(member.role);
  }

  cancelEditRole(): void {
    this.editingRoleId.set(null);
    this.editingRoleValue.set('');
  }

  saveRole(memberId: string): void {
    this.http
      .patch<{ data: Member }>(`/v1/org/members/${memberId}`, {
        role: this.editingRoleValue(),
      })
      .subscribe({
        next: () => {
          this.editingRoleId.set(null);
          this.loadData();
        },
        error: () => this.editingRoleId.set(null),
      });
  }

  confirmRemove(memberId: string): void {
    this.deletingMemberId.set(memberId);
  }

  cancelRemove(): void {
    this.deletingMemberId.set(null);
  }

  removeMember(memberId: string): void {
    this.http.delete(`/v1/org/members/${memberId}`).subscribe({
      next: () => {
        this.deletingMemberId.set(null);
        this.loadData();
      },
      error: () => this.deletingMemberId.set(null),
    });
  }

  roleBadgeClass(role: string): string {
    switch (role) {
      case 'owner':
        return 'badge badge--info';
      case 'admin':
        return 'badge badge--info';
      case 'member':
        return 'badge badge--neutral';
      case 'viewer':
        return 'badge badge--neutral';
      default:
        return 'badge badge--neutral';
    }
  }

  memberInitial(name: string): string {
    return name ? name.charAt(0).toUpperCase() : '?';
  }
}
