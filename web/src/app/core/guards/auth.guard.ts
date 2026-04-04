import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { map, catchError, of } from 'rxjs';
import { AuthService } from '../auth/auth.service';

export const authGuard: CanActivateFn = () => {
  const authService = inject(AuthService);
  const router = inject(Router);

  if (authService.isAuthenticated()) {
    return true;
  }

  if (authService.token) {
    return authService.loadMe().pipe(
      map(() => true),
      catchError(() => {
        router.navigate(['/login']);
        return of(false);
      }),
    );
  }

  router.navigate(['/login']);
  return false;
};
