import { Component, OnInit, inject } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { AuthService } from '../../../core/auth/auth.service';

@Component({
  selector: 'app-callback',
  standalone: true,
  templateUrl: './callback.component.html',
  styleUrl: './callback.component.scss',
})
export class CallbackComponent implements OnInit {
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly authService = inject(AuthService);

  ngOnInit(): void {
    // Token is passed in the URL fragment (#token=...) to prevent it from
    // appearing in Referer headers, server logs, or proxy logs.
    const fragment = this.route.snapshot.fragment || '';
    const params = new URLSearchParams(fragment);
    const token = params.get('token');

    if (token) {
      this.authService.handleCallback(token);
    } else {
      this.router.navigate(['/login']);
    }
  }
}
